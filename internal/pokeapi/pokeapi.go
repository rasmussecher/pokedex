package pokeapi

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type ListResponse struct {
	Count    int32  `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// TODO - These request functions should be generic instead.
func (c *Client) GetList(url string) ListResponse {
	if val, ok := c.cache.Get(url); ok {
		locationsResp := ListResponse{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			log.Fatalf("err")
			return ListResponse{}
		}

		return locationsResp
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	listRes := ListResponse{}
	err = json.Unmarshal(body, &listRes)
	if err != nil {
		log.Fatal(err)
	}
	return listRes
}

func (listRes *ListResponse) ExtractNames() []string {
	names := []string{}
	for i := range listRes.Results {
		names = append(names, listRes.Results[i].Name)
	}
	return names
}

type PokemonEncounterList struct {
	Encounters []struct {
		Pokemon struct {
			Name string `json:"name"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
}

func (c *Client) GetPokemonsForArea(url string) PokemonEncounterList {
	if val, ok := c.cache.Get(url); ok {
		locationsResp := PokemonEncounterList{}
		err := json.Unmarshal(val, &locationsResp)
		if err != nil {
			log.Fatalf("err")
			return PokemonEncounterList{}
		}

		return locationsResp
	}

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	}
	listRes := PokemonEncounterList{}
	err = json.Unmarshal(body, &listRes)
	if err != nil {
		log.Fatal(err)
	}
	return listRes
}

func (c *Client) GetPokemon(pokemonName string) (Pokemon, error) {
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemonName

	if val, ok := c.cache.Get(url); ok {
		pokemonResp := Pokemon{}
		err := json.Unmarshal(val, &pokemonResp)
		if err != nil {
			return Pokemon{}, err
		}
		return pokemonResp, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return Pokemon{}, err
	}

	pokemonResp := Pokemon{}
	err = json.Unmarshal(dat, &pokemonResp)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, dat)

	return pokemonResp, nil
}
