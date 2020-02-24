package motionsensor

import (
	"context"

	pb "github.com/nicofeals/hommy/rpc/motionsensor"
	"github.com/twitchtv/twirp"
	"go.uber.org/zap"
)

type LightController interface {
	EnterRoom(ctx context.Context) error
	LeaveRoom(ctx context.Context) error
}

type Server struct {
	l                *zap.Logger
	lc               LightController
	doorSensorActive bool
	deskSensorActive bool
	counter          uint32
}

func NewServer(l *zap.Logger, lc LightController) *Server {
	return &Server{
		l:                l,
		lc:               lc,
		doorSensorActive: false,
		deskSensorActive: false,
		counter:          0,
	}
}

func (s *Server) DetectMovement(ctx context.Context, r *pb.DetectMovementRequest) (res *pb.DetectMovementResponse, e error) {
	// Log request
	s.l.Info("DetectMovement",
		zap.String("position", r.GetPosition().String()),
	)
	res = &pb.DetectMovementResponse{Position: r.GetPosition()}

	switch r.GetPosition() {
	case pb.Position_DOOR:
		if !s.deskSensorActive {
			// Desk sensor is not active, means the person is coming in
			// Set door sensor to active and return
			s.doorSensorActive = true
			return
		}
		if s.counter == 0 {
			s.deskSensorActive = false
			// There's already no one in the room and someone has just gone out...
			return
		}

		// The person is going out
		s.counter--

		s.l.Info("Out",
			zap.Uint32("counter", s.counter),
		)
		// Deactivate desk sensor and decrease counter
		s.deskSensorActive = false

		if s.counter != 0 {
			return
		}
		if err := s.lc.LeaveRoom(ctx); err != nil {
			return nil, twirp.InternalErrorWith(err)
		}
	case pb.Position_DESK:
		if !s.doorSensorActive {
			// Door sensor is not active, means the person is going out
			// Set desk sensor to active and return
			s.deskSensorActive = true
			return
		}

		// The person is coming in
		s.counter++
		s.l.Info("In",
			zap.Uint32("counter", s.counter),
		)
		// Deactivate desk sensor and decrease counter
		s.doorSensorActive = false

		if s.counter != 1 {
			return
		}
		if err := s.lc.EnterRoom(ctx); err != nil {
			return nil, twirp.InternalErrorWith(err)
		}
	}

	return
}
