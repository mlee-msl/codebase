package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"time"

	errs "github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// https://mp.weixin.qq.com/s/aC5BZWuO7bRJdc_xUmMTTw
func ScopeCheck() { // 作用域验证
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Int()
	value := (randomNum % 4) + 1
	if a := 1; value == a {
		fmt.Println(a, randomNum)
	} else if b := 2; value == 1 {
		fmt.Println(a, b, randomNum)
	} else if c := 3; value == c {
		fmt.Println(a, b, c, randomNum)
	} else {
		d := 4
		fmt.Println(a, b, c, d, randomNum)
	}
}

func testMap() {
	a := map[int]int{1: 30, 2: 20, 3: 10}
	b := make(map[int]int, len(a))
	for k, v := range a {
		if k == 1 {
			a[k] = 300
			a[2] = 200
			a[30] = 2000
		}
		b[k] = v
	}
	fmt.Println(a, b)
	// map[1:300 2:200 3:10 30:2000] map[1:30 2:200 3:10 30:2000]
	// map[1:300 2:200 3:10 30:2000] map[1:30 2:200 3:10]
	// 以上两种情况都有可能
}

func testSlice1() {
	var a = []int{1, 2, 3, 4, 5}
	r := make([]int, len(a))

	fmt.Println("original a =", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
			a = append(a, 1000)
		}
		r[i] = v
	}

	fmt.Println("after for range loop, r =", r)
	fmt.Println("after for range loop, a =", a)
}

func testSlice2() {
	var a = []int{1, 2, 3, 4, 5}
	r := make([]int, len(a))

	fmt.Println("original a =", a)

	for i := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
			a = append(a, 1000)
		}
		r[i] = a[i]
	}

	fmt.Println("after for range loop, r =", r)
	fmt.Println("after for range loop, a =", a)
}

var _ = func() {
	var _ context.Context

	var _ sync.Cond
	var _ sync.WaitGroup
	var _ sync.Once
	var _ sync.Locker
	var _ sync.Map
	var _ sync.Mutex
	var _ sync.Pool
	var _ sync.RWMutex
	_ = time.After(5 * time.Second)
}

type MyStruct struct {
	noCopy noCopy
	// 其他字段...
}

// https://www.jianshu.com/p/ed7c0b028695
func testNoCopy() {
	// 创建一个 MyStruct 实例
	ms1 := MyStruct{}

	// 尝试复制 MyStruct 会导致编译错误
	// 因为复制会复制 noCopy，而 noCopy 有一个私有的 lock() 方法
	ms2 := ms1 // 这会导致编译错误

	// 如果我们注释掉 noCopy 的嵌入，那么上面的复制将工作正常
	// 但是，通常我们不想复制这样的结构体，因为它们可能包含不应该被复制的资源或状态

	fmt.Println("MyStruct created", ms1)
	fmt.Println("MyStruct created", ms1, ms2)
}

// noCopy may be added to structs which must not be copied
// after the first use.
//
// See https://golang.org/issues/8005#issuecomment-190753527
// for details.
//
// Note that it must not be embedded, due to the Lock and Unlock methods.
type noCopy struct{}

// Lock is a no-op used by -copylocks checker from `go vet`.
func (*noCopy) Lock()   {}
func (*noCopy) Unlock() {}

type A1 interface {
	f1()
}

type B1 struct {
	b int
}

func (b *B1) f1() {
	return
}

// 编译期校验结构体是否实现了特定接口
var _ A1 = (*B1)(nil)

func ctxWorker1(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopped, cause %v\n", id, ctx.Err())
			return
		default:
			fmt.Printf("Worker %d working...\n", id)
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func ctxWorker2(ctx context.Context, id int) {
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("mleeeeeeeeeee")
	// 		break
	// 	}
	// }

	var once sync.Once
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopped, cause %v\n", id, ctx.Err())
			return
		default:
			once.Do(func() {
				fmt.Printf("Worker %d working...\n", id)
				// time.Sleep(100 * time.Millisecond)
				time.Sleep(2 * time.Second)
			})
		}
	}
}

func ctxWorker3(ctx context.Context, id int) {
	// for {
	// 	select {
	// 	case <-ctx.Done():
	// 		fmt.Println("mleeeeeeeeeee")
	// 		break
	// 	}
	// }

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Worker %d stopped, cause by %v\n", id, context.Cause(ctx))
			return
		default:
			fmt.Printf("Worker %d working...\n", id)
			time.Sleep(100 * time.Millisecond)
			// time.Sleep(2 * time.Second)
		}
	}
}

func testCtx1() {
	ctx, cancel := context.WithCancel(context.Background())

	go ctxWorker1(ctx, 1)
	go ctxWorker1(ctx, 2)

	time.Sleep(1 * time.Second)
	cancel() // 发送取消信号

	// 等待一段时间以确保 worker 接收到取消信号
	time.Sleep(500 * time.Millisecond)
}

func testCtx2() {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 20; i++ {
		go ctxWorker2(ctx, i)
	}

	time.Sleep(1 * time.Second)
	cancel() // 发送取消信号

	// 等待一段时间以确保 worker 接收到取消信号
	time.Sleep(500 * time.Millisecond)
}

func testCtx3() {
	ctx, cancel := context.WithCancelCause(context.Background())
	// ctx, cancel := context.WithCancel(context.Background())
	for i := 1; i <= 4; i++ {
		go ctxWorker3(ctx, i)
	}

	time.Sleep(1 * time.Second)
	// cancel() // 发送取消信号
	cancel(errors.New("mlee cancel")) // 发送取消信号
	// cancel(errors.New("mlee cancel")) // 发送取消信号

	// 等待一段时间以确保 worker 接收到取消信号
	fmt.Println("....................")
	time.Sleep(500 * time.Second)
	fmt.Println("end")
}

// 1. 是否可以使用任务池
// 2. 任务太多如何处理？
// 3. 使用sync.Cond进行关键任务失败尽早取消
// var _ sync.Pool

// https://blog.csdn.net/weiguang102/article/details/131008608
// 查看下errgroup的原理，确认也是否是出现一个错误全部goroutine都立即退出，还是？

// errgroup 是添加一个任务随即执行一个任务，这对于不需要收集执行结果的场景比较合适，对于需要收集每个任务的结果的情况，需要首先将所有任务收集起来，才能知道收集结果的变量的空间大小

func testCtx() {
	const shortDuration = 1 * time.Millisecond

	// ctx, _ := context.WithTimeout(context.Background(), shortDuration)
	ctx, cancel := context.WithTimeoutCause(context.Background(), shortDuration, errors.New("occurs err"))
	// ctx, cancel := context.WithCancelCause(context.Background())
	// cancel(errors.New("occurs err"))
	defer cancel() // 释放资源，避免context(关联的channel等)或和context关联的goroutine泄露

	select {
	case <-time.After(1 * time.Second):
		fmt.Println("overslept")
	case <-ctx.Done():
		t, f := ctx.Deadline()
		fmt.Println("aa", ctx.Err(), t, f, ctx.Value("ml"), context.Cause(ctx))
	}
}

func textCtxCancel1() {
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(errors.New("err1"))
	fmt.Println(context.Cause(ctx))
	cancel(errors.New("err2"))
	fmt.Println(context.Cause(ctx))
}

func textCtxCancel2() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errors.New("err0"))
	fmt.Println(context.Cause(ctx))
}

func testErrGroup() {
	// 创建一个带有上下文的 errgroup.Group
	_, _ = errgroup.WithContext(context.Background())
	group := WithContext(context.Background())

	// 启动第一个任务
	group.Go(func() error {
		// 模拟一个耗时的任务
		time.Sleep(2 * time.Second)
		fmt.Println("Task 1 completed")
		return nil
	})

	// 启动第二个任务
	group.Go(func() error {
		// 模拟一个有错误的任务
		time.Sleep(1 * time.Second)
		fmt.Println("Task 2 completed with error")
		return fmt.Errorf("Task 2 encountered an error")
		// return nil
	})

	// 启动第三个任务
	group.Go(func() error {
		// 模拟一个耗时的任务
		time.Sleep(3 * time.Second)
		fmt.Println("Task 3 completed")
		return nil
	})

	// 等待所有任务完成或发生错误
	err := group.Wait()

	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}
}

func testOnce() {
	var once sync.Once
	once.Do(func() {
		fmt.Println("func 1!")
	})
	once.Do(func() {
		fmt.Println("func 1!")
	})

	once.Do(func() {
		fmt.Println("func 2!")
	})
}

func testSortBool(bools []bool) {
	sort.Slice(bools, func(i, j int) bool {
		return bools[i] && !bools[j]
	})
	fmt.Println(bools)
}

func testError() {
	fmt.Println(errC())
	fmt.Println(errD())
	fmt.Println(errs.Cause(errC()), errs.Cause(errD()))
}

func errA() error {
	return errors.New("this is a errA")
}

func errB() error {
	err := errA()
	return errs.Wrap(err, "this is a errB")
}

func errC() error {
	err := errB()
	return errs.WithMessage(err, "this is a errC")
}

func errD() error {
	return errs.WithStack(errC())
}

func testMapSliCap() {
	m := make(map[uint]byte, 10)
	var m1 map[uint]byte
	sli := make([]byte, 0, 10)
	var sli1 []byte
	fmt.Println(len(m), len(m1))
	fmt.Println(len(sli), cap(sli), len(sli1), cap(sli1))
}

// recover必须在defer函数中，且不能再包装其他函数了，比如testPanic1，testPanic2都不能recover住
// panic后不能在其他协程中recover住
func testPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
			var buf [4096]byte
			n := runtime.Stack(buf[:], false)
			fmt.Printf("Stack trace:\n%s\n", buf[:n])
		}
	}()
	fmt.Println("mleeeee")
	panic(123)
}

func testPanic1() {
	defer func() {
		rec()
	}()
	fmt.Println("mleeeee")
	panic(123)
}

func rec() {
	if r := recover(); r != nil {
		fmt.Println(r)
	}
}

func testPanic2() {
	defer func() {
		func() {
			if r := recover(); r != nil {
				fmt.Println(r)
			}
		}()
	}()
	fmt.Println("mleeeee")
	panic(123)
}

type AIn interface {
	f1()
	f2()
}

type StA struct {
	a string
}

func (s StA)f1() {
	fmt.Println(s.a)
}

func (s StA)f2() {
	fmt.Println(s.a+"mleeee")
}

type StB  struct {
	StA // 通过嵌入结构体实现了继承(实际上，是一种组合)
}

func (s StB) f1() {
	fmt.Println(s.a+"msl")
}

func testInheritance(){
	s := StB{StA: StA{"111111"}}
	testInterface(s)
}

func testInterface(s AIn) { // 通过接口实现多态
	s.f1()
	s.f2()
}


