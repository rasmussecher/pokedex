package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rasmussecher/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(cfg *config, params []string) error
}

type config struct {
	pokeapiClient pokeapi.Client
	Next          string
	Previous      string
	Explore       string
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
		"explore": {
			name:        "explore <area_name>",
			description: "Explore an area",
			callback:    commandExplore,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, time.Minute*5)

	ctx := config{
		pokeapiClient: pokeClient,
		Next:          "https://pokeapi.co/api/v2/location-area",
		Previous:      "https://pokeapi.co/api/v2/location-area",
		Explore:       "https://pokeapi.co/api/v2/location-area/",
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := cleanInput(scanner.Text())
		if len(input) == 0 {
			fmt.Printf("You must input a command. Type \"help\" for a list of commands")
			continue
		}
		command, ok := commands[input[0]]
		if ok {
			command.callback(&ctx, input[1:])
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

func commandHelp(cfg *config, params []string) error {
	fmt.Print("Welcome to the Pokedex!\n")
	fmt.Print("Usage:\n\n")
	for _, c := range commands {
		fmt.Printf("%s, %s\n", c.name, c.description)
	}
	return nil
}

func commandMap(cfg *config, params []string) error {
	handleMap(cfg, cfg.Next)
	return nil
}

func commandMapb(cfg *config, params []string) error {
	handleMap(cfg, cfg.Previous)
	return nil
}

func commandExplore(cfg *config, params []string) error {
	if len(params) < 1 {
		fmt.Printf("You must enter an area!\n")
		return nil
	}
	area := params[0]
	encounters := cfg.pokeapiClient.GetPokemonsForArea(cfg.Explore + area)
	for _, e := range encounters.Encounters {
		fmt.Printf("%s\n", e.Pokemon.Name)
	}
	return nil
}

func commandExit(cfg *config, params []string) error {
	fmt.Print("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func handleMap(cfg *config, url string) {
	if url == "" {
		fmt.Printf("You must go further forward in the pagination.\n")
		return
	}
	location := cfg.pokeapiClient.GetList(url)
	cfg.Next = location.Next
	cfg.Previous = location.Previous
	for _, m := range location.ExtractNames() {
		fmt.Printf("%s\n", m)
	}
}
