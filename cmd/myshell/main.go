package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func handleInput(reader *bufio.Reader) {

	input, err := reader.ReadString('\n')

	input = strings.TrimSpace(input)

	cmd_args := strings.Split(input, " ")

	if err != nil {
		panic(err)
	}

	cmd := strings.ToLower(cmd_args[0])

	switch cmd {
	case "exit":
		exit_status, err := strconv.Atoi(cmd_args[1])
		if err != nil {
			fmt.Errorf("Can't convert '%s' to valid exit status.", cmd_args[1])
		}
		os.Exit(exit_status)
	default:
		fmt.Printf("%s: command not found\n", cmd)
	}
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		handleInput(bufio.NewReader(os.Stdin))
	}

}
