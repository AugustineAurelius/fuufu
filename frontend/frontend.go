package frontend

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed  dist/*
var frontendFS embed.FS

func RegisterFrontend(r *http.ServeMux) {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Panic(err.Error())
	}

	r.Handle("/", http.FileServer(http.FS(distFS)))
}
