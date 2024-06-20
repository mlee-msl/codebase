package main

import (
	"context"
	"errors"
	"fmt"
	"math/bits"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/axgle/mahonia"
	errs "github.com/pkg/errors"
	"github.com/shopspring/decimal"
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

	var m map[uint8]map[uint32]string
	fmt.Println("取nil map: ", m[1][12])

	var m1 map[uint8]map[uint32]map[byte]string
	fmt.Println("取nil map: ", m1[1][12][122])
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

type DoNotCopy [0]sync.Mutex // 也实现sync.Locker接口
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
	ms2 := ms1 // 使用`go vet`可检查出错误

	// 如果我们注释掉 noCopy 的嵌入，那么上面的复制将工作正常
	// 但是，通常我们不想复制这样的结构体，因为它们可能包含不应该被复制的资源或状态,比如复制了sync.WaitGroup对象后，那么对象的Add和Done这一对操作就不匹配了，就会导致死锁问题
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

// 实现了sync.Locker接口
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

// func testCtx4(){
// 	ctx := context.Background()
// 	ctx1 := testCtx4_(&ctx)
// 	fmt.Println("222", ctx.Value("key"))
// 	fmt.Println("333", ctx1.Value("key"))
// }
// func testCtx4_(ctx *context.Context) context.Context{
// 	val := (*ctx).Value("key")
//     if val == nil {
//         // 如果上下文值为空，则创建一个新的上下文值
//         val = "default"
//         ctx = context.WithValue(*ctx, "key", val)
//     } else {
//         // 如果上下文值不为空，则修改其值
//         val = "new value"
//         ctx = context.WithValue(ctx, "key", val)
//     }
//     // 在函数中使用修改后的上下文值
//     fmt.Println("111", ctx.Value("key"))
// 	return ctx
// }

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

// 两阶段延迟执行
func testDefer() {
	setup := func() {
		fmt.Println("setup")
	}
	f := func() func() {
		setup()
		return func() {
			fmt.Println("tear down")
		}
	}
	defer f()()
	fmt.Println("test defer")
}

func testDefer1() {
	a := testDefer1_()
	fmt.Println(a)
}

func testDefer1_() (a int) { // 这里需要声明一个具名的返回值变量，才能修改，对于匿名的临时变量无法被修改
	a = 10
	defer func() func() {
		fmt.Println("set up", a)
		fmt.Println("set up")
		a = 20
		return func() {
			fmt.Println("tear down", a)
			a = 30
			fmt.Println("tear down")
		}
	}()() // 注意一个括号就是调用一次，defer 仅会将最外层的函数推进调用栈中

	fmt.Println("testDefer1", a)
	return a + 100
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

func (s StA) f1() {
	fmt.Println(s.a)
}

func (s StA) f2() {
	fmt.Println(s.a + "mleeee")
}

type StB struct {
	StA // 通过嵌入结构体实现了继承(实际上，是一种组合)
}

func (s StB) f1() {
	fmt.Println(s.a + "msl")
}

func testInheritance() {
	s := StB{StA: StA{"111111"}}
	testInterface(s)
}

func testInterface(s AIn) { // 通过接口实现多态
	s.f1()
	s.f2()
}

func testCancelCtx1() {
	fmt.Println(testCancelCtx())
}
func testCancelCtx() error {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil)
	// cancel(errors.New("mlee1111")) // 仅获取到第一次ctx cancel的错误。后面再次取消不会有任何动作，直接返回
	// cancel(errors.New("mlee2222"))
	// cancel(nil)
	// fmt.Println(context.Cause(ctx), ctx.Err())
	// fmt.Println(context.Cause(ctx), ctx.Err())
	// cancel(nil)
	// fmt.Println(<-ctx.Done(), context.Cause(ctx), ctx.Err())
	// fmt.Println(<-ctx.Done(), context.Cause(ctx), ctx.Err())
	// cancel(errors.New("mlee"))
	return context.Cause(ctx)
}

func checkChannelClosed() {
	var m1 map[string]int
	if m1 == nil {
		fmt.Println("mleeeeeeeeeeee")
	}
	m := make(map[string]int, 10)
	m["mlee"] = 1
	fmt.Println(len(m))
	ch := make(chan int, 10)
	ch <- 1
	fmt.Println(len(ch), cap(ch))
	close(ch)
	fmt.Println(len(ch), cap(ch))
	// close(ch)
	// select {
	// case ch <- 1:
	//     fmt.Println("Sent successfully")
	// default:
	//     fmt.Println("Channel is closed")
	// }
	// time.Sleep(2*time.Second)
	// ch <- 1
}

func testMapSort() {
	m := map[int][]int{
		1: []int{1, 2, -1, -10, 3},
		4: []int{10, 2, 1, -20, 3},
		2: []int{100, 20000, 3, -10, -100},
	}
	for _, v := range m {
		sort.Slice(v, func(i, j int) bool {
			return v[i] > v[j]
		})
	}
	fmt.Println(m)
}

func testAtomic() {
	var v atomic.Value
	v.Store(1)
}

type complexData1 map[int]string
type complexData2 map[string]complexData1

func testComplexData() {
	c := make(complexData2, 12)
	c["mlee1"] = complexData1{10: "ok", 11: "bad"}
	c["mlee2"] = complexData1{1022: "ok22", 1122: "bad222"}
	fmt.Println(c)
}

func testContainer() {
	// list.List
	// heap.Fix
	// ring.New(3)
}

func testNilMapOrSlice() {
	var m map[int]int
	fmt.Println(m[11])
	m1 := map[int]int(nil)
	val, ok := m1[11]
	fmt.Println(val, ok)

	// var sli []int
	// fmt.Println(sli[0]) // panic
	sli1 := []int(nil)
	fmt.Println(sli1[0]) // panic
}

func testCap() {
	var b bool
	var b1 byte
	var i int
	fmt.Println(unsafe.Sizeof(b), unsafe.Sizeof(b1), unsafe.Sizeof(i)) // 1 1 8
}

// 实现下bitset
func testBits() {
	x := uint(0b1010101010) // 最高位1的位数
	len := bits.Len(x)
	fmt.Printf("The length of %b is %d\n", x, len)

	x1 := uint(0b0100000111)
	len1 := bits.Len(x1)
	fmt.Printf("The length of %b is %d\n", x1, len1)
}

func testStringSplit() {
	const s = "12_"
	ss := strings.Split(s, "_")
	fmt.Println(len(ss), ss)
	fmt.Println(strconv.Atoi(""))
}

func testDecimal() {
	v := decimal.New(2, 1e7)
	v1 := decimal.New(2, 7)
	fmt.Println(v, v1)
}

func testLangTrans() {
	enc := mahonia.NewEncoder("zh-tw")
	//converts a  string from UTF-8 to gbk encoding.
	fmt.Println(enc.ConvertString("hello,世界"))
}

func Test_if() {
	var a int = 10
	a0 := _if(true, (*int)(nil), &a).(*int)
	fmt.Println(a0)
	// a1 := _if(true, nil, &a).(*int) // panic
	// fmt.Println(a1)
	a2 := _if(false, nil, &a).(*int)
	fmt.Println(*a2)
}

var _if = func(cond bool, a, b interface{}) interface{} {
	if cond {
		return a
	}
	return b
}

func If(cond bool, f1, f2 func() error) error {
	f := func(f func() error) error {
		if f == nil {
			return nil
		}
		return f()
	}
	if cond {
		return f(f1)
	}
	return f(f2)
}

func TestAppend() {
	var a, b []int
	if a == nil {
		fmt.Println("aaa")
	}
	c := append(a, b...)
	c = append(c, nil...)
	fmt.Println(len(c), c)
}

func TestChannel() {
	results := make(chan int, 10)
	// go func(){
	// 	results <- 12
	// 	close(results)
	// }()

	for result := range results {
		fmt.Println(result)
	}
}

func testInt() {
	a := 2_040_051_011 // 使用下划线分隔，便于查看
	fmt.Println(a, a+301)
	fmt.Printf("%[0]s, %[1]s, %[0]s, %d", 12, 13, 14, 15)
}

func testEmptyArray() {
	var a [0]byte
	m := map[int][0]byte{1: [0]byte{}}
	m1 := map[int][0]struct{}{1: [0]struct{}{}}
	fmt.Println(a, m, m1)
}

// https://blog.csdn.net/u013272009/article/details/135876694
func TestSkills() {
	// https://mytechshares.com/2021/12/14/gopher-should-know-struct-ops/
	type NoUnkeyedLiterals struct{}
	type ab struct {
		aa int
	}
	type a struct {
		// _ NoUnkeyedLiterals // 忽略该字段，就算赋值了也不会有值，这样要求声明 a 结构体就需要显示指定每个字段的值了，采用默认的顺序赋值的方式可能有隐患，比如,将结构体中同类型字段顺序修改，就会有问题（编译发现不了）
		A int
		B string
		C string
		ab
	}
	a1 := a{A: 1, B: "12", ab: ab{}}
	// a2 := a{1, "12", ab{}} // too few values in struct literal of type a
	fmt.Println(a1)
}

// Processor 定义处理函数类型
type Processor func(in chan interface{}) chan interface{}

// CreatePipeline 创建并启动处理流水线
func CreatePipeline(procs ...Processor) (inChan chan interface{}, outChan chan interface{}) {
	// 创建一个初始的输入通道和最终的输出通道
	inChan = make(chan interface{})
	outChan = make(chan interface{})

	// 使用WaitGroup等待所有处理器启动
	var wg sync.WaitGroup
	wg.Add(len(procs))

	// 逐个串联处理器
	var currentChan chan interface{} = inChan
	for _, proc := range procs {
		nextChan := make(chan interface{}, 1) // 创建下一个处理器的输入通道

		go func(proc Processor, inputChan chan interface{}, outputChan chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done() // 通知WaitGroup当前处理器已启动
			for data := range inputChan {
				proc(inputChan)    // 执行当前处理器逻辑，可能不会用到传入的data，具体取决于处理器实现
				outputChan <- data // 将数据传递到下一个处理器
			}
			close(outputChan) // 当前处理器处理完数据后，关闭输出通道
		}(proc, currentChan, nextChan, &wg)

		currentChan = nextChan // 更新当前通道为下一个处理器的输入通道
	}

	// 等待所有处理器启动
	wg.Wait()

	// 返回输入通道和最终输出通道
	return inChan, outChan
}

// 示例处理器：打印接收到的数据，并原样发送到输出通道
func printProcessor() Processor {
	return func(in chan interface{}) chan interface{} {
		out := make(chan interface{}, 1)
		go func() {
			for data := range in {
				fmt.Printf("Processing: %v\n", data)
				out <- data
			}
			close(out)
		}()
		return out
	}
}

// https://juejin.cn/post/7202153645441318967
func TestPipeline() {
	inChan, outChan := CreatePipeline(printProcessor(), printProcessor(), printProcessor())

	// 向流水线发送数据
	inChan <- "Hello, Pipeline!"
	inChan <- "Another message."
	close(inChan) // 关闭输入通道，表示没有更多的数据了

	// 读取并打印最终结果
	for result := range outChan {
		fmt.Printf("Result: %v\n", result)
	}
}
type iface interface{
	f()
}
type st struct {
	a int
}

func(st) f(){}
func isNil() {


	m := map[int]st{}

	f1 := func(a int)(iface, bool){
		v1, ok := m[a]
		return v1, ok
	}
	v, ok := f1(1)
	fmt.Println(v, ok)
	fmt.Println(IsNilPtr(v))

	// fmt.Println(reflect.ValueOf(v).IsNil())
}

// IsNilPtr 判断接口是否为空指针
func IsNilPtr(i any) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}

// IsNilPtr 判断接口是否为空指针(nil说明已经是引用类型或者指针类型了，比如struct类型就会panic)
func _IsNilPtr(i any) bool {
	if i == nil {
		return true
	}

	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Pointer, reflect.UnsafePointer:
		return v.IsNil()
	default:
		return false
	}
}
