package handler

import (
	"encoding/json"
	"fmt"
	"gobject-storage/meta"
	"gobject-storage/util"
	"io"
	"net/http"
	"os"
	"time"
)

// UploadHandler: Deal with file upload
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// return upload html page
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))

	} else if r.Method == "POST" {
		// receive file stream
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err:%s\n", err.Error())
		}
		defer file.Close()

		// create a local newFile in buffer
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file, err:%s\n", err.Error())
			return
		}
		defer newFile.Close()

		// copy the file stream to the newFile just created
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Print("Failed to save data into f")
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler: generate response when the upload is finished.
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler: get meta data of a file
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler: download file with file hash value
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fshar1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fshar1)
	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// tell the client that content being sent is a binary data stream (often used when sending files)
	w.Header().Set("Content-Type", "application/octect-stream")
	// tell the client that the file should be downloaded
	w.Header().Set("Content-Disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)
}