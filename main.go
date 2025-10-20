package main

import (
	"pwner/elf"
)

func main() {

	e := elf.ELF("/home/tsuneki/dc/ctf/qnqsec/notez/notez")
	Dump(e.Sym("main"))
}
