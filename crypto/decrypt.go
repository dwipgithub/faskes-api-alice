package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/hkdf"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ============================================================
// Decrypt
// ============================================================
//
// Bob menggunakan:
//
//     Bob Private Key
//           +
//     Alice Ephemeral Public Key
//           ↓
//          ECDH
//           ↓
//     Shared Secret
//           ↓
//          HKDF
//           ↓
//     Encryption Key
//           ↓
//       AES-256-GCM
//           ↓
//        Plaintext
//
// ============================================================

func Decrypt(
	bobPrivateKey *ecdh.PrivateKey,
	encryptedMessage *EncryptedMessage,
) ([]byte, error) {

	// ========================================================
	// STEP 1
	// Decode Ephemeral Public Key
	// ========================================================

	ephemeralPublicKeyBytes, err := hex.DecodeString(
		encryptedMessage.EphemeralPublicKey,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid ephemeral public key: %w",
			err,
		)
	}

	// ========================================================
	// STEP 2
	// Convert menjadi ECDH Public Key
	// ========================================================

	ephemeralPublicKey, err := ecdh.P256().NewPublicKey(
		ephemeralPublicKeyBytes,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid ephemeral public key: %w",
			err,
		)
	}

	// ========================================================
	// STEP 3
	// ECDH
	// ========================================================

	sharedSecret, err := bobPrivateKey.ECDH(
		ephemeralPublicKey,
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

	encryptionKey, err := hkdf.Expand(
		sha256.New,
		prk,
		"ecc-learning encryption",
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
	// AES
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
	// Decode Nonce
	// ========================================================

	nonce, err := hex.DecodeString(
		encryptedMessage.Nonce,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid nonce: %w",
			err,
		)
	}

	// Pastikan nonce memiliki ukuran yang benar.

	if len(nonce) != gcm.NonceSize() {

		return nil, fmt.Errorf(
			"invalid nonce size",
		)
	}

	// ========================================================
	// STEP 9
	// Decode Ciphertext
	// ========================================================

	ciphertext, err := hex.DecodeString(
		encryptedMessage.Ciphertext,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"invalid ciphertext: %w",
			err,
		)
	}

	// ========================================================
	// STEP 10
	// Decrypt + Authentication
	// ========================================================
	//
	// Open() melakukan:
	//
	//     Authentication Tag verification
	//                 +
	//             Decryption
	//
	// Jika ciphertext diubah:
	//
	//     Open() → error
	//
	// ========================================================

	plaintext, err := gcm.Open(
		nil,
		nonce,
		ciphertext,
		[]byte(encryptedMessage.KeyID),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"decryption failed: %w",
			err,
		)
	}

	return plaintext, nil
}
