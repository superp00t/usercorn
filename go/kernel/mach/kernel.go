package mach

import (
	co "github.com/superp00t/usercorn/go/kernel/common"
)

type MachKernel struct {
	*co.KernelBase
}

func NewKernel() *MachKernel {
	return &MachKernel{&co.KernelBase{}}
}
