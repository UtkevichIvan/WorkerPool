package main

import (
	"VKTest/workerpool"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

const (
	CreateCommand = iota
	DeleteCommand
	ExitCommand
)

func Gen(n int, jobs chan string, stopGen chan struct{}, stop chan struct{}) {
	defer close(jobs)
	defer close(stop)

	for {
		select {
		case <-time.After(time.Second / 2):
			b := make([]byte, n)
			for i := range b {
				b[i] = letterBytes[rand.Intn(len(letterBytes))]
			}
			select {
			case jobs <- string(b):
			case <-stopGen:
				return
			}
		case <-stopGen:
			return
		}
	}
}

func main() {
	commands := []int{CreateCommand, CreateCommand, CreateCommand, DeleteCommand, DeleteCommand, CreateCommand, ExitCommand}
	jobs := make(chan string)
	stopGen := make(chan struct{})
	genStopSignal := make(chan struct{})

	workerPool := workerpool.NewPool(jobs)
	defer workerPool.Close()

	go Gen(10, jobs, stopGen, genStopSignal)
	defer func() { <-genStopSignal }()

	for _, word := range commands {
		switch word {
		case CreateCommand:
			workerPool.Add()
		case DeleteCommand:
			workerPool.StopOne()
		case ExitCommand:
			close(stopGen)
			return
		default:
			panic("unsupported command")
		}
		time.Sleep(time.Second * 3)
	}
}
