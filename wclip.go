package main

import (
	// "bytes"
	"fmt"

	"github.com/JamesHovious/w32"

	// "io"
	"os"
	"reflect"
	"syscall"

	// "unicode/utf16"
	"unsafe"
)

func main() {
	w32.OpenClipboard(0)
	defer w32.CloseClipboard()
	hData := w32.GetClipboardData(w32.CF_UNICODETEXT)
	if hData == 0 {
		fmt.Fprint(os.Stderr, "Clipboard connot be converted to text.")
		return
	}
	fmt.Print(ClipboardToString(hData))
}

func ClipboardToString(hData w32.HANDLE) string {

	hgData := (w32.HGLOBAL)(hData)
	w32.GlobalLock(hgData)
	defer w32.GlobalUnlock(hgData)

	// Ideally we would find the len of HGLOBAL with (W32API) GlobalSize
	// But it appears to be unimplemented in the w32 package as of writing.
	wcLen := FindUint16Null(unsafe.Pointer(hData)) / 2

	buffer := cArray2Uint16Slice(unsafe.Pointer(hData), int(wcLen))
	return syscall.UTF16ToString(buffer)
}

// cArray2Uint16Slice creates a uint16 slice from an unsafe.Pointer
//from https://gist.github.com/nasitra/98bb59421be49a518c4a
func cArray2Uint16Slice(array unsafe.Pointer, length int) (list []uint16) {
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = length
	sliceHeader.Len = length
	sliceHeader.Data = uintptr(array)
	return list
}

// FindUint16Null returns the byte offset from inPtr to the first byteword aligned uint16(0) - [0x00 0x00]
func FindUint16Null(inPtr unsafe.Pointer) uint {
	match := uint16(0)
	for i := uint(0); ; i += 2 {
		matchBWord := *(*uint16)(unsafe.Pointer(uintptr(inPtr) + uintptr(i)))
		// fmt.Print(*(*[2]byte)(unsafe.Pointer(&matchBWord)), " ")
		if matchBWord == match {
			// fmt.Println("\ninPtr:", inPtr, "\noffset:", i, "\nEnd:", unsafe.Pointer(uintptr(inPtr)+uintptr(i)))
			return i
		}

	}
}
