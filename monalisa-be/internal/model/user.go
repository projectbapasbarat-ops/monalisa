package model

type UserWithRoles struct {
	ID      string   `json:"id"`
	NIP     string   `json:"nip"`
	Nama    string   `json:"nama"`
	Jabatan string   `json:"jabatan"`
	Roles   []string `json:"roles"`
}
