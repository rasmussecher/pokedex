package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
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
	caughtPokemon map[string]pokeapi.Pokemon
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
		"catch": {
			name:        "catch <pokemon_name>",
			description: "Attempt to catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect <pokemon_name>",
			description: "Inspect a Pokemon in your inventory",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Print all the Pokemons in your Pokedex",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func main() {
	pokeClient := pokeapi.NewClient(5*time.Second, 5*time.Minute)

	ctx := config{
		pokeapiClient: pokeClient,
		caughtPokemon: map[string]pokeapi.Pokemon{},
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

func commandCatch(cfg *config, params []string) error {
	if len(params) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	name := params[0]
	pokemon, err := cfg.pokeapiClient.GetPokemon(name)
	if err != nil {
		return err
	}

	res := rand.Intn(pokemon.BaseExperience)

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemon.Name)
	if res > 40 {
		fmt.Printf("%s escaped!\n", pokemon.Name)
		return nil
	}

	fmt.Printf("%s was caught!\n", pokemon.Name)

	cfg.caughtPokemon[pokemon.Name] = pokemon
	return nil
}

func commandInspect(cfg *config, params []string) error {
	if len(params) != 1 {
		return errors.New("you must provide a pokemon name")
	}

	name := params[0]
	pokemon, ok := cfg.caughtPokemon[name]
	if !ok {
		fmt.Printf("you have not cought that pokemon\n")
	}
	printPokemon(pokemon)
	return nil
}

func commandPokedex(cfg *config, params []string) error {
	fmt.Printf("Your Pokedex:\n")
	for _, p := range cfg.caughtPokemon {
		fmt.Printf(" - %s\n", p.Name)
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

func printPokemon(p pokeapi.Pokemon) {
	fmt.Printf("Name: %s\nHeight: %d\nWeight: %d\nStats:\n", p.Name, p.Height, p.Weight)
	for _, s := range p.Stats {
		fmt.Printf("  -%s: %d\n", s.Stat.Name, s.BaseStat)
	}
	fmt.Printf("Types:\n")
	for _, t := range p.Types {
		fmt.Printf("  - %s\n", t.Type.Name)
	}
}
