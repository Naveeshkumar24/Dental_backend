package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Invalid file", 400)
		return
	}
	defer file.Close()

	filename := time.Now().Format("20060102150405") + "_" + handler.Filename
	path := filepath.Join("uploads/blogs", filename)

	out, _ := os.Create(path)
	defer out.Close()
	io.Copy(out, file)

	w.Write([]byte("/uploads/blogs/" + filename))
}
