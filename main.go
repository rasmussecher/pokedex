package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	pokeAPI "github.com/rasmussecher/pokedex/internal/api"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *config) error
}

type config struct {
	Next     string
	Previous string
}

var commands map[string]cliCommand

func init() {
	commands = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Show the next 20 areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous 20 areas",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {
	ctx := config{
		Next:     "https://pokeapi.co/api/v2/location-area",
		Previous: "https://pokeapi.co/api/v2/location-area",
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		command, ok := commands[input[0]]
		if ok {
			command.callback(&ctx)
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

func commandHelp(cfg *config) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s, %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config) error {
	handleMap(cfg, cfg.Next)
	return nil
}

func commandMapb(cfg *config) error {
	handleMap(cfg, cfg.Previous)
	return nil
}

func commandExit(cfg *config) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func handleMap(cfg *config, url string) {
	if url == "" {
		fmt.Printf("You must go further forward in the pagination.\n")
		return
	}
	location := pokeAPI.GetList(url)
	cfg.Next = location.Next
	cfg.Previous = location.Previous
	for _, m := range location.ExtractNames() {
		fmt.Printf("%s\n", m)
	}
}
