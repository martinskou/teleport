package config

import (
	"errors"
	"path"
	col "teleport/pkg/color"
	"teleport/pkg/util"
	"time"
)

type Config struct {
	Server    string        // Which server will teleport upload/download from
	Port      int           // Server port
	AuthToken string        // A token which must match on client and server
	TmpFolder string        // Where to place files
	TimeOut   time.Duration // Delete files after this duration
}

func Load(paths ...string) (Config, error) {
	// Load config.json places in any of paths.
	// If not found, create a default and return an error.
	cfg := Config{Server: "0.0.0.0", Port: 31345, TmpFolder: "tmp", AuthToken: "1234", TimeOut: 3600}
	found := false
	for _, pospath := range paths {
		cfile := path.Join(pospath, "config.json")
		c, err := util.LoadJSON[Config](cfile)
		if err != nil {
			continue
		}
		col.CM.Printf("[purple]Config in use: %s[res]\n", cfile)
		cfg = c
		found = true
		break
	}
	if !found {
		cfile := path.Join(paths[0], "config.json")
		util.SaveJSON(cfile, cfg)
		return cfg, errors.New("not found")
	}

	return cfg, nil
}
