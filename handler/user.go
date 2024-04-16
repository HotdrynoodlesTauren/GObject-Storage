package handler

import (

	// "io/ioutil"

	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	dblayer "gobject-storage/db"
	"gobject-storage/util"
)

const (
	// user for salt in encryption
	pwdSalt = "*#890"
)

// SignupHandler : Handle the signup request from an user
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))

		//http.Redirect(w, r, "/static/view/signup.html", http.StatusFound)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	// encrypt password via sha1
	encPasswd := util.Sha1([]byte(passwd + pwdSalt))
	// save user identity into the db
	suc := dblayer.UserSignup(username, encPasswd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SignInHandler: Login interface
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Serve the signin.html page
		data, err := os.ReadFile("./static/view/signin.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	}

	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwdSalt))

	// 1. Validate username and password
	pwdChecked := dblayer.UserSignin(username, encPasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	// 2. Generate access token
	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	// 3. Redirect to the homepage after successful login
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	//http.Redirect(w, r, "/static/view/home.html", http.StatusTemporaryRedirect)
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())

}

func GenToken(username string) string {
	// md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]

}

// UserInfoHandler: Query user information
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Parse request parameters
	r.ParseForm()
	username := r.Form.Get("username")
	//	token := r.Form.Get("token")

	// // 2. Validate if the token is valid
	// isValidToken := IsTokenValid(token)
	// if !isValidToken {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }
	
	// 3. Query user information
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	// 4. Assemble and respond with user data
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return true
}