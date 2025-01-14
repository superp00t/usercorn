package x86_16

import (
	"fmt"
	"strings"
	"syscall"

	"github.com/lunixbochs/argjoy"
	"github.com/pkg/errors"

	co "github.com/superp00t/usercorn/go/kernel/common"
	"github.com/superp00t/usercorn/go/models"
	uc "github.com/unicorn-engine/unicorn/bindings/go/unicorn"
)

const (
	STACK_BASE = 0x8000
	STACK_SIZE = 0x1000
	NUM_FDS    = 256

	// Registers
	AH    = uc.X86_REG_AH
	AL    = uc.X86_REG_AL
	AX    = uc.X86_REG_AX
	BH    = uc.X86_REG_BH
	BL    = uc.X86_REG_BL
	BP    = uc.X86_REG_BP
	BX    = uc.X86_REG_BX
	CH    = uc.X86_REG_CH
	CL    = uc.X86_REG_CL
	CS    = uc.X86_REG_CS
	CX    = uc.X86_REG_CX
	DH    = uc.X86_REG_DH
	DI    = uc.X86_REG_DI
	DL    = uc.X86_REG_DL
	DS    = uc.X86_REG_DS
	DX    = uc.X86_REG_DX
	ES    = uc.X86_REG_ES
	FS    = uc.X86_REG_FS
	GS    = uc.X86_REG_GS
	IP    = uc.X86_REG_IP
	SI    = uc.X86_REG_SI
	SP    = uc.X86_REG_SP
	SS    = uc.X86_REG_SS
	FLAGS = uc.X86_REG_EFLAGS
)

func (k *DosKernel) reg16(enum int) uint16 {
	r, _ := k.U.RegRead(enum)
	return uint16(r)
}
func (k *DosKernel) reg8(enum int) uint8 {
	r, _ := k.U.RegRead(enum)
	return uint8(r)
}
func (k *DosKernel) wreg16(enum int, val uint16) {
	k.U.RegWrite(enum, uint64(val))
}
func (k *DosKernel) wreg8(enum int, val uint8) {
	k.U.RegWrite(enum, uint64(val))
}
func (k *DosKernel) setFlagC(set bool) {
	// TODO: Write setFlagX with enum for each flag
	// Unicorn doesn't have the non-extended FLAGS register, so we're
	// dealing with 32 bits here
	flags, _ := k.U.RegRead(FLAGS)
	if set {
		flags |= 1 // CF = 1
	} else {
		flags &= 0xfffffffe // CF = 0
	}
	k.U.RegWrite(FLAGS, flags)
}

var dosSysNum = map[int]string{
	0x00: "terminate",
	0x01: "char_in",
	0x02: "char_out",
	0x09: "display",
	0x30: "get_dos_version",
	0x3C: "create_or_truncate",
	0x3D: "open",
	0x3E: "close",
	0x3F: "read",
	0x40: "write",
	0x41: "unlink",
	0x4C: "terminate_with_code",
}

type abi []int

// ABI to syscall number mapping
var abiMap = map[*abi][]int{
	&abi{BX, DX, CX}: {0x00, 0x30, 0x3E, 0x3F, 0x40},
	&abi{DX, CX}:     {0x01, 0x02, 0x09, 0x3C, 0x41},
	&abi{DX, AL}:     {0x3D},
	&abi{AL}:         {0x4C},
}

var syscallAbis = map[int][]int{}

type PSP struct {
	CPMExit                     [2]byte
	FirstFreeSegment            uint16
	Reserved1                   uint8
	CPMCall5Compat              [5]byte
	OldTSRAddress               uint32
	OldBreakAddress             uint32
	CriticalErrorHandlerAddress uint32
	CallerPSPSegment            uint16
	JobFileTable                [20]byte
	EnvironmentSegment          uint16
	INT21SSSP                   uint32
	JobFileTableSize            uint16
	JobFileTablePointer         uint32
	PreviousPSP                 uint32
	Reserved2                   uint32
	DOSVersion                  uint16
	Reserved3                   [14]byte
	DOSFarCall                  [3]byte
	Reserved4                   uint16
	ExtendedFCB1                [7]byte
	FCB1                        [16]byte
	FCB2                        [20]byte
	CommandLineLength           uint8
	CommandLine                 [127]byte
}

type DosKernel struct {
	*co.KernelBase
	fds map[int]int
}

func initPsp(argc int, argv []string) *PSP {
	psp := &PSP{
		CPMExit:    [2]byte{0xcd, 0x20},       // int 0x20
		DOSFarCall: [3]byte{0xcd, 0x21, 0xcd}, // int 0x21 + retf
	}

	psp.FCB1[0] = 0x01
	psp.FCB1[1] = 0x20

	// Combine all args into one string
	commandline := strings.Join(argv, " ")
	copy(psp.CommandLine[:126], commandline)
	if len(commandline) > 126 {
		psp.CommandLineLength = 126
	} else {
		psp.CommandLineLength = uint8(len(commandline))
	}

	return psp
}

func (k *DosKernel) readUntilChar(addr uint64, c byte) []byte {
	var mem []byte
	var i uint64
	var char byte = 0

	// TODO: Read ahead? This'll be slow
	for i = 1; char != c || i == 1; i++ {
		mem, _ = k.U.MemRead(addr, i)
		char = mem[i-1]
	}
	return mem[:i-2]
}

func (k *DosKernel) getFd(fd int) (uint16, error) {
	for i := 0; i < NUM_FDS; i++ {
		if _, ok := k.fds[i]; !ok {
			k.fds[i] = fd
			return uint16(i), nil
		}
	}
	return 0xFFFF, errors.New("DOS FD table exhausted")
}

func (k *DosKernel) freeFd(fd int) (int, error) {
	realfd, ok := k.fds[fd]
	if !ok {
		return 0xFFFF, errors.New("FD not found in FD table")
	}
	delete(k.fds, fd)
	return realfd, nil
}

func (k *DosKernel) Terminate() {
	k.U.Exit(models.ExitStatus(0))
}

func (k *DosKernel) CharIn(buf co.Buf) byte {
	var char byte
	fmt.Scanf("%c", &char)
	k.U.MemWrite(buf.Addr, []byte{char})
	return char
}

func (k *DosKernel) CharOut(char uint16) byte {
	fmt.Printf("%c", byte(char&0xFF))
	return byte(char & 0xFF)
}

func (k *DosKernel) Display(buf co.Buf) int {
	mem := k.readUntilChar(buf.Addr, '$')

	syscall.Write(1, mem)
	k.wreg8(AL, 0x24)
	return 0x24
}

func (k *DosKernel) GetDosVersion() int {
	k.wreg16(AX, 0x7)
	return 0x7
}

func (k *DosKernel) openFile(filename string, mode int) co.Fd {
	realfd, err := syscall.Open(filename, mode, 0666)
	if err != nil {
		k.wreg16(AX, 0xFFFF)
		k.setFlagC(true)
		return 0xFFFF
	}

	// Find an internal fd number
	dosfd, err := k.getFd(realfd)
	if err != nil {
		k.wreg16(AX, dosfd)
		k.setFlagC(true)
		return 0xFFFF
	}
	k.setFlagC(false)
	k.wreg16(AX, dosfd)
	return co.Fd(dosfd)
}

func (k *DosKernel) CreateOrTruncate(buf co.Buf, attr int) co.Fd {
	filename := string(k.readUntilChar(buf.Addr, '$'))
	return k.openFile(filename, syscall.O_CREAT|syscall.O_TRUNC|syscall.O_RDWR)
}

func (k *DosKernel) Open(filename string, mode int) co.Fd {
	return k.openFile(filename, mode)
}

func (k *DosKernel) Close(fd co.Fd) {
	// Find and free the internal fd
	realfd, _ := k.freeFd(int(fd))
	err := syscall.Close(realfd)
	if err != nil {
		k.setFlagC(true)
		// TODO: Set AX to error code
		k.wreg16(AX, 0xFFFF)
	}
	k.setFlagC(false)
	k.wreg16(AX, 0)
}

func (k *DosKernel) Read(fd co.Fd, buf co.Obuf, len co.Len) int {
	mem := make([]byte, len)
	n, err := syscall.Read(int(fd), mem)
	if err != nil {
		k.setFlagC(true)
		// TODO: Set AX to error code
		k.wreg16(AX, 0xFFFF)
	}
	k.U.MemWrite(buf.Addr, mem)
	k.setFlagC(false)
	k.wreg16(AX, uint16(n))
	return n
}

func (k *DosKernel) Write(fd co.Fd, buf co.Buf, n co.Len) int {
	mem, _ := k.U.MemRead(buf.Addr, uint64(n))
	realfd, ok := k.fds[int(fd)]
	if !ok {
		k.setFlagC(true)
		// TODO: Set AX to error code
		k.wreg16(AX, 0xFFFF)
	}
	written, err := syscall.Write(realfd, mem)
	if err != nil {
		k.setFlagC(true)
		// TODO: Set AX to error code
		k.wreg16(AX, 0xFFFF)
	}
	k.setFlagC(false)
	k.wreg16(AX, uint16(written))
	return written
}

func (k *DosKernel) Unlink(filename string, attr int) int {
	err := syscall.Unlink(filename)
	if err != nil {
		k.setFlagC(true)
		k.wreg16(AX, 0xFFFF)
		return 0xFFFF
	}
	k.setFlagC(false)
	k.wreg16(AX, 0)
	return 0
}

func (k *DosKernel) TerminateWithCode(code int) {
	k.U.Exit(models.ExitStatus(code))
}

func NewKernel() *DosKernel {
	k := &DosKernel{
		KernelBase: &co.KernelBase{},
		fds: map[int]int{
			0: 0,
			1: 1,
			2: 2,
		},
	}
	// Invert the syscall map
	for abi, syscalls := range abiMap {
		for _, syscall := range syscalls {
			syscallAbis[syscall] = *abi
		}
	}
	k.Argjoy.Register(k.getDosArgCodec())
	return k
}

func DosInit(u models.Usercorn, args, env []string) error {
	// Setup PSP
	// TODO: Setup args
	psp := initPsp(0, nil)
	u.StrucAt(0).Pack(psp)

	u.SetEntry(0x100)
	return u.MapStack(STACK_BASE, STACK_SIZE, false)
}

func DosSyscall(u models.Usercorn) {
	num, _ := u.RegRead(AH)
	name, _ := dosSysNum[int(num)]
	// TODO: How are registers numbered from here?
	u.Syscall(int(num), name, dosArgs(u, int(num)))
	// TODO: Set error
}

func (k *DosKernel) getDosArgCodec() func(interface{}, []interface{}) error {
	return func(arg interface{}, vals []interface{}) error {
		// DOS takes address as DS+DX
		if reg, ok := vals[0].(uint64); ok && len(vals) > 1 {
			ds, _ := k.U.RegRead(DS)
			reg += ds
			switch v := arg.(type) {
			case *co.Buf:
				*v = co.NewBuf(k, reg)
			case *co.Obuf:
				*v = co.Obuf{co.NewBuf(k, reg)}
			case *co.Ptr:
				*v = co.Ptr(reg)
			case *string:
				s, err := k.U.Mem().ReadStrAt(reg)
				if err != nil {
					return errors.Wrapf(err, "ReadStrAt(%#x) failed", reg)
				}
				*v = s
			default:
				return argjoy.NoMatch
			}
			return nil
		}
		return argjoy.NoMatch
	}
}

func dosArgs(u models.Usercorn, num int) func(int) ([]uint64, error) {
	// Return a closure over the correct arglist based on the syscall number
	return co.RegArgs(u, syscallAbis[num])
}

func DosInterrupt(u models.Usercorn, cause uint32) {
	intno := cause & 0xFF
	if intno == 0x21 {
		DosSyscall(u)
	} else if intno == 0x20 {
		u.Syscall(0, "terminate", func(int) ([]uint64, error) { return []uint64{}, nil })
	} else {
		panic(fmt.Sprintf("unhandled X86 interrupt %#X", intno))
	}
}
func DosKernels(u models.Usercorn) []interface{} {
	return []interface{}{NewKernel()}
}

func init() {
	Arch.RegisterOS(&models.OS{
		Name:      "DOS",
		Init:      DosInit,
		Interrupt: DosInterrupt,
		Kernels:   DosKernels,
	})
}
