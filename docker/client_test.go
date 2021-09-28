package docker

import (
	"context"
	"testing"
)

func TestClientVersion(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	t.Log(cli.ClientVersion())
}

func TestBuildDockerImage(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	//workDir := "C:\\Users\\rick\\go\\src\\awesomeProject"
	workDir := "E:\\build-test"
	//buildCtx, _ := archive.TarWithOptions("C:\\Users\\rick\\go\\src\\awesomeProject", &archive.TarOptions{})
	tags := []string{"test:v1"}

	t.Log(cli.BuildDockerImage(ctx, workDir, tags...))
}

func TestPullDockerImage(t *testing.T) {
	cli, err := NewClient()
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	alpine := "vimodetect"

	t.Log(cli.PullDockerImage(ctx, alpine, "", NewWriter(nil, 0)))
}
