package crypto

import (
	"crypto/ecdh"
	"encoding/base64"
	"fmt"
)

// ============================================================
// ParsePublicKey
// ============================================================
//
// Mengubah Public Key dari format Base64 menjadi
// *ecdh.PublicKey.
//
// Public Key berasal dari database.
//
// Alur:
//
//     MySQL
//        ↓
//     Base64 String
//        ↓
//     ParsePublicKey()
//        ↓
//     *ecdh.PublicKey
//
// ============================================================

func ParsePublicKey(
	publicKey string,
) (*ecdh.PublicKey, error) {

	// --------------------------------------------------------
	// Validasi
	// --------------------------------------------------------

	if publicKey == "" {
		return nil, fmt.Errorf(
			"public key is empty",
		)
	}

	// --------------------------------------------------------
	// Decode Base64
	// --------------------------------------------------------

	publicKeyBytes, err := base64.StdEncoding.DecodeString(
		publicKey,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode public key: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// Convert bytes menjadi ECDH Public Key
	// --------------------------------------------------------

	key, err := ecdh.P256().NewPublicKey(
		publicKeyBytes,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create public key: %w",
			err,
		)
	}

	return key, nil
}
