package main

import (
	"encoding/json"
	"os"
)

// data we want for one layer in our list of layers
type TilemapLayerJson struct {
	Data   []int `json:"data"`
	Width  int   `json:"width"`
	Height int   `json:"height"`
}

// all layers in a tilemap
type TilemapJSON struct {
	Layers []TilemapLayerJson `json:"layers"`
}

func NewTilemapJSON(filepath string) (*TilemapJSON, error) {
	// read the file
	contents, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var tilemapJSON TilemapJSON
	err = json.Unmarshal(contents, &tilemapJSON)
	if err != nil {
		return nil, err
	}

	return &tilemapJSON, nil
}
