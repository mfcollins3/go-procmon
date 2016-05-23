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

/*
Package debugger implements tools to assist debugging Go programs in
development, testing, and production environments.

Initially, this package is focused on building debugging support for programs
running on the Microsoft Windows platforms. Microsoft Windows has an API
function named OutputDebugString that programs can use to send a short message
to an attached debugger to be included in the debugger's log or shown in the
debug console. There are also tools such as DebugView
(https://technet.microsoft.com/en-us/Library/bb896647.aspx) that is distributed
as part of Microsoft's Sysinternals Suite that can intercept debug messages at
runtime without the need to attach a debugger. This package implements support
for sending debug messages to OutputDebugString at runtime.
*/
package debugger
