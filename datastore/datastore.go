package datastore

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"cli_kanban/task"
	"cli_kanban/typedef"
)

type Todo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type InProgress struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Done struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// JSON-to-Go is amazeballs https://mholt.github.io/json-to-go/
type Columns struct {
	Todo       []Todo       `json:"todo"`
	InProgress []InProgress `json:"inProgress"`
	Done       []Done       `json:"done"`
}

func ReadFromStorage() Columns {
	var columns Columns

	storageFile := "datastore/storage.json"

	if _, err := os.Stat(storageFile); err == nil {
		jsonFile, err := os.Open(storageFile)
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

func WriteToStorage(m typedef.Model) {
	var columns Columns

	for _, element := range m.Lists[typedef.Todo].Items() {
		columns.Todo = append(columns.Todo, Todo{element.(*task.Task).Title(), element.(*task.Task).Description()})
	}

	for _, element := range m.Lists[typedef.InProgress].Items() {
		columns.InProgress = append(columns.InProgress, InProgress{element.(*task.Task).Title(), element.(*task.Task).Description()})
	}

	for _, element := range m.Lists[typedef.Done].Items() {
		columns.Done = append(columns.Done, Done{element.(*task.Task).Title(), element.(*task.Task).Description()})
	}

	storageFile := "datastore/storage.json"

	if _, err := os.Stat(storageFile); err == nil {
		e := os.Remove(storageFile)
		if e != nil {
			fmt.Println(e)
		}
	}

	_, e := os.Create(storageFile)
	if e != nil {
		fmt.Println(e)
	}

	jsonString, _ := json.MarshalIndent(columns, "", " ")
	_ = ioutil.WriteFile(storageFile, jsonString, 0644)
}
