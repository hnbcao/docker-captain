package main

import (
	"context"
	"docker-captain/common"
	"docker-captain/docker"
	"docker-captain/handler"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	Debug = "debug"
	Info  = "info"
)

//var RegistryAuth = map[string]string{"123":buildAuth("admin", "fOciCYZ4EolU1uHR"),}

/**
DOCKER_SOURCE_USER=admin
DOCKER_SOURCE_PASSWORD=fOciCYZ4EolU1uHR
DOCKER_TARGET_USER=admin
DOCKER_TARGET_PASSWORD=cisdigital-12345
DEBUG_LEVEL=debug
APP_HOST=0.0.0.0
APP_PORT=8080
*/
func main() {
	app := cli.NewApp()
	app.Name = "docker client"
	app.Usage = "download file from minio"
	app.Action = run
	app.Version = "v0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "dockerHost",
			Usage:  "docker client ENV",
			EnvVar: "DOCKER_HOST",
			Value:  "",
		},
		cli.StringFlag{
			Name:   "debugLevel",
			Usage:  "debugLevel",
			EnvVar: "DEBUG_LEVEL",
		},
		cli.StringFlag{
			Name:   "appHost",
			Usage:  "host",
			EnvVar: "APP_HOST",
			Value:  "0.0.0.0",
		},
		cli.IntFlag{
			Name:   "appPort",
			Usage:  "appPort",
			EnvVar: "APP_PORT",
			Value:  8080,
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.WithError(err).Fatal("application execute failed")
	}

}

func run(c *cli.Context) error {
	// set log mode
	debugLevel := c.String("debugLevel")
	setMode(debugLevel)
	//sourceAuth, err := buildAuth(c.String("dockerSourceUser"), c.String("dockerSourcePassword"))
	//if err != nil {
	//	logrus.WithError(err).Error("create auth config error")
	//	return err
	//}
	//
	//targetAuth, err := buildAuth(c.String("dockerTargetUser"), c.String("dockerTargetPassword"))
	//if err != nil {
	//	logrus.WithError(err).Error("create auth config error")
	//	return err
	//}
	registryAuth := buildRegistryAuth()

	// init docker
	dockerClient, err := docker.NewClient()
	if err != nil {
		logrus.WithError(err).Error("create docker client error")
		return err
	}

	// init http server
	addr := net.JoinHostPort(c.String("appHost"), strconv.Itoa(c.Int("appPort")))
	engine := createEngine()
	ctx := common.WithContext(context.Background())
	requestHandler := handler.NewRequestAPIHandler(dockerClient, registryAuth)
	apiHandler := handler.New(engine, requestHandler)
	return apiHandler.Run(ctx, addr)
}

func buildRegistryAuth() map[string]string {
	registryAuth := map[string]string{
		"hub.d.cisdigital.cn":       "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJmT2NpQ1laNEVvbFUxdUhSIn0=",
		"harbor.test.cisdigital.cn": "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJjaXNkaWdpdGFsLTEyMzQ1In0=",
		"harbor.dev.cisdigital.cn":  "eyJ1c2VybmFtZSI6ImFkbWluIiwicGFzc3dvcmQiOiJER0xhOHJoeXpDV25uRFFmIn0=",
	}
	return registryAuth
}

func setMode(debugLevel string) {
	switch debugLevel {
	case Debug:
		logrus.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
	case Info:
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	default:
		logrus.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}
}

func createEngine() *gin.Engine {
	gin.ForceConsoleColor()
	engine := gin.New()
	engine.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	engine.Use(gin.Recovery())
	return engine
}

func buildAuth(userName string, password string) (string, error) {
	authConfig := types.AuthConfig{
		Username: userName,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
