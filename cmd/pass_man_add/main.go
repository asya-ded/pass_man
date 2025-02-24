package main

import (
	"pass_man/internal/pkg/helper"
	"pass_man/internal/pkg/service/crypto"
	"pass_man/internal/pkg/service/vault"
)

func main() {
	mp := helper.GetMasterPassword()

	entry := helper.GetEntry()

	vaultService := vault.New()
	vault, err := vaultService.Getvault()
	if err != nil {
		panic(err.Error())
	}

	cryptoService := crypto.New()
	kek := cryptoService.GenerateKEK(mp, vault.KDF)
	dek, err := cryptoService.DecryptDEK(kek, vault.EncryptedDEK)
	if err != nil {
		panic(err.Error())
	}

	encryptedEntry := cryptoService.EncryptEntity(dek, entry)

	vault.Entries = append(vault.Entries, encryptedEntry)

	vaultService.SaveVault(vault)
}
