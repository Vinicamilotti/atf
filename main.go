package main

import (
	"fmt"
	"os"

	"github.com/Vinicamilotti/archlinux_terraformer/cli"
	"github.com/Vinicamilotti/archlinux_terraformer/funcs"
	"github.com/Vinicamilotti/archlinux_terraformer/repository"
)

const terraformFileLocation = "atf_file.json"

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: archlinux_terraformer <add|remove|terraform> [args...]")
		os.Exit(1)
	}

	repo := repository.NewTerraformFileRepository(terraformFileLocation)

	router := cli.CommandRouter{
		Commands: map[string]cli.CliFunc{},
	}

	router.AddCommand("add", funcs.NewAddToTerraformFile(repo))
	router.AddCommand("remove", funcs.NewRemoveFromTerraformFile(repo))
	router.AddCommand("terraform", funcs.NewTerraformCommand(repo))

	command, err := router.GetCommand(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := command.Exec(os.Args[2:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
