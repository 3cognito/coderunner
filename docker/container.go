package docker

import (
	"3cognito/coderunner/types"
	"3cognito/coderunner/utils"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
)

type ContainerConfig struct {
	Image       string
	Command     []string
	WorkingDir  string
	Env         []string
	MemoryLimit int64
	CPULimit    int64
	Timeout     int // seconds
}

var (
	ErrUnsupportedLanguage = errors.New("unsupported language")
)

func (c *Client) CreateContainer(ctx context.Context, language string) (string, error) {
	config, err := c.getContainerConfigs(language)
	if err != nil {
		return "", err
	}

	containerConfig := &container.Config{
		Image:      config.Image,
		Cmd:        config.Command,
		WorkingDir: config.WorkingDir,
		Env:        config.Env,
	}

	hostConfig := &container.HostConfig{
		Resources: container.Resources{
			Memory:   config.MemoryLimit,
			NanoCPUs: config.CPULimit,
		},
		NetworkMode: "none",
		CapDrop:     []string{"ALL"},
	}

	resp, err := c.cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *Client) RunContainer(ctx context.Context, fileData types.FileData, containerID string) (string, string, error) {
	if err := c.copyToContainer(ctx, fileData, containerID); err != nil {
		return "", "", err
	}

	if err := c.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
		return "", "", err
	}

	statusCh, errCh := c.cli.ContainerWait(ctx, containerID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		return "", "", err
	case <-statusCh:
	}

	stdout, stderr, err := c.getLogs(ctx, containerID)
	if err != nil {
		return "", "", err
	}

	return stdout, stderr, nil
}

func (c *Client) RemoveContainer(ctx context.Context, containerID string) error {
	if err := c.cli.ContainerRemove(ctx, containerID, container.RemoveOptions{
		Force:         true,
		RemoveVolumes: true,
	}); err != nil {
		fmt.Println("Error removing container with id", containerID, "error", err)
		return err
	}
	return nil
}

func (c *Client) getContainerConfigs(language string) (ContainerConfig, error) {
	runtime, err := GetRuntime(language)
	if err != nil {
		return ContainerConfig{}, err
	}

	config := ContainerConfig{
		Image:      runtime.Image,
		Command:    runtime.Command,
		WorkingDir: "/app",
		// Env:         runtime.Env,
		MemoryLimit: 256 * 1024 * 1024,
		CPULimit:    550_000_000,
		Timeout:     10,
	}

	return config, nil
}

func (c *Client) copyToContainer(ctx context.Context, fileData types.FileData, containerID string) error {
	runtime, err := GetRuntime(fileData.Language)
	if err != nil {
		return ErrUnsupportedLanguage
	}

	tarBytes, err := utils.TarFile(runtime.FileName, fileData.Content)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(tarBytes)
	if err = c.cli.CopyToContainer(ctx, containerID, "/app", reader, container.CopyToContainerOptions{}); err != nil {
		return err
	}

	return nil
}

func (c *Client) ResetContainer(ctx context.Context, containerID string) error {
	res, err := c.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return err
	}

	if !res.State.Running {
		log.Printf("container %s is not running, restarting it for reset...", containerID)
		if err := c.cli.ContainerStart(ctx, containerID, container.StartOptions{}); err != nil {
			log.Printf("failed to start container %s: %v", containerID, err)
			return err
		}
	}

	exec, err := c.cli.ContainerExecCreate(ctx, containerID, container.ExecOptions{
		Cmd: []string{"rm", "-rf", "/tmp/*", "/app/*"},
	})
	if err != nil {
		log.Printf("error: failed to create exec command in container %s: %v", containerID, err)
		return err
	}

	if err := c.cli.ContainerExecStart(ctx, exec.ID, container.ExecStartOptions{}); err != nil {
		log.Printf("error: failed to start exec command in container %s: %v", containerID, err)
		return err
	}

	return nil
}

func (c *Client) getLogs(ctx context.Context, containerID string) (string, string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		// Follow:     true,
	}

	logs, err := c.cli.ContainerLogs(ctx, containerID, options)
	if err != nil {
		return "", "", err
	}
	defer logs.Close()

	var stdoutBuf, stderrBuf bytes.Buffer

	_, err = stdcopy.StdCopy(&stdoutBuf, &stderrBuf, logs)
	if err != nil {
		return "", "", fmt.Errorf("failed to copy log output: %w", err)
	}

	return stdoutBuf.String(), stderrBuf.String(), nil
}
