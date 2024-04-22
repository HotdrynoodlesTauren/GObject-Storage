package handler

import (
	"encoding/json"
	"fmt"
	cmn "gobject-storage/common"
	cfg "gobject-storage/config"
	dblayer "gobject-storage/db"
	"gobject-storage/meta"
	"gobject-storage/mq"
	"gobject-storage/store/oss"
	"gobject-storage/util"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
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
		//data, _ := io.ReadAll(newFile)

		fileMeta.FileSha1 = util.FileSha1(newFile)

		ossPath := "oss/" + fileMeta.FileSha1
		/*err = oss.Bucket().PutObjectFromFile(ossPath, fileMeta.Location)
		if err != nil {
			fmt.Println(err.Error())
			w.Write([]byte("Upload failed!"))
			return
		}
		fileMeta.Location = ossPath*/
		data := mq.TransferData{
			FileHash:      fileMeta.FileSha1,
			CurLocation:   fileMeta.Location,
			DestLocation:  ossPath,
			DestStoreType: cmn.StoreOSS,
		}
		pubData, _ := json.Marshal(data)
		suc := mq.Publish(cfg.TransExchangeName,
			cfg.TransOSSRoutingKey,
			pubData)
		if !suc {

		}

		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)

		// update tbl_user_file
		r.ParseForm()
		username := r.Form.Get("username")
		fmt.Println(username + " is uploading")
		suc = dblayer.OnUserFileUploadFinished(username, fileMeta.FileSha1, fileMeta.FileName, fileMeta.FileSize)
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

// FileQueryHandler : batch query user file meta
func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	username := r.Form.Get("username")
	//fileMetas, _ := meta.GetLastFileMetasDB(limitCnt)
	userFiles, err := dblayer.QueryUserFileMetas(username, limitCnt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(userFiles)
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
	fm, _ := meta.GetFileMetaDB(fshar1)
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

// TryFastUploadHandler:
func TryFastUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filename := r.Form.Get("filename")
	filesize, _ := strconv.Atoi(r.Form.Get("filesize"))

	// Get the file record with the same hash from tbl_file
	fileMeta, err := meta.GetFileMetaDB(filehash)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// If not record then
	if fileMeta == nil {
		resp := util.RespMsg{
			Code: -1,
			Msg:  "Fast upload failed, visit ordinary endpoint",
		}
		w.Write(resp.JSONBytes())
		return
	}

	// If there is a record, then write file meta to tbl_user_file, return true
	suc := dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(filesize))
	if suc {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "Fast upload successed",
		}
		w.Write(resp.JSONBytes())
		return
	} else {
		resp := util.RespMsg{
			Code: -2,
			Msg:  "Fast upload failed, try later",
		}
		w.Write(resp.JSONBytes())
		return
	}
}

// DownloadURLHandler: Generates the download URL for a file
func DownloadURLHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	// Retrieve record from the file table
	row, _ := dblayer.GetFileMeta(filehash)
	// TODO: Determine if the file exists in OSS, Ceph, or locally
	if strings.HasPrefix(row.FileAddr.String, "/tmp") {
		username := r.Form.Get("username")
		token := r.Form.Get("token")
		tmpUrl := fmt.Sprintf("http://%s/file/download?filehash=%s&username=%s&token=%s",
			r.Host, filehash, username, token)

		w.Write([]byte(tmpUrl))
	} else if strings.HasPrefix(row.FileAddr.String, "/ceph") {
		// TODO: Ceph download URL
	} else if strings.HasPrefix(row.FileAddr.String, "oss/") {
		// OSS download URL
		signedURL := oss.DownloadURL(row.FileAddr.String)
		w.Write([]byte(signedURL))
	}
}
