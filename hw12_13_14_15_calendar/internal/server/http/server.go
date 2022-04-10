package internalhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	App Application
	//	Logg     Logger
	HttpConf string
	MyHTTP   http.Server
	//	MyHandler MyHandler
}

//func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintf(w, "Hello! ", time.Now(), r)
//}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	myTimeStart := time.Now()
	fmt.Fprintf(w, "Hello world!")
	loggingMiddleware(s, r, myTimeStart)

	//	s.App.Print(fmt.Sprintf("%#v", w))
}

/*
func (s *Server) HelloW(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
*/
//type Logger interface { // TODO
//}

type Application interface {
	CreateEvent(ctx context.Context, title, startDateStr, details string, userID int) error
	UpdateEvent(ctx context.Context, eventID, title, startDateStr, details string, userID int) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEventByDate(ctx context.Context, startDateStr string) ([]st.Event, error)
	GetEventMonth(ctx context.Context, startDateStr string) ([]st.Event, error)
	GetEventWeek(ctx context.Context, startDateStr string) ([]st.Event, error)

	Trace(args ...interface{})
	Debug(args ...interface{})
	Info(args ...interface{})
	Print(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

func NewServer( /*logger Logger,*/ app Application, httpConf string) *Server {
	return &Server{App: app /*Logg: logger,*/, HttpConf: httpConf}
}

func (s *Server) Start(ctx context.Context) error {
	fmt.Printf("serv=%#v", s.HttpConf)
	//	mux := http.NewServeMux()
	//	mux.HandleFunc("/hello", loggingMiddleware(s) s.HelloW))
	s.MyHTTP.Addr = s.HttpConf
	s.MyHTTP.Handler = s //loggingMiddleware(s, s) //s//loggingMiddleware(s)
	//	fmt.Printf("MyHTTP=%#v", s.MyHTTP)

	s.MyHTTP.ListenAndServe()

	/*	http.HandleFunc("/", sayhello)                     // Устанавливаем роутер
		myEr := http.ListenAndServe(s.HttpConf, nil) // устанавливаем порт веб-сервера
		if myEr != nil {
			log.Fatal("ListenAndServe: ", myEr)
		}
	*/
	<-ctx.Done()
	return nil
}

/*
func sayhello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Привет! ", time.Now(), r)
}
*/
func (s *Server) Stop(ctx context.Context) error {
	s.MyHTTP.Shutdown(ctx)
	return nil
}
