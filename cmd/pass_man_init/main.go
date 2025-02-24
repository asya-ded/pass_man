package main

import (
	"pass_man/internal/pkg/helper"
	"pass_man/internal/pkg/service/crypto"
)

func main() {
	mp := helper.GetMasterPassword()

	service := crypto.New()
	service.Init(mp)
}
