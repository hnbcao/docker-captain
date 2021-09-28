package process

import (
	"context"
	"docker-captain/core"
	"docker-captain/docker"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
)

func CleanImages(context context.Context, dockerClient docker.Client, images ...string) {
	for _, image := range images {
		if err := dockerClient.RemoveDockerImage(context, image); err != nil {
			logrus.WithError(err).Error(fmt.Sprintf("remove image %s error", image))
		}
	}
}

func CreateNewImage(context context.Context, request *core.Image, dockerClient docker.Client, writer io.Writer) error {
	err := dockerClient.PullDockerImage(context, request.SourceImage, request.SourceAuth, writer)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("cannot pull image %s with auth %s", request.SourceImage, request.SourceAuth))
		return err
	}
	err = dockerClient.TagDockerImage(context, request.SourceImage, request.TargetImage, writer)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("cannot tag image %s to %s", request.SourceImage, request.TargetImage))
		return err
	}
	err = dockerClient.PushDockerImage(context, request.TargetImage, request.TargetAuth, writer)
	if err != nil {
		logrus.WithError(err).Error(fmt.Sprintf("cannot push image %s with auth %s", request.TargetImage, request.TargetAuth))
		return err
	}
	return err
}
