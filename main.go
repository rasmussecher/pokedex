package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command, ok := commands[input[0]]
		if ok {
			command.callback()
		} else {
			fmt.Print("Unknown command\n")
		}
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	for i := range words {
		words[i] = strings.ToLower(words[i])
	}
	return words
}

func commandHelp() error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s, %s\n", c.name, c.description)
	}
	return nil
}

func commandExit() error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}
