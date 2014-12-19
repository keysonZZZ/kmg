package kmgTapTun

import (
	"os"
	"strings"
	"syscall"
	"unsafe"
)

// Interface is a TUN/TAP interface.
type iInterface struct {
	deviceType DeviceType
	file       *os.File
	name       string
}

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTap(ifName string) (ifce Interface, err error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(file.Fd(), ifName, cIFF_TAP|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	ifce = &iInterface{deviceType: DeviceTypeTap, file: file, name: name}
	return
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTun(ifName string) (ifce Interface, err error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}
	name, err := createInterface(file.Fd(), ifName, cIFF_TUN|cIFF_NO_PI)
	if err != nil {
		return nil, err
	}
	ifce = &iInterface{deviceType: DeviceTypeTun, file: file, name: name}
	return
}

func (ifce *iInterface) GetDeviceType() DeviceType {
	return ifce.deviceType
}

// Returns the interface name of ifce, e.g. tun0, tap1, etc..
func (ifce *iInterface) Name() string {
	return ifce.name
}

// Implement io.Writer interface.
func (ifce *iInterface) Write(p []byte) (n int, err error) {
	return ifce.file.Write(p)
}

// Implement io.Reader interface.
func (ifce *iInterface) Read(p []byte) (n int, err error) {
	return ifce.file.Read(p)
}
func (ifce *iInterface) Close() (err error) {
	return ifce.file.Close()
}

const (
	cIFF_TUN   = 0x0001
	cIFF_TAP   = 0x0002
	cIFF_NO_PI = 0x1000
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

func createInterface(fd uintptr, ifName string, flags uint16) (createdIFName string, err error) {
	var req ifReq
	req.Flags = flags
	copy(req.Name[:], ifName)
	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, fd, uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		err = errno
		return
	}
	createdIFName = strings.Trim(string(req.Name[:]), "\x00")
	return
}
