package models

type Hospital struct {
	Kode                        *string    `json:"kode"`
	KodeBPJS                    *string    `json:"kodeBpjs"`
	Nama                        *string    `json:"nama"`
	Jenis                       *string    `json:"jenis"`
	Kelas                       *string    `json:"kelas"`
	Telepon                     *string    `json:"telepon"`
	Email                       *string    `json:"email"`
	Website                     *string    `json:"website"`
	StatusBLU                   *string    `json:"statusBLU"`
	NoSuratIjinOperasional      *string    `json:"noSuratIjinOperasional"`
	TanggalSuratIjinOperasional *string    `json:"tanggalSuratIjinOperasional"`
	TMO                         *string    `json:"tmo"`
	Direktur                    *string    `json:"direktur"`
	KetersediaanSIMRS           *string    `json:"ketersediaanSIMRS"`
	AksesInternet               *string    `json:"aksesInternet"`
	LuasTanah                   *string    `json:"luasTanah"`
	LuasBangunan                *string    `json:"luasBangunan"`
	Kepemilikan                 *string    `json:"kepemilikan"`
	URLTarif                    *string    `json:"urlTarif"`
	StatusValidasiTarif         *string    `json:"statusValidasiTarif"`
	Alamat                      *string    `json:"alamat"`
	ProvinsiID                  *string    `json:"provinsi_id"`
	ProvinsiNama                *string    `json:"provinsiNama"`
	KabKotaID                   *string    `json:"kab_kota_id"`
	KabKotaNama                 *string    `json:"kabKotaNama"`
	KecamatanID                 *string    `json:"kecamatan_id"`
	KecamatanNama               *string    `json:"kecamatanNama"`
	KelurahanID                 *string    `json:"kelurahan_id"`
	KelurahanNama               *string    `json:"kelurahanNama"`
	Longitude                   *string    `json:"longitude"`
	Latitude                    *string    `json:"latitude"`
	NoSTRPJ                     *string    `json:"no_str_pj"`
	StatusAktivasi              *int       `json:"statusAktivasi"`
	URLFotoDepan                *string    `json:"urlFotoDepan"`
	ModifiedAt                  *string    `json:"modified_at"`
	SatuSehat                   *SatuSehat `json:"satuSehat"`
}

type SatuSehat struct {
	OrganizationID *string `json:"organizationId"`
	NamaPIC        *string `json:"namaPic"`
	EmailIntegrasi *string `json:"emailIntegrasi"`
	TelpPIC        *string `json:"telpPic"`
}
