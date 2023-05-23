package darwin

import (
	"github.com/lunixbochs/argjoy"
	"syscall"

	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/kernel/darwin/unpack"
	"github.com/superp00t/usercorn/go/native/enum"
)

func Unpack(k co.Kernel, arg interface{}, vals []interface{}) error {
	// TODO: this is the exact same preamble as linux
	reg0 := vals[0].(uint64)
	// null pointer guard
	if reg0 == 0 {
		return nil
	}
	buf := co.NewBuf(k, reg0)
	switch v := arg.(type) {
	case *syscall.Sockaddr:
		*v = unpack.Sockaddr(buf, int(vals[1].(uint64)))
	case *enum.OpenFlag:
		*v = unpack.OpenFlag(reg0)
	case *enum.MmapFlag:
		*v = unpack.MmapFlag(reg0)
	case *enum.MmapProt:
		*v = unpack.MmapProt(reg0)
	default:
		return argjoy.NoMatch
	}
	return nil
}

func registerUnpack(d *DarwinKernel) {
	d.Argjoy.Register(func(arg interface{}, vals []interface{}) error {
		return Unpack(d, arg, vals)
	})
}
