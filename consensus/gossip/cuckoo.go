package gossip

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"hash"
	"math"
	"math/rand"
	"sync"
)

const (
	entryCountPerBucket uint = 4
	maxRetries          uint = 500
)

var defaultHasherBuilder = sha256.New

type fingerprint []byte

type bucket []fingerprint

type Cuckoo struct {
	buckets             []bucket
	bucketCount         uint
	entryCountPerBucket uint
	fingerprintLength   uint
	capacity            uint
	hasher              hash.Hash
	mut                 sync.Mutex
}

func NewCuckoo(capacity uint, falsePositiveRate float64, hasherBuilder func() hash.Hash) *Cuckoo {
	fingerprintLength := fingerprintLength(entryCountPerBucket, falsePositiveRate)
	bucketCount := nextPower(capacity / fingerprintLength * 8)
	buckets := make([]bucket, bucketCount)
	for i := uint(0); i < bucketCount; i++ {
		buckets[i] = make(bucket, entryCountPerBucket)
	}

	if hasherBuilder == nil {
		hasherBuilder = defaultHasherBuilder
	}

	return &Cuckoo{
		buckets:             buckets,
		bucketCount:         bucketCount,
		entryCountPerBucket: entryCountPerBucket,
		fingerprintLength:   fingerprintLength,
		capacity:            capacity,
		hasher:              hasherBuilder(),
	}
}

func (c *Cuckoo) Insert(input []byte) error {
	i1, i2, f := c.hashes(input)

	b1 := c.buckets[i1%c.bucketCount]
	if i, err := b1.nextIndex(); err == nil {
		b1[i] = f
		return nil
	}

	b2 := c.buckets[i2%c.bucketCount]
	if i, err := b2.nextIndex(); err == nil {
		b2[i] = f
		return nil
	}

	// else we need to start relocating items
	i := i1
	for retries := 0; retries < int(maxRetries); retries++ {
		index := i % c.bucketCount
		entryIndex := rand.Intn(int(c.entryCountPerBucket))
		f, c.buckets[index][entryIndex] = c.buckets[index][entryIndex], f
		i = i ^ uint(binary.BigEndian.Uint32(c.hash(f)))
		b := c.buckets[i%c.bucketCount]
		if idx, err := b.nextIndex(); err == nil {
			b[idx] = f
			return nil
		}
	}

	return errors.New("could not insert item ; filter is probably full")
}

func (c *Cuckoo) Delete(needle []byte) {
	i1, i2, f := c.hashes(needle)
	// try to remove from f1
	b1 := c.buckets[i1%c.bucketCount]
	if ind, ok := b1.contains(f); ok {
		b1[ind] = nil
		return
	}

	b2 := c.buckets[i2%c.bucketCount]
	if ind, ok := b2.contains(f); ok {
		b2[ind] = nil
		return
	}
}

func (c *Cuckoo) Lookup(needle []byte) bool {
	i1, i2, f := c.hashes(needle)
	_, b1 := c.buckets[i1%c.bucketCount].contains(f)
	_, b2 := c.buckets[i2%c.bucketCount].contains(f)
	return b1 || b2
}

func (c *Cuckoo) hashes(data []byte) (uint, uint, fingerprint) {
	h := c.hash(data)
	f := h[0:c.fingerprintLength]
	i1 := uint(binary.BigEndian.Uint32(h))
	i2 := i1 ^ uint(binary.BigEndian.Uint32(c.hash(f)))
	return i1, i2, fingerprint(f)
}

func (c *Cuckoo) hash(data []byte) []byte {
	c.mut.Lock()
	defer c.mut.Unlock()

	c.hasher.Write([]byte(data))
	hash := c.hasher.Sum(nil)
	c.hasher.Reset()

	return hash
}

func (b bucket) contains(f fingerprint) (int, bool) {
	for i, x := range b {
		if bytes.Equal(x, f) {
			return i, true
		}
	}
	return -1, false
}

func (b bucket) nextIndex() (int, error) {
	for i, f := range b {
		if f == nil {
			return i, nil
		}
	}
	return -1, errors.New("bucket is full")
}

func fingerprintLength(b uint, e float64) uint {
	f := uint(math.Ceil(math.Log(2 * float64(b) / e)))
	f /= 8
	if f < 1 {
		return 1
	}
	return f
}

func nextPower(i uint) uint {
	i--
	i |= i >> 1
	i |= i >> 2
	i |= i >> 4
	i |= i >> 8
	i |= i >> 16
	i |= i >> 32
	i++
	return i
}
