package crypto

import (
	"crypto/ecdh"
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type GeneratedKey struct {
	KeyID      string
	PrivateKey string
	PublicKey  string
}

func GenerateKeyPair(
	keyID string,
) (*GeneratedKey, error) {

	// --------------------------------------------------------
	// Validasi Key ID
	// --------------------------------------------------------

	if keyID == "" {
		return nil, fmt.Errorf("key_id is required")
	}

	// --------------------------------------------------------
	// Generate P-256 Private Key
	// --------------------------------------------------------

	curve := ecdh.P256()

	privateKey, err := curve.GenerateKey(
		rand.Reader,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to generate private key: %w",
			err,
		)
	}

	// --------------------------------------------------------
	// Ambil Public Key
	// --------------------------------------------------------

	publicKey := privateKey.PublicKey()

	// --------------------------------------------------------
	// Encode Base64
	// --------------------------------------------------------

	privateKeyBase64 := base64.StdEncoding.EncodeToString(
		privateKey.Bytes(),
	)

	publicKeyBase64 := base64.StdEncoding.EncodeToString(
		publicKey.Bytes(),
	)

	// --------------------------------------------------------
	// Return Key Pair
	// --------------------------------------------------------

	return &GeneratedKey{
		KeyID:      keyID,
		PrivateKey: privateKeyBase64,
		PublicKey:  publicKeyBase64,
	}, nil
}
