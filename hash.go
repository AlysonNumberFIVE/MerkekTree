package main

import "crypto/sha256"

// Hash128 implements SHA-256 encryption on incoming data.
func Hash128(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:16]
}

// compareHash compares the individual data within hashes to ensure
// they are equivalent. This is done for explicit comparison over
// string(hash1) == string(hash2)
func compareHash(hash1, hash2 []byte) bool {
	for index, data1 := range hash1 {
		if data1 != hash2[index] {
			return false
		}
	}

	return true
}

func GenerateHash(LeftChildHash, RightChidlHash []byte) []byte {
	joinHashes := string(LeftChildHash) + string(RightChidlHash)
	return []byte(joinHashes)
}
