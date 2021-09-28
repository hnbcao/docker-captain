package handler

import (
	"context"
	"docker-captain/docker"
	"docker-captain/docs"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	ginSwaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
)

type apiHandler struct {
	engine            *gin.Engine
	requestAPIHandler ImageAPIHandler
}

type APIHandler interface {
	Run(ctx context.Context, addr ...string) error
}

func (self *apiHandler) initHTTPAPIHandler(basePath string) {
	api := self.engine.Group(basePath)
	self.initSwaggerRoute(basePath)
	self.requestAPIHandler.registeredImageHttpHandler("image", api)

	self.engine.GET("/healthz", func(context *gin.Context) {
		context.String(http.StatusOK, "Health Ok")
	})
	self.engine.LoadHTMLFiles("static/index.html")

	self.engine.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	self.engine.StaticFile("/index.html", "static/index.html")
	self.engine.Static("/static", "static")
}

func (self *apiHandler) Run(ctx context.Context, addr ...string) error {
	var g errgroup.Group
	self.initHTTPAPIHandler("/api/v1")
	address := resolveAddress(addr)
	logrus.Info("http server listening to :", addr)
	srv := &http.Server{
		Addr:    address,
		Handler: self.engine,
	}

	g.Go(func() error {
		select {
		case <-ctx.Done():
			logrus.Info("shutdown http server")
			return srv.Shutdown(ctx)
		}
	})
	g.Go(func() error {
		return srv.ListenAndServe()
	})
	return g.Wait()
}

func (self *apiHandler) initSwaggerRoute(basePath string) {
	docs.SwaggerInfo.BasePath = basePath
	self.engine.GET("/swagger/*any", ginSwagger.WrapHandler(ginSwaggerFiles.Handler))
}

func NewRequestAPIHandler(dockerClient docker.Client, registryAuth map[string]string) ImageAPIHandler {
	return &imageAPIHandler{
		dockerClient: dockerClient,
		registryAuth: registryAuth,
		wsUpgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

func New(engine *gin.Engine, requestAPIHandler ImageAPIHandler) APIHandler {
	return &apiHandler{
		engine:            engine,
		requestAPIHandler: requestAPIHandler,
	}
}

func resolveAddress(addr []string) string {
	switch len(addr) {
	case 0:
		if port := os.Getenv("PORT"); port != "" {
			logrus.Debug("Environment variable PORT=\"%s\"", port)
			return ":" + port
		}
		logrus.Debug("Environment variable PORT is undefined. Using port :8080 by default")
		return ":8080"
	case 1:
		return addr[0]
	default:
		panic("too many parameters")
	}
}
