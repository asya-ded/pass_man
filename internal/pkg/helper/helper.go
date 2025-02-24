package helper

import (
	"bufio"
	"fmt"
	"os"
	"pass_man/internal/pkg/model"
)

func GetMasterPassword() []byte {
	fmt.Println()
	fmt.Print("Enter master password:")
	reader := bufio.NewReader(os.Stdin)

	mp, err := reader.ReadString('\n')
	if err != nil {
		panic("cannot read master password")
	}

	if len(mp) == 0 {
		panic("master password is empty")
	}

	return []byte(mp)
}

func SaveToFile(fileName string, content []byte) error {
	file, err := os.Create(model.Path + fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write(content)
	return nil
}

func GetEntry() *model.Entry {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print("name: ")
	name, _ := reader.ReadString('\n')
	fmt.Print("service: ")
	service, _ := reader.ReadString('\n')
	fmt.Print("login: ")
	login, _ := reader.ReadString('\n')
	fmt.Print("password:")
	pw, _ := reader.ReadString('\n')
	fmt.Print("notes: ")
	notes, _ := reader.ReadString('\n')

	return &model.Entry{
		Name:     name,
		Login:    login,
		Password: pw,
		Service:  service,
		Notes:    notes,
	}
}

func GetEntryName() string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println()
	fmt.Print("name: ")
	name, _ := reader.ReadString('\n')

	return name
}
