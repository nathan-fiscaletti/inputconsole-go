# InputConsole

[![GoDoc](https://godoc.org/github.com/nathan-fiscaletti/inputconsole-go?status.svg)](https://godoc.org/github.com/nathan-fiscaletti/inputconsole-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/nathan-fiscaletti/inputconsole-go)](https://goreportcard.com/report/github.com/nathan-fiscaletti/inputconsole-go)

**InputConsole** is a console for Go that will keep all output above the input line without interrupting the input line.

[Looking for the Python version?](https://github.com/nathan-fiscaletti/inputconsole)

## Install

```sh
$ go get github.com/nathan-fiscaletti/inputconsole
```

## Usage

```go
import(
    "time"

    "github.com/nathan-fiscaletti/inputconsole-go"
)

func main() {
    console := inputconsole.NewInputConsole()

    // Register a command.
    // Runtime exceptions caused by commands are automatically caught
    // and an error message will be written to the inputconsole.
    console.RegisterCommand("help", func(params []string) {
        console.Writef("I don't want to help you %s", params[0])
    })

    // Start listening for input on a new thread
    // Input line will always stay at the bottom
    console.ListenForInput("> ")

    // Set unknown command handler, Return 'true' for command handled
    // or 'false' for command not handled.
    console.SetUnknownCommandHandler(func(command string) bool {
        console.Writef("Unknown command: %s\n", command)
        return true
    })

    // Generate random output to keep the output thread active.
    go func() {
        var cnt int = 0
        for {
            console.Writef("This is a test message: %d\n", cnt)
            time.Sleep(time.Second)
            cnt = cnt + 1
        }
    }()

    // Keep the process alive
    for true {
        time.Sleep(time.Second)
    }
}
```