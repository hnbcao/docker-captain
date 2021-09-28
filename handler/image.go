package handler

import (
	"context"
	"docker-captain/core"
	"docker-captain/docker"
	"docker-captain/process"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

type (
	imageAPIHandler struct {
		dockerClient docker.Client
		registryAuth map[string]string
		wsUpgrader   websocket.Upgrader
	}

	ImageAPIHandler interface {
		registeredImageHttpHandler(relativePath string, api *gin.RouterGroup)
	}
)

//var buildStatus = make(map[int64]int)

func (handler *imageAPIHandler) registeredImageHttpHandler(relativePath string, api *gin.RouterGroup) {
	imageGroup := api.Group(relativePath)
	imageGroup.POST("/tag", handler.tagImageHandler)
	imageGroup.GET("/tag/ws", func(ctx *gin.Context) {
		handler.tagImageWsHandler(ctx.Writer, ctx.Request)
	})
}

// @BasePath /api/v1
// tagImageHandler godoc
// @Summary 推送镜像
// @Description 推送镜像
// @Tags image
// @Accept  json
// @Produce  json
// @Param image body core.Image true "image info"
// @Success 200 {string} true
// @Router /image/tag [post]
func (self *imageAPIHandler) tagImageHandler(ctx *gin.Context) {
	image := &core.Image{}
	if err := ctx.ShouldBindBodyWith(image, binding.JSON); err != nil {
		ctx.JSON(http.StatusOK, Failure("500", "参数错误"))
		return
	}
	logrus.Info("image param is: ", *image)
	image.BuildRegistryAuth(self.registryAuth)

	err := process.CreateNewImage(context.Background(), image, self.dockerClient, docker.NewWriter(nil, 0))
	if err != nil {
		ctx.JSON(http.StatusOK, Failure("500", "推送镜像失败"))
		return
	}
	ctx.JSON(http.StatusOK, Success(true))
}

func (handler *imageAPIHandler) tagImageWsHandler(writer gin.ResponseWriter, request *http.Request) {
	conn, err := handler.wsUpgrader.Upgrade(writer, request, nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to set websocket upgrade")
		return
	}
	id := uuid.New()
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		image := &core.Image{}
		json.Unmarshal(msg, image)
		logrus.Info("image param is: ", *image)
		conn.WriteMessage(messageType, []byte(fmt.Sprintf("receive param %s", *image)))
		image.BuildRegistryAuth(handler.registryAuth)

		out := fmt.Sprintf("success=====>>>%s:%s", id.String(), string(msg))

		writer := docker.NewWriter(conn, messageType)

		if err != nil {
			break
		}
		err = process.CreateNewImage(context.Background(), image, handler.dockerClient, writer)
		if err != nil {
			conn.WriteMessage(messageType, []byte("PROCESS_MESSAGE>>error: "+err.Error()))
		} else {
			conn.WriteMessage(messageType, []byte("PROCESS_MESSAGE>>success: "+out))
		}
		go process.CleanImages(context.Background(), handler.dockerClient, image.SourceImage, image.TargetImage)
	}
}

type Result struct {
	Data interface{} `json:"data"`

	Success bool `json:"success"`

	RequestId string `json:"requestId"`

	ResultCode string `json:"resultCode"`

	ResultMsg string `json:"resultMsg"`

	Solution string `json:"solution"`
}

func Success(data interface{}) Result {
	if data == nil {
		data = gin.H{}
	}
	res := Result{}
	res.Success = true
	res.Data = data
	res.ResultCode = "200"
	res.RequestId = ""
	res.ResultMsg = ""
	res.Solution = ""
	return res
}

func Failure(code string, message string) Result {
	res := Result{}
	res.Success = false
	res.Data = false
	res.ResultCode = code
	res.RequestId = ""
	res.ResultMsg = message
	res.Solution = ""
	return res
}
