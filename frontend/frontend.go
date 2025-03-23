package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed  dist/*
var frontendFS embed.FS

func RegisterFrontend(r *http.ServeMux) {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		panic(err.Error())
	}

	r.Handle("/", http.FileServer(http.FS(distFS)))
}
