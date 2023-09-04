package ton

import (
	"crypto/sha256"
	"encoding/json"
	"testing"
)

func TestBits256JSON(t *testing.T) {

	hash := Bits256(sha256.Sum256([]byte("test")))
	j, err := json.Marshal(hash)
	if err != nil {
		t.Fatal(err)
	}
	if string(j) != `"9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08"` {
		t.Fatal("invalid sum", string(j))
	}
	var hash2 Bits256
	err = json.Unmarshal(j, &hash2)
	if err != nil {
		t.Fatal(err)
	}
	if hash2 != hash {
		t.Fatal("mismatch hashes", hash, hash2)
	}
}
