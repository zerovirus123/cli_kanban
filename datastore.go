package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// JSON-to-Go is amazeballs https://mholt.github.io/json-to-go/
type Columns struct {
	Todo []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"todo"`
	InProgress []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"inProgress"`
	Done []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	} `json:"done"`
}

func ReadFromStorage() Columns {

	var columns Columns

	if _, err := os.Stat("storage.json"); err == nil {
		jsonFile, err := os.Open("storage.json")

		if err != nil {
			fmt.Println("JSON IO ERROR: " + err.Error())
		}

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'columns' which we defined above
		json.Unmarshal(byteValue, &columns)
	}

	return columns
}
