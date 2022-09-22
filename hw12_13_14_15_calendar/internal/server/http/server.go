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

func getBodyRow(req_body io.ReadCloser) []byte {
	raw, err := ioutil.ReadAll(req_body)
	if err != nil {
		return nil
	}
	defer req_body.Close()
	return raw
}

func (s *Server) CreateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("CreateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := ForCreate{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			myErr := s.App.CreateEvent(s.myCtx, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID))
			if myErr != nil {
				s.App.Error(myErr)
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
	return
}

func (s *Server) UpdateEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("UpdateEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := ForUpdate{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			fmt.Println("myStruct=", myStruct)
			myErr := s.App.UpdateEvent(s.myCtx, myStruct.EventID, myStruct.Title, myStruct.StartDate, myStruct.Details, int(myStruct.UserID))
			if myErr != nil {
				s.App.Error(myErr)
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
	return
}

func (s *Server) DeleteEventFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("DeleteEvent")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := ForDelete{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			myErr := s.App.DeleteEvent(s.myCtx, myStruct.EventID)
			if myErr != nil {
				s.App.Error(myErr)
				rw.WriteHeader(http.StatusInternalServerError)
			}
		}
	}
	return
}

func (s *Server) GetEventByDateFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("GetEventByDate")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := StartDate{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			eventList, myErr := s.App.GetEventByDate(s.myCtx, myStruct.StartDateStr)
			if myErr == nil {
				rawResp, err3 := json.Marshal(eventList)
				if err3 == nil {
					rw.Write(rawResp)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
					s.App.Error(err3)
				}
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
				s.App.Error(myErr)

			}
		}
	}
	return
}

func (s *Server) GetEventMonthFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("GetEventMonth")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := StartDate{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			eventList, myErr := s.App.GetEventMonth(s.myCtx, myStruct.StartDateStr)
			if myErr == nil {
				rawResp, err3 := json.Marshal(eventList)
				if err3 == nil {
					rw.Write(rawResp)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
					s.App.Error(err3)
				}
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
				s.App.Error(myErr)

			}
		}
	}
	return
}

func (s *Server) GetEventByWeekFunc(rw http.ResponseWriter, req *http.Request) {
	s.App.Info("GetEventByWeekFunc")
	myRaw := getBodyRow(req.Body)
	if myRaw == nil {
		s.App.Error("Ошибка обработки тела запроса")
		rw.WriteHeader(http.StatusInternalServerError)
	} else {
		myStruct := StartDate{}
		if err := json.Unmarshal(myRaw, &myStruct); err != nil {
			s.App.Error("Ошибка перевода json в структуру")
			rw.WriteHeader(http.StatusInternalServerError)
		} else {
			eventList, myErr := s.App.GetEventWeek(s.myCtx, myStruct.StartDateStr)
			if myErr == nil {
				rawResp, err3 := json.Marshal(eventList)
				if err3 == nil {
					rw.Write(rawResp)
				} else {
					rw.WriteHeader(http.StatusInternalServerError)
					s.App.Error(err3)
				}
			} else {
				rw.WriteHeader(http.StatusInternalServerError)
				s.App.Error(myErr)

			}
		}
	}
	return
}

func (s *Server) Start(ctx context.Context) error {
	fmt.Printf("serv=%#v", s.HTTPConf)
	s.MyHTTP.Addr = s.HTTPConf
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
	http.ListenAndServe(s.MyHTTP.Addr, s.loggingMiddleware(mux))
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.MyHTTP.Shutdown(ctx)
	return nil
}
