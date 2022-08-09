package internalhttpgrpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	st "github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/internal/storage"
	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	App      Application
	HTTPConf string
	MyHTTP   *grpc.Server
}

func writeLogResult(ctx context.Context, myMethod string, app Application) {
	headers, ok := metadata.FromIncomingContext(ctx)
	var myIP, myClient string
	if !ok {
		IP := headers.Get(":authority")
		if len(IP) > 0 {
			myIP = IP[0]
		}
		Client := headers.Get("grpc-client")
		if len(Client) > 0 {
			myClient = Client[0]
		}
	}
	myTimeStart, ok := ctx.Value("TimeStart").(time.Time)
	if !ok {
		myTimeStart = time.Now()
		app.Error("Erorr getting myTimeStart from context:", myTimeStart)
	}
	app.Info(myIP, myClient, myMethod, myTimeStart, time.Since(myTimeStart))
}

type Service struct {
	pb.UnimplementedMyServServer
	myServer Server
}

func eventToRes(procEvents []st.Event) []*pb.OneEvent {
	myEvents := []*pb.OneEvent{}
	for i, v := range procEvents {
		myEvents = append(myEvents, &pb.OneEvent{})
		myEvents[i].Details = v.Details
		myEvents[i].EventID = v.ID
		myEvents[i].Title = v.Title
		myEvents[i].UserID = uint32(v.UserID)
		myEvents[i].StartDate = v.StartDate.Format("02.01.2006")
	}
	return myEvents
}

func (s *Service) CreateEvent(ctx context.Context, req *pb.ForCreate) (*pb.Error, error) {
	myErr := s.myServer.App.CreateEvent(ctx, req.Title, req.StartDate, req.Details, int(req.UserID))
	Resp := &pb.Error{}
	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}
	writeLogResult(ctx, "CreateEvent", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
}

func (s *Service) UpdateEvent(ctx context.Context, req *pb.ForUpdate) (*pb.Error, error) {
	myErr := s.myServer.App.UpdateEvent(ctx, req.EventID, req.Title, req.StartDate, req.Details, int(req.UserID))
	Resp := &pb.Error{}
	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}
	writeLogResult(ctx, "UpdateEvent", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
}

func (s *Service) DeleteEvent(ctx context.Context, req *pb.ForDelete) (*pb.Error, error) {
	Resp := &pb.Error{}
	myErr := s.myServer.App.DeleteEvent(ctx, req.EventID)
	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}
	writeLogResult(ctx, "DeleteEvent", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
}

func (s *Service) GetEventByDate(ctx context.Context, req *pb.StartDate) (*pb.Events, error) {
	Resp := &pb.Events{}
	myRes, myErr := s.myServer.App.GetEventByDate(ctx, req.StartDateStr)
	Resp.Result = eventToRes(myRes)

	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}

	writeLogResult(ctx, "GetEventByDate", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
}

func (s *Service) GetEventWeek(ctx context.Context, req *pb.StartDate) (*pb.Events, error) {
	Resp := &pb.Events{}
	myRes, myErr := s.myServer.App.GetEventWeek(ctx, req.StartDateStr)
	Resp.Result = eventToRes(myRes)
	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}
	writeLogResult(ctx, "GetEventWeek", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
}

func (s *Service) GetEventMonth(ctx context.Context, req *pb.StartDate) (*pb.Events, error) {
	Resp := &pb.Events{}
	myRes, myErr := s.myServer.App.GetEventMonth(ctx, req.StartDateStr)
	Resp.Result = eventToRes(myRes)
	if myErr == nil {
		Resp.ErrorText = "OK"
	} else {
		Resp.ErrorText = myErr.Error()
	}
	writeLogResult(ctx, "GetEventMonth", s.myServer.App)
	s.myServer.App.Print("Response=", Resp)
	return Resp, nil
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

func NewServer( /*logger Logger,*/ app Application, httpConf string) *Server {
	return &Server{App: app /*Logg: logger,*/, HTTPConf: httpConf}
}

func (s *Server) Start(ctx context.Context) error {
	lsn, err := net.Listen("tcp", s.HTTPConf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("lsn=", lsn)
	s.App.Info("serv=", s.HTTPConf)
	s.MyHTTP = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			myMiddleware(ValidateReq),
		),
	)
	fmt.Println("server=", s.MyHTTP)
	fmt.Println("*s=", *s)
	myService := &Service{myServer: *s}
	pb.RegisterMyServServer(s.MyHTTP, myService)
	log.Printf("starting server on %s", lsn.Addr().String())
	if err := s.MyHTTP.Serve(lsn); err != nil {
		log.Fatal(err)
	}
	log.Printf("starting server on %s", lsn.Addr().String())
	fmt.Printf("serv=%#v", s.HTTPConf)
	s.App.Info("qqq")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.MyHTTP.Stop()
	return nil
}
