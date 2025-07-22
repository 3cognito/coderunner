package executor

import (
	"3cognito/coderunner/cache"
	"3cognito/coderunner/docker"
	"3cognito/coderunner/types"
	"3cognito/coderunner/utils"
	"context"
)

type Executor struct {
	containerPool   ContainerPoolInterface
	containerClient docker.ClientInterface
	cache           cache.CacheInterface
}

type ExecutorInterface interface {
	Execute(ctx context.Context, fileData types.FileData) (types.ExecutionOutput, error)
}

func NewExecutor(pool ContainerPoolInterface, client docker.ClientInterface, cache cache.CacheInterface) ExecutorInterface {
	return &Executor{
		containerPool:   pool,
		containerClient: client,
		cache:           cache,
	}
}

func (e *Executor) Execute(ctx context.Context, fileData types.FileData) (types.ExecutionOutput, error) {
	var output types.ExecutionOutput
	cacheKey := utils.HashData(fileData)

	if result, err := e.cache.Get(cacheKey); err == nil {
		return result, nil
	}

	containerID, err := e.containerPool.GetContainer(ctx, fileData.Language)
	if err != nil {
		return output, err
	}

	stdout, stderr, err := e.containerClient.RunContainer(ctx, fileData, containerID)
	if err != nil {
		return output, err
	}

	output.Stdout = stdout
	output.Stderr = stderr

	go e.cache.Set(cacheKey, output)

	go e.containerPool.CleanUp(ctx, containerID, fileData.Language)

	return output, nil
}
