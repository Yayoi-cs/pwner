package utils

import (
	"encoding/hex"
	"log"
	"strconv"
)

func Encode(s string) []byte {
	return []byte(s)
}

func Decode(data []byte) string {
	return string(data)
}

func Hex(data []byte) string {
	return hex.EncodeToString(data)
}

func Strtou64(s string, base int) uint64 {
	val, err := strconv.ParseUint(s, base, 64)
	if err != nil {
		log.Fatalf("ParseUint failed: %v", err)
	}
	return val
}

func Bintou64(s string) uint64 {
	return Strtou64(s, 2)
}

func Octtou64(s string) uint64 {
	return Strtou64(s, 8)
}

func Dectou64(s string) uint64 {
	return Strtou64(s, 10)
}

func Hextou64(s string) uint64 {
	return Strtou64(s, 16)
}

func Strtou32(s string, base int) uint32 {
	val, err := strconv.ParseUint(s, base, 32)
	if err != nil {
		log.Fatalf("ParseUint failed: %v", err)
	}
	return uint32(val)
}

func Bintou32(s string) uint32 {
	return Strtou32(s, 2)
}

func Octtou32(s string) uint32 {
	return Strtou32(s, 8)
}

func Dectou32(s string) uint32 {
	return Strtou32(s, 10)
}

func Hextou32(s string) uint32 {
	return Strtou32(s, 16)
}

func Strtou16(s string, base int) uint16 {
	val, err := strconv.ParseUint(s, base, 16)
	if err != nil {
		log.Fatalf("ParseUint failed: %v", err)
	}
	return uint16(val)
}

func Bintou16(s string) uint16 {
	return Strtou16(s, 2)
}

func Octtou16(s string) uint16 {
	return Strtou16(s, 8)
}

func Dectou16(s string) uint16 {
	return Strtou16(s, 10)
}

func Hextou16(s string) uint16 {
	return Strtou16(s, 16)
}

func Strtou8(s string, base int) uint8 {
	val, err := strconv.ParseUint(s, base, 8)
	if err != nil {
		log.Fatalf("ParseUint failed: %v", err)
	}
	return uint8(val)
}

func Bintou8(s string) uint8 {
	return Strtou8(s, 2)
}

func Octtou8(s string) uint8 {
	return Strtou8(s, 8)
}

func Dectou8(s string) uint8 {
	return Strtou8(s, 10)
}

func Hextou8(s string) uint8 {
	return Strtou8(s, 16)
}
