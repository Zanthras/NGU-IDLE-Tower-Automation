// +build windows

package main

import (
	"log"
	"syscall"
	"unsafe"
)

type HWND uintptr

var (
	user32, _        = syscall.LoadLibrary("user32.dll")
	findWindowW, _   = syscall.GetProcAddress(user32, "FindWindowW")
	getWindowRect, _ = syscall.GetProcAddress(user32, "GetWindowRect")
)

func getNGULocation() RECT {
	defer syscall.FreeLibrary(user32)

	hwnd := FindWindowByTitle("NGU Idle")

	if hwnd > 0 {
		rect := GetWindowDimensions(hwnd)
		TOP = int(rect.Top) + BAR_OFFSET_TOP
		LEFT = int(rect.Left) + BAR_OFFSET_LEFT
		return *rect
	}
	log.Fatal("Cant find a window named 'NGU Idle'")
	return RECT{}
}

func FindWindowByTitle(title string) HWND {
	ret, _, _ := syscall.Syscall(
		findWindowW,
		2,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))),
		0,
	)
	return HWND(ret)
}

func GetWindowDimensions(hwnd HWND) *RECT {
	var rect RECT

	syscall.Syscall(
		getWindowRect,
		2,
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&rect)),
		0,
	)

	return &rect
}
