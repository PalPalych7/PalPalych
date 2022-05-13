package internalhttp

import (
	"net/http"
	"time"
)

func loggingMiddleware(server /*http.Handler*/ *Server, r *http.Request, myTimeStart time.Time) /*http.Handler*/ {
	myRec := *r
	server.App.Info(myRec.RemoteAddr, myTimeStart, myRec.Method, myRec.Proto, time.Since(myTimeStart), myRec.UserAgent())
	//	return Server //http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
}
