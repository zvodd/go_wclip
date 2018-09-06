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
	wclen := FindUint16Null(unsafe.Pointer(hData)) / 2
	// fmt.Println("\nLength of WCHAR string:", wclen)

	buffer := cArray2Uint16Slice(unsafe.Pointer(hData), int(wclen))
	// fmt.Println(buffer)
	return syscall.UTF16ToString(buffer)
}

//from https://gist.github.com/nasitra/98bb59421be49a518c4a
func cArray2Uint16Slice(array unsafe.Pointer, len int) (list []uint16) {
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(array)
	return list
}

func FindUint16Null(inPtr unsafe.Pointer) uint {
	buf := make([]byte, 2)
	for i := uint(0); ; i += 1 {
		for j := uint(0); j < 2; j += 1 {
			onebyte := *(*[1]byte)(unsafe.Pointer(uintptr(inPtr) + uintptr(i+j)))
			buf[j] = onebyte[0]
		}
		// fmt.Println("i:", i, " : ", buf)

		if buf[0] == 0 && buf[1] == 0 {
			return i + 1
		}
	}
}

// // Seems unneccisary
// func UnsafeFindByteSequence(inPtr unsafe.Pointer, sequence []byte) uint {
//  size := uint(len(sequence))
//  buf := make([]byte, size)
//  for i := uint(0); ; i++ {
//      for j := uint(0); j < size; j++ {
//          offset := uintptr(i + j)
//          ptr := uintptr(inPtr)
//          onebyte := *(*byte)(unsafe.Pointer(ptr + offset))
//          buf[j] = onebyte
//      }
//      if bytes.Equal(buf, sequence) {
//          return i + size
//      }
//  }
// }

// //from https://gist.github.com/nasitra/98bb59421be49a518c4a
// func cArray2ByteSlice(array unsafe.Pointer, len int) (list []byte) {
//     sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
//     sliceHeader.Cap = len
//     sliceHeader.Len = len
//     sliceHeader.Data = uintptr(array)
//     return list
// }

// // Only works for 2 byte aligned codepoints
// func FprintUTF16CString(output io.Writer, hData w32.HANDLE) (strlen uint) {
// 	hgData := (w32.HGLOBAL)(hData)
// 	w32.GlobalLock(hgData)

// 	// Should use winapi GlobalSize to set the limit of this loop.
// 	for i := uint(0); ; i += 4 {
// 		ptrAdder := uintptr(unsafe.Pointer(hData)) + uintptr(i)
// 		u16codepoint := *(*[2]uint16)(unsafe.Pointer(ptrAdder))
// 		if u16codepoint[0] == 0 {
// 			break
// 		} else if u16codepoint[1] == 0 {
// 			strlen++
// 			fmt.Fprint(output, string(utf16.Decode(u16codepoint[:1])))
// 			break
// 		} else {
// 			fmt.Fprint(output, string(utf16.Decode(u16codepoint[:])))
// 			strlen += 2
// 		}
// 	}
// 	// fmt.Println("\nLength of WCHAR string:", strlen)
// 	fmt.Println("")
// 	w32.GlobalUnlock(hgData)
// 	return
// }
