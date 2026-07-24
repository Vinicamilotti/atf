package funcs

import (
	"fmt"

	"github.com/Vinicamilotti/archlinux_terraformer/repository"
)

type AddToTerraformFile struct {
	Repository repository.TerraformFileRepository
}

func NewAddToTerraformFile(repo repository.TerraformFileRepository) AddToTerraformFile {
	return AddToTerraformFile{
		Repository: repo,
	}
}

func (a AddToTerraformFile) Exec(args ...[]string) error {
	if len(args) == 0 || len(args[0]) == 0 {
		return fmt.Errorf("missing arguments: expected \"pacman <package>\" or \"custom <name> <instructions...>\"")
	}

	params := args[0]
	target := params[0]

	switch target {
	case "pacman":
		if len(params) < 2 {
			return fmt.Errorf("missing package name for pacman")
		}

		return a.Repository.AddToPacman(params[1])
	case "custom":
		if len(params) < 3 {
			return fmt.Errorf("missing name or instructions for custom")
		}

		name := params[1]
		instructions := params[2:]

		return a.Repository.AddToCustom(name, instructions)
	default:
		return fmt.Errorf("unknown target %q: expected \"pacman\" or \"custom\"", target)
	}
}
