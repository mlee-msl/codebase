package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Fan-Out/Fan-In 模式(扇出/扇入模式)
// 1. 所有操作是否需要考虑并发安全
// 2. 是否需要有`NewTaskGroup()`方法，这样可以在这个方法中做一次一次性操作，比如初始化`fNOs`

// TaskGroup 表示可将多个任务进行安全并发执行的一个对象
type TaskGroup struct {
	once       sync.Once
	fNOs       map[uint32]struct{}
	workerNums uint32  // 工作组数量（协程数）
	tasks      []*Task // 待执行的任务集合
}

type Task struct {
	fNO         uint32                      // 任务编号
	f           func() (interface{}, error) // 任务方法
	mustSuccess bool                        // 任务必须执行成功，否则整个任务组将会立即结束，且失败
}

// NewTask 创建一个任务
func NewTask(fNO uint32 /* 任务唯一标识 */, mustSuccess bool /* 标识任务是否必须执行成功 */, f func() (interface{}, error) /* 任务执行方法 */) *Task {
	return &Task{fNO, f, mustSuccess}
}

// SetWorkerNums 设置任务所需的协程数
// Worker数量的选择：
// 硬件资源：系统上的 CPU 核心数量、内存大小和网络带宽等因素会限制可以并行运行的 worker 的数量。如果 worker 数量超过硬件资源能够支持的程度，那么增加更多的 worker 并不会提高整体性能，反而可能因为上下文切换和资源争用而降低性能
// 任务的性质：任务可能是 I/O 密集型（如网络请求或磁盘读写）或 CPU 密集型（如复杂的数学计算）。对于 I/O 密集型任务，增加 worker 数量可以更有效地利用等待时间，因为当一个 worker 在等待 I/O 操作完成时，其他 worker 可以继续执行。然而，对于 CPU 密集型任务，增加 worker 数量可能会导致 CPU 资源耗尽，从而降低性能。
// 任务的粒度：任务的粒度（即每个任务所需的时间）也会影响 worker 的数量。如果任务粒度很小，那么可能需要更多的 worker 来确保 CPU 始终保持忙碌状态。但是，如果任务粒度很大，那么少量的 worker 就足以处理所有任务，增加 worker 数量可能是不必要的。
func (tg *TaskGroup) SetWorkerNums(workerNums uint32) *TaskGroup {
	tg.workerNums = workerNums
	if tg.workerNums == 0 {
		tg.workerNums = uint32(runtime.NumCPU())
	}
	return tg
}

// AddTask 向任务组中添加若干待执行的任务
func (tg *TaskGroup) AddTask(tasks ...*Task) *TaskGroup {
	tg.once.Do(func() {
		var preAllocatedCapacity = 2*len(tasks) + 1
		tg.fNOs = make(map[uint32]struct{}, preAllocatedCapacity)
		tg.tasks = make([]*Task, 0, preAllocatedCapacity)
	})

	for _, task := range tasks {
		if task == nil {
			continue
		}
		if _, exist := tg.fNOs[task.fNO]; exist {
			panic(fmt.Sprintf("AddTask: Already have the same task %d", task.fNO)) // 已经有相同的任务了
		}
		if task.f != nil {
			tg.fNOs[task.fNO] = struct{}{}
			tg.tasks = append(tg.tasks, task)
		}
	}
	return tg
}

// Run 启动并运行任务组中的所有任务
func (tg *TaskGroup) Run() map[uint32]*taskResult {
	if len(tg.tasks) == 0 {
		return nil
	}

	var (
		tasks   = make(chan *Task, len(tg.tasks))
		results = make(chan *taskResult, len(tg.tasks))

		wg sync.WaitGroup
	)
	// 工作协程数不得多余待执行任务总数，否则，因多余协程不会做任务，反而会由于创建或销毁这些协程而带来额外不必要的性能消耗
	if tg.workerNums > uint32(len(tg.tasks)) {
		tg.workerNums = uint32(len(tg.tasks))
	}
	// Start workers
	for i := 1; i <= int(tg.workerNums); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tg.worker(tasks, results)
		}()
	}

	// Send tasks to workers
	for i := 0; i < len(tg.tasks); i++ {
		tasks <- tg.tasks[i]
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results from workers
	taskResults := make(map[uint32]*taskResult, len(tg.tasks))
	for result := range results {
		taskResults[result.fNO] = result
	}
	return taskResults
}

func (tg *TaskGroup) worker(tasks chan *Task, results chan *taskResult) {
	for task := range tasks {
		result, err := task.f() // 每个任务逐一被执行
		results <- &taskResult{task.fNO, result, err}
	}
}

type taskResult struct {
	fNO    uint32
	result interface{}
	err    error
}

// Result 获取任务执行结果
func (tr *taskResult) Result() interface{} {
	return tr.result
}

// Error 获取任务执行状态
func (tr *taskResult) Error() error {
	return tr.err
}
