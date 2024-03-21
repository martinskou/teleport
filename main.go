package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	c "teleport/pkg/color"
	"teleport/pkg/config"
	"teleport/pkg/server"
	"teleport/pkg/util"

	"github.com/urfave/cli/v2"
)

var version = "0.1"

func main() {

	c.CM.Printf("[green]Teleport [yellow](v.%s)[res]\n", version)

	wpath, err := util.FindExecPath()
	if err != nil {
		fmt.Println(err)
		return
	}

	cfg, err := config.Load(wpath, "~/.config/teleport")
	if err != nil {
		c.CM.Printf("[red]Error loading config, default config created (please check config.json)[res]\n")
		return
	}
	os.MkdirAll(cfg.TmpFolder, os.ModePerm)

	app := &cli.App{
		Name:  "teleport",
		Usage: "Teleport",
		Action: func(*cli.Context) error {
			fmt.Println("Teleport anything anywhere")
			fmt.Println("For help, type \"teleport help\".")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:        "server",
				Aliases:     []string{},
				Usage:       "teleport server",
				UsageText:   "starts a server",
				Description: "starts a server",
				Action: func(cCtx *cli.Context) error {
					return server.Server(cfg)
				},
			},
			{
				Name:        "upload",
				Aliases:     []string{"u"},
				Usage:       "teleport upload <folder>",
				UsageText:   "send a folder",
				Description: "send a folder",
				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() > 0 {
						sourcepath, err := filepath.Abs(cCtx.Args().Get(0))
						if err != nil {
							return err
						}
						return Send(cfg, sourcepath)
					} else {
						return errors.New("missing argument <folder>")
					}
				},
			},
			{
				Name:        "download",
				Aliases:     []string{"d"},
				Usage:       "teleport download <code> <destination>",
				UsageText:   "receive a folder",
				Description: "receive a folder",
				Action: func(cCtx *cli.Context) error {
					code := ""
					if cCtx.NArg() > 0 {
						code = cCtx.Args().Get(0)
					} else {
						return errors.New("missing argument <code>")
					}

					destpath := cfg.TmpFolder
					if cCtx.NArg() > 1 {
						destpath = cCtx.Args().Get(1)
					}
					destpath, err := filepath.Abs(destpath)
					if err != nil {
						return err
					}
					return Receive(cfg, code, destpath)
				},
			},
			{
				Name:        "test",
				Aliases:     []string{"t"},
				Usage:       "teleport test ",
				UsageText:   "test",
				Description: "test",
				Action: func(cCtx *cli.Context) error {
					fmt.Println(util.GeneratePassword())
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}

}

func Send(cfg config.Config, sourcepath string) error {
	sourceabs, err := filepath.Abs(sourcepath)
	if err != nil {
		return err
	}
	if !util.ExistsPath(sourceabs) {
		return fmt.Errorf(c.CM.Sprintf("[red]Source path [yellow]%s[red] not found[res]", sourceabs))
	}

	code := util.GeneratePassword()
	filename := path.Join(os.TempDir(), code+".zip")
	util.ZipFolder(sourcepath, filename)
	url := fmt.Sprintf("http://%s:%d/upload/", cfg.Server, cfg.Port)
	err = UploadFile(filename, url, cfg.AuthToken)
	if err != nil {
		return errors.New(c.CM.Sprintf("[red]Unable to upload file (%s)[res]", err.Error()))
	} else {
		c.CM.Printf("[green]Folder [yellow]%s[green] uploaded with retrieval code [yellow]%s[res]\n", sourcepath, code)
	}
	os.Remove(filename)
	return nil
}

func Receive(cfg config.Config, code string, destpath string) error {

	url := fmt.Sprintf("http://%s:%d/download/%s/", cfg.Server, cfg.Port, code)

	var b bytes.Buffer
	req, err := http.NewRequest("GET", url, &b)
	if err != nil {
		c.CM.Printf("[red]Auth error %s[res]\n", err.Error())
		return err
	}
	req.Header.Set("Auth-Token", cfg.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		c.CM.Printf("[red]Request error %s[res]\n", err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {

		tmp := path.Join(cfg.TmpFolder, "tmp.zip")
		out, err := os.Create(tmp)
		if err != nil {
			c.CM.Printf("[red]Create error %s[res]\n", err.Error())
			return err
		}
		defer out.Close()

		n, err := io.Copy(out, resp.Body)
		if err != nil {
			c.CM.Printf("[red]Copy error %s[res]\n", err.Error())
			return err
		}

		err = util.Unzip(tmp, destpath)
		if err != nil {
			c.CM.Printf("[red]Unzip error %s[res]\n", err.Error())
		} else {
			c.CM.Printf("[green]Downloaded [yellow]%s[green] (%d bytes) to [yellow]%s[green][res]\n", code, n, destpath)

		}
	} else {
		c.CM.Printf("[red]Error status %s[res]\n", resp.Status)
	}

	return nil
}

func UploadFile(filePath, targetURL, auth_token string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fw, err := w.CreateFormFile("file", filePath)
	if err != nil {
		return err
	}
	if _, err = io.Copy(fw, file); err != nil {
		return err
	}
	w.Close()

	req, err := http.NewRequest("POST", targetURL, &b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Auth-Token", auth_token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	return nil
}
