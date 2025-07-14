package docker

import (
	"context"
	"log"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

type ClientInterface interface {
	InitImages(ctx context.Context) error

	Execute(ctx context.Context, fileData FileData) (stdout string, stderr string, err error)

	CreateContainer(ctx context.Context, config ContainerConfig) (containerID string, err error)
	RunContainer(ctx context.Context, fileData FileData, containerID string) (stdout string, stderr string, err error)
	RemoveContainer(ctx context.Context, containerID string) error

	Close() error
}

func NewClient() ClientInterface {
	cli, err := client.NewClientWithOpts()
	if err != nil {
		log.Println("Error establishing connection to docker daemon", err)
		panic(err)
	}

	if _, err = cli.Ping(context.Background()); err != nil {
		log.Println("Failed to ping Docker daemon:", err)
		panic(err)
	}

	log.Println("Successfully connected to Docker daemon.")

	return &Client{
		cli: cli,
	}
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) InitImages(ctx context.Context) error {
	for _, data := range SupportedRuntimes {
		_, err := c.cli.ImagePull(context.Background(), data.Image, image.PullOptions{})
		if err != nil {
			log.Printf("unable to pull image '%s': %v\n", data.Image, err)
			return err
		}
	}
	return nil
}
