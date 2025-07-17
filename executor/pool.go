package executor

import (
	"3cognito/coderunner/docker"
	"context"
	"sync"
)

type ContainerPool struct {
	containerClient    docker.ClientInterface
	languageChannelMap map[string]chan string
	mutex              sync.RWMutex
}

type ContainerPoolInterface interface {
	GetContainer(ctx context.Context, language string) (string, error)
	CleanUp(ctx context.Context, containerID, language string) error
}

func NewContainerPool(containerClient docker.ClientInterface, langChanMap map[string]chan string) ContainerPoolInterface {
	return &ContainerPool{
		containerClient:    containerClient,
		languageChannelMap: langChanMap,
		mutex:              sync.RWMutex{},
	}
}

func (cp *ContainerPool) getorCreate(language string) chan string {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()
	channel, exists := cp.languageChannelMap[language]
	if !exists {
		channel = make(chan string, 10)
		cp.languageChannelMap[language] = channel
		return channel
	}

	return channel
}

func (cp *ContainerPool) GetContainer(ctx context.Context, language string) (string, error) {
	channel := cp.getorCreate(language)

	select {
	case containerID := <-channel:
		return containerID, nil
	default:
		return cp.containerClient.CreateContainer(ctx, language)
	}
}

func (cp *ContainerPool) CleanUp(ctx context.Context, containerID, language string) error {
	if err := cp.containerClient.ResetContainer(ctx, containerID); err != nil {
		return err
	}

	channel := cp.getorCreate(language)

	select {
	case channel <- containerID:
	default:
		go cp.containerClient.RemoveContainer(ctx, containerID)
	}

	return nil
}
