package model

type Vault struct {
	KDF          *KDF     `json:"kdf"`
	EncryptedDEK string   `json:"encrypted_dek"`
	Entries      []*Entry `json:"entries"`
}

type KDF struct {
	Salt        string `json:"salt"`
	Iterations  uint32 `json:"iterations"`
	Memory      uint32 `json:"memory"`
	Parallelism uint8  `json:"parallelism"`
}

type Entry struct {
	Name     string `json:"name"`
	Service  string `json:"service"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Notes    string `json:"notes"`
}
