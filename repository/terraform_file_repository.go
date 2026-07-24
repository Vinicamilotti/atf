package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Vinicamilotti/archlinux_terraformer/models"
)

type TerraformFileRepository interface {
	ReadFromFile() (models.ConfigFile, error)
	AddToPacman(pkg string) error
	AddToCustom(name string, instructions []string) error
	RemoveFromPacman(pkg string) error
	RemoveFromCustom(name string) error
	MarkCustomInstalled(name string) error
	EnsureFileExists() error
}

type TerraformFileRepositoryImpl struct {
	fileLocation string
}

func NewTerraformFileRepository(fileLocation string) TerraformFileRepository {
	return TerraformFileRepositoryImpl{
		fileLocation: fileLocation,
	}
}

func (r TerraformFileRepositoryImpl) ReadFromFile() (models.ConfigFile, error) {
	config := models.ConfigFile{
		PacmanInstallations: []string{},
		CustomInstalations:  map[string]models.CustomInstallation{},
	}

	data, err := os.ReadFile(r.fileLocation)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil
		}
		return config, err
	}

	if len(data) == 0 {
		return config, nil
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, err
	}

	return config, nil
}

func (r TerraformFileRepositoryImpl) writeToFile(config models.ConfigFile) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.fileLocation, data, 0644)
}

func (r TerraformFileRepositoryImpl) AddToPacman(pkg string) error {
	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	for _, existing := range config.PacmanInstallations {
		if existing == pkg {
			return fmt.Errorf("package %q is already present", pkg)
		}
	}

	config.PacmanInstallations = append(config.PacmanInstallations, pkg)

	return r.writeToFile(config)
}

func (r TerraformFileRepositoryImpl) AddToCustom(name string, instructions []string) error {
	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	if config.CustomInstalations == nil {
		config.CustomInstalations = map[string]models.CustomInstallation{}
	}

	config.CustomInstalations[name] = models.CustomInstallation{
		Instructions: instructions,
		Installed:    false,
	}

	return r.writeToFile(config)
}

func (r TerraformFileRepositoryImpl) RemoveFromPacman(pkg string) error {
	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	index := -1
	for i, existing := range config.PacmanInstallations {
		if existing == pkg {
			index = i
			break
		}
	}

	if index == -1 {
		return fmt.Errorf("package %q not found", pkg)
	}

	config.PacmanInstallations = append(config.PacmanInstallations[:index], config.PacmanInstallations[index+1:]...)

	return r.writeToFile(config)
}

func (r TerraformFileRepositoryImpl) RemoveFromCustom(name string) error {
	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	if _, ok := config.CustomInstalations[name]; !ok {
		return fmt.Errorf("custom instruction %q not found", name)
	}

	delete(config.CustomInstalations, name)

	return r.writeToFile(config)
}

func (r TerraformFileRepositoryImpl) MarkCustomInstalled(name string) error {
	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	installation, ok := config.CustomInstalations[name]
	if !ok {
		return fmt.Errorf("custom instruction %q not found", name)
	}

	installation.Installed = true
	config.CustomInstalations[name] = installation

	return r.writeToFile(config)
}

func (r TerraformFileRepositoryImpl) EnsureFileExists() error {
	if _, err := os.Stat(r.fileLocation); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}

	config, err := r.ReadFromFile()
	if err != nil {
		return err
	}

	return r.writeToFile(config)
}
