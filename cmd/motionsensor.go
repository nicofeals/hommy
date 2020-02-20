package main

import (
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/nicofeals/hommy/common/config"
	"github.com/nicofeals/hommy/common/hooks"
	pb "github.com/nicofeals/hommy/rpc/motionsensor"
	"github.com/nicofeals/hommy/service/motionsensor"
	"github.com/twitchtv/twirp"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

func startMotionSensorServerHandler(c *cli.Context) {
	rand.Seed(time.Now().UnixNano())
	env := getEnvironment(c)
	log := config.NewZapLogger(getLogLevel(c))

	log.Info("loading config",
		zap.String("env", env),
	)

	errHooks := hooks.ErrorLoggerHooks(log)

	lights := map[motionsensor.LightPosition]string{
		motionsensor.Corner:    "192.168.1.28:55443",
		motionsensor.Ceiling:   "192.168.1.29:55443",
		motionsensor.Ambilight: "192.168.1.32:55443",
		motionsensor.Bedside:   "192.168.1.33:55443",
	}
	lc := motionsensor.NewLighter(log, getLightsOffDelay(c), lights)
	server := motionsensor.NewServer(log, lc)

	handler := pb.NewMotionSensorServer(server, twirp.ChainHooks(errHooks))

	log.Info("start")

	host := "localhost"
	if strings.EqualFold(env, "production") {
		host = "0.0.0.0"
	}
	port := strings.TrimSpace(c.String("port"))

	if err := http.ListenAndServe(net.JoinHostPort(host, port), handler); err != nil {
		log.Error("listen", zap.Error(err))
	}
}
