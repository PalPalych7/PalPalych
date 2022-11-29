package internalhttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	App      Application
	HTTPConf string
	MyHTTP   http.Server
	myCtx    context.Context
}

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

func NewServer(app Application, httpConf string) *Server {
	return &Server{App: app, HTTPConf: httpConf}
}

func getBodyRow(reqBody io.ReadCloser) []byte {
	raw, err := ioutil.ReadAll(reqBody)
	if err != nil {
		return nil
	}
	defer reqBody.Close()
	return raw
}

func (s *Server) CreateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("CreateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForCreate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Info(myRaw)
		s.App.Error("Ошибка перевода json в структуру - " + err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.App.CreateEvent(s.myCtx, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID))
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) UpdateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("UpdateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForUpdate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println("myStruct=", myStruct)
	myErr := s.App.UpdateEvent(s.myCtx, myStruct.EventID, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID)) //nolint
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) DeleteEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("DeleteEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := ForDelete{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myErr := s.App.DeleteEvent(s.myCtx, myStruct.EventID)
	if myErr != nil {
		s.App.Error(myErr)
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) GetEventByDateFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventByDate")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventByDate(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}

func (s *Server) GetEventMonthFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventMonth")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса!")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру!")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventMonth(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}

func (s *Server) GetEventByWeekFunc(rw http.ResponseWriter, req *http.Request) { //nolint:dupl
	s.App.Info("GetEventByWeekFunc")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	myStruct := StartDate{}
	if err := json.Unmarshal(myRaw, &myStruct); err != nil {
		s.App.Error("Ошибка перевода json в структуру")
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	eventList, myErr := s.App.GetEventWeek(s.myCtx, myStruct.StartDateStr)
	if myErr == nil {
		rawResp, err3 := json.Marshal(eventList)
		if err3 == nil {
			rw.Write(rawResp)
		} else {
			rw.WriteHeader(http.StatusInternalServerError)
			s.App.Error(err3)
			return
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
		s.App.Error(myErr)
	}
}

func (s *Server) Start(ctx context.Context) error {
	fmt.Printf("serv=%#v", s.HTTPConf)
	s.MyHTTP.Addr = s.HTTPConf
	//	s.myHTTP.Addr = ":" + s.HTTPConf.Port
	s.myCtx = ctx
	fmt.Println("serv=", s.HTTPConf)
	s.App.Info("serv=", s.HTTPConf)

	mux := http.NewServeMux()
	mux.HandleFunc("/CreateEvent", s.CreateEventFunc)
	mux.HandleFunc("/UpdateEvent", s.UpdateEventFunc)
	mux.HandleFunc("/DeleteEvent", s.DeleteEventFunc)
	mux.HandleFunc("/GetEventByDate", s.GetEventByDateFunc)
	mux.HandleFunc("/GetEventMonth", s.GetEventMonthFunc)
	mux.HandleFunc("/GetEventWeek", s.GetEventByWeekFunc)
	http.ListenAndServe(s.MyHTTP.Addr, s.loggingMiddleware(mux)) //nolint
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.MyHTTP.Shutdown(ctx)
	return nil
}
