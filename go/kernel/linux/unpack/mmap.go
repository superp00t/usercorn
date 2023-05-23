package unpack

import (
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"

	"github.com/superp00t/usercorn/go/native/enum"
)

var mmapProtMap = map[int]int{
	0: uc.PROT_NONE,
	1: uc.PROT_READ,
	2: uc.PROT_WRITE,
	4: uc.PROT_EXEC,
	// FIXME?
	// 0x01000000: syscall.PROT_GROWSDOWN,
	// 0x02000000: syscall.PROT_GROWSUP,
}

func MmapProt(reg uint64) enum.MmapProt {
	var out enum.MmapProt
	for a, b := range mmapProtMap {
		if int(reg)&a == a {
			out |= enum.MmapProt(b)
		}
	}
	return out
}

func MmapFlag(reg uint64) enum.MmapFlag {
	var out enum.MmapFlag
	for a, b := range mmapFlagMap {
		if int(reg)&a == a {
			out |= enum.MmapFlag(b)
		}
	}
	return out
}
