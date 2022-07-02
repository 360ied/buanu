package sums

import (
	"buanu/util"
	"hash"
	"sort"
	"strings"

	"github.com/minio/sha256-simd"
	"github.com/zeebo/blake3"
)

var Hashes = map[string]func() hash.Hash{
	"sha256": sha256.New,
	"blake3": func() hash.Hash { return blake3.New() },
}

func GetHashNames() string {
	keysSlice := util.MapKeysSlice(Hashes)
	sort.Strings(keysSlice)
	return strings.Join(keysSlice, ", ")
}
