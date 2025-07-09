package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/peterh/liner"
)

var builtins = []string{"cd", "exit", "help"}
var ErrNoPath = errors.New("path required")

func main() {

	line := liner.NewLiner()
	defer line.Close()

	line.SetCtrlCAborts(true)

	allCommands := getAllCommands()

	line.SetCompleter(func(line string) []string {
		var completions []string

		for _, cmd := range allCommands {
			if strings.HasPrefix(cmd, line) {
				completions = append(completions, cmd)
			}
		}

		return completions
	})

	for {
		cwd, _ := os.Getwd()

		input, err := line.Prompt(fmt.Sprintf("[%s] > ", cwd))
		if err != nil {
			if err == liner.ErrPromptAborted {
				continue
			}
			break
		}
		input = strings.TrimSpace(input)

		if input == "" {
			continue
		}

		line.AppendHistory(input)

		if err := execIn(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

	}
}

func execIn(input string) error {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil
	}

	switch args[0] {
	case "cd":
		if len(args) < 2 {
			return ErrNoPath
		}
		return os.Chdir(args[1])
	case "exit":
		os.Exit(0)
	case "help":
		fmt.Println("Built-in commands: cd, exit, help")
		return nil
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("command not found: %s", args[0])
		}
	}
	return err
}

func getAllCommands() []string {
	paths := strings.Split(os.Getenv("PATH"), ":")
	commands := make(map[string]bool)

	for _, dir := range paths {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, file := range files {
			commands[file.Name()] = true
		}
	}

	for _, cmd := range builtins {
		commands[cmd] = true
	}

	var list []string
	for cmd := range commands {
		list = append(list, cmd)
	}
	return list
}
