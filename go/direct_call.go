package usercorn

import (
	"bytes"
	"fmt"

	"github.com/superp00t/usercorn/go/models"
)

func (u *Usercorn) asmCallSetupCdecl(void bool, addr uint64, arguments []uint64) (text string, err error) {
	call_S := bytes.NewBuffer(nil)
	for i := len(arguments) - 1; i >= 0; i-- {
		arg := arguments[i]
		fmt.Fprintf(call_S, "push 0x%016X\n", arg)
	}

	fmt.Fprintf(call_S, "mov eax, 0x%016X\n", addr)
	fmt.Fprintf(call_S, "call eax\n")
	text = call_S.String()
	return
}

// Begin emulation at a certain addr, packing in arguments according to a certain calling convention
func (u *Usercorn) StartDirectCall(conv models.Callconvention, void bool, addr uint64, arguments []uint64) (retval uint64, err error) {
	var callsetup string
	switch conv {
	case models.Cdecl:
		callsetup, err = u.asmCallSetupCdecl(void, addr, arguments)
	default:
		err = fmt.Errorf("unsupported call conv %d", conv)
	}

	if err != nil {
		return
	}

	fmt.Println("Running ASM!", callsetup)

	err = u.RunAsm(0, callsetup, nil, nil)
	return
}
