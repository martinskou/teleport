package config

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	col "teleport/pkg/color"
	"time"
)

type Config struct {
	Name      string
	Server    string
	Port      int
	AuthToken string
	TmpFolder string
	TimeOut   time.Duration
}

func Load(paths ...string) (Config, error) {
	cfg := Config{Server: "0.0.0.0", Port: 31345, TmpFolder: "tmp", AuthToken: "1234", TimeOut: 3600}
	found := false
	for _, pospath := range paths {
		cfile := path.Join(pospath, "config.json")
		c, err := LoadJSON[Config](cfile)
		if err != nil {
			continue
		}
		col.CM.Printf("[purple]Config in use: %s[res]\n", cfile)
		cfg = c
		found = true
	}
	if !found {
		cfile := path.Join(paths[0], "config.json")
		SaveJSON(cfile, cfg)
		return cfg, errors.New("not found")
	}

	return cfg, nil
}

func LoadJSON[T any](path string) (T, error) {
	var data T

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}

	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return data, err
	}

	return data, nil
}

func SaveJSON[T any](path string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
