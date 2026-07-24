package funcs

import (
	"fmt"

	"github.com/Vinicamilotti/archlinux_terraformer/repository"
)

type RemoveFromTerraformFile struct {
	Repository repository.TerraformFileRepository
}

func NewRemoveFromTerraformFile(repo repository.TerraformFileRepository) RemoveFromTerraformFile {
	return RemoveFromTerraformFile{
		Repository: repo,
	}
}

func (r RemoveFromTerraformFile) Exec(args ...[]string) error {
	if len(args) == 0 || len(args[0]) == 0 {
		return fmt.Errorf("missing arguments: expected \"pacman <package>\" or \"custom <name>\"")
	}

	params := args[0]
	target := params[0]

	switch target {
	case "pacman":
		if len(params) < 2 {
			return fmt.Errorf("missing package name for pacman")
		}

		return r.Repository.RemoveFromPacman(params[1])
	case "custom":
		if len(params) < 2 {
			return fmt.Errorf("missing name for custom")
		}

		return r.Repository.RemoveFromCustom(params[1])
	default:
		return fmt.Errorf("unknown target %q: expected \"pacman\" or \"custom\"", target)
	}
}
