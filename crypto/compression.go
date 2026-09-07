package crypto

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
)

// Compress melakukan kompresi data menggunakan GZIP.
func Compress(data []byte) ([]byte, error) {

	var buffer bytes.Buffer

	writer := gzip.NewWriter(&buffer)

	_, err := writer.Write(data)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to compress data: %w",
			err,
		)
	}

	err = writer.Close()

	if err != nil {
		return nil, fmt.Errorf(
			"failed to close gzip writer: %w",
			err,
		)
	}

	return buffer.Bytes(), nil
}

// Decompress mengembalikan data GZIP menjadi data asli.
func Decompress(data []byte) ([]byte, error) {

	reader, err := gzip.NewReader(
		bytes.NewReader(data),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create gzip reader: %w",
			err,
		)
	}

	defer reader.Close()

	decompressed, err := io.ReadAll(reader)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decompress data: %w",
			err,
		)
	}

	return decompressed, nil
}
