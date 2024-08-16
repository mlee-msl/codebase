package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
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
	"unicode/utf8"
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
	fmt.Println(time.Date(2012+100, time.December, 1, 1, 1, 1, 1, time.Local).Unix())
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
	fmt.Printf("取nil map: <%v>\n", m1[1][12][122])
}

func TestBoolean() {
	a := 12
	if !(a <= 9 && 0 <= a || a == 12 || a == 100 || a == -100) {
		fmt.Println("bad1")
		return
	}
	fmt.Println("good1")

	if !((a <= 9 && 0 <= a) || a == 12 || a == 100 || a == -100) {
		fmt.Println("bad2")
		return
	}
	fmt.Println("good2")
}

type iface1 interface {
	free()
}

func testIface() {
	f := iface1(nil)
	f.free()
}
func testStr() {
	name := "mlee 富途"
	fmt.Println(len(name)) // 11个字节
	for i := 0; i < len(name); i++ {
		fmt.Printf("%v, %T, %T\n", string(name[i]), name[i], 'm')
		// if name[i] == 'e' {
		// 	fmt.Println("aaaa: ", name[i])
		// }
	}
	fmt.Println("-------------")
	for _, c := range name {
		fmt.Println(c, string(c))
	}
	s := "世界你好"
	r, size := utf8.DecodeRuneInString(s)
	fmt.Println(string(r), size)
}

func TestSli() {
	tags := map[string]int{"mlee1": 12, "mlee2": 1222}
	sli := make([]int, len(tags))
	for _, tag := range tags {
		sli = append(sli, tag)
	}
	fmt.Println(len(sli), cap(sli), sli)
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

// slice切片的底层结构。
// note: a[start:end],其中start>=0, end <= cap(a), 不指定就是len(a)
// type slice struct {
//     array unsafe.Pointer
//     len   int
//     cap   int
// }

func testSlice3() {
	f1 := func(a []int) {
		a = append(a, 1)
		a = append(a, 1)
		a = append(a, 1)
		a = append(a, 1)
	}
	_ = func(a []int) {
		a = append(a, 1)
		a = append(a, 10)
	}
	a := make([]int, 0, 2)
	a = append(a, 110)
	fmt.Println(a)
	f1(a)
	fmt.Println(a, len(a))
	fmt.Println(a[:cap(a)], a[:])
	fmt.Print()
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
	var _ strings.Builder
	var _ time.Timer
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

func testSortTime() {
	var (
		now = time.Now()
		ts  = []time.Time{now.AddDate(0, 0, 10), now.AddDate(0, 0, 8), now.AddDate(0, 0, 2), now.AddDate(0, 0, 5)}
	)
	sort.Slice(ts, func(i int, j int) bool {
		return ts[i].After(ts[j])
	})
	fmt.Println(ts)
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
	defer func() func() {
		setup()
		return func() {
			fmt.Println("tear down")
		}
	}()()
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
	var m2 map[int]map[string][0]func()
	val1, ok := m2[1]["mlee"]
	fmt.Println(val1, ok)

	// var sli []int
	// fmt.Println(sli[0]) // panic
	// sli1 := []int(nil)
	// fmt.Println(sli1[0]) // panic
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
	// fmt.Printf("%[0]s, %[1]s, %[0]s, %d", 12, 13, 14, 15)
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

type iface interface {
	f()
}
type st struct {
	a int
}

func (st) f() {}
func isNil() {

	m := map[int]st{}

	f1 := func(a int) (iface, bool) {
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

func testTypeSize() {
	var a int
	fmt.Println("size int", unsafe.Sizeof(a))
}

func testTicker() {
	_ = time.NewTimer(1 * time.Second).C // time.After(1*time.Second)
	// t := time.NewTicker(100 * time.Millisecond)
	t := time.NewTicker(1 * time.Second)
	defer func() {
		t.Stop()
	}()
	for {
		select {
		case t_ := <-time.After(3 * time.Second):
			fmt.Println("-----end-----", t_)
			return
			// case t_ := <-t.C: // ? 为啥不能命中第一个case
		// fmt.Println(t_)
		default: // ? 为啥不能命中第一个case
			// 	time.Sleep(1 * time.Second)
			// 	fmt.Println("default")
		}
	}
}

func testIOTA() {
	// iota的值每一行自动加1,后面的行计算方式与前面的行完全一样
	const (
		_ = iota // iota的值每一行自动加1

		_
		a
		b
		c = iota
		d
		e = 10
		f1
		f2
		f3
		g = iota + 1 // 后面都是iota+1，只不过iota每一行自动加1
		h            // iota+1
	)
	const (
		v0 = 1000
		v1 = 1024 // 后面的自动和v1保持一样，除非另外定义
		v2
		v3
		v4 = 100 // 后面的自动和v4保持一样，除非另外定义
		v5
		v6
	)
	const (
		s1 = iota
		s2 = iota
		_

		_  = iota
		s3 = iota
		_  = iota + 1 // 5+1
		s4            // 6+1
	)
	const (
		t1 = 1 << iota
		t2
		t3
		t4
		t5
		t6
		t0  = t6 - iota
		t10 = t0 + iota
		t11
		t12
	)
	const (
		One      = 1 << iota         // 1
		Two                          // 2
		Four                         // 4
		Eight                        // 8
		Nine     = One | Eight       // 9
		Ten      = Two | Eight       // 10
		Eleven   = One | Two | Eight // 11
		Twelve   = 1 << iota         // 12
		Thirteen                     // 13
		Fourteen                     // 14
		Fifteen                      // 15
		Sixteen                      // 16
	)

	type FinancialDataType uint32

	const (
		_ FinancialDataType = iota // 零值占位，不采用
		// 以下为可按需获取的数据类型：采用位掩码的方式初始化
		BalanceSheetFinancialData = 1 << (iota - 1) // 资产负债表数据
		BuybackFinancialData                        // 回购数据
		CashFlowFinancialData                       // 现金流量表数据
		IncomeFinancialData                         // 利润表数据
		MainIndexFinancialData                      // 关键指标数据
		placeholder
		AllOnDemandFinancialData = placeholder - 1 // 表示所有的关键数据类型
		_
		resetFinancialDataType = AllOnDemandFinancialData - iota
		// 以下为非可按需获取的数据类型采用线性加1的方式初始化
		OperationFinancialData = resetFinancialDataType + iota // 财报的运营数据
		CrawlerFinancialData                                   // 财报的爬虫数据
		BrightRiskData                                         // 亮点和风险点数据
		FinancialExtraData                                     // 财报相关的其他数据

		FinancialDataTypeTotalNums = iota
	)
	fmt.Println(a, b, c, d, e, f1, f2, f3, g, h)
	fmt.Println(v0, v1, v2, v3, v4, v5, v6)
	fmt.Println(s1, s2, s3, s4)
	fmt.Println(t1, t2, t3, t4, t5, t6, t10, t11, t12)
	fmt.Println(One, Two, Four, Eight, Nine, Ten, Eleven, Twelve, Thirteen, Fourteen, Fifteen, Sixteen)
	fmt.Println(BalanceSheetFinancialData, BuybackFinancialData, CashFlowFinancialData, IncomeFinancialData, MainIndexFinancialData, "-----", AllOnDemandFinancialData, "---",
		OperationFinancialData, CrawlerFinancialData, BrightRiskData, FinancialExtraData, FinancialDataTypeTotalNums)
	aaa := MainIndexFinancialData | IncomeFinancialData
	bbb := aaa | AllOnDemandFinancialData
	fmt.Println(aaa&BalanceSheetFinancialData > 0, aaa&BuybackFinancialData > 0, aaa&CashFlowFinancialData > 0, aaa&IncomeFinancialData > 0, aaa&MainIndexFinancialData > 0, AllOnDemandFinancialData, "-----",
		bbb&BalanceSheetFinancialData > 0, bbb&BuybackFinancialData > 0, bbb&CashFlowFinancialData > 0, bbb&IncomeFinancialData > 0, bbb&MainIndexFinancialData > 0)

	fmt.Println("----", 0&AllOnDemandFinancialData == 0, BalanceSheetFinancialData&AllOnDemandFinancialData == 0, (BalanceSheetFinancialData|MainIndexFinancialData)&AllOnDemandFinancialData == 0)
}

func TestObjSize() {
	s := f10HandlerHK{}
	fmt.Println(unsafe.Sizeof(&s))

	s1 := new(int)
	fmt.Println(unsafe.Sizeof(s1))

	s2 := new(string)
	fmt.Println(unsafe.Sizeof(s2))
}

const gNums = 200

func TestSyncPool() {
	p := sync.Pool{
		New: func() any {
			return f10HandlerHK{}
		},
	}

	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)
	for i := 1; i <= gNums; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Millisecond)
			defer func() {
				wg.Done()
			}()
			a_ := p.Get().(f10HandlerHK)
			defer func() {
				// fmt.Println("111", a_.s)
				p.Put(a_)
			}()

			// if a_.s == "mlee" {
			_ = counter.Add(1)
			// 	fmt.Printf("<%+v>\n", a_.s)
			// }
			a_.s = "mlee"
			// fmt.Println(a_.s)
		}()
	}
	wg.Wait()
	// fmt.Println("end", counter.Load())
}

func TestSyncPool2() {
	p := sync.Pool{
		New: func() any {
			return new(f10HandlerHK)
		},
	}

	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)
	for i := 1; i <= gNums; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Millisecond)
			defer func() {
				wg.Done()
			}()
			a_ := p.Get().(*f10HandlerHK)
			defer func() {
				// fmt.Println("111", a_.s)
				a_.f10FinanceIndicatorHK = f10FinanceIndicatorHK{}
				a_.renderDependency = renderDependency{}
				a_.s = ""
				p.Put(a_)
			}()

			// if a_.s == "mlee" {
			_ = counter.Add(1)
			// 	fmt.Printf("<%+v>\n", a_.s)
			// }
			a_.s = "mlee"
			// fmt.Println(a_.s)
		}()
	}
	wg.Wait()
	// fmt.Println("end", counter.Load())
}

func TestSyncPool3() {
	// p := sync.Pool{
	// 	New: func() any {
	// 		return new(f10HandlerHK)
	// 	},
	// }

	p := sync.Pool{
		New: func() any {
			return &f10HandlerHK{
				f10FinanceIndicatorHK: f10FinanceIndicatorHK{
					mainIndexHK:    mainIndexHK{mainIndexDetails: make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)},
					balanceSheetHK: balanceSheetHK{balanceSheetDetails: make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)},
				},
			}
		},
	}

	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)
	for i := 1; i <= gNums; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Millisecond)
			defer func() {
				wg.Done()
			}()
			a_ := p.Get().(*f10HandlerHK)
			defer func() {
				// fmt.Println("111", a_.s)
				// a_.f10FinanceIndicatorHK = f10FinanceIndicatorHK{}
				// a_.renderDependency = renderDependency{}
				a_.mainIndexDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
				a_.balanceSheetDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
				a_.s = ""
				p.Put(a_)
			}()
			a_.mainIndexDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.mainIndexDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["3"] = new(HK_BalanceSheetGEHK)
			// if a_.s == "mlee" {
			_ = counter.Add(1)
			// 	fmt.Printf("<%+v>\n", a_.s)
			// }
			a_.s = "mlee"
			// fmt.Println(a_.s)
		}()
	}
	wg.Wait()
	// fmt.Println("end", counter.Load())
}

func TestSyncPool4() {
	p := sync.Pool{
		New: func() any {
			return new(f10HandlerHK)
		},
	}

	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)
	for i := 1; i <= gNums; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(3)+1) * time.Millisecond)
			defer func() {
				wg.Done()
			}()
			a_ := p.Get().(*f10HandlerHK)
			defer func() {
				// fmt.Println("111", a_.s)
				a_.f10FinanceIndicatorHK = f10FinanceIndicatorHK{}
				a_.renderDependency = renderDependency{}
				a_.s = ""
				p.Put(a_)
			}()
			a_.mainIndexDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
			a_.balanceSheetDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
			a_.mainIndexDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.mainIndexDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["3"] = new(HK_BalanceSheetGEHK)
			// if a_.s == "mlee" {
			_ = counter.Add(1)
			// 	fmt.Printf("<%+v>\n", a_.s)
			// }
			a_.s = "mlee"
			// fmt.Println(a_.s)
		}()
	}
	wg.Wait()
	// fmt.Println("end", counter.Load())
	fmt.Print()
}

func TestNonPool() {
	var (
		counter atomic.Int32
		wg      sync.WaitGroup
	)
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Intn(4)+1) * time.Millisecond)
			defer func() {
				wg.Done()
			}()
			a_ := new(f10HandlerHK)
			a_.mainIndexDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
			a_.balanceSheetDetails = make(map[FinancialUniqueKey]*HK_BalanceSheetGEHK, 20)
			a_.mainIndexDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.mainIndexDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["1"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["2"] = new(HK_BalanceSheetGEHK)
			a_.balanceSheetDetails["3"] = new(HK_BalanceSheetGEHK)
			_ = counter.Add(1)
			a_.s = "mlee"
			// fmt.Println(a_.s)
		}()
	}
	wg.Wait()
	// fmt.Println("end", counter.Load())
}

func init() {
	rand.Seed(time.Now().Unix())
}

type f10HandlerHK struct {
	f10FinanceIndicatorHK

	renderDependency
	s string
}

type renderDependency struct {
	indicatorResults AllIndicatorResults

	currencyDetails map[int32]struct{ currencyUnit, currencyCode string }
	brightRiskData  BrightRiskDataDetail
	extraData       FinancialExtraDataDetail
}

type FinancialExtraDataDetail struct {
	PlateId     uint64  // 个股所属的行业板块
	MarketValue float64 // 个股市值（扩大1000， 单位已转为港币）
}

type BrightRiskDataDetail struct {
	_ NoUnkeyedLiterals

	FinancialType   uint8
	BrightIndicator map[int32]Score
	RiskIndicator   map[int32]Score
}

// Score 风险提示与投资亮点得分
type Score struct {
	_ NoUnkeyedLiterals

	TotalScore      int64 // 指标的总分数
	HorizontalScore uint8 // 指标的横向分数
	VerticalScore   uint8 // 指标的纵向分数
}

type f10FinanceIndicatorHK struct {
	mainIndexHK
	incomeHK
	balanceSheetHK
	cashFlowHK
	buybackHK
	operationAndCrawlerHK
}

type FinancialUniqueKey string
type mainIndexHK struct {
	mainIndexDetails map[FinancialUniqueKey]*HK_BalanceSheetGEHK

	assistCalculation
}

type incomeHK struct {
	incomeDetails map[FinancialUniqueKey]*HK_BalanceSheetGEHK

	assistCalculation
}

type assistCalculation struct {
	mergeKeysByFinancialType map[int32][]FinancialUniqueKey
	allKeys                  []FinancialUniqueKey
}

type HK_BalanceSheetGEHK struct {
	UniqueKey                 *string  `protobuf:"bytes,142,opt,name=unique_key,json=uniqueKey" json:"unique_key,omitempty"`
	F10FinancialYear          *uint32  `protobuf:"varint,143,opt,name=f10_financial_year,json=f10FinancialYear" json:"f10_financial_year,omitempty"`
	F10FinancialType          *uint32  `protobuf:"varint,144,opt,name=f10_financial_type,json=f10FinancialType" json:"f10_financial_type,omitempty"`
	Type                      *uint32  `protobuf:"varint,131,opt,name=type" json:"type,omitempty"`
	EndDate                   *uint32  `protobuf:"varint,132,opt,name=EndDate" json:"EndDate,omitempty"`
	InfoSource                *string  `protobuf:"bytes,133,opt,name=InfoSource" json:"InfoSource,omitempty"`
	CompanyNature             *uint32  `protobuf:"varint,134,opt,name=CompanyNature" json:"CompanyNature,omitempty"`
	CurrencyUnit              *string  `protobuf:"bytes,135,opt,name=CurrencyUnit" json:"CurrencyUnit,omitempty"`
	CurrencyCode              *string  `protobuf:"bytes,141,opt,name=CurrencyCode" json:"CurrencyCode,omitempty"`
	AccountingStandards       *string  `protobuf:"bytes,136,opt,name=AccountingStandards" json:"AccountingStandards,omitempty"`
	PeriodMark                *uint32  `protobuf:"varint,137,opt,name=PeriodMark" json:"PeriodMark,omitempty"`
	Mark                      *uint32  `protobuf:"varint,138,opt,name=Mark" json:"Mark,omitempty"`
	InfoPublDate              *uint32  `protobuf:"varint,139,opt,name=InfoPublDate" json:"InfoPublDate,omitempty"`
	FiscalYear                *uint32  `protobuf:"varint,140,opt,name=FiscalYear" json:"FiscalYear,omitempty"`
	Inventories               *float64 `protobuf:"fixed64,1,opt,name=Inventories" json:"Inventories,omitempty"`
	DeveAndForSalePro         *float64 `protobuf:"fixed64,2,opt,name=DeveAndForSalePro" json:"DeveAndForSalePro,omitempty"`
	AccountReceivables        *float64 `protobuf:"fixed64,3,opt,name=AccountReceivables" json:"AccountReceivables,omitempty"`
	BillReceivable            *float64 `protobuf:"fixed64,4,opt,name=BillReceivable" json:"BillReceivable,omitempty"`
	AssociateFundRece         *float64 `protobuf:"fixed64,5,opt,name=AssociateFundRece" json:"AssociateFundRece,omitempty"`
	CWrksCliMonReceCA         *float64 `protobuf:"fixed64,6,opt,name=CWrksCliMonReceCA" json:"CWrksCliMonReceCA,omitempty"`
	InterestReceivables       *float64 `protobuf:"fixed64,7,opt,name=InterestReceivables" json:"InterestReceivables,omitempty"`
	InsOtherReceCA            *float64 `protobuf:"fixed64,8,opt,name=InsOtherReceCA" json:"InsOtherReceCA,omitempty"`
	OtherAccounetrece         *float64 `protobuf:"fixed64,9,opt,name=OtherAccounetrece" json:"OtherAccounetrece,omitempty"`
	TaxReceivable             *float64 `protobuf:"fixed64,10,opt,name=TaxReceivable" json:"TaxReceivable,omitempty"`
	PrepaidRentCA             *float64 `protobuf:"fixed64,11,opt,name=PrepaidRentCA" json:"PrepaidRentCA,omitempty"`
	Cash                      *float64 `protobuf:"fixed64,12,opt,name=Cash" json:"Cash,omitempty"`
	ShortTermDeposit          *float64 `protobuf:"fixed64,13,opt,name=ShortTermDeposit" json:"ShortTermDeposit,omitempty"`
	FixedDepositCA            *float64 `protobuf:"fixed64,14,opt,name=FixedDepositCA" json:"FixedDepositCA,omitempty"`
	DepositInCentralBank      *float64 `protobuf:"fixed64,15,opt,name=DepositInCentralBank" json:"DepositInCentralBank,omitempty"`
	MortagageDeposit          *float64 `protobuf:"fixed64,16,opt,name=MortagageDeposit" json:"MortagageDeposit,omitempty"`
	AdvancesTCusts            *float64 `protobuf:"fixed64,17,opt,name=AdvancesTCusts" json:"AdvancesTCusts,omitempty"`
	LendCapital               *float64 `protobuf:"fixed64,18,opt,name=LendCapital" json:"LendCapital,omitempty"`
	ShortTermInvest           *float64 `protobuf:"fixed64,19,opt,name=ShortTermInvest" json:"ShortTermInvest,omitempty"`
	HForSaleAssetsCA          *float64 `protobuf:"fixed64,20,opt,name=HForSaleAssetsCA" json:"HForSaleAssetsCA,omitempty"`
	FinAetAtFValTPLCA         *float64 `protobuf:"fixed64,21,opt,name=FinAetAtFValTPLCA" json:"FinAetAtFValTPLCA,omitempty"`
	DerFinInstsCA             *float64 `protobuf:"fixed64,22,opt,name=DerFinInstsCA" json:"DerFinInstsCA,omitempty"`
	OtherCurrentAssets        *float64 `protobuf:"fixed64,23,opt,name=OtherCurrentAssets" json:"OtherCurrentAssets,omitempty"`
	CAExcepItems              *float64 `protobuf:"fixed64,24,opt,name=CAExcepItems" json:"CAExcepItems,omitempty"`
	CAAdjItems                *float64 `protobuf:"fixed64,25,opt,name=CAAdjItems" json:"CAAdjItems,omitempty"`
	TotalCurrentAssets        *float64 `protobuf:"fixed64,26,opt,name=TotalCurrentAssets" json:"TotalCurrentAssets,omitempty"`
	FixedAssets               *float64 `protobuf:"fixed64,27,opt,name=FixedAssets" json:"FixedAssets,omitempty"`
	WorkshopAndEquipment      *float64 `protobuf:"fixed64,28,opt,name=WorkshopAndEquipment" json:"WorkshopAndEquipment,omitempty"`
	InvestProperty            *float64 `protobuf:"fixed64,29,opt,name=InvestProperty" json:"InvestProperty,omitempty"`
	ConstruInProcess          *float64 `protobuf:"fixed64,30,opt,name=ConstruInProcess" json:"ConstruInProcess,omitempty"`
	LandUsufruct              *float64 `protobuf:"fixed64,31,opt,name=LandUsufruct" json:"LandUsufruct,omitempty"`
	AdvancePayment            *float64 `protobuf:"fixed64,32,opt,name=AdvancePayment" json:"AdvancePayment,omitempty"`
	PrepaidRentNCA            *float64 `protobuf:"fixed64,33,opt,name=PrepaidRentNCA" json:"PrepaidRentNCA,omitempty"`
	LongtermReceivableAccount *float64 `protobuf:"fixed64,34,opt,name=LongtermReceivableAccount" json:"LongtermReceivableAccount,omitempty"`
	CWrksCliMonReceNCA        *float64 `protobuf:"fixed64,35,opt,name=CWrksCliMonReceNCA" json:"CWrksCliMonReceNCA,omitempty"`
	InsOtherReceNCA           *float64 `protobuf:"fixed64,36,opt,name=InsOtherReceNCA" json:"InsOtherReceNCA,omitempty"`
	DevelopmentExpenditure    *float64 `protobuf:"fixed64,37,opt,name=DevelopmentExpenditure" json:"DevelopmentExpenditure,omitempty"`
	SubCompanyEquity          *float64 `protobuf:"fixed64,38,opt,name=SubCompanyEquity" json:"SubCompanyEquity,omitempty"`
	CoBusinessEquity          *float64 `protobuf:"fixed64,39,opt,name=CoBusinessEquity" json:"CoBusinessEquity,omitempty"`
	SuppCompEquity            *float64 `protobuf:"fixed64,40,opt,name=SuppCompEquity" json:"SuppCompEquity,omitempty"`
	JointVenturesEquity       *float64 `protobuf:"fixed64,41,opt,name=JointVenturesEquity" json:"JointVenturesEquity,omitempty"`
	FixedDepositNCA           *float64 `protobuf:"fixed64,42,opt,name=FixedDepositNCA" json:"FixedDepositNCA,omitempty"`
	MortagageDepositNCA       *float64 `protobuf:"fixed64,43,opt,name=MortagageDepositNCA" json:"MortagageDepositNCA,omitempty"`
	LTInvestments             *float64 `protobuf:"fixed64,44,opt,name=LTInvestments" json:"LTInvestments,omitempty"`
	SecuInvestment            *float64 `protobuf:"fixed64,45,opt,name=SecuInvestment" json:"SecuInvestment,omitempty"`
	FinAetAtFValTPLNCA        *float64 `protobuf:"fixed64,46,opt,name=FinAetAtFValTPLNCA" json:"FinAetAtFValTPLNCA,omitempty"`
	DerFinInstsNCA            *float64 `protobuf:"fixed64,47,opt,name=DerFinInstsNCA" json:"DerFinInstsNCA,omitempty"`
	HoldForSaleAssetsNCA      *float64 `protobuf:"fixed64,48,opt,name=HoldForSaleAssetsNCA" json:"HoldForSaleAssetsNCA,omitempty"`
	OtherInvestment           *float64 `protobuf:"fixed64,49,opt,name=OtherInvestment" json:"OtherInvestment,omitempty"`
	IntangibleAssets          *float64 `protobuf:"fixed64,50,opt,name=IntangibleAssets" json:"IntangibleAssets,omitempty"`
	GoodWill                  *float64 `protobuf:"fixed64,51,opt,name=GoodWill" json:"GoodWill,omitempty"`
	NegaGoodWill              *float64 `protobuf:"fixed64,52,opt,name=NegaGoodWill" json:"NegaGoodWill,omitempty"`
	DeferredTaxAssets         *float64 `protobuf:"fixed64,53,opt,name=DeferredTaxAssets" json:"DeferredTaxAssets,omitempty"`
	OtherNonCurrentAssets     *float64 `protobuf:"fixed64,54,opt,name=OtherNonCurrentAssets" json:"OtherNonCurrentAssets,omitempty"`
	NCAExcepItems             *float64 `protobuf:"fixed64,55,opt,name=NCAExcepItems" json:"NCAExcepItems,omitempty"`
	NCAAdjItems               *float64 `protobuf:"fixed64,56,opt,name=NCAAdjItems" json:"NCAAdjItems,omitempty"`
	TotalNonCurrentAssets     *float64 `protobuf:"fixed64,57,opt,name=TotalNonCurrentAssets" json:"TotalNonCurrentAssets,omitempty"`
	OtherAssets               *float64 `protobuf:"fixed64,58,opt,name=OtherAssets" json:"OtherAssets,omitempty"`
	TotalAssets               *float64 `protobuf:"fixed64,59,opt,name=TotalAssets" json:"TotalAssets,omitempty"`
	AccountsPayable           *float64 `protobuf:"fixed64,60,opt,name=AccountsPayable" json:"AccountsPayable,omitempty"`
	NotesPayable              *float64 `protobuf:"fixed64,61,opt,name=NotesPayable" json:"NotesPayable,omitempty"`
	TaxesPayable              *float64 `protobuf:"fixed64,62,opt,name=TaxesPayable" json:"TaxesPayable,omitempty"`
	DividendPayable           *float64 `protobuf:"fixed64,63,opt,name=DividendPayable" json:"DividendPayable,omitempty"`
	FAssociateFundRecCL       *float64 `protobuf:"fixed64,64,opt,name=FAssociateFundRecCL" json:"FAssociateFundRecCL,omitempty"`
	OtherFeesPayable          *float64 `protobuf:"fixed64,65,opt,name=OtherFeesPayable" json:"OtherFeesPayable,omitempty"`
	AdvanceReceipts           *float64 `protobuf:"fixed64,66,opt,name=AdvanceReceipts" json:"AdvanceReceipts,omitempty"`
	CustomerDeposits          *float64 `protobuf:"fixed64,67,opt,name=CustomerDeposits" json:"CustomerDeposits,omitempty"`
	ShortTermLoan             *float64 `protobuf:"fixed64,68,opt,name=ShortTermLoan" json:"ShortTermLoan,omitempty"`
	BankLoansAndOverdraft     *float64 `protobuf:"fixed64,69,opt,name=BankLoansAndOverdraft" json:"BankLoansAndOverdraft,omitempty"`
	FOtherLoanCL              *float64 `protobuf:"fixed64,70,opt,name=FOtherLoanCL" json:"FOtherLoanCL,omitempty"`
	UnearnedPremiumReserve    *float64 `protobuf:"fixed64,71,opt,name=UnearnedPremiumReserve" json:"UnearnedPremiumReserve,omitempty"`
	NotDecidedReservesCL      *float64 `protobuf:"fixed64,72,opt,name=NotDecidedReservesCL" json:"NotDecidedReservesCL,omitempty"`
	NotMatuRiskReserves       *float64 `protobuf:"fixed64,73,opt,name=NotMatuRiskReserves" json:"NotMatuRiskReserves,omitempty"`
	DerFinInstsCL             *float64 `protobuf:"fixed64,74,opt,name=DerFinInstsCL" json:"DerFinInstsCL,omitempty"`
	InveContLiaCL             *float64 `protobuf:"fixed64,75,opt,name=InveContLiaCL" json:"InveContLiaCL,omitempty"`
	BFinInstDepBorMoney       *float64 `protobuf:"fixed64,76,opt,name=BFinInstDepBorMoney" json:"BFinInstDepBorMoney,omitempty"`
	LFromOthBanksCL           *float64 `protobuf:"fixed64,77,opt,name=LFromOthBanksCL" json:"LFromOthBanksCL,omitempty"`
	SBbSecuProceeds           *float64 `protobuf:"fixed64,78,opt,name=SBbSecuProceeds" json:"SBbSecuProceeds,omitempty"`
	FAccruedBadDebtCL         *float64 `protobuf:"fixed64,79,opt,name=FAccruedBadDebtCL" json:"FAccruedBadDebtCL,omitempty"`
	FFinanceLeaseOwesCL       *float64 `protobuf:"fixed64,80,opt,name=FFinanceLeaseOwesCL" json:"FFinanceLeaseOwesCL,omitempty"`
	DeferredProceedsCL        *float64 `protobuf:"fixed64,81,opt,name=DeferredProceedsCL" json:"DeferredProceedsCL,omitempty"`
	OtherCurrentLiability     *float64 `protobuf:"fixed64,82,opt,name=OtherCurrentLiability" json:"OtherCurrentLiability,omitempty"`
	CLExcepItems              *float64 `protobuf:"fixed64,83,opt,name=CLExcepItems" json:"CLExcepItems,omitempty"`
	CLAdjItems                *float64 `protobuf:"fixed64,84,opt,name=CLAdjItems" json:"CLAdjItems,omitempty"`
	TotalCurrentLiability     *float64 `protobuf:"fixed64,85,opt,name=TotalCurrentLiability" json:"TotalCurrentLiability,omitempty"`
	NetCurrentLiability       *float64 `protobuf:"fixed64,86,opt,name=NetCurrentLiability" json:"NetCurrentLiability,omitempty"`
	AssetLessCLiability       *float64 `protobuf:"fixed64,87,opt,name=AssetLessCLiability" json:"AssetLessCLiability,omitempty"`
	LongtermLoan              *float64 `protobuf:"fixed64,88,opt,name=LongtermLoan" json:"LongtermLoan,omitempty"`
	NFOtherLoanNCL            *float64 `protobuf:"fixed64,89,opt,name=NFOtherLoanNCL" json:"NFOtherLoanNCL,omitempty"`
	LFromOthBanksNCL          *float64 `protobuf:"fixed64,90,opt,name=LFromOthBanksNCL" json:"LFromOthBanksNCL,omitempty"`
	LTAccountPayable          *float64 `protobuf:"fixed64,91,opt,name=LTAccountPayable" json:"LTAccountPayable,omitempty"`
	LongSalariesPay           *float64 `protobuf:"fixed64,92,opt,name=LongSalariesPay" json:"LongSalariesPay,omitempty"`
	NFAssociateFundRecNCL     *float64 `protobuf:"fixed64,93,opt,name=NFAssociateFundRecNCL" json:"NFAssociateFundRecNCL,omitempty"`
	NFFinanceLeaseOwesNCL     *float64 `protobuf:"fixed64,94,opt,name=NFFinanceLeaseOwesNCL" json:"NFFinanceLeaseOwesNCL,omitempty"`
	DeferredTaxLiability      *float64 `protobuf:"fixed64,95,opt,name=DeferredTaxLiability" json:"DeferredTaxLiability,omitempty"`
	NFDeferredProceedsNCL     *float64 `protobuf:"fixed64,96,opt,name=NFDeferredProceedsNCL" json:"NFDeferredProceedsNCL,omitempty"`
	NFAccruedBadDebtNCL       *float64 `protobuf:"fixed64,97,opt,name=NFAccruedBadDebtNCL" json:"NFAccruedBadDebtNCL,omitempty"`
	ConBillAndBond            *float64 `protobuf:"fixed64,98,opt,name=ConBillAndBond" json:"ConBillAndBond,omitempty"`
	DebtInstruIssued          *float64 `protobuf:"fixed64,99,opt,name=DebtInstruIssued" json:"DebtInstruIssued,omitempty"`
	DerFinInstsNCL            *float64 `protobuf:"fixed64,100,opt,name=DerFinInstsNCL" json:"DerFinInstsNCL,omitempty"`
	RetBfitsResp              *float64 `protobuf:"fixed64,101,opt,name=RetBfitsResp" json:"RetBfitsResp,omitempty"`
	NotDecidedReservesNCL     *float64 `protobuf:"fixed64,102,opt,name=NotDecidedReservesNCL" json:"NotDecidedReservesNCL,omitempty"`
	InveContLiaNCL            *float64 `protobuf:"fixed64,103,opt,name=InveContLiaNCL" json:"InveContLiaNCL,omitempty"`
	InsurAccPayableNCL        *float64 `protobuf:"fixed64,104,opt,name=InsurAccPayableNCL" json:"InsurAccPayableNCL,omitempty"`
	OtherNonCurrentLiab       *float64 `protobuf:"fixed64,105,opt,name=OtherNonCurrentLiab" json:"OtherNonCurrentLiab,omitempty"`
	NCLExcepItems             *float64 `protobuf:"fixed64,106,opt,name=NCLExcepItems" json:"NCLExcepItems,omitempty"`
	NCLAdjItems               *float64 `protobuf:"fixed64,107,opt,name=NCLAdjItems" json:"NCLAdjItems,omitempty"`
	TotalNonCurrentLiab       *float64 `protobuf:"fixed64,108,opt,name=TotalNonCurrentLiab" json:"TotalNonCurrentLiab,omitempty"`
	OtherLiability            *float64 `protobuf:"fixed64,109,opt,name=OtherLiability" json:"OtherLiability,omitempty"`
	TotalLiability            *float64 `protobuf:"fixed64,110,opt,name=TotalLiability" json:"TotalLiability,omitempty"`
	AssetLessTLiability       *float64 `protobuf:"fixed64,111,opt,name=AssetLessTLiability" json:"AssetLessTLiability,omitempty"`
	TotalIntANCTLiability     *float64 `protobuf:"fixed64,112,opt,name=TotalIntANCTLiability" json:"TotalIntANCTLiability,omitempty"`
	ShareCapital              *float64 `protobuf:"fixed64,113,opt,name=ShareCapital" json:"ShareCapital,omitempty"`
	OtherEquityinstruments    *float64 `protobuf:"fixed64,114,opt,name=OtherEquityinstruments" json:"OtherEquityinstruments,omitempty"`
	Reserve                   *float64 `protobuf:"fixed64,115,opt,name=Reserve" json:"Reserve,omitempty"`
	StockPremium              *float64 `protobuf:"fixed64,116,opt,name=StockPremium" json:"StockPremium,omitempty"`
	ReserveFund               *float64 `protobuf:"fixed64,117,opt,name=ReserveFund" json:"ReserveFund,omitempty"`
	CapitalReserveFund        *float64 `protobuf:"fixed64,118,opt,name=CapitalReserveFund" json:"CapitalReserveFund,omitempty"`
	RevaluationReserve        *float64 `protobuf:"fixed64,119,opt,name=RevaluationReserve" json:"RevaluationReserve,omitempty"`
	ExchangeReserve           *float64 `protobuf:"fixed64,120,opt,name=ExchangeReserve" json:"ExchangeReserve,omitempty"`
	OtherReserve              *float64 `protobuf:"fixed64,121,opt,name=OtherReserve" json:"OtherReserve,omitempty"`
	HoldProfit                *float64 `protobuf:"fixed64,122,opt,name=HoldProfit" json:"HoldProfit,omitempty"`
	SimulantAllotDividend     *float64 `protobuf:"fixed64,123,opt,name=SimulantAllotDividend" json:"SimulantAllotDividend,omitempty"`
	RetainedProfit            *float64 `protobuf:"fixed64,124,opt,name=RetainedProfit" json:"RetainedProfit,omitempty"`
	SEExcepItems              *float64 `protobuf:"fixed64,125,opt,name=SEExcepItems" json:"SEExcepItems,omitempty"`
	SEAdjItems                *float64 `protobuf:"fixed64,126,opt,name=SEAdjItems" json:"SEAdjItems,omitempty"`
	ShareholderEquity         *float64 `protobuf:"fixed64,127,opt,name=ShareholderEquity" json:"ShareholderEquity,omitempty"`
	MinorityInterests         *float64 `protobuf:"fixed64,128,opt,name=MinorityInterests" json:"MinorityInterests,omitempty"`
	TotalInterests            *float64 `protobuf:"fixed64,129,opt,name=TotalInterests" json:"TotalInterests,omitempty"`
	TotalIntATotalLiab        *float64 `protobuf:"fixed64,130,opt,name=TotalIntATotalLiab" json:"TotalIntATotalLiab,omitempty"`
	XXX_NoUnkeyedLiteral      struct{} `json:"-"`
	XXX_unrecognized          []byte   `json:"-"`
	XXX_sizecache             int32    `json:"-"`
}

type balanceSheetHK struct {
	balanceSheetDetails map[FinancialUniqueKey]*HK_BalanceSheetGEHK

	assistCalculation
}

type HK_CashFlowStatementHK struct {
	UniqueKey            *string  `protobuf:"bytes,120,opt,name=unique_key,json=uniqueKey" json:"unique_key,omitempty"`
	F10FinancialYear     *uint32  `protobuf:"varint,121,opt,name=f10_financial_year,json=f10FinancialYear" json:"f10_financial_year,omitempty"`
	F10FinancialType     *uint32  `protobuf:"varint,122,opt,name=f10_financial_type,json=f10FinancialType" json:"f10_financial_type,omitempty"`
	Type                 *uint32  `protobuf:"varint,109,opt,name=type" json:"type,omitempty"`
	EndDate              *uint32  `protobuf:"varint,110,opt,name=EndDate" json:"EndDate,omitempty"`
	InfoSource           *string  `protobuf:"bytes,111,opt,name=InfoSource" json:"InfoSource,omitempty"`
	CompanyType          *uint32  `protobuf:"varint,112,opt,name=CompanyType" json:"CompanyType,omitempty"`
	CurrencyUnit         *string  `protobuf:"bytes,113,opt,name=CurrencyUnit" json:"CurrencyUnit,omitempty"`
	CurrencyCode         *string  `protobuf:"bytes,119,opt,name=CurrencyCode" json:"CurrencyCode,omitempty"`
	AccountingStandards  *string  `protobuf:"bytes,114,opt,name=AccountingStandards" json:"AccountingStandards,omitempty"`
	PeriodMark           *uint32  `protobuf:"varint,115,opt,name=PeriodMark" json:"PeriodMark,omitempty"`
	Mark                 *uint32  `protobuf:"varint,116,opt,name=Mark" json:"Mark,omitempty"`
	InfoPublDate         *uint32  `protobuf:"varint,117,opt,name=InfoPublDate" json:"InfoPublDate,omitempty"`
	FiscalYear           *uint32  `protobuf:"varint,118,opt,name=FiscalYear" json:"FiscalYear,omitempty"`
	EarningBeforeTax     *float64 `protobuf:"fixed64,1,opt,name=EarningBeforeTax" json:"EarningBeforeTax,omitempty"`
	InterestIncomeAD     *float64 `protobuf:"fixed64,2,opt,name=InterestIncomeAD" json:"InterestIncomeAD,omitempty"`
	InterestExpAD        *float64 `protobuf:"fixed64,3,opt,name=InterestExpAD" json:"InterestExpAD,omitempty"`
	DividendIncomeAD     *float64 `protobuf:"fixed64,4,opt,name=DividendIncomeAD" json:"DividendIncomeAD,omitempty"`
	InvestProfAloss      *float64 `protobuf:"fixed64,5,opt,name=InvestProfAloss" json:"InvestProfAloss,omitempty"`
	AffilCompProfAloss   *float64 `protobuf:"fixed64,6,opt,name=AffilCompProfAloss" json:"AffilCompProfAloss,omitempty"`
	DevalAndAccBadDebt   *float64 `protobuf:"fixed64,7,opt,name=DevalAndAccBadDebt" json:"DevalAndAccBadDebt,omitempty"`
	DevalofProPlEquip    *float64 `protobuf:"fixed64,8,opt,name=DevalofProPlEquip" json:"DevalofProPlEquip,omitempty"`
	DevalofAvaForSaleInv *float64 `protobuf:"fixed64,9,opt,name=DevalofAvaForSaleInv" json:"DevalofAvaForSaleInv,omitempty"`
	DevalofInventories   *float64 `protobuf:"fixed64,10,opt,name=DevalofInventories" json:"DevalofInventories,omitempty"`
	DevalofTradeRece     *float64 `protobuf:"fixed64,11,opt,name=DevalofTradeRece" json:"DevalofTradeRece,omitempty"`
	DevalofGoodwill      *float64 `protobuf:"fixed64,12,opt,name=DevalofGoodwill" json:"DevalofGoodwill,omitempty"`
	DevalofOthers        *float64 `protobuf:"fixed64,13,opt,name=DevalofOthers" json:"DevalofOthers,omitempty"`
	RevaluationSurplus   *float64 `protobuf:"fixed64,14,opt,name=RevaluationSurplus" json:"RevaluationSurplus,omitempty"`
	CInFVofInvPropert    *float64 `protobuf:"fixed64,15,opt,name=CInFVofInvPropert" json:"CInFVofInvPropert,omitempty"`
	CInFVofDerFinInst    *float64 `protobuf:"fixed64,16,opt,name=CInFVofDerFinInst" json:"CInFVofDerFinInst,omitempty"`
	CInFVofOtherAssets   *float64 `protobuf:"fixed64,17,opt,name=CInFVofOtherAssets" json:"CInFVofOtherAssets,omitempty"`
	ProfitDispOfAssets   *float64 `protobuf:"fixed64,18,opt,name=ProfitDispOfAssets" json:"ProfitDispOfAssets,omitempty"`
	ProfDispOfAFSaleInv  *float64 `protobuf:"fixed64,19,opt,name=ProfDispOfAFSaleInv" json:"ProfDispOfAFSaleInv,omitempty"`
	ProfDispOfAffCEqu    *float64 `protobuf:"fixed64,20,opt,name=ProfDispOfAffCEqu" json:"ProfDispOfAffCEqu,omitempty"`
	ProfDispOfProPlEquip *float64 `protobuf:"fixed64,21,opt,name=ProfDispOfProPlEquip" json:"ProfDispOfProPlEquip,omitempty"`
	ProfDispOfOthAssets  *float64 `protobuf:"fixed64,22,opt,name=ProfDispOfOthAssets" json:"ProfDispOfOthAssets,omitempty"`
	DepDividerSale       *float64 `protobuf:"fixed64,23,opt,name=DepDividerSale" json:"DepDividerSale,omitempty"`
	Depreciation         *float64 `protobuf:"fixed64,24,opt,name=Depreciation" json:"Depreciation,omitempty"`
	IntangibleAssetAmort *float64 `protobuf:"fixed64,25,opt,name=IntangibleAssetAmort" json:"IntangibleAssetAmort,omitempty"`
	OtherDepDividerSale  *float64 `protobuf:"fixed64,26,opt,name=OtherDepDividerSale" json:"OtherDepDividerSale,omitempty"`
	FinancialExpense     *float64 `protobuf:"fixed64,27,opt,name=FinancialExpense" json:"FinancialExpense,omitempty"`
	ExchangeIncome       *float64 `protobuf:"fixed64,28,opt,name=ExchangeIncome" json:"ExchangeIncome,omitempty"`
	UnrealExchangeIncome *float64 `protobuf:"fixed64,29,opt,name=UnrealExchangeIncome" json:"UnrealExchangeIncome,omitempty"`
	SpeItemsManageAdj    *float64 `protobuf:"fixed64,30,opt,name=SpeItemsManageAdj" json:"SpeItemsManageAdj,omitempty"`
	AdjItemsManageAdj    *float64 `protobuf:"fixed64,31,opt,name=AdjItemsManageAdj" json:"AdjItemsManageAdj,omitempty"`
	OpeProBefChgInOpeCap *float64 `protobuf:"fixed64,32,opt,name=OpeProBefChgInOpeCap" json:"OpeProBefChgInOpeCap,omitempty"`
	InventoryChange      *float64 `protobuf:"fixed64,33,opt,name=InventoryChange" json:"InventoryChange,omitempty"`
	ProUnderDevelChange  *float64 `protobuf:"fixed64,34,opt,name=ProUnderDevelChange" json:"ProUnderDevelChange,omitempty"`
	AccReceivChange      *float64 `protobuf:"fixed64,35,opt,name=AccReceivChange" json:"AccReceivChange,omitempty"`
	AccPayableChange     *float64 `protobuf:"fixed64,36,opt,name=AccPayableChange" json:"AccPayableChange,omitempty"`
	AdvReceiptsChange    *float64 `protobuf:"fixed64,37,opt,name=AdvReceiptsChange" json:"AdvReceiptsChange,omitempty"`
	AdvPaymentChange     *float64 `protobuf:"fixed64,38,opt,name=AdvPaymentChange" json:"AdvPaymentChange,omitempty"`
	FAAFValOnPLChange    *float64 `protobuf:"fixed64,39,opt,name=FAAFValOnPLChange" json:"FAAFValOnPLChange,omitempty"`
	FLAFValOnPLChange    *float64 `protobuf:"fixed64,40,opt,name=FLAFValOnPLChange" json:"FLAFValOnPLChange,omitempty"`
	DerFinInstChange     *float64 `protobuf:"fixed64,41,opt,name=DerFinInstChange" json:"DerFinInstChange,omitempty"`
	InsuReceivChange     *float64 `protobuf:"fixed64,42,opt,name=InsuReceivChange" json:"InsuReceivChange,omitempty"`
	InsuContLiabChange   *float64 `protobuf:"fixed64,43,opt,name=InsuContLiabChange" json:"InsuContLiabChange,omitempty"`
	AccRePayableChange   *float64 `protobuf:"fixed64,44,opt,name=AccRePayableChange" json:"AccRePayableChange,omitempty"`
	BBackSFAssetsChange  *float64 `protobuf:"fixed64,45,opt,name=BBackSFAssetsChange" json:"BBackSFAssetsChange,omitempty"`
	SpeItemsWCapChange   *float64 `protobuf:"fixed64,46,opt,name=SpeItemsWCapChange" json:"SpeItemsWCapChange,omitempty"`
	AdjItemsWCapChange   *float64 `protobuf:"fixed64,47,opt,name=AdjItemsWCapChange" json:"AdjItemsWCapChange,omitempty"`
	BankDepositChange    *float64 `protobuf:"fixed64,48,opt,name=BankDepositChange" json:"BankDepositChange,omitempty"`
	LoansAAdvanChange    *float64 `protobuf:"fixed64,49,opt,name=LoansAAdvanChange" json:"LoansAAdvanChange,omitempty"`
	BFAAFValOnPLChange   *float64 `protobuf:"fixed64,50,opt,name=BFAAFValOnPLChange" json:"BFAAFValOnPLChange,omitempty"`
	SpeItemsOpeAchange   *float64 `protobuf:"fixed64,51,opt,name=SpeItemsOpeAchange" json:"SpeItemsOpeAchange,omitempty"`
	BorFromCBChange      *float64 `protobuf:"fixed64,52,opt,name=BorFromCBChange" json:"BorFromCBChange,omitempty"`
	CusDepositsChange    *float64 `protobuf:"fixed64,53,opt,name=CusDepositsChange" json:"CusDepositsChange,omitempty"`
	BFLAFValOnPLChange   *float64 `protobuf:"fixed64,54,opt,name=BFLAFValOnPLChange" json:"BFLAFValOnPLChange,omitempty"`
	SpeItemsOpeLchange   *float64 `protobuf:"fixed64,55,opt,name=SpeItemsOpeLchange" json:"SpeItemsOpeLchange,omitempty"`
	CashReceiptsFOpe     *float64 `protobuf:"fixed64,56,opt,name=CashReceiptsFOpe" json:"CashReceiptsFOpe,omitempty"`
	HKProfitsTaxPaid     *float64 `protobuf:"fixed64,57,opt,name=HKProfitsTaxPaid" json:"HKProfitsTaxPaid,omitempty"`
	ChinaIncomeTaxPaid   *float64 `protobuf:"fixed64,58,opt,name=ChinaIncomeTaxPaid" json:"ChinaIncomeTaxPaid,omitempty"`
	OtherTaxes           *float64 `protobuf:"fixed64,59,opt,name=OtherTaxes" json:"OtherTaxes,omitempty"`
	DividendsRecBO       *float64 `protobuf:"fixed64,60,opt,name=DividendsRecBO" json:"DividendsRecBO,omitempty"`
	DividendPaidBO       *float64 `protobuf:"fixed64,61,opt,name=DividendPaidBO" json:"DividendPaidBO,omitempty"`
	InterestRecBO        *float64 `protobuf:"fixed64,62,opt,name=InterestRecBO" json:"InterestRecBO,omitempty"`
	InterestPaidBO       *float64 `protobuf:"fixed64,63,opt,name=InterestPaidBO" json:"InterestPaidBO,omitempty"`
	SpeItemsOpeBusi      *float64 `protobuf:"fixed64,64,opt,name=SpeItemsOpeBusi" json:"SpeItemsOpeBusi,omitempty"`
	AdjItemsOpeBusi      *float64 `protobuf:"fixed64,65,opt,name=AdjItemsOpeBusi" json:"AdjItemsOpeBusi,omitempty"`
	NetOpeCFlow          *float64 `protobuf:"fixed64,66,opt,name=NetOpeCFlow" json:"NetOpeCFlow,omitempty"`
	FinanceAndSpeItems   *float64 `protobuf:"fixed64,67,opt,name=FinanceAndSpeItems" json:"FinanceAndSpeItems,omitempty"`
	InterestRecIB        *float64 `protobuf:"fixed64,68,opt,name=InterestRecIB" json:"InterestRecIB,omitempty"`
	DividendsRecIB       *float64 `protobuf:"fixed64,69,opt,name=DividendsRecIB" json:"DividendsRecIB,omitempty"`
	RestrictCashChange   *float64 `protobuf:"fixed64,70,opt,name=RestrictCashChange" json:"RestrictCashChange,omitempty"`
	LoanReceivableChange *float64 `protobuf:"fixed64,71,opt,name=LoanReceivableChange" json:"LoanReceivableChange,omitempty"`
	DepositChange        *float64 `protobuf:"fixed64,72,opt,name=DepositChange" json:"DepositChange,omitempty"`
	VendCapitalAssents   *float64 `protobuf:"fixed64,73,opt,name=VendCapitalAssents" json:"VendCapitalAssents,omitempty"`
	PurCapitalAssents    *float64 `protobuf:"fixed64,74,opt,name=PurCapitalAssents" json:"PurCapitalAssents,omitempty"`
	VendIntassets        *float64 `protobuf:"fixed64,75,opt,name=VendIntassets" json:"VendIntassets,omitempty"`
	PurIntassets         *float64 `protobuf:"fixed64,76,opt,name=PurIntassets" json:"PurIntassets,omitempty"`
	VendAffCompanies     *float64 `protobuf:"fixed64,77,opt,name=VendAffCompanies" json:"VendAffCompanies,omitempty"`
	PurAffCompanies      *float64 `protobuf:"fixed64,78,opt,name=PurAffCompanies" json:"PurAffCompanies,omitempty"`
	DisinvestmentCash    *float64 `protobuf:"fixed64,79,opt,name=DisinvestmentCash" json:"DisinvestmentCash,omitempty"`
	InvestPaymentCash    *float64 `protobuf:"fixed64,80,opt,name=InvestPaymentCash" json:"InvestPaymentCash,omitempty"`
	InvestAdjustedOther  *float64 `protobuf:"fixed64,81,opt,name=InvestAdjustedOther" json:"InvestAdjustedOther,omitempty"`
	InvestAdjustedItems  *float64 `protobuf:"fixed64,82,opt,name=InvestAdjustedItems" json:"InvestAdjustedItems,omitempty"`
	NetInvbusiCFlow      *float64 `protobuf:"fixed64,83,opt,name=NetInvbusiCFlow" json:"NetInvbusiCFlow,omitempty"`
	CashAndOtherBefFin   *float64 `protobuf:"fixed64,84,opt,name=CashAndOtherBefFin" json:"CashAndOtherBefFin,omitempty"`
	NetCashBeforFinance  *float64 `protobuf:"fixed64,85,opt,name=NetCashBeforFinance" json:"NetCashBeforFinance,omitempty"`
	NewLoan              *float64 `protobuf:"fixed64,86,opt,name=NewLoan" json:"NewLoan,omitempty"`
	Refund               *float64 `protobuf:"fixed64,87,opt,name=Refund" json:"Refund,omitempty"`
	IssueShares          *float64 `protobuf:"fixed64,88,opt,name=IssueShares" json:"IssueShares,omitempty"`
	IssueBonds           *float64 `protobuf:"fixed64,89,opt,name=IssueBonds" json:"IssueBonds,omitempty"`
	InterestPaidFB       *float64 `protobuf:"fixed64,90,opt,name=InterestPaidFB" json:"InterestPaidFB,omitempty"`
	DividendPaidFB       *float64 `protobuf:"fixed64,91,opt,name=DividendPaidFB" json:"DividendPaidFB,omitempty"`
	AbsorbInvestIncome   *float64 `protobuf:"fixed64,92,opt,name=AbsorbInvestIncome" json:"AbsorbInvestIncome,omitempty"`
	IssExpAPayOfRedSecu  *float64 `protobuf:"fixed64,93,opt,name=IssExpAPayOfRedSecu" json:"IssExpAPayOfRedSecu,omitempty"`
	PledgedDepositChange *float64 `protobuf:"fixed64,94,opt,name=PledgedDepositChange" json:"PledgedDepositChange,omitempty"`
	FinanceAdjustedOther *float64 `protobuf:"fixed64,95,opt,name=FinanceAdjustedOther" json:"FinanceAdjustedOther,omitempty"`
	FinanceAdjustedItems *float64 `protobuf:"fixed64,96,opt,name=FinanceAdjustedItems" json:"FinanceAdjustedItems,omitempty"`
	NetCashFromFinance   *float64 `protobuf:"fixed64,97,opt,name=NetCashFromFinance" json:"NetCashFromFinance,omitempty"`
	EffectOfRate         *float64 `protobuf:"fixed64,98,opt,name=EffectOfRate" json:"EffectOfRate,omitempty"`
	OtherItemsAffectNC   *float64 `protobuf:"fixed64,99,opt,name=OtherItemsAffectNC" json:"OtherItemsAffectNC,omitempty"`
	NetCash              *float64 `protobuf:"fixed64,100,opt,name=NetCash" json:"NetCash,omitempty"`
	BeginPeriodCash      *float64 `protobuf:"fixed64,101,opt,name=BeginPeriodCash" json:"BeginPeriodCash,omitempty"`
	ItemsPeriod          *float64 `protobuf:"fixed64,102,opt,name=ItemsPeriod" json:"ItemsPeriod,omitempty"`
	CashEndPer           *float64 `protobuf:"fixed64,103,opt,name=CashEndPer" json:"CashEndPer,omitempty"`
	CashABankBalances    *float64 `protobuf:"fixed64,104,opt,name=CashABankBalances" json:"CashABankBalances,omitempty"`
	BankDeposits         *float64 `protobuf:"fixed64,105,opt,name=BankDeposits" json:"BankDeposits,omitempty"`
	InterestRecCB        *float64 `protobuf:"fixed64,106,opt,name=InterestRecCB" json:"InterestRecCB,omitempty"`
	InterestPaidCB       *float64 `protobuf:"fixed64,107,opt,name=InterestPaidCB" json:"InterestPaidCB,omitempty"`
	CashCashEquival      *float64 `protobuf:"fixed64,108,opt,name=CashCashEquival" json:"CashCashEquival,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type cashFlowHK struct {
	cashFlowDetails map[FinancialUniqueKey]*HK_CashFlowStatementHK

	assistCalculation
}

type HK_BuyBackitems struct {
	PublDate             *uint64  `protobuf:"varint,1,opt,name=publ_date,json=publDate" json:"publ_date,omitempty"`
	EndDate              *uint64  `protobuf:"varint,2,opt,name=EndDate" json:"EndDate,omitempty"`
	BuybackMoney         *float64 `protobuf:"fixed64,3,opt,name=buyback_money,json=buybackMoney" json:"buyback_money,omitempty"`
	BuybackSum           *int64   `protobuf:"varint,4,opt,name=buyback_sum,json=buybackSum" json:"buyback_sum,omitempty"`
	Percentage           *float64 `protobuf:"fixed64,5,opt,name=percentage" json:"percentage,omitempty"`
	HighPrice            *float64 `protobuf:"fixed64,6,opt,name=high_price,json=highPrice" json:"high_price,omitempty"`
	LowPrice             *float64 `protobuf:"fixed64,7,opt,name=low_price,json=lowPrice" json:"low_price,omitempty"`
	CumulativeSum        *int64   `protobuf:"varint,8,opt,name=cumulative_sum,json=cumulativeSum" json:"cumulative_sum,omitempty"`
	CumulativeSumToTS    *float64 `protobuf:"fixed64,9,opt,name=cumulative_sumToTS,json=cumulativeSumToTS" json:"cumulative_sumToTS,omitempty"`
	BuybackMoneyCode     *string  `protobuf:"bytes,10,opt,name=BuybackMoneyCode" json:"BuybackMoneyCode,omitempty"`
	ShareType            *string  `protobuf:"bytes,11,opt,name=share_type,json=shareType" json:"share_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

type buybackHK struct {
	buybackDetails []*HK_BuyBackitems
}

type NoUnkeyedLiterals struct{}

type IndicatorResult struct {
	_ NoUnkeyedLiterals

	IsInvalid bool               // 标识指标计算结果无效
	Source    byte               // 指标的数据来源(f10, 运营后台，爬虫)
	Key       FinancialUniqueKey // 指标的唯一key(由财年和财报类型决定)
	Value     float64            // 指标对应的取值
	YoY       float64            // 同比
	MoM       float64            // 环比
	EndDate   int64              // 当期财报统计截止日时间戳
}

// IndicatorResults 特定指标所有财务周期的计算数据
type IndicatorResults map[FinancialUniqueKey]IndicatorResult

// AllIndicatorResults 各指标所有财务周期的计算数据
type AllIndicatorResults map[int]IndicatorResults

type operationAndCrawlerHK struct {
	operationDetails AllIndicatorResults // 运营方式的财务指标数据
	crawlerDetails   AllIndicatorResults // 爬虫获取的财务指标
}

var buffers = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset() // 这里存在问题，因为底层切片的len为0了，但是cap可能还很大(https://github.com/CodeFish-xiao/go_concurrent_notes/blob/master/1.%E5%9F%BA%E6%9C%AC%E5%B9%B6%E5%8F%91%E5%8E%9F%E8%AF%AD/1.10%EF%BC%9APool%EF%BC%9A%E6%80%A7%E8%83%BD%E6%8F%90%E5%8D%87%E5%A4%A7%E6%9D%80%E5%99%A8/10.00-Pool%EF%BC%9A%E6%80%A7%E8%83%BD%E6%8F%90%E5%8D%87%E5%A4%A7%E6%9D%80%E5%99%A8.md)
	buffers.Put(buf)

	fmt.Print() // buf.Reset()
}

// DiffSli 实现集合的差集
func DiffSli(sli1, sli2 []FinancialUniqueKey) []FinancialUniqueKey {
	var (
		diff    = make([]FinancialUniqueKey, 0, len(sli1))
		sli1Map = make(map[FinancialUniqueKey][0]struct{}, len(sli1))
		sli2Map = make(map[FinancialUniqueKey][0]struct{}, len(sli2))
	)
	for _, ele := range sli1 { // 对`sli1`去重
		sli1Map[ele] = [0]struct{}{}
	}
	for _, ele := range sli2 {
		sli2Map[ele] = [0]struct{}{}
	}
	for ele := range sli1Map {
		if _, has := sli2Map[ele]; !has {
			diff = append(diff, ele)
		}
	}
	return diff
}

func TestDiff() {
	r := DiffSli(append([]FinancialUniqueKey{"mlee1", "mlee2", "mlee2", "mlee3", "mlee1"}, []FinancialUniqueKey{"mlee1000", "mlee121", "mlee3", "mlee2", "mlee2"}...), nil)
	fmt.Println(r)
}

func fff() {
	var sli = make([]int, 0, 10)
	f := func(a []int) {
		sli = append(sli, a...)
	}
	f([]int{1, 2, 3})
	f([]int{1, 20, 30})
	fmt.Println(sli)
}

func testSize() {
	type a struct {
		f1 string
		f2 int32
		f3 []int32
	}
	runtime.NumCPU()
	fmt.Println(unsafe.Sizeof(a{}), unsafe.Sizeof(new(a)), unsafe.Sizeof(new(string)), unsafe.Sizeof(new([]int)))
	v1 := "aaa"
	v1_ := strings.Repeat("a", 100)
	v2 := []int{}
	v3 := new(string)
	v4 := []int{1, 2}

	fmt.Println(unsafe.Sizeof(v1), unsafe.Sizeof(v1_), unsafe.Sizeof(v2), unsafe.Sizeof(*v3), unsafe.Sizeof(v4))
}

func getAllMonthLastDay(year int) {
	for _, month := range []time.Month{
		time.January, time.February, time.March,
		time.April, time.May, time.June,
		time.July, time.August, time.September,
		time.October, time.November, time.December,
	} {
		t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC) // 计算下一个月的第0天，即本月的最后一天
		fmt.Println(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
}

// func EqualSli(sli1, sli2 []int){
// 	slices.Compare[]()
// 	if reflect.DeepEqual(sli1, sli2) {
// 		fmt.Println("equal")
// 	} else {
// 		fmt.Println("not equal")
// 	}
// }

func testStringFunc() {
	ss1 := strings.Split("_", "_")
	fmt.Println(len(ss1), ss1, len(ss1[0]), len(ss1[1]))
	ss2 := strings.Split("a_", "_")
	fmt.Println(len(ss2), ss2, len(ss2[0]), len(ss2[1]))
	ss3 := strings.Split("_b", "_")
	fmt.Println(len(ss3), ss3, len(ss3[0]), len(ss3[1]))
	ss4 := strings.Split("aa_bb", "_")
	fmt.Println(len(ss4), ss4, len(ss4[0]), len(ss4[1]))
}

func testDeferOrder() {
	defer fmt.Println("defer first")
	defer fmt.Println("defer second")
	fmt.Println("third")
}

func testWaitGroup() {
	fmt.Println("begin")
	a := []int{}
	var wg sync.WaitGroup
	for range a {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("----")
		}()
	}
	wg.Wait()
	fmt.Println("end")
}

type testMM map[int]map[string][]int

func (d testMM) String() string {
	var b strings.Builder
	for key1, m1 := range d {
		for key2, value := range m1 {
			b.WriteString(fmt.Sprintln(key1, key2, value))
		}
	}
	return b.String()
}

func (d testMM) add1(key1 int, key2 string, value []int) testMM {
	if d == nil {
		d = make(testMM)
	}
	if _, has := d[key1]; !has {
		d[key1] = make(map[string][]int)
	}
	d[key1][key2] = value
	return d
}

func testMMAdd1() {
	var d testMM
	d = d.add1(1, "a", []int{1, 2}).
		add1(2, "a", []int{3, 4}).
		add1(1, "b", []int{5, 6}).
		add1(1, "a", []int{10, 20, 30}).
		add1(3, "c", []int{3, 4}).
		add1(4, "d", []int{13, 14})
	fmt.Println("add1:\n", d)
}

func (d testMM) add2(key1 int, key2 string, value []int) {
	if d == nil {
		d = make(testMM)
	}
	if _, has := d[key1]; !has {
		d[key1] = make(map[string][]int)
	}
	d[key1][key2] = value
}

func testMMAdd2() {
	d := new(testMM) // 解引用后，会再拷贝给接收器临时对象
	d.add2(1, "a", []int{1, 2})
	d.add2(2, "a", []int{3, 4})
	d.add2(1, "b", []int{5, 6})
	d.add2(1, "a", []int{10, 20, 30})
	d.add2(3, "c", []int{3, 4})
	d.add2(4, "d", []int{13, 14})
	fmt.Println("add2:\n", d)
}

func testMMAdd2_() {
	var d testMM // 会拷贝给接收器临时对象
	d.add2(1, "a", []int{1, 2})
	d.add2(2, "a", []int{3, 4})
	d.add2(1, "b", []int{5, 6})
	d.add2(1, "a", []int{10, 20, 30})
	d.add2(3, "c", []int{3, 4})
	d.add2(4, "d", []int{13, 14})
	fmt.Println("add2_:\n", d)
}

func (d *testMM) add3(key1 int, key2 string, value []int) {
	if d == nil {
		d = new(testMM)
	}
	if *d == nil {
		*d = make(testMM)
	}

	if _, has := (*d)[key1]; !has {
		(*d)[key1] = make(map[string][]int)
	}
	(*d)[key1][key2] = value
}

func testMMAdd3() {
	var d testMM
	d.add3(1, "a", []int{1, 2})
	d.add3(2, "a", []int{3, 4})
	d.add3(1, "b", []int{5, 6})
	d.add3(1, "a", []int{10, 20, 30})
	d.add3(3, "c", []int{3, 4})
	d.add3(4, "d", []int{13, 14})
	fmt.Println("add3:\n", d)
}

func testMMAdd3_() {
	var d *testMM
	d.add3(1, "a", []int{1, 2})
	d.add3(2, "a", []int{3, 4})
	d.add3(1, "b", []int{5, 6})
	d.add3(1, "a", []int{10, 20, 30})
	d.add3(3, "c", []int{3, 4})
	d.add3(4, "d", []int{13, 14})
	fmt.Println("add3_:\n", d)
}

func (d *testMM) add4(key1 int, key2 string, value []int) *testMM {
	if d == nil {
		d = new(testMM)

	}
	if *d == nil {
		*d = make(testMM)
	}
	if _, has := (*d)[key1]; !has {
		(*d)[key1] = make(map[string][]int)
	}
	(*d)[key1][key2] = value
	return d
}

func testMMAdd4() {
	var d testMM
	d.add4(1, "a", []int{1, 2}).
		add4(2, "a", []int{3, 4}).
		add4(1, "b", []int{5, 6}).
		add4(1, "a", []int{10, 20, 30}).
		add4(3, "c", []int{3, 4}).
		add4(4, "d", []int{13, 14})
	fmt.Println("add4:\n", d)
}

func testMMAdd4_() {
	var d *testMM
	d = d.add4(1, "a", []int{1, 2}).
		add4(2, "a", []int{3, 4}).
		add4(1, "b", []int{5, 6}).
		add4(1, "a", []int{10, 20, 30}).
		add4(3, "c", []int{3, 4}).
		add4(4, "d", []int{13, 14})
	fmt.Println("add4_:\n", d)
}

func testMMAdd() {
	var d testMM
	if d == nil {
		fmt.Println("nil nil")
	} else {
		fmt.Println("non-nil")
	}
	if d_ := new(testMM); d_ == nil {
		fmt.Println("nil nil nil")
	} else {
		fmt.Println("non-nil non-nil")
	}
	if d_ := new(testMM); *d_ == nil {
		fmt.Println("nil nil nil nil")
	} else {
		fmt.Println("non-nil non-nil non-nil")
	}
	testMMAdd1()
	testMMAdd2()
	testMMAdd2_()
	testMMAdd3()
	testMMAdd3_()
	testMMAdd4()
	testMMAdd4_()
}

func testMapInit() {
	f := func(m map[int]int) {
		fmt.Println(m == nil, len(m)) // map类型没有cap函数
	}
	var m0 map[int]int
	f(m0)                   // case 0
	f(make(map[int]int))    // case 2
	f(make(map[int]int, 2)) // case 3
	// 如何区分出 case2 和 case3
}

func testCalculation() {
	a := 1.0 * 2 * 3.0
	fmt.Println(18/a, 18/2*3.0)
	v := 18 / a
	fmt.Printf("%f, %v, %+v, %T\n", v, v, v, v)
}

type MarketType = byte

const (
	MarketUndefined MarketType = iota // 占位零值，不采用
	MarketCN                          // A股市场
	MarketHK                          // 港股市场
	MarketUS                          // 美股市场
	MarketSG                          // 新加坡市场
	MarketCA                          // 加拿大市场
	MarketAU                          // 澳大利亚市场
	MarketJP                          // 日本市场
	MarketMY                          // 马来西亚市场
)

// MarketTimezone 市场和其对应的时区映射
var MarketTimezone = map[MarketType]string{
	MarketCN: "Asia/Shanghai",
	MarketHK: "Asia/Hong_Kong",
	MarketUS: "America/New_York",
	MarketSG: "Asia/Singapore",
	MarketCA: "America/Toronto",
	MarketAU: "Australia/Sydney",
	MarketJP: "Asia/Tokyo",
	MarketMY: "Asia/Kuala_Lumpur",
}

func test_getLastDayOfYearMonth(year int) {
	for m := time.January; m <= time.December; m++ {
		getLastDayOfYearMonth(MarketCN, year, m)
	}
}

func getLastDayOfYearMonth(typ MarketType, year int, month time.Month) {
	loc := time.UTC
	if dstTZStr, has := MarketTimezone[typ]; has {
		if loc_, err := time.LoadLocation(dstTZStr); err == nil {
			loc = loc_
		}
	}
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, loc)
	fmt.Println(t, t.Unix()) // 计算下一个月的第0天，即本月的最后一天
}

func testSyncMap() {
	var (
		isNotEmpty1 bool
		isNotEmpty2 bool
		isNotEmpty3 bool
		m           sync.Map
	)
	m.Range(func(_, _ any) bool {
		fmt.Println("isNotEmpty1 range")
		isNotEmpty1 = true
		return false // 停止迭代
	})
	m.Store(1, [0]struct{}{})
	m.Store(1, [0]struct{}{})
	m.Store(11, [0]struct{}{})
	m.Store(111, [0]struct{}{})
	m.Range(func(_, _ any) bool {
		fmt.Println("isNotEmpty2,3 range")
		isNotEmpty2 = true
		isNotEmpty3 = true
		return false
	})
	_, ok := m.Load(11)
	fmt.Println("11", ok)
	fmt.Println(isNotEmpty1, isNotEmpty2, isNotEmpty3)

	clearSyncMap(m)

	m.Range(func(_, _ any) bool {
		isNotEmpty3 = true
		fmt.Println("isNotEmpty3 range")
		return false // 停止迭代
	})
	fmt.Println(isNotEmpty1, isNotEmpty2, isNotEmpty3)
	_, ok = m.Load(11)
	fmt.Println("11", ok)
}

func clearSyncMap(m sync.Map) {
	m.Range(func(key, _ any) bool {
		m.Delete(key)
		return true // 继续迭代
	})
}

func testFormatTime() {
	// func (t Time) Format(layout string) string  标准库中这个方法（包）是真的复杂，圈复杂度爆表！！！
	fmt.Println(time.Now().Format("200601"), time.Now().Format(time.DateOnly))
}

func testSyncSafetyForSlice() {
	type s struct {
		a int
		b string
	}
	const nums = 10
	var (
		sli = make([]s, nums)
		wg  sync.WaitGroup
	)
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			sli[index].a = i * 10
			sli[index].b = strconv.Itoa(i * 10)
		}(i)
	}
	wg.Wait()

	fmt.Println(sli)
}

func TestTime() {
	now := time.Now()
	fmt.Println(now)
	fmt.Println(now.Truncate(24 * time.Hour))
	fmt.Println(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()))

	t := now.Add(2 * time.Hour)
	fmt.Println(now, t)
	t1, t2 := TimeOnlyContainsYMD(MarketCN, now), TimeOnlyContainsYMD(MarketCN, t)
	fmt.Println(t1, t2, t1.Equal(t2))

}

var marketTimezone = map[MarketType]string{
	MarketCN: "Asia/Shanghai",
	MarketHK: "Asia/Hong_Kong",
	MarketUS: "America/New_York",
	MarketSG: "Asia/Singapore",
	MarketCA: "America/Toronto",
	MarketAU: "Australia/Sydney",
	MarketJP: "Asia/Tokyo",
	MarketMY: "Asia/Kuala_Lumpur",
}

func TimeOnlyContainsYMD(typ MarketType, t time.Time) time.Time {
	loc, _ := time.LoadLocation(marketTimezone[typ])
	if loc == nil {
		loc = time.UTC
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func testRound() {
	fmt.Println(math.Round(0.3), math.Round(0.5), math.Round(0.7), math.Round(2.3), math.Round(2.5), math.Round(2.6))
}

func testJSON() {
	type a struct {
		A1 *int
		A2 *string
		A3 int
		A4 string
		A5 int
	}
	d1, d2 := 0, "mleee"
	d := &a{A1: &d1, A4: d2, A3: 0}
	bs, _ := json.Marshal(d)
	var dd *a
	fmt.Println(&dd) // nil 变量可以取其地址
	json.Unmarshal(bs, &dd)
	fmt.Printf("%+v\n%+v\n", dd, *dd)
}

func TestOverflow() {
	f := 18346963000000000.000000
	i := int64(f)
	ii := uint64(f)
	fmt.Println(f, i, i*1e3, ii*1e3)
	fmt.Println(2596039000.000000/183469630000000.000000, (2596039000000.000000-2406557000.000000)/2406557000.000000)
	fmt.Println(math.MaxInt64)
	
	// https://blog.csdn.net/raoxiaoya/article/details/129158263
	fmt.Println(math.Round(0.4999999), math.Round(23.51), math.Round(0.51), math.RoundToEven(0.51), math.RoundToEven(0.5), math.RoundToEven(0.61))
}
