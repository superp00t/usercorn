package darwin

import (
	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/native/enum"
)

func (k *DarwinKernel) Pread(fd co.Fd, buf co.Obuf, size co.Len, offset int64) uint64 {
	return k.PosixKernel.Pread64(fd, buf, size, offset)
}

func (k *DarwinKernel) PreadNocancel(fd co.Fd, buf co.Obuf, size co.Len, offset int64) uint64 {
	// TODO: will it be possible to cancel syscalls here?
	// what conditions can cancel a syscall in the real world?
	return k.Pread(fd, buf, size, offset)
}

func (k *DarwinKernel) Pwrite(fd co.Fd, buf co.Buf, size co.Len, offset int64) uint64 {
	return k.PosixKernel.Pwrite64(fd, buf, size, offset)
}

func (k *DarwinKernel) Fstat64(fd co.Fd, buf co.Obuf) uint64 {
	return k.PosixKernel.Fstat(fd, buf)
}

func (k *DarwinKernel) Stat64(path string, buf co.Obuf) uint64 {
	return k.PosixKernel.Stat(path, buf)
}

func (k *DarwinKernel) OpenNocancel(path string, flags enum.OpenFlag, mode uint64) uint64 {
	return k.PosixKernel.Open(path, flags, mode)
}

func (k *DarwinKernel) ReadNocancel(fd co.Fd, buf co.Obuf, size co.Len) uint64 {
	return k.PosixKernel.Read(fd, buf, size)
}

func (k *DarwinKernel) CloseNocancel(fd co.Fd) uint64 {
	return k.PosixKernel.Close(fd)
}
