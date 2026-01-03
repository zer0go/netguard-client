package handler

import (
	"unsafe"

	"github.com/pkg/errors"
	"github.com/zer0go/netguard-client/internal/config"
	"golang.org/x/sys/windows"
)

var (
	modWireGuard                 = windows.NewLazyDLL("wireguard.dll")
	procWireGuardCreate          = modWireGuard.NewProc("WireGuardCreateAdapter")
	procWireGuardClose           = modWireGuard.NewProc("WireGuardCloseAdapter")
	procWireGuardSetAdapterState = modWireGuard.NewProc("WireGuardSetAdapterState")
)

func startNetworkAdapter() (err error, close func()) {
	interfaceName := config.Get().InterfaceName
	namePtr, _ := windows.UTF16PtrFromString(interfaceName)
	typePtr, _ := windows.UTF16PtrFromString("WireGuard")

	handle, _, err := procWireGuardCreate.Call(
		uintptr(unsafe.Pointer(namePtr)),
		uintptr(unsafe.Pointer(typePtr)),
		0,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create WireGuard adapter"), nil
	}

	_, _, err = procWireGuardSetAdapterState.Call(
		handle,
		uintptr(1),
	)

	if err != nil {
		return errors.Wrap(err, "failed to set WireGuard adapter state"), nil
	}

	return nil, func() {
		_, _, _ = procWireGuardClose.Call(handle)
	}
}
