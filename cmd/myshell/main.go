package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
			fmt.Fprint(os.Stdout, "$ ")

			// Wait for user input
			cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')

			cmd = strings.TrimSpace(cmd)

			if err != nil {
				panic(err)
			}

			fmt.Printf("%s: command not found\n", cmd)
		}

}
