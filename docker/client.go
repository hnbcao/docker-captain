package docker

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

type Client interface {
	BuildDockerImage(ctx context.Context, workDir string, tags ...string) error
	ClientVersion() string
	PullDockerImage(ctx context.Context, image, registryAuth string, writer io.Writer) error
	TagDockerImage(ctx context.Context, source, target string, writer io.Writer) error
	PushDockerImage(ctx context.Context, image, registryAuth string, writer io.Writer) error
	RemoveDockerImage(ctx context.Context, image string) error
}

type dockerClient struct {
	client *client.Client
}

func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logrus.WithError(err).Error("create client error")
		return nil, err
	}
	return &dockerClient{
		client: cli,
	}, nil
}

func (client *dockerClient) ClientVersion() string {
	return client.client.ClientVersion()
}

// BuildDockerImage build镜像
func (client *dockerClient) BuildDockerImage(ctx context.Context, workDir string, tags ...string) error {
	buildCtx, _ := archive.TarWithOptions(workDir, &archive.TarOptions{})
	ops := types.ImageBuildOptions{
		Tags: tags,
	}
	res, err := client.client.ImageBuild(ctx, buildCtx, ops)
	if err == nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		errStr := buf.String()
		results := strings.Split(errStr, "\r\n")
		logrus.Debug(results)
		for _, result := range results {
			if strings.Contains(result, "errorDetail") {
				logrus.Info("build image error", result)
				return errors.New(result)
			}
		}
	}
	return err
}

// PullDockerImage 镜像仓库拉取基础镜像
func (client *dockerClient) PullDockerImage(ctx context.Context, image, registryAuth string, writer io.Writer) error {
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>begin to pull image %s", image)))
	out, err := client.client.ImagePull(ctx, image, types.ImagePullOptions{RegistryAuth: registryAuth})
	if err != nil {
		logrus.WithError(err).Error("pull base image error")
		return err
	}
	if _, err := io.Copy(writer, out); err != nil {
		logrus.WithError(err).Error("copy base image error ")
		return err
	}
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>success to pull image %s", image)))
	return nil
}

// TagDockerImage 镜像仓库拉取基础镜像
func (client *dockerClient) TagDockerImage(ctx context.Context, source, target string, writer io.Writer) error {
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>begin to tag image %s to %s", source, target)))
	err := client.client.ImageTag(ctx, source, target)
	if err != nil {
		logrus.WithError(err).Error("pull base image error")
		return err
	}
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>success to tag image %s to %s", source, target)))
	return nil
}

// PushDockerImage push镜像到镜像仓库
func (client *dockerClient) PushDockerImage(ctx context.Context, image, registryAuth string, writer io.Writer) error {
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>begin to push image %s", image)))
	out, err := client.client.ImagePush(ctx, image, types.ImagePushOptions{RegistryAuth: registryAuth})
	if err != nil {
		logrus.WithError(err).Error("pull base image error")
		return err
	}
	scanner := bufio.NewScanner(out)
	msg := ""
	for scanner.Scan() {
		msg = scanner.Text()
		if strings.Contains(msg, "errorDetail") {
			return errors.New(msg)
		}
		writer.Write([]byte(msg))
	}
	writer.Write([]byte(fmt.Sprintf("PROCESS_MESSAGE>>=======>>>>success to push image %s", image)))
	return err
}

func (client *dockerClient) RemoveDockerImage(ctx context.Context, image string) error {
	_, err := client.client.ImageRemove(ctx, image, types.ImageRemoveOptions{
		Force:         false,
		PruneChildren: false,
	})
	return err
}
