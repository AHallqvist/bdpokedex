package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func startRepl() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			break
		}
		words := cleanInput(scanner.Text())
		if len(words) == 0 {
			continue
		}
		fmt.Printf("Your command was: %s\n", words[0])
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read input:", err)
	}
}

func cleanInput(text string) []string {
	return strings.Fields(strings.ToLower(text))
}
