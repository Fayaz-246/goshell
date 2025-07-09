package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		cwd, _ := os.Getwd()
		fmt.Printf("%s > ", cwd)
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		if err := execIn(input); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

var ErrNoPath = errors.New("path required")

func execIn(input string) error {
	input = strings.TrimSpace(input)

	if input == "" {
		return nil
	}

	args := strings.Split(input, " ")

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
	}

	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			return fmt.Errorf("command not found: %s", args[0])
		}
		return err
	}
	return err
}
