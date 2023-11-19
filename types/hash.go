package types

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

type HashTopic string

func HashType[T any]() (HashTopic, error) {
	ptr := new(T)
	var hash string
	var err error
	if hash, err = AsSha256(*ptr); err != nil {
		return "", err
	}
	return HashTopic(hash), nil
}

func AsSha256(s interface{}) (string, error) {
	// Step 1: Serialize the struct to a byte slice.
	bytes, err := json.Marshal(s)
	if err != nil {
		return "", err
	}

	// Step 2: Create a SHA-256 hash from the byte slice.
	hash := sha256.Sum256(bytes)

	// Step 3: Convert the hash to a string.
	hashString := hex.EncodeToString(hash[:])

	return hashString, nil
}
