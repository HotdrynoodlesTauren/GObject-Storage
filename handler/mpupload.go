package handler

import (
	"fmt"
	rPool "gobject-storage/cache/redis"
	dblayer "gobject-storage/db"
	"gobject-storage/util"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gomodule/redigo/redis"
)

// MultipartUploadInfo: initialization information
type MultipartUploadInfo struct {
	FileHash string
	FileSize int
	UploadID string
	ChunkSize int
	ChunkCount int
}

// InitialMultipartUploadHandler: Initialization
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse parameters from user request
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}

	// get a redis connection
	rConn := rPool.RedisPool().Get() // .Get() returns the reredis client
	defer rConn.Close()

	//  generate initialization information for multi part upload
	upInfo := MultipartUploadInfo {
		FileHash: filehash,
		FileSize: filesize,
		UploadID: username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize: 5 * 1024 * 1024, // 5MB
		ChunkCount: int(math.Ceil(float64(filesize) / (5 * 1024 * 1024))),
	}

	// write the initialization info to redis cache
	rConn.Do("HSET", "MP_" + upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_" + upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_" + upInfo.UploadID, "filesize", upInfo.FileSize)

	// return the response of initialization to client
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

// UploadPartHandler: Upload the blocked file
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	// parse parameters from user request
	r.ParseForm()
	// username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")

	// get a redis connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// obtain the handle to the file, used to store the block content
	fpath := "./test/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "Upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024 * 1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	// update the redis cache
	rConn.Do("HSET", "MP_" + uploadID, "chkidx_" + chunkIndex, 1)

	// return the response to client
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

// CompleteUploadHandler: notify upload and merge
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// parse parameters from user request
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	// get a redis connection
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	// check and determine whether all chunks are uploaded via uploadid
	data, err := redis.Values(rConn.Do("HGETALL", "MP_" + uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount += 1
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	// TODO: merge the blocks

	// update the tbl_file and tbl_user_file
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "" )
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))

	// return the response to client
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())

}