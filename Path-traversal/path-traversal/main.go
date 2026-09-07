package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	http.Handle("/", http.FileServer(http.Dir("public")))

	http.HandleFunc("/api/content", func(w http.ResponseWriter, r *http.Request) {
		fileParam := r.URL.Query().Get("file")
		if fileParam == "" {
			http.Error(w, "Missing 'file' parameter", http.StatusBadRequest)
			return
		}

		basedir := "data"
		r_base, err := filepath.Abs(basedir)
		resolv, err := filepath.Abs(filepath.Join(basedir, fileParam))

		real_base, err := filepath.EvalSymlinks(r_base)
		real_target, err := filepath.EvalSymlinks(resolv)

		fmt.Println(real_base)

		//condition
		if real_base != real_target || !strings.HasPrefix(resolv, r_base+"/") {
			http.Error(w, "Not allowed", http.StatusForbidden)
			return
		}

		resolvedPath := filepath.Join(r_base, fileParam)
		file, err := os.Open(resolvedPath)
		if err != nil {
			fmt.Printf("Error opening file: %v\n", err)
			http.Error(w, "Error opening a file", http.StatusNotFound)
			return
		}
		defer file.Close()

		io.Copy(w, file)
	})

	fmt.Println("Starting server on http://127.0.0.1:3000")
	fmt.Println("Warning: This server is intentionally vulnerable to LFI.")

	log.Fatal(http.ListenAndServe("127.0.0.1:3000", nil))
}
