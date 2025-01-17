package executors

import "time"

const defaultChunkSize = 1024 * 1024 // 1M

type (
	ChunkOption func(options *chunkOptions)

	ChunkExecutor struct {
		executor  *IntervalExecutor
		container *chunkContainer
	}

	chunkOptions struct {
		chunkSize     int
		flushInterval time.Duration
	}
)

func NewChunkExecutor(execute Execute, opts ...ChunkOption) *ChunkExecutor {
	options := newChunkOptions()
	for _, opt := range opts {
		opt(&options)
	}

	container := &chunkContainer{
		execute:      execute,
		maxChunkSize: options.chunkSize,
	}
	executor := &ChunkExecutor{
		executor:  NewIntervalExecutor(options.flushInterval, container),
		container: container,
	}

	return executor
}

func (ce *ChunkExecutor) Add(task any, size int) error {
	ce.executor.Add(chunk{
		val:  task,
		size: size,
	})
	return nil
}

func (ce *ChunkExecutor) Flush() {
	ce.executor.Flush()
}

func (ce *ChunkExecutor) Wait() {
	ce.executor.Wait()
}

func WithChunkBytes(size int) ChunkOption {
	return func(options *chunkOptions) {
		options.chunkSize = size
	}
}

func WithFlushInterval(duration time.Duration) ChunkOption {
	return func(options *chunkOptions) {
		options.flushInterval = duration
	}
}

func newChunkOptions() chunkOptions {
	return chunkOptions{
		chunkSize:     defaultChunkSize,
		flushInterval: defaultFlushInterval,
	}
}

type chunkContainer struct {
	tasks        []any
	execute      Execute
	size         int
	maxChunkSize int
}

func (bc *chunkContainer) AddItem(task any) bool {
	ck := task.(chunk)
	bc.tasks = append(bc.tasks, ck.val)
	bc.size += ck.size
	return bc.size >= bc.maxChunkSize
}

func (bc *chunkContainer) Execute(tasks any) {
	vals := tasks.([]any)
	bc.execute(vals)
}

func (bc *chunkContainer) RemoveAll() any {
	tasks := bc.tasks
	bc.tasks = nil
	bc.size = 0
	return tasks
}

type chunk struct {
	val  any
	size int
}
