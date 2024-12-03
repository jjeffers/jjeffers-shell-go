package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type shell_func func(cmd_map map[string]shell_func, args []string)

func exitCmd(command_map map[string]shell_func, args []string) {

	exit_status, err := strconv.Atoi(args[0])
	if err != nil {
		panic(fmt.Errorf("Can't convert '%s' to valid exit status.", args[0]))

	}
	os.Exit(exit_status)
}

func echoCmd(command_map map[string]shell_func, args []string) {
	echo_args := strings.Join(args, " ")
	fmt.Printf("%s\n", echo_args)
}

func IsExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

func typeCmd(command_map map[string]shell_func, args []string) {
	cmd := args[0]

	if _, exists := command_map[cmd]; exists {
		fmt.Printf("%s is a shell builtin\n", cmd)
	} else {
		path := os.Getenv("PATH")

		for _, dir := range strings.Split(path, ":") {
			files, _ := os.ReadDir(dir)
			for _, file := range files {

				info, _ := file.Info()

				if file.Name() == cmd && !info.IsDir() && IsExecAny(info.Mode()) {
					fmt.Printf("%s is %s\n", cmd, filepath.Join(dir, cmd))
					return
				}
			}
		}
		fmt.Printf("%s: not found\n", cmd)
	}
}

func handleInput(reader *bufio.Reader) {

	commands := map[string]shell_func{
		"exit": exitCmd,
		"echo": echoCmd,
		"type": typeCmd,
	}

	input, err := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	cmd_args := strings.Split(input, " ")

	if err != nil {
		panic(err)
	}

	cmd := strings.ToLower(cmd_args[0])

	if matched_cmd, exists := commands[cmd]; exists {

		matched_cmd(commands, cmd_args[1:])

	} else {
		fmt.Printf("%s: command not found\n", cmd)
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		handleInput(bufio.NewReader(os.Stdin))
	}

}
