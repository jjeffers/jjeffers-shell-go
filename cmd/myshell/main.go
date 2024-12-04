package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

		filepath, err := exec.LookPath(cmd)

		if err != nil {
			fmt.Printf("%s: not found\n", cmd)
		} else {
			fmt.Printf("%s is %s\n", cmd, filepath)
		}
	}
}

func execCommand(cmd string, cmd_args []string) {

	filepath, err := exec.LookPath(cmd)

	if err != nil {
		fmt.Printf("%s: command not found\n", cmd)
	} else {
		proc := exec.Command(filepath, cmd_args...)

		stdout, err := proc.Output()

		if err != nil {
			panic(err)
		}

		fmt.Print(string(stdout))
	}
}

func pwdCmd(_ map[string]shell_func, cmd_args []string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
}

func cdCmd(_ map[string]shell_func, cmd_args []string) {

	target_dir := cmd_args[0]

	if target_dir[0] == '~' {
		target_dir = os.Getenv("HOME")
	}

	err := os.Chdir(target_dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", cmd_args[0])
	}
}

func handleInput(reader *bufio.Reader) {

	commands := map[string]shell_func{
		"exit": exitCmd,
		"echo": echoCmd,
		"type": typeCmd,
		"pwd":  pwdCmd,
		"cd": cdCmd,
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
		execCommand(cmd, cmd_args[1:])
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		handleInput(bufio.NewReader(os.Stdin))
	}

}
