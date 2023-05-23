package linux

import (
	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/native/enum"
)

func (k *LinuxKernel) Mmap2(addrHint, size uint64, prot enum.MmapProt, flags enum.MmapFlag, fd co.Fd, off co.Off) uint64 {
	return k.Mmap(addrHint, size, prot, flags, fd, off*0x1000)
}
