package arm64

import (
	ks "github.com/keystone-engine/keystone/bindings/go/keystone"
	cs "github.com/lunixbochs/capstr"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"

	"github.com/superp00t/usercorn/go/cpu"
	"github.com/superp00t/usercorn/go/cpu/unicorn"
	"github.com/superp00t/usercorn/go/models"
)

var Arch = &models.Arch{
	Name:   "arm64",
	Bits:   64,
	Radare: "arm64",

	Cpu: &unicorn.Builder{Arch: uc.ARCH_ARM64, Mode: uc.MODE_ARM},
	Dis: &cpu.Capstr{Arch: cs.ARCH_ARM64, Mode: cs.MODE_ARM},
	Asm: &cpu.Keystone{Arch: ks.ARCH_ARM64, Mode: ks.MODE_LITTLE_ENDIAN},

	PC: uc.ARM64_REG_PC,
	SP: uc.ARM64_REG_SP,
	Regs: map[string]int{
		"x0":  uc.ARM64_REG_X0,
		"x1":  uc.ARM64_REG_X1,
		"x2":  uc.ARM64_REG_X2,
		"x3":  uc.ARM64_REG_X3,
		"x4":  uc.ARM64_REG_X4,
		"x5":  uc.ARM64_REG_X5,
		"x6":  uc.ARM64_REG_X6,
		"x7":  uc.ARM64_REG_X7,
		"x8":  uc.ARM64_REG_X8,
		"x9":  uc.ARM64_REG_X9,
		"x10": uc.ARM64_REG_X10,
		"x11": uc.ARM64_REG_X11,
		"x12": uc.ARM64_REG_X12,
		"x13": uc.ARM64_REG_X13,
		"x14": uc.ARM64_REG_X14,
		"x15": uc.ARM64_REG_X15,
		"x16": uc.ARM64_REG_X16,
		"x17": uc.ARM64_REG_X17,
		"x18": uc.ARM64_REG_X18,
		"x19": uc.ARM64_REG_X19,
		"x20": uc.ARM64_REG_X20,
		"x21": uc.ARM64_REG_X21,
		"x22": uc.ARM64_REG_X22,
		"x23": uc.ARM64_REG_X23,
		"x24": uc.ARM64_REG_X24,
		"x25": uc.ARM64_REG_X25,
		"x26": uc.ARM64_REG_X26,
		"x27": uc.ARM64_REG_X27,
		"x28": uc.ARM64_REG_X28,
		"fp":  uc.ARM64_REG_FP,
		"lr":  uc.ARM64_REG_LR,
		"sp":  uc.ARM64_REG_SP,
		"pc":  uc.ARM64_REG_PC,
	},
	DefaultRegs: []string{
		"x0", "x1", "x2", "x3", "x4", "x5", "x6", "x7", "x8",
		"x9", "x10", "x11", "x12", "x13", "x14", "x15", "x16",
		"x17", "x18", "x19", "x20", "x21", "x22", "x23", "x24",
		"x25", "x26", "x27", "x28",
	},
}

func EnableFPU(u models.Usercorn) error {
	val, err := u.RegRead(uc.ARM64_REG_CPACR_EL1)
	if err != nil {
		return err
	}
	val |= 0x300000 // set the FPEN bits
	return u.RegWrite(uc.ARM64_REG_CPACR_EL1, val)
}
