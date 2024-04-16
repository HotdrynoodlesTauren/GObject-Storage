package handler

import (
	"net/http"
)

// HTTPInterceptor: HTTP request Interceptor
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc {
	// 视频里写的会报错
	// return http.HandlerFunc(
		return func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			
			// 
			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h(w, r)
		}
	// )
}