package internalhttp

import (
	"net/http"
	"time"
)

func loggingMiddleware(Server /*http.Handler*/ *Server, r *http.Request, myTimeStart time.Time) /*http.Handler*/ {
	myRec := *r
	Server.App.Info(myRec.RemoteAddr, myTimeStart, myRec.Method /*myRec.Response.StatusCode,*/, myRec.Proto, time.Since(myTimeStart), myRec.UserAgent())
	//	return Server //http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

}
