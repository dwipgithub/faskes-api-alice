package models

type HospitalFilter struct {
	Page            int
	Limit           int
	PelayananNama   string
	ProvinsiID      string
	KabKotaID       string
	Nama            string
	Aktive          string
	StartModifiedAt string
	EndModifiedAt   string
}
