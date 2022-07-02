package cmdsum

import (
	"fmt"
	"hash"
	"io"
	"log"
	"os"

	"github.com/minio/sha256-simd"
	"github.com/zeebo/blake3"
)

func Run() {
	hashes := map[string]func() hash.Hash{
		"sha256": sha256.New,
		"blake3": func() hash.Hash { return blake3.New() },
	}
	if len(os.Args) < 3 {
		log.Fatal("Hash function not specified.")
	}
	hashName := os.Args[2]
	hashWriterFunc, found := hashes[hashName]
	if !found {
		log.Fatalf("%q is not a valid hash function.", hashName)
	}
	hashWriter := hashWriterFunc()
	if _, err := io.Copy(hashWriter, os.Stdin); err != nil {
		panic(err)
	}
	sum := hashWriter.Sum(nil)
	_, err := fmt.Printf("%x", sum)
	if err != nil {
		panic(err)
	}
}
