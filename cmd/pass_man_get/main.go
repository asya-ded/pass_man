package main

import (
	"fmt"
	"pass_man/internal/pkg/helper"
	"pass_man/internal/pkg/service/crypto"
	"pass_man/internal/pkg/service/vault"
)

func main() {
	mp := helper.GetMasterPassword()

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

	name := helper.GetEntryName()
	entry := vaultService.GetEntryByName(vault, name)

	if entry != nil {
		decryptedEntry := cryptoService.DecryptEntity(dek, entry)
		fmt.Printf("%+v\n", decryptedEntry)
	}
}
