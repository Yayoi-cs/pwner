package elf

import (
	"debug/elf"
	"log"
	"sync"
)

type ELFer struct {
	base   uint64
	path   string
	sym    map[string]uint64
	rawSym map[string]uint64
}

func ELF(path string) *ELFer {
	file, err := elf.Open(path)
	if err != nil {
		log.Fatalf("failed to open ELF: %v", err)
	}
	defer file.Close()

	e := &ELFer{
		base:   0,
		path:   path,
		sym:    make(map[string]uint64),
		rawSym: make(map[string]uint64),
	}

	var symSymbols, dynSymbols []elf.Symbol
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if symbols, err := file.Symbols(); err == nil {
			symSymbols = symbols
		}
	}()

	go func() {
		defer wg.Done()
		if symbols, err := file.DynamicSymbols(); err == nil {
			dynSymbols = symbols
		}
	}()

	wg.Wait()

	for _, s := range symSymbols {
		if s.Name != "" && s.Value != 0 {
			e.rawSym[s.Name] = s.Value
			e.sym[s.Name] = s.Value
		}
	}

	for _, s := range dynSymbols {
		if s.Name != "" && s.Value != 0 {
			if _, exists := e.rawSym[s.Name]; !exists {
				e.rawSym[s.Name] = s.Value
				e.sym[s.Name] = s.Value
			}
		}
	}

	return e
}

func (e *ELFer) Base(base uint64) {
	e.base = base
	for name, rawAddr := range e.rawSym {
		e.sym[name] = e.base + rawAddr
	}
}

func (e *ELFer) Sym(name string) uint64 {
	if addr, exists := e.sym[name]; exists {
		return addr
	}
	log.Fatalf("symbol '%s' not found", name)
	return 0
}

func (e *ELFer) Plt(name string) uint64 {
	pltName := name + "@plt"
	if addr, exists := e.sym[pltName]; exists {
		return addr
	}
	log.Fatalf("PLT entry '%s' not found", name)
	return 0
}

func (e *ELFer) Got(name string) uint64 {
	file, err := elf.Open(e.path)
	if err != nil {
		log.Fatalf("failed to open ELF: %v", err)
	}
	defer file.Close()

	gotSection := file.Section(".got.plt")
	if gotSection == nil {
		gotSection = file.Section(".got")
	}
	if gotSection == nil {
		log.Fatalf("GOT section not found")
	}

	relPltSection := file.Section(".rela.plt")
	if relPltSection == nil {
		relPltSection = file.Section(".rel.plt")
	}
	if relPltSection == nil {
		log.Fatalf("relocation section not found")
	}

	dynSymbols, err := file.DynamicSymbols()
	if err != nil {
		log.Fatalf("failed to get dynamic symbols: %v", err)
	}

	data, _ := relPltSection.Data()
	var entrySize int
	if file.Class == elf.ELFCLASS64 {
		if relPltSection.Type == elf.SHT_RELA {
			entrySize = 24
		} else {
			entrySize = 16
		}
	} else {
		if relPltSection.Type == elf.SHT_RELA {
			entrySize = 12
		} else {
			entrySize = 8
		}
	}

	entryCount := len(data) / entrySize

	for i := 0; i < entryCount; i++ {
		entryData := data[i*entrySize : (i+1)*entrySize]

		var offset uint64
		var symIndex uint32

		if file.Class == elf.ELFCLASS64 {
			offset = uint64(entryData[0]) | uint64(entryData[1])<<8 |
				uint64(entryData[2])<<16 | uint64(entryData[3])<<24 |
				uint64(entryData[4])<<32 | uint64(entryData[5])<<40 |
				uint64(entryData[6])<<48 | uint64(entryData[7])<<56

			info := uint64(entryData[8]) | uint64(entryData[9])<<8 |
				uint64(entryData[10])<<16 | uint64(entryData[11])<<24 |
				uint64(entryData[12])<<32 | uint64(entryData[13])<<40 |
				uint64(entryData[14])<<48 | uint64(entryData[15])<<56

			symIndex = uint32(info >> 32)
		} else {
			offset = uint64(entryData[0]) | uint64(entryData[1])<<8 |
				uint64(entryData[2])<<16 | uint64(entryData[3])<<24

			info := uint32(entryData[4]) | uint32(entryData[5])<<8 |
				uint32(entryData[6])<<16 | uint32(entryData[7])<<24

			symIndex = info >> 8
		}

		if symIndex < uint32(len(dynSymbols)) {
			if dynSymbols[symIndex].Name == name {
				return e.base + offset
			}
		}
	}

	log.Fatalf("GOT entry '%s' not found", name)
	return 0
}
