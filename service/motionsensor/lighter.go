package motionsensor

import (
	"context"
	"sync"
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
	log            *zap.Logger
	lightsOffDelay time.Duration
	lights         []*Light
	scheduledOff   bool
	cancel         context.CancelFunc
}

func NewLighter(log *zap.Logger, lightsOffDelay time.Duration, addrByPosition map[LightPosition]string) *Lighter {
	lighter := &Lighter{
		log:            log,
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

// EnterRoom turns the lights on when someone enters the room for the first time
func (s *Lighter) EnterRoom(ctx context.Context) error {
	s.log.Info("EnterRoom")

	if s.scheduledOff {
		// If scheduled to be switched off, simply cancel the schedule
		s.cancel()
		s.scheduledOff = false
		return nil
	}

	var wg sync.WaitGroup
	errs := make(chan error, len(s.lights))
	// Turn lights on
	for _, light := range s.lights {
		wg.Add(1)
		go func(light *Light, wg *sync.WaitGroup) {
			s.log.Info("Turn light on",
				zap.Any("light", light),
			)
			defer wg.Done()
			errs <- light.Bulb.SetPower(true)
		}(light, &wg)
	}

	wg.Wait()
	// If an error has been added to the channel, return it
	if err := <-errs; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// LeaveRoom schedules the lights to be turned when the last person leaves the room
func (s *Lighter) LeaveRoom(ctx context.Context) error {
	s.log.Info("LeaveRoom")
	s.scheduledOff = true

	lightCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel
	go func() {
		_ = s.scheduleLightsOff(lightCtx)
	}()

	return nil
}

func (s *Lighter) scheduleLightsOff(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.lightsOffDelay):
			for _, light := range s.lights {
				s.log.Info("Turn light off",
					zap.Any("light", light),
				)
				if err := light.Bulb.SetPower(false); err != nil {
					return errors.Wrapf(err, "%d light", light.Position)
				}
			}
			s.scheduledOff = false
			return nil
		}
	}
}
