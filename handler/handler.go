package handler

import (
	"encoding/json"
	"fmt"
	dblayer "gobject-storage/db"
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
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		// update tbl_user_file
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println(username + " is uploading")
		suc := dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
		if suc {
			http.Redirect(w, r, "/static/view/home.html", http.StatusFound)
		} else {
			w.Write([]byte("Upload Failed"))
		}
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
	//fMeta := meta.GetFileMeta(filehash)
	fMeta, err1 := meta.GetFileMetaDB(filehash)
	data, err2 := json.Marshal(fMeta)
	if err2 != nil && err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 第一章的视频好像漏掉了这个部分, 暂时先贴上, 不知道后面有没有用
// FileQueryHandler : 查询批量的文件元信息
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	// r.ParseForm()

	// limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	// username := r.Form.Get("username")
	// //fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	// userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	// data, err := json.Marshal(userFiles)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// w.Write(data)
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

// FileMetaUpdateHandler: modify file meta info
func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta((fileSha1))
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// FileDeleteHandler: modify file meta info
func FileDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
