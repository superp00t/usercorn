package linux

import (
	"github.com/lunixbochs/argjoy"
	"syscall"

	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/native"
)

func Pack(buf co.Buf, i interface{}) error {
	switch v := i.(type) {
	case *syscall.Statfs_t:
		return buf.Pack(native.StatfsToLinux(v))
	default:
		return argjoy.NoMatch
	}
}
