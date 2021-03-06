package usb

import (
	"errors"
	"fmt"
	"io"
)

const (
	VendorT1            = 0x534c
	ProductT1Firmware   = 0x0001
	VendorT2            = 0x1209
	ProductT2Bootloader = 0x53C0
	ProductT2Firmware   = 0x53C1

	VendorArchos               = 0x0E79
	ProductSafeTminiBootloader = 0x6001
	ProductSafeTminiFirmware   = 0x6000
)

var (
	ErrNotFound = fmt.Errorf("device not found")
)

type Info struct {
	Path      string
	VendorID  int
	ProductID int
}

type Device interface {
	io.ReadWriteCloser
}

type Bus interface {
	Enumerate() ([]Info, error)
	Connect(path string) (Device, error)
	Has(path string) bool
}

type USB struct {
	buses []Bus
}

func Init(buses ...Bus) *USB {
	return &USB{
		buses: buses,
	}
}

func (b *USB) Enumerate() ([]Info, error) {
	var infos []Info

	for _, b := range b.buses {
		l, err := b.Enumerate()
		if err != nil {
			return nil, err
		}
		infos = append(infos, l...)
	}
	return infos, nil
}

func (b *USB) Connect(path string) (Device, error) {
	for _, b := range b.buses {
		if b.Has(path) {
			return b.Connect(path)
		}
	}
	return nil, ErrNotFound
}

var errDisconnect = errors.New("Device disconnected during action")
var errClosedDevice = errors.New("Closed device")
