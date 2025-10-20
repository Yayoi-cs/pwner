package utils

import (
	"encoding/binary"
	"log"
)

type Endian int

const (
	LittleEndian Endian = iota
	BigEndian
)

func P8(v uint8) []byte {
	return []byte{v}
}

func P16(v uint16, endian ...Endian) []byte {
	buf := make([]byte, 2)
	e := getEndian(endian...)
	e.PutUint16(buf, v)
	return buf
}

func P32(v uint32, endian ...Endian) []byte {
	buf := make([]byte, 4)
	e := getEndian(endian...)
	e.PutUint32(buf, v)
	return buf
}

func P64(v uint64, endian ...Endian) []byte {
	buf := make([]byte, 8)
	e := getEndian(endian...)
	e.PutUint64(buf, v)
	return buf
}

func U8(data []byte) uint8 {
	if len(data) < 1 {
		log.Fatalf("U8: insufficient data length, need 1 byte, got %d", len(data))
	}
	return data[0]
}

func U16(data []byte, endian ...Endian) uint16 {
	if len(data) < 2 {
		log.Fatalf("U16: insufficient data length, need 2 bytes, got %d", len(data))
	}
	e := getEndian(endian...)
	return e.Uint16(data[:2])
}

func U32(data []byte, endian ...Endian) uint32 {
	if len(data) < 4 {
		log.Fatalf("U32: insufficient data length, need 4 bytes, got %d", len(data))
	}
	e := getEndian(endian...)
	return e.Uint32(data[:4])
}

func U64(data []byte, endian ...Endian) uint64 {
	if len(data) < 8 {
		log.Fatalf("U64: insufficient data length, need 8 bytes, got %d", len(data))
	}
	e := getEndian(endian...)
	return e.Uint64(data[:8])
}

func getEndian(endian ...Endian) binary.ByteOrder {
	if len(endian) > 0 && endian[0] == BigEndian {
		return binary.BigEndian
	}
	return binary.LittleEndian
}

func P16LE(v uint16) []byte { return P16(v, LittleEndian) }
func P32LE(v uint32) []byte { return P32(v, LittleEndian) }
func P64LE(v uint64) []byte { return P64(v, LittleEndian) }

func P16BE(v uint16) []byte { return P16(v, BigEndian) }
func P32BE(v uint32) []byte { return P32(v, BigEndian) }
func P64BE(v uint64) []byte { return P64(v, BigEndian) }

func U16LE(data []byte) uint16 { return U16(data, LittleEndian) }
func U32LE(data []byte) uint32 { return U32(data, LittleEndian) }
func U64LE(data []byte) uint64 { return U64(data, LittleEndian) }

func U16BE(data []byte) uint16 { return U16(data, BigEndian) }
func U32BE(data []byte) uint32 { return U32(data, BigEndian) }
func U64BE(data []byte) uint64 { return U64(data, BigEndian) }
