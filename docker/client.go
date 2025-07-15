package docker

import (
	"3cognito/coderunner/types"
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

	CreateContainer(ctx context.Context, language string) (containerID string, err error)
	RunContainer(ctx context.Context, fileData types.FileData, containerID string) (stdout string, stderr string, err error)
	RemoveContainer(ctx context.Context, containerID string) error
	ResetContainer(ctx context.Context, containerID string) error

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
	for _, runtime := range SupportedRuntimes {
		imageName := runtime.Image

		images, err := c.cli.ImageList(ctx, image.ListOptions{})
		if err != nil {
			log.Printf("error listing images: %v\n", err)
			return err
		}

		imageExists := false
		for _, img := range images {
			for _, tag := range img.RepoTags {
				if tag == imageName {
					imageExists = true
					break
				}
			}
			if imageExists {
				break
			}
		}
		if !imageExists {
			log.Printf("image not found locally: %s — pulling...\n", imageName)
			_, err := c.cli.ImagePull(ctx, imageName, image.PullOptions{})
			if err != nil {
				log.Printf("failed to pull image '%s': %v\n", imageName, err)
				return err
			}
			log.Printf("successfully pulled image: %s\n", imageName)
		} else {
			log.Printf("image already exists locally: %s\n", imageName)
		}
	}
	return nil
}
