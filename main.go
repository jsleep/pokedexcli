package main

import (
	"bufio"
	"fmt"
	"internal/pokeapi"
	"internal/pokecache"
	"os"
	"strings"
	"time"
)

func cleanInput(text string) []string {
	words := strings.Fields(text)
	res := make([]string, 0, len(words))

	for _, word := range words {
		res = append(res, strings.ToLower(word))
	}

	return res
}

func commandExit([]string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return fmt.Errorf("this should never happen")
}

func commandHelp([]string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range cliMap {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}

	return nil
}

var offset int = 0

func commandMap([]string) error {
	err := pokeapi.GetLocationAreas(offset, cache)
	offset += 20
	return err
}

func commandMapBack([]string) error {
	if offset == 0 {
		fmt.Println("You're on the first page")
		return nil
	} else {
		offset -= 20
		return pokeapi.GetLocationAreas(offset, cache)
	}
}

func commandExplore(a []string) error {
	if len(a) == 0 {
		fmt.Println("Please enter a location to explore")
		return nil
	}

	location := a[0]
	fmt.Printf("Exploring %s...\n", location)
	return pokeapi.GetLocationAreaDetails(location, cache)
}

func commandCatch(a []string) error {
	if len(a) == 0 {
		fmt.Println("Please enter a pokemon to catch")
		return nil
	}

	pokemon := a[0]
	return pokeapi.CatchPokemon(pokemon, cache, caught)
}

func commandInspect(a []string) error {
	if len(a) == 0 {
		fmt.Println("Please enter a pokemon to catch")
		return nil
	}

	pokemon := a[0]
	return pokeapi.InspectPokemon(pokemon, caught)
}

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

var cliMap map[string]cliCommand
var cache *pokecache.Cache
var caught map[string]pokeapi.Pokemon

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	cliMap = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Displays next 20 locations on map",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays previous 20 locations on map",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "explore the location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect a pokemon",
			callback:    commandInspect,
		},
	}
	cache = pokecache.NewCache(time.Minute * 5)
	caught = make(map[string]pokeapi.Pokemon)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()

		input := scanner.Text()
		args := cleanInput(input)

		if len(args) == 0 {
			fmt.Println("Empty command, Please enter a command")
			continue
		}

		command := args[0]
		if cmdObj, ok := cliMap[command]; !ok {
			fmt.Printf("Unknown command: %s\n", command)
			continue
		} else {
			err := cmdObj.callback(args[1:])
			if err != nil {
				fmt.Printf("Error executing command %s: %v\n", command, err)
			}
		}

	}
}
