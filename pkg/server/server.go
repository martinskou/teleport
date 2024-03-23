package server

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	c "teleport/pkg/color"
	"teleport/pkg/config"
	"teleport/pkg/util"
	"time"
)

type Handler struct {
	cfg config.Config
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Auth-Token") != h.cfg.AuthToken {
		http.Error(w, "Authentication error", http.StatusMethodNotAllowed)
		return
	}
	filelst, err := os.ReadDir(h.cfg.TmpFolder)
	if err != nil {
		c.CM.Printf("[red]Unable to read file list (%s)[res]\n", err.Error())
		http.Error(w, "error", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, f := range filelst {
		if !f.IsDir() {
			afile := path.Join(h.cfg.TmpFolder, f.Name())
			s, err := os.Stat(afile)
			if err == nil {
				age := time.Since(s.ModTime())
				w.Write([]byte(fmt.Sprintf("%s (%s)\n", util.FileNameWithoutExt(f.Name()), age.String())))
			}
		}
	}
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Auth-Token") != h.cfg.AuthToken {
		http.Error(w, "Authentication error", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	filename := r.PathValue("file")
	c.CM.Printf("[green]Request for [yellow]%s[res]\n", filename)

	absfile, err := filepath.Abs(path.Join(h.cfg.TmpFolder, filename+".zip"))
	if err != nil {
		c.CM.Printf("[red]Unable for understand request for [yellow]%s[red] (%s)[res]\n", filename, err.Error())
		http.Error(w, "error", http.StatusNotFound)
		return
	}
	if !util.ExistsPath(absfile) {
		c.CM.Printf("[red]Not found [yellow]%s[red][res]\n", filename)
		http.Error(w, "not found", http.StatusNotFound)
		return
	} else {
		c.CM.Printf("[green]Serving file [yellow]%s[res]\n", filename)
		http.ServeFile(w, r, absfile)
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method != "POST" {
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Auth-Token") != h.cfg.AuthToken {
		http.Error(w, "Authentication error", http.StatusMethodNotAllowed)
		return
	}

	// Parse the multipart form, 100 << 20 specifies a maximum upload of 100 MB files.
	err := r.ParseMultipartForm(100 << 20)
	if err != nil {
		c.CM.Printf("[red]Error Retrieving the File (%s)[res]\n", err.Error())
		http.Error(w, "Parse error.", http.StatusBadRequest)
		return
	}

	// Retrieve the file from posted form-data
	file, handler, err := r.FormFile("file")
	if err != nil {
		c.CM.Printf("[red]Error retrieving the file (%s)[res]\n", err.Error())
		http.Error(w, "Form file error.", http.StatusBadRequest)
		return
	}
	defer file.Close()

	destabs, err := filepath.Abs(path.Join(h.cfg.TmpFolder, handler.Filename))
	if err != nil {
		c.CM.Printf("[red]Error saving the file (%s)[res]\n", err.Error())
		http.Error(w, "Save file error.", http.StatusBadRequest)
		return
	}

	dst, err := os.Create(destabs)
	if err != nil {
		c.CM.Printf("[red]Error saving the file (%s)[res]\n", err.Error())
		http.Error(w, "Save file error.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the uploaded file to the created file on the filesystem
	n, err := io.Copy(dst, file)
	if err != nil {
		c.CM.Printf("[red]Error saving the file (%s)[res]\n", err.Error())
		http.Error(w, "Save file error.", http.StatusInternalServerError)
		return
	}

	c.CM.Printf("[green]Received file [yellow]%s[green] (%d bytes) succesfully saved[res]\n", handler.Filename, n)

	w.Write([]byte("success"))
}

func cleanup(cfg config.Config) {
	for {
		files, err := os.ReadDir(cfg.TmpFolder)
		if err != nil {
			c.CM.Printf("[red]Cleanup error (%s)[res]\n", err.Error())

		} else {
			for _, f := range files {
				if !f.IsDir() {
					afile := path.Join(cfg.TmpFolder, f.Name())
					s, err := os.Stat(afile)
					if err == nil {
						if time.Since(s.ModTime()) > cfg.TimeOut*time.Second {
							c.CM.Printf("[pink]Cleanup deleting %s[res]\n", f.Name())
							os.Remove(afile)
						}
					}

				}
			}
		}

		time.Sleep(10 * time.Second)
	}
}

func Server(cfg config.Config) error {
	c.CM.Printf("[green]Starting server on [yellow]0.0.0.0:%d[res]\n", cfg.Port)

	// Delete all files in tmp

	h := Handler{cfg: cfg}

	mux := http.NewServeMux()

	mux.HandleFunc("/list/", h.List)
	mux.HandleFunc("/upload/", h.Upload)
	mux.HandleFunc("/download/{file}/", h.Download)

	srv := &http.Server{
		Handler: mux,
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", cfg.Port))
	if err != nil {
		panic(err)
	}

	go cleanup(cfg)

	return srv.Serve(ln)

}
