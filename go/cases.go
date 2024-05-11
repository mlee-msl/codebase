package main

import (
	"fmt"
	"math/rand"
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
