package gocraft

import (
	"fmt"
	gocweb "github.com/gocraft/web"
	"net/http"
	"path/filepath"
)

func StaticMiddleware(path string) func(gocweb.ResponseWriter, *gocweb.Request, gocweb.NextMiddlewareFunc) {
	dir := http.Dir(path)
	return func(w gocweb.ResponseWriter, req *gocweb.Request, next gocweb.NextMiddlewareFunc) {
		if req.Method != "GET" && req.Method != "HEAD" {
			next(w, req)
			return
		}

		file := req.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			fmt.Println("Error opening Path: ", req.URL.Path, err)
			next(w, req)
			return
		}
		defer f.Close()

		fmt.Println("Here.....")
		fi, err := f.Stat()
		if err != nil {
			next(w, req)
			return
		}

		// Try to serve index.html
		if fi.IsDir() {
			file = filepath.Join(file, "index.html")
			f, err = dir.Open(file)
			if err != nil {
				next(w, req)
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				next(w, req)
				return
			}
		}

		http.ServeContent(w, req.Request, file, fi.ModTime(), f)
	}
}
