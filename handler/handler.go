package handler

import (
	"io"
	"net/http"
	"os"
)

// UploadHandler: Deal with file upload
func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		// return upload html page
		data, err := os.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	} else if r.Method == "POST"{
		// receive file stream and save to local cataglog
	}
}