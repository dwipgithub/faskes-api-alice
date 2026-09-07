package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	hospital "go-api-learning/alice/models/hospital"
)

type HospitalRepository struct {
	DB *sql.DB
}

func (r *HospitalRepository) GetAll(
	filter hospital.HospitalFilter,
) ([]hospital.Hospital, int, error) {

	// ========================================================
	// Pagination
	// ========================================================

	if filter.Page < 1 {
		filter.Page = 1
	}

	if filter.Limit < 1 {
		filter.Limit = 100
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}

	offset := (filter.Page - 1) * filter.Limit

	// ========================================================
	// SELECT
	// ========================================================

	sqlSelect := `
		SELECT
			db_fasyankes.data.Propinsi AS kode,
			db_fasyankes.data.RUMAH_SAKIT AS nama,
			db_fasyankes.m_jenis.alias AS jenis,
			db_fasyankes.m_kelas.kelas AS kelas,
			db_fasyankes.data.TELEPON AS telepon,
			db_fasyankes.data.EMAIL AS email,
			db_fasyankes.data.WEBSITE AS website,
			db_fasyankes.m_blu.blu AS statusBLU,
			db_fasyankes.data.NO_SURAT_IJIN AS noSuratIjinOperasional,
			db_fasyankes.data.TANGGAL_SURAT_IJIN AS tanggalSuratIjinOperasional,
			db_fasyankes.data.DIREKTUR_RS AS direktur,
			db_fasyankes.m_simrs.simrs AS ketersediaanSIMRS,
			db_fasyankes.data.LUAS_TANAH AS luasTanah,
			db_fasyankes.data.LUAS_BANGUNAN AS luasBangunan,
			db_fasyankes.m_kepemilikan.kepemilikan AS kepemilikan,
			db_fasyankes.t_dok_tariflayanan_rs.url AS urlTarif,
			db_fasyankes.t_dok_tariflayanan_rs.status_validasi AS statusValidasiTarif,
			db_fasyankes.data.ALAMAT AS alamat,
			db_fasyankes.data.provinsi_id,
			db_fasyankes.provinsi.nama AS provinsiNama,
			db_fasyankes.data.kab_kota_id,
			db_fasyankes.kab_kota.nama AS kabKotaNama,
			db_fasyankes.koordinat.long AS longitude,
			db_fasyankes.koordinat.alt AS latitude,
			db_fasyankes.data.aktive AS statusAktivasi,
			derivedtable2.url AS urlFotoDepan,
			db_fasyankes.data.TANGGAL_UPDATE AS modified_at
	`

	// ========================================================
	// FROM + JOIN
	// ========================================================

	sqlFrom := `
		FROM
		(
			SELECT
				db_fasyankes.data.Propinsi AS faskesId
			FROM db_fasyankes.data

			LEFT OUTER JOIN db_fasyankes.t_pelayanan
				ON db_fasyankes.t_pelayanan.koders =
					db_fasyankes.data.Propinsi

			LEFT JOIN db_fasyankes.m_pelayanan
				ON db_fasyankes.m_pelayanan.kode_pelayanan =
					db_fasyankes.t_pelayanan.kode_pelayanan

			WHERE db_fasyankes.m_pelayanan.pelayanan LIKE ?

			GROUP BY db_fasyankes.data.Propinsi
		) derivedTable1

		INNER JOIN db_fasyankes.data
			ON db_fasyankes.data.Propinsi =
				derivedTable1.faskesId

		LEFT OUTER JOIN db_fasyankes.provinsi
			ON db_fasyankes.provinsi.id =
				db_fasyankes.data.provinsi_id

		LEFT OUTER JOIN db_fasyankes.kab_kota
			ON db_fasyankes.kab_kota.id =
				db_fasyankes.data.kab_kota_id

		LEFT OUTER JOIN db_fasyankes.m_jenis
			ON db_fasyankes.m_jenis.id_jenis =
				db_fasyankes.data.JENIS

		LEFT OUTER JOIN db_fasyankes.m_kelas
			ON db_fasyankes.m_kelas.id_kelas =
				db_fasyankes.data.KLS_RS

		LEFT OUTER JOIN db_fasyankes.m_kepemilikan
			ON db_fasyankes.m_kepemilikan.id_kepemilikan =
				db_fasyankes.data.PENYELENGGARA

		LEFT OUTER JOIN db_fasyankes.m_blu
			ON db_fasyankes.m_blu.id_blu =
				db_fasyankes.data.blu

		LEFT OUTER JOIN db_fasyankes.koordinat
			ON db_fasyankes.koordinat.koders =
				db_fasyankes.data.propinsi

		LEFT OUTER JOIN db_fasyankes.m_simrs
			ON db_fasyankes.m_simrs.id_simrs =
				db_fasyankes.data.simrs

		LEFT OUTER JOIN db_fasyankes.t_dok_tariflayanan_rs
			ON db_fasyankes.t_dok_tariflayanan_rs.koders =
				db_fasyankes.data.Propinsi

		LEFT OUTER JOIN
		(
			SELECT
				db_fasyankes.t_images.koders,
				db_fasyankes.t_images.url
			FROM db_fasyankes.t_images
			WHERE db_fasyankes.t_images.keterangan = "depan"
		) derivedtable2
			ON derivedtable2.koders =
				db_fasyankes.data.Propinsi
	`

	// ========================================================
	// WHERE
	// ========================================================

	where := []string{
		`db_fasyankes.data.Propinsi NOT IN ("9999999", "7371435", "7371121", "")`,
		`db_fasyankes.data.JENIS <> 20`,
	}

	// Parameter untuk derivedTable1
	args := []interface{}{
		"%" + filter.PelayananNama + "%",
	}

	// ========================================================
	// Dynamic filter
	// ========================================================

	if filter.ProvinsiID != "" {
		where = append(
			where,
			"db_fasyankes.data.provinsi_id = ?",
		)

		args = append(
			args,
			filter.ProvinsiID,
		)
	}

	if filter.KabKotaID != "" {
		where = append(
			where,
			"db_fasyankes.data.kab_kota_id = ?",
		)

		args = append(
			args,
			filter.KabKotaID,
		)
	}

	if filter.Nama != "" {
		where = append(
			where,
			"db_fasyankes.data.RUMAH_SAKIT LIKE ?",
		)

		args = append(
			args,
			"%"+filter.Nama+"%",
		)
	}

	if filter.Aktive != "" {
		where = append(
			where,
			"db_fasyankes.data.aktive = ?",
		)

		args = append(
			args,
			filter.Aktive,
		)
	}

	if filter.StartModifiedAt != "" {
		where = append(
			where,
			"db_fasyankes.data.TANGGAL_UPDATE >= ?",
		)

		args = append(
			args,
			filter.StartModifiedAt,
		)
	}

	if filter.EndModifiedAt != "" {
		where = append(
			where,
			"db_fasyankes.data.TANGGAL_UPDATE <= ?",
		)

		args = append(
			args,
			filter.EndModifiedAt,
		)
	}

	whereSQL := " WHERE " + strings.Join(where, " AND ")

	// ========================================================
	// ORDER
	// ========================================================

	sqlOrder := `
		ORDER BY db_fasyankes.data.RUMAH_SAKIT
	`

	// ========================================================
	// DATA QUERY
	// ========================================================

	dataQuery := sqlSelect +
		sqlFrom +
		whereSQL +
		sqlOrder +
		" LIMIT ? OFFSET ?"

	dataArgs := append(
		append([]interface{}{}, args...),
		filter.Limit,
		offset,
	)

	rows, err := r.DB.Query(
		dataQuery,
		dataArgs...,
	)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to query hospitals: %w",
			err,
		)
	}

	defer rows.Close()

	// ========================================================
	// Scan
	// ========================================================

	var hospitals []hospital.Hospital

	for rows.Next() {

		var hospital hospital.Hospital

		err := rows.Scan(
			&hospital.Kode,
			&hospital.Nama,
			&hospital.Jenis,
			&hospital.Kelas,
			&hospital.Telepon,
			&hospital.Email,
			&hospital.Website,
			&hospital.StatusBLU,
			&hospital.NoSuratIjinOperasional,
			&hospital.TanggalSuratIjinOperasional,
			&hospital.Direktur,
			&hospital.KetersediaanSIMRS,
			&hospital.LuasTanah,
			&hospital.LuasBangunan,
			&hospital.Kepemilikan,
			&hospital.URLTarif,
			&hospital.StatusValidasiTarif,
			&hospital.Alamat,
			&hospital.ProvinsiID,
			&hospital.ProvinsiNama,
			&hospital.KabKotaID,
			&hospital.KabKotaNama,
			&hospital.Longitude,
			&hospital.Latitude,
			&hospital.StatusAktivasi,
			&hospital.URLFotoDepan,
			&hospital.ModifiedAt,
		)

		if err != nil {
			return nil, 0, fmt.Errorf(
				"failed to scan hospital: %w",
				err,
			)
		}

		hospitals = append(
			hospitals,
			hospital,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	// ========================================================
	// COUNT QUERY
	// ========================================================

	countQuery :=
		"SELECT COUNT(db_fasyankes.data.Propinsi) " +
			sqlFrom +
			whereSQL

	var totalRowCount int

	err = r.DB.QueryRow(
		countQuery,
		args...,
	).Scan(&totalRowCount)

	if err != nil {
		return nil, 0, fmt.Errorf(
			"failed to count hospitals: %w",
			err,
		)
	}

	return hospitals, totalRowCount, nil
}
