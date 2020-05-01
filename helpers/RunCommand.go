package helpers

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"text/template"
)

func RunCommand(commandTemplate string, arguments interface{}) (int, string, error) {
	commandString := MustCreateString(commandTemplate, arguments)
	// Run the command and return the exit code
	commandPieces := strings.Fields(commandString)
	var args []string
	if len(commandString) > 0 {
		for _, value := range commandPieces {
			args = append(args, value)
		}
		name, args := args[0], args[1:]
		cmd := exec.Command(name, args...)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					return status.ExitStatus(), "", err
				}
			}
			return 0, "", err
		}
		return 0, out.String(), nil
	}
	return 0, "", errors.New("invalid command")
}

func MustCreateString(str string, vars interface{}) string {
	if len(str) == 0 {
		log.Fatal("Must pass a string of longer length.")
	}
	// Create template
	tpl, err := template.New("string").Parse(str)
	if err != nil {
		log.Fatal(err)
	}
	// Use vars to create the string
	var strBuffer bytes.Buffer
	err = tpl.Execute(&strBuffer, vars)
	if err != nil {
		log.Fatal(err)
	}
	return strBuffer.String()
}
