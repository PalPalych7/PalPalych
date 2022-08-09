package internalhttpgrpc

import (
	"context"
	"errors"
	"time"

	"github.com/PalPalych7/PalPalych/hw12_13_14_15_calendar/pb"
	"google.golang.org/grpc"
)

type Validator func(req interface{}) error

func myMiddleware(validator Validator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) { //nolint:lll
		startTime := time.Now()
		ctx = context.WithValue(ctx, "TimeStart", startTime) //nolint:staticcheck
		if err := validator(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func ValidateReq(req interface{}) error {
	switch r := req.(type) {
	case *pb.ForCreate:
		if r.UserID == 0 {
			return errors.New("middleware validator: UserID wrong")
		}
	case *pb.StartDate:
		_, err := time.Parse("2.1.2006", r.StartDateStr)
		if err != nil {
			return errors.New("middleware validator: StartDateStr wrong")
		}
	default:
	}
	return nil
}
