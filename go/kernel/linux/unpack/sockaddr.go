package unpack

import (
	"encoding/binary"
	"syscall"

	"github.com/superp00t/usercorn/go/kernel/common"
)

const (
	AF_LOCAL   = 1
	AF_INET    = 2
	AF_INET6   = 10
	AF_PACKET  = 17
	AF_NETLINK = 16
)

type SockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte
}

type SockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte
	Scope_id uint32
}

type SockaddrLinklayer struct {
	Family   uint16
	Protocol uint16
	Ifindex  int32
	Hatype   uint16
	Pkttype  uint8
	Halen    uint8
	Addr     [8]uint8
}

type SockaddrNetlink struct {
	Family uint16
	Pad    uint16
	Pid    uint32
	Groups uint32
}

type SockaddrUnix struct {
	Family uint16
	Path   [108]byte
}

func Sockaddr(buf common.Buf, length int) syscall.Sockaddr {
	var port [2]byte
	order := buf.K.U.ByteOrder()
	// TODO: handle insufficient length
	var family uint16
	if err := buf.Unpack(&family); err != nil {
		return nil
	}
	// TODO: handle errors?
	st := buf.Struc()
	switch family {
	case AF_LOCAL:
		var a SockaddrUnix
		st.Unpack(&a)
		return sockaddrToNative(&a)
	case AF_INET:
		var a SockaddrInet4
		st.Unpack(&a)
		order.PutUint16(port[:], a.Port)
		a.Port = binary.BigEndian.Uint16(port[:])
		return sockaddrToNative(&a)
	case AF_INET6:
		var a SockaddrInet6
		st.Unpack(&a)
		order.PutUint16(port[:], a.Port)
		a.Port = binary.BigEndian.Uint16(port[:])
		return sockaddrToNative(&a)
	case AF_PACKET:
		var a SockaddrLinklayer
		st.Unpack(&a)
		return sockaddrToNative(&a)
	case AF_NETLINK:
		var a SockaddrNetlink
		st.Unpack(&a)
		return sockaddrToNative(&a)
	}
	return nil
}
