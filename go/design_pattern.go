package main

import (
	"runtime"
	"sync"
)

// Fan-Out/Fan-In 模式(扇出/扇入模式)
// 1. 所有操作是否需要考虑并发安全
// 2. 是否需要有`NewTaskGroup()`方法，这样可以在这个方法中做一次一次性操作，比如初始化`fNOs`

type TaskGroup struct {
	once       sync.Once
	fNOs       map[uint32]struct{}
	workerNums uint32 // 工作组数量（协程数）
	tasks      []TaskDesc
}

type TaskDesc struct {
	fNO uint32 // 任务编号
	f   func() (interface{}, error)
}

type TaskResult struct {
	result interface{}
	err    error
}

type taskResult struct {
	fNO uint32
	tr  TaskResult
}

func (td *TaskGroup) AddTask(tasks ...TaskDesc) *TaskGroup {
	td.once.Do(func() {
		td.fNOs = make(map[uint32]struct{}, 2*len(tasks))
	})

	for _, task := range tasks {
		if _, exist := td.fNOs[task.fNO]; exist {
			// return errors.New("AddTask: Already have the same task") // 已经有相同的任务了
			panic("AddTask: Already have the same task") // 已经有相同的任务了
		}
		if task.f != nil {
			td.fNOs[task.fNO] = struct{}{}
			td.tasks = append(td.tasks, task)
		}
	}
	return td
}

// SetWorkersNums 设置任务所需的协程数
// Worker数量的选择：
// 硬件资源：系统上的 CPU 核心数量、内存大小和网络带宽等因素会限制可以并行运行的 worker 的数量。如果 worker 数量超过硬件资源能够支持的程度，那么增加更多的 worker 并不会提高整体性能，反而可能因为上下文切换和资源争用而降低性能
// 任务的性质：任务可能是 I/O 密集型（如网络请求或磁盘读写）或 CPU 密集型（如复杂的数学计算）。对于 I/O 密集型任务，增加 worker 数量可以更有效地利用等待时间，因为当一个 worker 在等待 I/O 操作完成时，其他 worker 可以继续执行。然而，对于 CPU 密集型任务，增加 worker 数量可能会导致 CPU 资源耗尽，从而降低性能。
// 任务的粒度：任务的粒度（即每个任务所需的时间）也会影响 worker 的数量。如果任务粒度很小，那么可能需要更多的 worker 来确保 CPU 始终保持忙碌状态。但是，如果任务粒度很大，那么少量的 worker 就足以处理所有任务，增加 worker 数量可能是不必要的。
func (td *TaskGroup) SetWorkersNums(workerNums uint32) *TaskGroup {
	td.workerNums = workerNums
	if td.workerNums == 0 {
		td.workerNums = uint32(runtime.NumCPU())
	}
	return td
}

func (td *TaskGroup) Run() map[uint32]TaskResult {
	var (
		tasks   = make(chan TaskDesc, len(td.tasks))
		results = make(chan taskResult, len(td.tasks))

		wg sync.WaitGroup
	)
	if td.workerNums > uint32(len(td.tasks)) {
		td.workerNums = uint32(len(td.tasks))
	}
	// Start workers
	for i := 1; i <= int(td.workerNums); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(tasks, results)
		}()
	}

	// Send tasks to workers
	for i := 0; i < len(td.tasks); i++ {
		tasks <- td.tasks[i]
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from workers
	taskResults := make(map[uint32]TaskResult, len(td.tasks))
	for result := range results {
		taskResults[result.fNO] = result.tr
	}
	return taskResults
}

func worker(tasks chan TaskDesc, results chan taskResult) {
	for task := range tasks {
		result, err := task.f()
		results <- taskResult{fNO: task.fNO, tr: TaskResult{result, err}}
	}
}