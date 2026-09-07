package consumer

import (
	"database/sql"
	"fmt"
	"time"
)

type ConsumerKeyRepository struct {
	DB *sql.DB
}

// GenerateKeyID membuat key_id baru.
func (r *ConsumerKeyRepository) GenerateKeyID() (string, error) {

	var lastID int

	err := r.DB.QueryRow(`
		SELECT COALESCE(MAX(id), 0)
		FROM db_api_auth.consumer_keys
	`).Scan(&lastID)

	if err != nil {
		return "", err
	}

	nextID := lastID + 1

	keyID := fmt.Sprintf(
		"KEY-%d-%03d",
		time.Now().Year(),
		nextID,
	)

	return keyID, nil
}

// DeactivateActiveKey menonaktifkan key aktif consumer.
func (r *ConsumerKeyRepository) DeactivateActiveKey(
	consumerID int,
) error {

	_, err := r.DB.Exec(`
		UPDATE db_api_auth.consumer_keys
		SET status = 'inactive'
		WHERE consumer_id = ?
			AND status = 'active'
	`,
		consumerID,
	)

	return err
}

// Create menyimpan public key consumer.
func (r *ConsumerKeyRepository) Create(
	consumerID int,
	keyID string,
	publicKey string,
) error {

	_, err := r.DB.Exec(`
		INSERT INTO db_api_auth.consumer_keys (
			consumer_id,
			key_id,
			public_key,
			status,
			created_at
		)
		VALUES (?, ?, ?, 'active', ?)
	`,
		consumerID,
		keyID,
		publicKey,
		time.Now(),
	)

	return err
}

// GetActiveKey mengambil public key aktif consumer.
func (r *ConsumerKeyRepository) GetActiveKey(
	consumerID int,
) (string, string, error) {

	var keyID string
	var publicKey string

	err := r.DB.QueryRow(`
		SELECT key_id, public_key
		FROM db_api_auth.consumer_keys
		WHERE consumer_id = ?
		AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`,
		consumerID,
	).Scan(
		&keyID,
		&publicKey,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			return "", "", fmt.Errorf(
				"active public key not found",
			)
		}

		return "", "", err
	}

	return keyID, publicKey, nil
}
