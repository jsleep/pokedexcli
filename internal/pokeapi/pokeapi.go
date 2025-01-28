package pokeapi

import (
	"encoding/json"
	"fmt"
	"internal/pokecache"
	"io"
	"math/rand"
	"net/http"
)

type MapResult struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous any            `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func GetLocationAreas(offset int, cache *pokecache.Cache) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/?offset=%d&limit=20", offset)

	var data []byte

	if cacheData, ok := cache.Get(url); ok {
		fmt.Println("Cache hit")
		data = cacheData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch location areas: %v", err)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed read bytes %v", err)
		}
		cache.Add(url, data)

		defer resp.Body.Close()
	}

	var mapResult MapResult
	err := json.Unmarshal(data, &mapResult)
	if err != nil {
		return fmt.Errorf("failed to decode location areas: %v", err)
	}

	locationAreas := mapResult.Results

	for _, area := range locationAreas {
		fmt.Println(area.Name)
	}

	return nil
}

type LocationAreaDetails struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

func GetLocationAreaDetails(name string, cache *pokecache.Cache) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/location-area/%s/", name)

	var data []byte

	if cacheData, ok := cache.Get(url); ok {
		fmt.Println("Cache hit")
		data = cacheData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch location areas: %v", err)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed read bytes %v", err)
		}
		cache.Add(url, data)

		defer resp.Body.Close()
	}

	var areaDetails LocationAreaDetails
	err := json.Unmarshal(data, &areaDetails)
	if err != nil {
		return fmt.Errorf("failed to decode location areas: %v", err)
	}

	encounters := areaDetails.PokemonEncounters
	fmt.Println("Found Pokemon:")
	for _, encounter := range encounters {
		fmt.Printf("- %s\n", encounter.Pokemon.Name)
	}

	return nil
}

type Pokemon struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	IsDefault      bool   `json:"is_default"`
	Order          int    `json:"order"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func CatchPokemon(pokemonName string, cache *pokecache.Cache, caught map[string]Pokemon) error {
	url := fmt.Sprintf("https://pokeapi.co/api/v2/pokemon/%s/", pokemonName)

	var data []byte

	throwVal := rand.Intn(800)

	if cacheData, ok := cache.Get(url); ok {
		fmt.Println("Cache hit")
		data = cacheData
	} else {
		resp, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch location areas: %v", err)
		}
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed read bytes %v", err)
		}
		cache.Add(url, data)

		defer resp.Body.Close()
	}

	var pokemon Pokemon
	err := json.Unmarshal(data, &pokemon)
	if err != nil {
		return fmt.Errorf("failed to decode pokemon: %v", err)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if throwVal >= pokemon.BaseExperience {
		fmt.Printf("%s was caught!\n", pokemonName)
		caught[pokemonName] = pokemon
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}

	return nil
}

func InspectPokemon(pokemonName string, caught map[string]Pokemon) error {

	if pokemon, ok := caught[pokemonName]; ok {
		fmt.Println("Name:", pokemon.Name)
		fmt.Println("Height:", pokemon.Height)
		fmt.Println("Weight:", pokemon.Weight)
		fmt.Println("Stats:")
		for _, stat := range pokemon.Stats {
			fmt.Printf("- %s: %d\n", stat.Stat.Name, stat.BaseStat)
		}
		fmt.Println("Types:")
		for _, typeData := range pokemon.Types {
			fmt.Printf("- %s\n", typeData.Type.Name)
		}
	} else {
		fmt.Println("You have not caught that pokemon")
	}

	return nil
}

func Pokedex(caught map[string]Pokemon) error {
	fmt.Println("Your Pokedex:")
	for name, _ := range caught {
		fmt.Printf("- %s\n", name)
	}

	return nil
}
