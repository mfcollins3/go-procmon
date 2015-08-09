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

import "io"

// ProcessMonitor is used to write debug messages to the Process Monitor log.
// Debug messages should be less than 2048 characters in order to display
// correctly in the log.
//
// ProcessMonitor is an io.Writer type and can be used with any code that
// supports writing to an io.Writer instance.
var ProcessMonitor Writer

func init() {
	ProcessMonitor = &nullProcessMonitor{}
}

// Writer defines an interface for writing debug messages to the Process
// Monitor log. Writer extends io.Writere and defines WriteString for
// efficiently wrriting strings directly to the Process Monitor log.
type Writer interface {
	io.Writer

	// Writes a string message to the Process Monitor log.
	WriteString(s string) (n int, err error)
}

// Default null implementation of the Process Monitor writer. This
// implementation is used on non-Windows platforms or on Windows hosts when
// Process Monitor is not installed.
type nullProcessMonitor struct{}

// Null implementation of i.Writer.Write. This function will return the length
// of the byte array, making the caller believe that the write was successful
// and that no error occurred.
func (pm *nullProcessMonitor) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Null implementation of WriteString. This function will return the length of
// the string, making the caller believe that the write was successful and that
// no error occurred.
func (pm *nullProcessMonitor) WriteString(s string) (n int, err error) {
	return len(s), nil
}
