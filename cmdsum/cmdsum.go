package cmdsum

import (
	"buanu/util"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/minio/sha256-simd"
	"github.com/zeebo/blake3"
)

func hashOne(hashWriter hash.Hash, filename string) error {
	var file io.Reader
	if filename == "-" {
		file = os.Stdin
	} else {
		var err error
		file, err = os.Open(filename)
		if err != nil {
			return err
		}
	}
	if _, err := io.Copy(hashWriter, file); err != nil {
		return err
	}
	sum := hashWriter.Sum(nil)
	_, err := fmt.Printf("%x  %s\n", sum, filename)
	return err
}

func Run() {
	hashes := map[string]func() hash.Hash{
		"sha256": sha256.New,
		"blake3": func() hash.Hash { return blake3.New() },
	}
	hashNamesFunc := func() string {
		keysSlice := util.MapKeysSlice(hashes)
		sort.Strings(keysSlice)
		return strings.Join(keysSlice, ", ")
	}
	if len(os.Args) < 3 {
		log.Fatalf("Hash function not specified. Valid hashes are: %s.", hashNamesFunc())
	}
	hashName := os.Args[2]
	hashWriterFunc, found := hashes[hashName]
	if !found {
		log.Fatalf("%q is not a valid hash function. Valid hashes are: %s.", hashName, hashNamesFunc())
	}
	var filenames []string
	if len(os.Args) < 4 {
		filenames = append(filenames, "-")
	} else {
		filenames = append(filenames, os.Args[3:]...)
	}
	hashWriter := hashWriterFunc()
	hadError := false
	for _, filename := range filenames {
		if err := hashOne(hashWriter, filename); err != nil {
			log.Printf("Error while trying to sum %q: %s", filename, err.Error())
			hadError = true
		}
	}
	if hadError {
		os.Exit(1)
	}
}
