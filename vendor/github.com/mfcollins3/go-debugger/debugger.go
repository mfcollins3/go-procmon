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
	"fmt"
	"io"
)

// Console is used to write messages to the debugger. Debugger messages are
// strings that should convey some description of an action that is being
// performed. The action could be the name of a function that was called, or
// a function that is being called, an algorithm that is being executed, the
// value of a parameter for a function, or a value that was returned by a
// function. Console implements the io.Writer interface and can be passed to
// code that outputs to an io.Writer.
//
// Console also implements WriteString and can be used with io.WriteString to
// efficiently write a string value without having to marshal the string to a
// byte array and back to a string.
var Console Writer

// Writer defines an interface for writing status messages to a debugger.
// Writer extends io.Writer and defines WriteString() for writing strings
// directly to the debugger as well as supporting writing strings that have
// been encoded into byte arrays like Buffer objects.
type Writer interface {
	io.Writer

	// Writes a string message to the debugger. WriteString is defined to make
	// writing messages more efficient than using Write because when using
	// Write the messages have to be converted back to a string and then written
	// to the debugger. Using WriteString or io.WriteString is more efficient
	// because the strings will be written directly to the debugger without
	// the extra marshaling process having to occur.
	WriteString(s string) (n int, err error)
}

func init() {
	Console = &nullDebuggerWriter{}
}

// Default null implementation of a debug writer. This implementation will be
// used on platforms that do not support writing to a debugger.
type nullDebuggerWriter struct{}

// Null implementation of io.Write. This function will return the length of the
// byte array, making the caller believe that the write was successful and that
// no error occurred.
func (w *nullDebuggerWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// Null implementation of WriteString. This function will return the length of
// the string, making the caller believe that the write was successful and that
// no error occurred.
func (w *nullDebuggerWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

// Println formats a message using the default formats for the operands and
// outputs the message to the debugger console. It returns the number of
// characters that were output and any error that occurred.
func Println(a ...interface{}) (n int, err error) {
	return fmt.Fprint(Console, a...)
}

// Printf formats a debug message using the format specifier and writes the
// formatted message to the debugger console. It returns the number of
// characters that were output and any error that occurred.
func Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(Console, format, a...)
}
