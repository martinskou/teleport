package config

import (
	"crypto/tls"
	"errors"
	"math/rand"
	"net/http"
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

	UseTLS          bool   // If true only transmit with TLS/SSL cert
	TLSCert         string // Certificate file
	TLSKey          string // Key file
	AllowSelfSigned bool   // If true and TLSCert and TLSKey are empty, a selfsigned cert is used
}

func Load(paths ...string) (Config, error) {
	// Load config.json places in any of paths.
	// If not found, create a default and return an error.
	cfg := Config{
		Server:          "0.0.0.0",
		Port:            10000 + rand.Intn(10000),
		TmpFolder:       "tmp",
		AuthToken:       util.GenerateRandomAuthToken(32),
		TimeOut:         3600,
		UseTLS:          true,
		AllowSelfSigned: true,
	}
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

func (cfg Config) GetClient() *http.Client {
	tr := &http.Transport{}
	if cfg.AllowSelfSigned {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &http.Client{Transport: tr}
}
func (cfg Config) GetProtocol() string {
	if cfg.UseTLS {
		return "https"
	}
	return "http"
}
