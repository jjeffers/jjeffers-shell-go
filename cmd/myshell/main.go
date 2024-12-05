package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func tokenizeArgumentString(arguments string) ([]string, error) {

	regex := `'(.*?)'|"(.*?)"|(\w+)`
	reg := regexp.MustCompile(regex)
	matches := reg.FindAllString(arguments, -1)

	var result []string

	for _, match := range matches {
		trimmed_match := strings.Trim(match, `'"`)
		result = append(result, trimmed_match)
	}

	return result, nil

}

type shell_func func(cmd_map map[string]shell_func, rest_of_input string)

func exitCmd(command_map map[string]shell_func, rest_of_input string) {

	args, _ := tokenizeArgumentString(rest_of_input)

	exit_status, err := strconv.Atoi(args[0])
	if err != nil {
		panic(fmt.Errorf("Can't convert '%s' to valid exit status.", args[0]))

	}
	os.Exit(exit_status)
}

func echoCmd(command_map map[string]shell_func, rest_of_input string) {
	args, _ := tokenizeArgumentString(rest_of_input)
	fmt.Printf("%s\n", strings.Join(args, " "))
}

func IsExecAny(mode os.FileMode) bool {
	return mode&0111 != 0
}

func typeCmd(command_map map[string]shell_func, rest_of_input string) {
	args, _ := tokenizeArgumentString(rest_of_input)
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

func execCommand(cmd string, rest_of_input string) {

	cmd_args, _ := tokenizeArgumentString(rest_of_input)
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

func pwdCmd(_ map[string]shell_func, rest_of_input string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println(dir)
}

func cdCmd(_ map[string]shell_func, rest_of_input string) {

	cmd_args := strings.Fields(rest_of_input)
	target_dir := cmd_args[0]

	if target_dir[0] == '~' {
		target_dir = os.Getenv("HOME")
	}

	err := os.Chdir(target_dir)
	if err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", target_dir)
	}
}

func handleInput(reader *bufio.Reader) {

	commands := map[string]shell_func{
		"exit": exitCmd,
		"echo": echoCmd,
		"type": typeCmd,
		"pwd":  pwdCmd,
		"cd":   cdCmd,
	}

	input, err := reader.ReadString('\n')

	if err != nil {
		panic(err)
	}

	cmd, rest_of_input, _ := strings.Cut(input, " ")
	cmd = strings.ToLower(strings.TrimSpace(cmd))

	if matched_cmd, exists := commands[cmd]; exists {
		matched_cmd(commands, rest_of_input)
	} else {
		execCommand(cmd, rest_of_input)
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		handleInput(bufio.NewReader(os.Stdin))
	}

}
