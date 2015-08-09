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

package procmon

import (
	"bytes"
	"fmt"
	"io"
	"unsafe"

	"github.com/mfcollins3/go-debugger"

	"golang.org/x/sys/windows"
)

func init() {
	io.WriteString(debugger.Console,
		"Detecting whether Process Monitor is installed")

	path, err := windows.UTF16PtrFromString(`\\.\Global\ProcmonDebugLogger`)
	if nil != err {
		fmt.Fprintf(debugger.Console, "init failed: %v", err)
		return
	}

	processMonitorHandle, err := windows.CreateFile(
		path,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE|windows.FILE_SHARE_DELETE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0)
	if nil != err {
		fmt.Fprintf(debugger.Console, "CreateFile failed: %v", err)
		return
	}

	if windows.InvalidHandle == processMonitorHandle {
		io.WriteString(debugger.Console, "Process Monitor is not installed")
		return
	}

	ProcessMonitor =
		&processMonitorWriter{&realProcessMonitor{processMonitorHandle}}
	io.WriteString(debugger.Console, "Process Monitor debug logging is enabled")
}

const debugOut uint32 = (0x9535 << 16) | (0x2 << 14) | (0x81 << 2)

type processMonitor interface {
	Write(message string) error
}

type realProcessMonitor struct {
	handle windows.Handle
}

func (pm *realProcessMonitor) Write(message string) error {
	msg, err := windows.UTF16PtrFromString(message)
	if nil != err {
		return err
	}

	var bytesReturned uint32
	err = windows.DeviceIoControl(pm.handle,
		debugOut,
		(*byte)(unsafe.Pointer(msg)),
		uint32((len(message)+1)*2),
		nil,
		0,
		&bytesReturned,
		nil)
	if nil != err {
		fmt.Fprintf(debugger.Console, "PROCESS MONITOR ERROR: %v", err)
	}

	return err
}

type processMonitorWriter struct {
	processMonitor processMonitor
}

func (w *processMonitorWriter) Write(p []byte) (n int, err error) {
	_, err = w.WriteString(bytes.NewBuffer(p).String())
	n = 0
	if nil == err {
		n = len(p)
	}

	return
}

func (w *processMonitorWriter) WriteString(s string) (n int, err error) {
	err = w.processMonitor.Write(s)
	n = 0
	if nil == err {
		n = len(s)
	}

	return
}
