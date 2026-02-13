package main

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"sync"
)

type RandBuffer struct {
	mu     sync.Mutex
	reader io.Reader
}

var randBuffer = &RandBuffer{reader: rand.Reader}

func generateUUID() (string, error) {
	uuid := [16]byte{}
	randBuffer.mu.Lock()
	defer randBuffer.mu.Unlock()

	if _, err := io.ReadFull(randBuffer.reader, uuid[:]); err != nil {
		return "", err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return formatUUID(uuid), nil
}

func formatUUID(uuid [16]byte) string {
	buf := make([]byte, 36)

	hex.Encode(buf[0:8], uuid[0:4])
	buf[8] = '-'
	hex.Encode(buf[9:13], uuid[4:6])
	buf[13] = '-'
	hex.Encode(buf[14:18], uuid[6:8])
	buf[18] = '-'
	hex.Encode(buf[19:23], uuid[8:10])
	buf[23] = '-'
	hex.Encode(buf[24:], uuid[10:])

	return string(buf)
}
