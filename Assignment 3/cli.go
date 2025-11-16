package main

import (
	"bufio"
	"fmt"
	"io"
	"reflect"

	"os"
	"strings"
	"unicode"
)

// CLI globals
// Run cli set to false will exit the program
// Prompt is either empty or contains current connected server name
var (
	runCli bool   = true
	prompt string = ""
)

// Tokenize string input to array of strings
func parseArg(cmd string) []string {
	curToken := ""
	tokens := []string{}

	for _, c := range cmd {
		if unicode.IsSpace(c) {
			if curToken == "" {
				continue
			}
			tokens = append(tokens, curToken)
			curToken = ""
		} else {
			curToken = curToken + string(c)
		}
	}

	if curToken != "" {
		tokens = append(tokens, curToken)
	}

	return tokens
}

func startCli() {
	// Read user input from stdin
	reader := bufio.NewReader(os.Stdin)

	// Loop until "exit" command is executed
	for runCli {
		fmt.Printf("%s$ ", prompt)
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error: %v\n", err)
			continue
		}

		line = strings.TrimSpace(line)

		tokens := parseArg(line)
		if len(tokens) == 0 {
			continue
		}

		// Fetch command
		f, ok := cmds[tokens[0]]
		if !ok {
			fmt.Printf("command '%s' does not exist. See 'help' for available commands\n", tokens[0])
			continue
		}

		// Exec command
		err = f(tokens[1:])
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

// Generic wrapper for CLI command to handle arbitary arguments
func CallbackFunc(f any) func([]string) error {
	fVal := reflect.ValueOf(f)
	fType := fVal.Type()

	if fType.Kind() != reflect.Func {
		panic("Callback must be function")
	}

	return func(args []string) error {
		//if len(args) != fType.NumIn() {
		//	return fmt.Errorf("expected %d args, got %d", fType.NumIn(), len(args))
		//}

		values := make([]reflect.Value, fType.NumIn())

		for i := 0; i < fType.NumIn(); i++ {
			var arg string = ""
			if len(args) > i {
				arg = args[i]
			}

			param := fType.In(i)

			switch param.Kind() {
			case reflect.String:
				values[i] = reflect.ValueOf(arg)
			default:
				return fmt.Errorf("parameter type unsupported: %s", param)
			}
		}

		// Exec wrapped function
		_ = fVal.Call(values)

		return nil
	}
}
