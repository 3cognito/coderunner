package docker

import (
	"context"
)

func (c *Client) Execute(ctx context.Context, fileData FileData) (string, string, error) {
	config, err := c.getContainerConfigs(fileData.Language)
	if err != nil {
		return "", "", err
	}

	containerID, err := c.CreateContainer(ctx, config)
	if err != nil {
		return "", "", err
	}

	stdout, stderr, err := c.RunContainer(ctx, fileData, containerID)

	go c.RemoveContainer(ctx, containerID)

	return stdout, stderr, err
}
