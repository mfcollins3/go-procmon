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
	"encoding/binary"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

// DebugMessage stores the information for a debug message that was output by a
// program for an attached or observing debugger.
type DebugMessage struct {
	// The process identifier for the process that published the message.
	ProcessID uint32

	// The text of the debug message.
	Message string
}

type closeableEvent interface {
	Close() error
}

type settableEvent interface {
	closeableEvent
	Set() error
}

type waitableEvent interface {
	closeableEvent
	Wait(timeout uint32) (uint32, error)
}

type windowsEvent struct {
	handle windows.Handle
}

func newWindowsEvent(name string) (*windowsEvent, error) {
	str, err := windows.UTF16PtrFromString(name)
	if nil != err {
		return nil, err
	}

	handle, err := windows.CreateEvent(nil, 0, 0, str)
	if nil != err {
		return nil, err
	}

	return &windowsEvent{handle: handle}, nil
}

func (e *windowsEvent) Close() error {
	return windows.CloseHandle(e.handle)
}

func (e *windowsEvent) Set() error {
	return windows.SetEvent(e.handle)
}

func (e *windowsEvent) Wait(timeout uint32) (uint32, error) {
	return windows.WaitForSingleObject(e.handle, timeout)
}

// OutputDebugStringReceiver implements the debugger-side of the windows
// OutputDebugString API. Programs can use OutputDebugStringReceiver to receive
// messages from other programs using OutputDebugString.
// OutputDebugStringReceiver will hold a channel and will publish received
// messages to the channel.
type OutputDebugStringReceiver struct {
	err              error
	buffer           windows.Handle
	view             uintptr
	bufferReadyEvent settableEvent
	dataReadyEvent   waitableEvent
	done             chan struct{}
	completed        chan struct{}
	messageChannel   chan DebugMessage
}

// NewOutputDebugStringReceiver creates an OutputDebugString receiver and
// begins receiving debug messages.
func NewOutputDebugStringReceiver(messageChannel chan DebugMessage) (*OutputDebugStringReceiver, error) {
	receiver := &OutputDebugStringReceiver{
		done:           make(chan struct{}),
		completed:      make(chan struct{}),
		messageChannel: messageChannel,
	}
	receiver.createMemoryMappedFile()
	receiver.createBufferReadyEvent()
	receiver.createDataReadyEvent()
	if nil != receiver.err {
		receiver.Close()
	} else {
		go receiver.receiveMessages()
	}

	return receiver, receiver.err
}

func (r *OutputDebugStringReceiver) createMemoryMappedFile() {
	if nil != r.err {
		return
	}

	var str *uint16
	str, r.err = windows.UTF16PtrFromString("DBWIN_BUFFER")
	if nil != r.err {
		return
	}

	r.buffer, r.err = windows.CreateFileMapping(
		windows.InvalidHandle,
		nil,
		0x4,
		0,
		4096,
		str)
	if nil != r.err {
		return
	}

	r.view, r.err = windows.MapViewOfFile(r.buffer, 0x4, 0, 0, 0)
}

func (r *OutputDebugStringReceiver) createBufferReadyEvent() {
	if nil != r.err {
		return
	}

	r.bufferReadyEvent, r.err = newWindowsEvent("DBWIN_BUFFER_READY")
}

func (r *OutputDebugStringReceiver) createDataReadyEvent() {
	if nil != r.err {
		return
	}

	r.dataReadyEvent, r.err = newWindowsEvent("DBWIN_DATA_READY")
}

// Close will release all of the objects and memory needed to process debug
// messages from other processes.
func (r *OutputDebugStringReceiver) Close() error {
	if nil == r.err {
		r.done <- struct{}{}
		<-r.completed
	}

	var result, err error
	result = r.dataReadyEvent.Close()

	err = r.bufferReadyEvent.Close()
	if nil != result {
		result = err
	}

	if 0 != r.view {
		err = windows.UnmapViewOfFile(r.view)
		if nil != result {
			result = err
		}
	}

	if windows.InvalidHandle != r.buffer {
		err = windows.CloseHandle(r.buffer)
		if nil != result {
			result = err
		}
	}

	return result
}

func (r *OutputDebugStringReceiver) receiveMessages() {
loop:
	for {
		r.err = r.bufferReadyEvent.Set()
		if nil != r.err {
			return
		}

		var event uint32
		event, r.err = r.dataReadyEvent.Wait(500)
		if nil != r.err {
			return
		}

		if 0 == event {
			b := (*[4096]byte)(unsafe.Pointer(r.view))
			i := 4
			for b[i] != 0 {
				i++
			}

			var message DebugMessage
			message.ProcessID = binary.LittleEndian.Uint32(b[0:4])
			message.Message = string(b[4:i])
			r.messageChannel <- message
		}

		select {
		case <-time.After(time.Millisecond * 200):
			continue
		case <-r.done:
			break loop
		}
	}

	r.completed <- struct{}{}
}
