package mips

import (
	ks "github.com/keystone-engine/keystone/bindings/go/keystone"
	cs "github.com/lunixbochs/capstr"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"

	"github.com/superp00t/usercorn/go/cpu"
	"github.com/superp00t/usercorn/go/cpu/unicorn"
	"github.com/superp00t/usercorn/go/models"
)

var Arch = &models.Arch{
	Name:   "mips",
	Bits:   32,
	Radare: "mips",

	Cpu: &unicorn.Builder{Arch: uc.ARCH_MIPS, Mode: uc.MODE_MIPS32 + uc.MODE_LITTLE_ENDIAN},
	Dis: &cpu.Capstr{Arch: cs.ARCH_MIPS, Mode: cs.MODE_MIPS32 + cs.MODE_LITTLE_ENDIAN},
	Asm: &cpu.Keystone{Arch: ks.ARCH_MIPS, Mode: ks.MODE_MIPS32 + ks.MODE_LITTLE_ENDIAN},

	PC: uc.MIPS_REG_PC,
	SP: uc.MIPS_REG_SP,
	Regs: map[string]int{
		"at": uc.MIPS_REG_AT,
		"v0": uc.MIPS_REG_V0,
		"v1": uc.MIPS_REG_V1,
		"a0": uc.MIPS_REG_A0,
		"a1": uc.MIPS_REG_A1,
		"a2": uc.MIPS_REG_A2,
		"a3": uc.MIPS_REG_A3,
		"t0": uc.MIPS_REG_T0,
		"t1": uc.MIPS_REG_T1,
		"t2": uc.MIPS_REG_T2,
		"t3": uc.MIPS_REG_T3,
		"t4": uc.MIPS_REG_T4,
		"t5": uc.MIPS_REG_T5,
		"t6": uc.MIPS_REG_T6,
		"t7": uc.MIPS_REG_T7,
		"t8": uc.MIPS_REG_T8,
		"t9": uc.MIPS_REG_T9,
		"s0": uc.MIPS_REG_S0,
		"s1": uc.MIPS_REG_S1,
		"s2": uc.MIPS_REG_S2,
		"s3": uc.MIPS_REG_S3,
		"s4": uc.MIPS_REG_S4,
		"s5": uc.MIPS_REG_S5,
		"s6": uc.MIPS_REG_S6,
		"s7": uc.MIPS_REG_S7,
		"s8": uc.MIPS_REG_S8,
		"k0": uc.MIPS_REG_K0,
		"k1": uc.MIPS_REG_K1,
		"gp": uc.MIPS_REG_GP,
		"sp": uc.MIPS_REG_SP,
		"ra": uc.MIPS_REG_RA,
	},
	DefaultRegs: []string{
		"at",
		"v0", "v1",
		"a0", "a1", "a2", "a3",
		"t0", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8", "t9",
		"s0", "s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8",
		"k0", "k1",
		"gp",
	},
}
