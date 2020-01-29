package motionsensor

import (
	"context"
	"time"

	"github.com/avarabyeu/yeelight"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type LightPosition int

const (
	Ceiling LightPosition = iota
	Corner
	Bedside
	Ambilight
)

type Yeelight interface {
}

type Light struct {
	Position LightPosition
	Bulb     *yeelight.Yeelight
}

type Lighter struct {
	l              *zap.Logger
	lightsOffDelay time.Duration
	lights         []*Light
	scheduledOff   bool
	cancel         context.CancelFunc
}

func NewLighter(l *zap.Logger, lightsOffDelay time.Duration, addrByPosition map[LightPosition]string) *Lighter {
	lighter := &Lighter{
		l:              l,
		lightsOffDelay: lightsOffDelay,
		scheduledOff:   false,
	}

	lighter.lights = make([]*Light, len(addrByPosition))
	i := 0
	for pos, addr := range addrByPosition {
		lighter.lights[i] = &Light{
			Position: pos,
			Bulb:     yeelight.New(addr),
		}
		i++
	}

	return lighter
}

func (s *Lighter) EnterRoom(ctx context.Context) error {
	s.l.Info("EnterRoom")

	if s.scheduledOff {
		// If scheduled to be switched off, simply cancel the schedule
		s.cancel()
		s.scheduledOff = false
		return nil
	}

	return nil
}

func (s *Lighter) LeaveRoom(ctx context.Context) error {
	s.l.Info("LeaveRoom")
	s.scheduledOff = true

	lightCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	s.scheduleLightsOff(lightCtx)

	return nil
}

func (s *Lighter) scheduleLightsOff(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.lightsOffDelay):
			for _, light := range s.lights {
				s.l.Info("Turn light off",
					zap.Any("light", light),
				)
				if err := light.Bulb.SetPower(false); err != nil {
					return errors.Wrapf(err, "%s light", light.Position)
				}
			}
		}
	}
}
