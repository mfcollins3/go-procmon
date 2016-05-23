Debugging Tools for Go
======================
This project implements tools that are helpful when debugging Go
programs running in development, testing, or production environments.

Initially, this project is focused on implementing debugging support for
Go programs running on Microsoft Windows. AS the project develops, I
will be looking at adding more support for Linux and Mac programs too.

Get the Code
------------
Getting the code is easy using the `go get` command. In a terminal or
command prompt window, execute the following command:

    $ go get github.com/mfcollins3/go-debugger

To use the debugging library in your own programs or libraries, add The
following `import` statement at the top of your programs:

```go
import debugger "github.com/mfcollins3/go-debugger"
```

Features
--------
### OutputDebugString Support (Microsoft Windows)

#### Writing Debug Messages

The `go-debugger` library supports writing messages at runtime to an
attached debugger or [DebugView](https://technet.microsoft.com/en-us/Library/bb896647.aspx)
using the Microsoft Windows [OutputDebugString](https://msdn.microsoft.com/en-us/library/windows/desktop/aa363362(v=vs.85).aspx)
API. This feature has been designed to support the [io.Writer](http://golang.org/pkg/io/#Writer)
interface because there is so much support within the Go framework and
third-party libaries for writing to `io.Writer` objects.

Writing to the debugger is easy:

```go
import debugger "github.com/mfcollins3/go-debugger"

func main() {
  debugger.Println("The program is starting")

  // TODO: implement the program

  debugger.Println("The program is ending")  
}
```

In addition, `debugger.Console` can be used with functions that rely on
`io.Writer` such as the `fmt` package:

```go
import debugger "github.com/mfcollins3/go-debugger"

func main() {
  fmt.Fprint(debugger.Console, "The program is starting")

  var firstName = "Michael", lastName = "Collins"
  fmt.Fprintf(debugger.Console, "Hello, %s %s", firstName, lastName)

  fmt.Fprint(debugger.Console, "The program is ending")
}
```

On unsupported platforms (currently anything that isn't Microsoft
Windows), `debugger.Console` is backed by a null implementation. This
code will run successfully with minimal performance impact on those
platforms.

#### Capturing Debugger Messages

In some environments, users may not have DebugView or a debugger
installed. To help out in these situations, the `go-debugger` library
implements the debugger side of the protocol. This allows you to capture
debug messages from running processes and log the messages or do
processing on the debugger messages.

To capture debug messages, you will use the OutputDebugStringReceiver
type and you will receive the messages on a channel.

```go
import (
  "log"
  "os"
  "os/signal"

  debugger "github.com/mfcollins3/go-debugger"
)

func main() {
  interruptChannel := make(chan os.Signal, 1)
  signal.Notify(interruptChannel, os.Interrupt)
  messageChannel := make(chan debugger.DebugMessage, 1)
  receiver, err := debugger.NewOutputDebugStringReceiver(messageChannel)
  if nil != err {
    log.Fatal(err)
  }

loop:
  for {
    select {
    case message := <-messageChannel:
      log.Printf("[%d] %s\r\n", message.ProcessID, message.Message)
    case <-interruptChannel:
      receiver.Close()
      break loop
    }
  }
}
```
