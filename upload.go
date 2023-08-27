package main

import (
	"io"
	"net/http"
	"os"
	"path"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Get form and file values
	userId := r.FormValue("userid")
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Read from file
	data, err := io.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// Get filename (path name + filename) and write to path
	filename := path.Join("avatars", userId+path.Ext(header.Filename))
	err = os.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	io.WriteString(w, "Successful")
}
