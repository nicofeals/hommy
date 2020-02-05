package motionsensor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type LighterTestSuite struct {
	suite.Suite
	log *zap.Logger
	s   *Lighter
}

func (s *LighterTestSuite) SetupTest() {
	log, _ := zap.NewDevelopment()
	s.log = log
	s.s = NewLighter(log, 0*time.Second, map[LightPosition]string{
		Corner:    "192.168.1.28:55443",
		Ceiling:   "192.168.1.29:55443",
		Ambilight: "192.168.1.32:55443",
		Bedside:   "192.168.1.33:55443",
	})
}

func TestLighterTestSuite(t *testing.T) {
	suite.Run(t, new(LighterTestSuite))
}

func (s *LighterTestSuite) TestEnterRoom_OK() {
	err := s.s.EnterRoom(context.Background())
	s.NotNil(err)
}
