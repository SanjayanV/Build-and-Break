package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"slices"
)

const uploadPath = "./uploads"

func main() {
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		log.Fatal("Failed to create upload directory:", err)
	}

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/upload", uploadHandler)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir(uploadPath))))

	fmt.Println("Vulnerable File Upload Server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Vulnerable File Upload</title>
		<style>
			body { font-family: 'Inter', sans-serif; margin: 0; padding: 0; background-color: #121212; color: #e0e0e0; display: flex; justify-content: center; align-items: center; height: 100vh; }
			.container { background-color: #1e1e1e; padding: 40px; border-radius: 12px; box-shadow: 0 10px 30px rgba(0,0,0,0.5); border: 1px solid #333; max-width: 500px; width: 100%; }
			h1 { color: #ff5252; margin-top: 0; }
			p { color: #aaa; line-height: 1.5; }
			.warning { background-color: rgba(255, 82, 82, 0.1); border-left: 4px solid #ff5252; padding: 10px; margin-bottom: 20px; font-size: 0.9em; }
			.upload-form { display: flex; flex-direction: column; gap: 15px; }
			input[type="file"] { padding: 10px; background-color: #2c2c2c; border: 1px dashed #555; border-radius: 6px; color: #fff; cursor: pointer; }
			.upload-btn { background-color: #4CAF50; color: white; border: none; padding: 12px; cursor: pointer; border-radius: 6px; font-weight: bold; font-size: 1em; transition: background-color 0.3s ease; }
			.upload-btn:hover { background-color: #45a049; }
		</style>
	</head>
	<body>
		<div class="container">
			<h1>Upload Simulator</h1>
			
			<form class="upload-form" action="/upload" method="post" enctype="multipart/form-data">
				<input type="file" name="uploadFile" required>
				<button type="submit" class="upload-btn">Upload File</button>
			</form>
		</div>
	</body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, handler, err := r.FormFile("uploadFile")

	allowlist := []string{".jpg", ".png", ".pdf", "docx"}
	ext := filepath.Ext(handler.Filename)
	var targetPath string
	if slices.Contains(allowlist, ext) {
		targetPath = uploadPath + "/" + handler.Filename //fetching the actual filename from the user input and appending it directly with path
	} else {
		http.Error(w, "Unauthorized file type", http.StatusBadRequest)
	}

	fmt.Println(ext)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving the file : %s", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	dst, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	successMsg := fmt.Sprintf(`
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<style>
			body { font-family: 'Inter', sans-serif; background-color: #121212; color: #e0e0e0; padding: 40px; }
			.success { color: #4CAF50; }
			a { color: #64b5f6; text-decoration: none; }
			a:hover { text-decoration: underline; }
		</style>
	</head>
	<body>
		<h2 class="success">Upload Successful!</h2>
		<p>File saved to: <code>%s</code></p>
		<p><a href="/uploads/%s">View your file</a></p>
		<p><a href="/">Upload another</a></p>
	</body>
	</html>
	`, targetPath, handler.Filename)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(successMsg))
}
