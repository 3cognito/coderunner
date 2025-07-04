package docker

import (
	"context"
	"log"

	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
}

type ClientInterface interface {
	Close() error
	Ping(ctx context.Context) error
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

func (c *Client) Ping(ctx context.Context) error {
	_, err := c.cli.Ping(ctx)
	return err
}
