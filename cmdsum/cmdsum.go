package cmdsum

import (
	"buanu/util"
	"context"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync/atomic"

	"github.com/minio/sha256-simd"
	"github.com/zeebo/blake3"
	"golang.org/x/sync/semaphore"
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
	hadError := uint32(0)
	ncpu := int64(runtime.NumCPU())
	sem := semaphore.NewWeighted(ncpu)
	for _, filename := range filenames {
		// loop variable filename captured by func literal
		filename := filename
		sem.Acquire(context.Background(), 1)
		go func() {
			hashWriter := hashWriterFunc()
			if err := hashOne(hashWriter, filename); err != nil {
				log.Printf("Error while trying to sum %q: %s", filename, err.Error())
				atomic.StoreUint32(&hadError, ^uint32(0))
			}
			sem.Release(1)
		}()
	}
	sem.Acquire(context.Background(), ncpu)
	if hadError != 0 {
		os.Exit(1)
	}
}
