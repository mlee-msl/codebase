package main

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// 功能测试
func TestTaskGroup(t *testing.T) {
	var (
		tg TaskGroup

		tasks = []*Task{
			NewTask(1, true, task1),
			NewTask(2, true, task2),
			NewTask(3, true, task3),
		}
	)

	taskResult := tg.SetWorkerNums(4).AddTask(tasks...).Run()
	fmt.Printf("**************TaskGroup************\n%+v\n", taskResult)
	for fno, result := range taskResult {
		fmt.Printf("FNO: %d, RESULT: %v , STATUS: %v\n", fno, result.Result(), result.Error())
	}
}

func getRandomNum() int {
	return rand.Int() % 1024
}

func task1() (interface{}, error) {
	const taskFlag = "TASK1"
	fmt.Println(taskFlag)
	return getRandomNum(), errors.New(fmt.Sprintf("%s err", taskFlag))
}

type task2Struct struct {
	a int
	b string
}

func task2() (interface{}, error) {
	const taskFlag = "TASK2"
	fmt.Println(taskFlag)
	return task2Struct{
		a: getRandomNum(),
		b: "mlee",
	}, nil
}

func task3() (interface{}, error) {
	const taskFlag = "TASK3"
	fmt.Println(taskFlag)
	return fmt.Sprintf("%s: The data is %d", taskFlag, getRandomNum()), errors.New(fmt.Sprintf("%s err", taskFlag))
}

// 性能测试
func BenchmarkTaskGroup(b *testing.B) {

}
