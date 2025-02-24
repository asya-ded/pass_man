package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"pass_man/internal/pkg/helper"
	"pass_man/internal/pkg/model"

	"golang.org/x/crypto/argon2"
)

type Service struct{}

func New() *Service {
	return &Service{}
}

func (s *Service) Init(masterPassword []byte) error {
	if _, err := os.Stat(model.Path); os.IsNotExist(err) {
		err = os.MkdirAll(model.Path, os.ModePerm)
		if err != nil {
			panic("cannot create storage directory")
		}
	}

	vault := &model.Vault{
		KDF: &model.KDF{
			Salt:        s.generateSalt(16),
			Iterations:  10,
			Memory:      256 * 1024,
			Parallelism: 4,
		},
		Entries: []*model.Entry{},
	}

	kek := s.GenerateKEK(masterPassword, vault.KDF)

	dek, err := s.GenerateDEK()
	if err != nil {
		panic(err.Error())
	}

	encryptedDEK, err := s.EncryptDEK(kek, dek)
	if err != nil {
		panic(err.Error())
	}

	vault.EncryptedDEK = encryptedDEK

	vaultJson, _ := json.Marshal(vault)

	helper.SaveToFile("vault.json", []byte(vaultJson))

	return nil
}

func (s *Service) GenerateKEK(masterPassword []byte, kdf *model.KDF) []byte {
	salt, _ := base64.StdEncoding.DecodeString(kdf.Salt)
	return argon2.IDKey(
		masterPassword,
		salt,
		kdf.Iterations,
		kdf.Memory,
		kdf.Parallelism,
		32,
	)
}

func (s *Service) generateSalt(size int) string {
	salt := make([]byte, size)
	rand.Read(salt)
	return base64.RawStdEncoding.EncodeToString(salt)
}

func (s *Service) GenerateDEK() ([]byte, error) {
	dek := make([]byte, 32)
	if _, err := rand.Read(dek); err != nil {
		return nil, err
	}
	return dek, nil
}

func (s *Service) EncryptDEK(kek, dek []byte) (string, error) {
	block, _ := aes.NewCipher(kek)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	encryptedDEK := gcm.Seal(nonce, nonce, dek, nil)

	return base64.StdEncoding.EncodeToString(encryptedDEK), nil
}

func (s *Service) DecryptDEK(kek []byte, encryptedDEK string) ([]byte, error) {
	block, _ := aes.NewCipher(kek)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	dek, _ := base64.StdEncoding.DecodeString(encryptedDEK)
	nonce, cipherText := dek[:nonceSize], dek[nonceSize:]

	return gcm.Open(nil, nonce, cipherText, nil)
}

func (s *Service) EncryptEntity(dek []byte, entry *model.Entry) *model.Entry {
	encrypedEntry := &model.Entry{
		Name: entry.Name,
	}

	block, err := aes.NewCipher(dek)
	if err != nil {
		fmt.Println(err.Error())
	}
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	encrypedEntry.Login = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(entry.Login), nil))
	encrypedEntry.Password = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(entry.Password), nil))
	encrypedEntry.Service = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(entry.Service), nil))
	encrypedEntry.Notes = base64.StdEncoding.EncodeToString(gcm.Seal(nonce, nonce, []byte(entry.Notes), nil))

	return encrypedEntry
}

func (s *Service) DecryptEntity(dek []byte, entry *model.Entry) *model.Entry {
	decrypedEntry := &model.Entry{
		Name: entry.Name,
	}

	block, _ := aes.NewCipher(dek)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()

	login, _ := base64.StdEncoding.DecodeString(entry.Login)
	nonce, ciphertext := login[:nonceSize], login[nonceSize:]
	decrypted, _ := gcm.Open(nil, nonce, ciphertext, nil)

	decrypedEntry.Login = string(decrypted)

	password, _ := base64.StdEncoding.DecodeString(entry.Password)
	nonce, ciphertext = password[:nonceSize], password[nonceSize:]
	decrypted, _ = gcm.Open(nil, nonce, ciphertext, nil)
	decrypedEntry.Password = string(decrypted)

	service, _ := base64.StdEncoding.DecodeString(entry.Service)
	nonce, ciphertext = service[:nonceSize], service[nonceSize:]
	decrypted, _ = gcm.Open(nil, nonce, ciphertext, nil)
	decrypedEntry.Service = string(decrypted)

	notes, _ := base64.StdEncoding.DecodeString(entry.Notes)
	nonce, ciphertext = notes[:nonceSize], notes[nonceSize:]
	decrypted, _ = gcm.Open(nil, nonce, ciphertext, nil)
	decrypedEntry.Notes = string(decrypted)

	return decrypedEntry
}
