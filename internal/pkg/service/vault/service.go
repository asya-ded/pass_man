package vault

import (
	"encoding/json"
	"os"
	"pass_man/internal/pkg/model"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Getvault() (*model.Vault, error) {
	data, _ := os.ReadFile(model.Path + model.VautlName)

	var vault model.Vault
	json.Unmarshal(data, &vault)

	return &vault, nil
}

func (s *Service) SaveVault(vault *model.Vault) error {
	data, _ := json.Marshal(vault)

	err := os.WriteFile(model.Path+model.VautlName, data, os.ModePerm)
	return err
}

func (s *Service) GetEntryByName(vault *model.Vault, name string) *model.Entry {
	for _, e := range vault.Entries {
		if e.Name == name {
			return e
		}
	}
	return nil
}
