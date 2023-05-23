package posix

import (
	"os"
	"syscall"

	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/models"
)

func (k *PosixKernel) Exit(code int) {
	k.U.Exit(models.ExitStatus(code))
}

func (k *PosixKernel) ExitGroup(code int) {
	k.Exit(code)
}

func (k *PosixKernel) Getpid() int {
	return os.Getpid()
}

func (k *PosixKernel) Getppid() int {
	return os.Getppid()
}

func (k *PosixKernel) Getpgid(pid int) uint64 {
	n, err := syscall.Getpgid(pid)
	if err != nil {
		return Errno(err)
	}
	return uint64(n)
}

func (k *PosixKernel) Getpgrp() int {
	return syscall.Getpgrp()
}

func (k *PosixKernel) Kill(pid, signal int) uint64 {
	// TODO: os-specific signal handling?
	return Errno(syscall.Kill(pid, syscall.Signal(signal)))
}

func (k *PosixKernel) Execve(path string, argvBuf, envpBuf co.Buf) uint64 {
	// TODO: put this function somewhere generic?
	readStrArray := func(buf co.Buf) []string {
		var out []string
		st := buf.Struc()
		for {
			var addr uint64
			if k.U.Bits() == 64 {
				st.Unpack(&addr)
			} else {
				var addr32 uint32
				st.Unpack(&addr32)
				addr = uint64(addr32)
			}
			if addr == 0 {
				break
			}
			s, _ := k.U.Mem().ReadStrAt(addr)
			out = append(out, s)
		}
		return out
	}
	argv := readStrArray(argvBuf)
	envp := readStrArray(envpBuf)
	return Errno(syscall.Exec(path, argv, envp))
}
