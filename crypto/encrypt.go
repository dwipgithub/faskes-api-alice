package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type EncryptedMessage struct {
	KeyID              string `json:"key_id"`
	EphemeralPublicKey string `json:"ephemeral_public_key"`
	Nonce              string `json:"nonce"`
	Ciphertext         string `json:"ciphertext"`
}

// ============================================================
// Encrypt
// ============================================================
//
// recipientPublicKey
//     Public Key milik penerima.
//
// keyID
//     ID dari key milik penerima.
//
// plaintext
//     Data yang akan dienkripsi.
//
// Proses:
//
//     Ephemeral Key Pair
//             ↓
//            ECDH
//             ↓
//       Shared Secret
//             ↓
//       HKDF Extract
//             ↓
//       HKDF Expand
//             ↓
//       Encryption Key
//             ↓
//        AES-256-GCM
//
// ============================================================

func Encrypt(
	recipientPublicKey *ecdh.PublicKey,
	keyID string,
	plaintext []byte,
) (*EncryptedMessage, error) {

	// ========================================================
	// STEP 1
	// Generate Ephemeral Private Key
	// ========================================================

	curve := ecdh.P256()

	ephemeralPrivateKey, err := curve.GenerateKey(
		rand.Reader,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate ephemeral key: %w",
			err,
		)
	}

	// ========================================================
	// STEP 2
	// Ambil Ephemeral Public Key
	// ========================================================

	ephemeralPublicKey := ephemeralPrivateKey.PublicKey()

	// ========================================================
	// STEP 3
	// ECDH
	// ========================================================
	//
	// Ephemeral Private Key Alice
	//             +
	// Public Key Bob
	//             ↓
	//       Shared Secret
	//
	// ========================================================

	sharedSecret, err := ephemeralPrivateKey.ECDH(
		recipientPublicKey,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"ECDH failed: %w",
			err,
		)
	}

	// ========================================================
	// STEP 4
	// HKDF Extract
	// ========================================================

	prk, err := hkdf.Extract(
		sha256.New,
		sharedSecret,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"HKDF extract failed: %w",
			err,
		)
	}

	// ========================================================
	// STEP 5
	// HKDF Expand
	// ========================================================
	//
	// Info harus sama dengan yang digunakan Bob
	// ketika melakukan decrypt.
	//
	// ========================================================

	const HKDFInfo = "hospital-data-encryption-v1"

	encryptionKey, err := hkdf.Expand(
		sha256.New,
		prk,
		HKDFInfo,
		32,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"HKDF expand failed: %w",
			err,
		)
	}

	// ========================================================
	// STEP 6
	// AES-256
	// ========================================================

	block, err := aes.NewCipher(
		encryptionKey,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"AES initialization failed: %w",
			err,
		)
	}

	// ========================================================
	// STEP 7
	// AES-GCM
	// ========================================================

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return nil, fmt.Errorf(
			"GCM initialization failed: %w",
			err,
		)
	}

	// ========================================================
	// STEP 8
	// Generate Nonce
	// ========================================================

	nonce := make(
		[]byte,
		gcm.NonceSize(),
	)

	_, err = rand.Read(nonce)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate nonce: %w",
			err,
		)
	}

	// ========================================================
	// STEP 9
	// AES-GCM Encryption
	// ========================================================
	//
	// Seal menghasilkan:
	//
	//     ciphertext + authentication tag
	//
	// keyID digunakan sebagai Additional Authenticated Data
	// sehingga key_id juga dilindungi integritasnya.
	//
	// ========================================================

	ciphertext := gcm.Seal(
		nil,
		nonce,
		plaintext,
		[]byte(keyID),
	)

	// ========================================================
	// STEP 10
	// Encode Binary → Base64
	// ========================================================
	//
	// Semua data binary yang dikirim melalui JSON
	// menggunakan Base64.
	//
	// ========================================================

	return &EncryptedMessage{
		KeyID: keyID,

		EphemeralPublicKey: base64.StdEncoding.EncodeToString(
			ephemeralPublicKey.Bytes(),
		),

		Nonce: base64.StdEncoding.EncodeToString(
			nonce,
		),

		Ciphertext: base64.StdEncoding.EncodeToString(
			ciphertext,
		),
	}, nil
}
