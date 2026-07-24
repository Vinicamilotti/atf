package funcs

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Vinicamilotti/archlinux_terraformer/models"
	"github.com/Vinicamilotti/archlinux_terraformer/repository"
)

const terraformLogFile = "terraform.log"

type TerraformCommand struct {
	Repository repository.TerraformFileRepository
}

func NewTerraformCommand(repo repository.TerraformFileRepository) TerraformCommand {
	return TerraformCommand{
		Repository: repo,
	}
}

func (t TerraformCommand) Exec(args ...[]string) error {
	target := ""
	if len(args) > 0 && len(args[0]) > 0 {
		target = args[0][0]
	}

	if target != "" && target != "pacman" && target != "custom" {
		return fmt.Errorf("unknown target %q: expected \"pacman\" or \"custom\"", target)
	}

	config, err := t.Repository.ReadFromFile()
	if err != nil {
		return err
	}

	var errs []error

	if target == "" || target == "pacman" {
		if err := t.runPacman(config.PacmanInstallations); err != nil {
			errs = append(errs, err)
		}
	}

	if target == "" || target == "custom" {
		errs = append(errs, t.runCustom(config.CustomInstalations)...)
	}

	if len(errs) > 0 {
		t.warnAndLog(errs)
	}

	return nil
}

func (t TerraformCommand) runPacman(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	cmdArgs := append([]string{"-S"}, packages...)
	cmd := exec.Command("pacman", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("pacman -S %s: %w", strings.Join(packages, " "), err)
	}

	return nil
}

func (t TerraformCommand) runCustom(customInstallations map[string]models.CustomInstallation) []error {
	var errs []error

	for name, installation := range customInstallations {
		if installation.Installed {
			continue
		}

		failed := false
		for _, instruction := range installation.Instructions {
			cmd := exec.Command("sh", "-c", instruction)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Stdin = os.Stdin

			if err := cmd.Run(); err != nil {
				errs = append(errs, fmt.Errorf("custom %q: %q: %w", name, instruction, err))
				failed = true
			}
		}

		if failed {
			continue
		}

		if err := t.Repository.MarkCustomInstalled(name); err != nil {
			errs = append(errs, fmt.Errorf("custom %q: mark installed: %w", name, err))
		}
	}

	return errs
}

func (t TerraformCommand) warnAndLog(errs []error) {
	fmt.Fprintf(os.Stderr, "warning: %d error(s) occurred while terraforming, see %s\n", len(errs), terraformLogFile)

	file, err := os.OpenFile(terraformLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open log file %s: %v\n", terraformLogFile, err)
		return
	}
	defer file.Close()

	timestamp := time.Now().Format(time.RFC3339)
	for _, err := range errs {
		fmt.Fprintf(file, "[%s] %v\n", timestamp, err)
	}
}
