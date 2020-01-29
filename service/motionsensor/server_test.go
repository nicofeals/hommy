package motionsensor

import (
	"context"
	"testing"

	pb "github.com/nicofeals/hommy/rpc/motionsensor"
	"github.com/nicofeals/hommy/service/motionsensor/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type ServerTestSuite struct {
	suite.Suite
	log *zap.Logger
	lc  *mocks.LightController
	s   *Server
}

func (s *ServerTestSuite) SetupTest() {
	log, _ := zap.NewDevelopment()
	s.log = log
	s.lc = new(mocks.LightController)
	s.s = NewServer(s.log, s.lc)
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) TestDetectMotion_ComingIn() {
	s.lc.On("TurnLightOn", mock.Anything).Return(nil)

	res, err := s.s.DetectMovement(context.Background(), &pb.DetectMovementRequest{Position: pb.DetectMovementRequest_DOOR})
	s.NoError(err)
	s.Empty(res)
	s.True(s.s.doorSensorActive)
	s.False(s.s.deskSensorActive)
	s.Equal(uint32(0), s.s.counter)

	res, err = s.s.DetectMovement(context.Background(), &pb.DetectMovementRequest{Position: pb.DetectMovementRequest_DESK})
	s.NoError(err)
	s.Empty(res)
	s.False(s.s.doorSensorActive)
	s.False(s.s.deskSensorActive)
	s.Equal(uint32(1), s.s.counter)

	mock.AssertExpectationsForObjects(s.T(), s.lc)
}

func (s *ServerTestSuite) TestDetectMotion_GoingOut() {
	s.lc.On("TurnLightOff", mock.Anything).Return(nil)

	s.s.counter = 4
	res, err := s.s.DetectMovement(context.Background(), &pb.DetectMovementRequest{Position: pb.DetectMovementRequest_DESK})
	s.NoError(err)
	s.Empty(res)
	s.True(s.s.deskSensorActive)
	s.False(s.s.doorSensorActive)
	s.Equal(uint32(4), s.s.counter)

	res, err = s.s.DetectMovement(context.Background(), &pb.DetectMovementRequest{Position: pb.DetectMovementRequest_DOOR})
	s.NoError(err)
	s.Empty(res)
	s.False(s.s.doorSensorActive)
	s.False(s.s.deskSensorActive)
	s.Equal(uint32(3), s.s.counter)

	mock.AssertExpectationsForObjects(s.T(), s.lc)
}
