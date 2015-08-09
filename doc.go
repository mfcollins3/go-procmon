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
Package procmon adds support for writing debug messages to Process Monitor
when running on a Microsoft Windows host.

Process Monitor (https://technet.microsoft.com/en-us/Library/bb896645.aspx)
is an invaluable tool for monitoring running applications and diagnosing
problems in testing and production environments where a debugger is not
available. For customer support scenarios, Process Monitor is extremely
valuable because remote customers can capture Process Monitor logs and
send the logs to developers to analyze offline.

To make Process Monitor log events more helpful, it is possible for
programs to output debug messages providing context to the events. This
makes it easier to understand where in a program an action is happening
or why a specific registry entry, file, or device is being accessed.
Using this library, Go program developers can output debug messages to
the Process Monitor log to provide this contextual information.
*/
package procmon
