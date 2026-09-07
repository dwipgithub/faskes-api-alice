package consumer

import (
	"database/sql"
	"fmt"
)

// ConsumerRepository mengakses data consumer.
type ConsumerRepository struct {
	DB *sql.DB
}

// GetByCode mencari consumer berdasarkan consumer_code.
func (r *ConsumerRepository) GetByCode(
	consumerCode string,
) (int, error) {

	var consumerID int

	err := r.DB.QueryRow(`
		SELECT id
		FROM db_api_auth.user
		WHERE user_name = ?
			AND is_active = 1
	`,
		consumerCode,
	).Scan(&consumerID)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf(
				"consume not found: %s",
				consumerCode,
			)
		}

		return 0, err
	}

	return consumerID, nil
}

// GetEncryptionStatus mengambil status encryption
// berdasarkan consumer ID.
func (r *ConsumerRepository) GetEncryptionStatus(
	consumerID int,
) (bool, error) {

	var encryptionEnabled bool

	err := r.DB.QueryRow(`
		SELECT encryption_enabled
		FROM db_api_auth.user
		WHERE id = ?
			AND is_active = 1
	`,
		consumerID,
	).Scan(
		&encryptionEnabled,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf(
				"consume not found: %d",
				consumerID,
			)
		}

		return false, err
	}

	return encryptionEnabled, nil
}

// GetByCodeWithPassword mencari consumer berdasarkan consumer_code
// dan mengembalikan data yang diperlukan untuk proses login.
func (r *ConsumerRepository) GetByUserName(
	userName string,
) (int, string, string, error) {

	var (
		consumerID   int
		passwordHash string
		status       string
	)

	err := r.DB.QueryRow(`
		SELECT
			id,
			password,
			is_active
		FROM db_api_auth.user
		WHERE user_name = ?
	`,
		userName,
	).Scan(
		&consumerID,
		&passwordHash,
		&status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", "", fmt.Errorf(
				"username not found: %s",
				userName,
			)
		}

		return 0, "", "", err
	}

	return consumerID, passwordHash, status, nil
}
