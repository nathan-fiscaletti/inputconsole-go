package inputconsole

import (
    "fmt"
    "strings"

    "github.com/eiannone/keyboard"
)

// InputConsole represents a console in which the input text is kept
// separate from the output text. You will be able to type on the
// input line without the output lines interfering with it.
type InputConsole struct {
    inputString string
    prompt string
    unknownCommandHandler func(string)bool
    commands map[string]func([]string)
}

// NewInputConsole creates a new instance of an InputConsole.
func NewInputConsole() *InputConsole {
    return &InputConsole {
        inputString: "",
        prompt: "> ",
        unknownCommandHandler: nil,
        commands: map[string]func([]string){},
    }
}

// Writef will write a message to the InputConsole using the specified
// format and arguments.
func (ic *InputConsole) Writef(format string, vargs ...interface{}) {
    text := fmt.Sprintf(format, vargs...)
    output := fmt.Sprintf(
        "\r\033[K%s\n", strings.TrimRight(text, "\n"),
    )
    output = fmt.Sprintf("%s\r%s", output, ic.prompt)
    output = ic.parseInputString(output)
    fmt.Printf(output)
}

// ListenForInput will listen for input on a new thread using the
// specified prompt.
func (ic *InputConsole) ListenForInput(prompt string) {
    ic.prompt = prompt
    go ic.inputThread()
}

// RegisterCommand will register a command with the specified name and
// action. The action callback should take an array of strings that
// will represent the arguments passed to the command.
func (ic *InputConsole) RegisterCommand(
    name string, action func([]string),
) {
    ic.commands[name] = action
}

// SetUnknownCommandHandler will set the handler to use for unknown
// commands. The handler will be passed the full command string and
// should return a boolean indicating whether or not the command was
// handled by the handler.
func (ic *InputConsole) SetUnknownCommandHandler(
    handler func(string)bool,
) {
    ic.unknownCommandHandler = handler
}

func (ic *InputConsole) inputThread() {
    fmt.Printf("%s", ic.prompt)
    err := keyboard.Open()
    defer func() {
		_ = keyboard.Close()
	}()
    if err != nil {
        panic(err)
    }

    inputloop:
    for true {
        var r rune
        var key keyboard.Key
        var err error
        for {
            r, key, err = keyboard.GetKey()
            if err != nil {
                panic(err)
            }

            if key == keyboard.KeyBackspace || 
               key == keyboard.KeyBackspace2 {
                if ic.inputString != "" {
                    lastCharPos := len(ic.inputString)-1
                    ic.inputString = ic.inputString[:lastCharPos]
                    fmt.Printf("\b\033[K")
                }
            } else if key == keyboard.KeyArrowUp ||
                    key == keyboard.KeyArrowDown ||
                    key == keyboard.KeyArrowLeft ||
                    key == keyboard.KeyArrowRight {
                fmt.Printf("\033[C")
            } else if key == keyboard.KeySpace {
                ic.inputString = fmt.Sprintf("%s ", ic.inputString)
            } else if r != '\x00' {
                ic.inputString = fmt.Sprintf(
                    "%s%c", ic.inputString, r,
                )
            }

            fmt.Printf(ic.parseInputString(
                fmt.Sprintf("\r%s", ic.prompt),
            ))

            if key == keyboard.KeyEnter || key == keyboard.KeyCtrlC {
                break
            }
        }

        if key == keyboard.KeyCtrlC {
            // _ = keyboard.Close()
            break inputloop
        }

        if key == keyboard.KeyEnter {
            ret := ic.inputString
            key = 0
            ic.inputString = ""
            ic.handleCommand(strings.TrimRight(ret, "\r\n"))
        }
    }
}

func (ic *InputConsole) handleCommand(command string) {
    commandComponents := strings.Split(command, " ")
    commandName := commandComponents[0]
    for name,action := range ic.commands {
        if name == commandName {
            defer func() {
                if r := recover(); r != nil {
                    ic.Writef(
                        "Failed to run command '%s': %v",
                        commandName, r,
                    )
                }
            }()
            action(commandComponents[1:])
            return
        }
    }

    if ic.unknownCommandHandler != nil {
        if ic.unknownCommandHandler(command) {
            return
        }
    }

    ic.Writef(fmt.Sprintf("Unknown command: %s\n", command))
}

func (ic *InputConsole) parseInputString(prefix string) string {
    printableInputString := ic.inputString
    _inputStringComponents := strings.Split(ic.inputString, " ")

    // Replicate the same style of string splitting that we use in
    // python to make this a little easier
    var inputStringComponents []string
    for _,str := range _inputStringComponents {
        if str != "" {
            inputStringComponents = append(inputStringComponents, str)
        }
    }

    if len(inputStringComponents) > 0 {
        spaceChar := ""
        if len(inputStringComponents) > 1 {
            spaceChar = " "
        }

        hasCommand := false
        for name,_ := range ic.commands {
            if name == inputStringComponents[0] {
                hasCommand = true
            }
        }

        if hasCommand {
            printableInputString = fmt.Sprintf(
                "\033[92m%s\033[0m%s%s",
                inputStringComponents[0],
                spaceChar,
                strings.Join(inputStringComponents[1:], " "),
            )
        } else {
            printableInputString = fmt.Sprintf(
                "\033[91m%s\033[0m%s%s",
                inputStringComponents[0],
                spaceChar,
                strings.Join(inputStringComponents[1:], " "),
            )
        }

        if strings.HasSuffix(ic.inputString, " ") {
            printableInputString = fmt.Sprintf(
                "%s ", printableInputString,
            )
        }
    }

    return fmt.Sprintf("%s%s", prefix, printableInputString)
}