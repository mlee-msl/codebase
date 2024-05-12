package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
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

	select {
	case <-ctx.Done():
		fmt.Printf("Worker %d stopped, cause %v\n", id, ctx.Err())
		return
	default:
		fmt.Printf("Worker %d working...\n", id)
		time.Sleep(100 * time.Millisecond)
		// time.Sleep(2 * time.Second)
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
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= 20; i++ {
		go ctxWorker3(ctx, i)
	}

	time.Sleep(1 * time.Second)
	cancel() // 发送取消信号
	cancel() // 发送取消信号

	// 等待一段时间以确保 worker 接收到取消信号
	time.Sleep(500 * time.Millisecond)
	fmt.Println("end")
}

// 1. 是否可以使用任务池
// 2. 任务太多如何处理？
// 3. 使用sync.Cond进行关键任务失败尽早取消
// var _ sync.Pool


// https://blog.csdn.net/weiguang102/article/details/131008608
// 查看下errgroup的原理，确认也是否是出现一个错误全部goroutine都立即退出，还是？
