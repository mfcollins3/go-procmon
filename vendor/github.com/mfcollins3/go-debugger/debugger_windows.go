/*
The MIT License (MIT)

Copyright (c) 2015 Michael F. Collins, III

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including but without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is furnished
to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package debugger

import (
	"bytes"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernelDLL              = windows.NewLazyDLL("kernel32.dll")
	outputDebugStringWProc = kernelDLL.NewProc("OutputDebugStringW")
)

func init() {
	Console = &windowsDebuggerWriter{&realOutputDebugString{}}
}

type outputDebugString interface {
	Write(message string) error
}

type realOutputDebugString struct{}

func (d *realOutputDebugString) Write(message string) error {
	msg, err := windows.UTF16PtrFromString(message)
	if nil != err {
		return err
	}

	_, _, err = outputDebugStringWProc.Call(uintptr(unsafe.Pointer(msg)))
	return err
}

type windowsDebuggerWriter struct {
	outputDebugString outputDebugString
}

func (w *windowsDebuggerWriter) Write(p []byte) (n int, err error) {
	_, err = w.WriteString(bytes.NewBuffer(p).String())
	n = 0
	if nil == err {
		n = len(p)
	}

	return
}

func (w *windowsDebuggerWriter) WriteString(s string) (n int, err error) {
	n = 0
	err = w.outputDebugString.Write(s)
	if nil == err {
		n = len(s)
	}

	return
}
