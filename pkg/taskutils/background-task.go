package taskutils

import (
	"fmt"
	"log"
	"sync"

	"greenlight/pkg/jsonlog"
)

type LogFunc func(v any)

var Logger *jsonlog.Logger

var (
	logFunc LogFunc
	wg      = sync.WaitGroup{}
)

func init() {
	logFunc = func(v any) { log.Default().Print(v) }
}

// Set the function to use for logging on panic recovery
//
// By default, it uses log.Default().Print to log the panic
func SetLogger(f LogFunc) {
	logFunc = f
}

// Setup Background task with with recover, and assures graceful exit on program exit
// when calling WaitAll()
//
// WARNING: This function by itself does not use a separete goroutine. Use `go` to run it concurrently
//
// Technically a little bit slower than using a goroutine, if you know it will not panic and don't care about
// a graceful exit, you should use a goroutine.
func BackgroundTask(task func()) {
	wg.Add(1)
	defer func() {
		defer wg.Done()
		if err := recover(); err != nil {
			if Logger != nil {
				Logger.PrintError(fmt.Errorf("%s", err), nil)
			} else {
				log.Print(err)
			}
		}
	}()

	task()
}

func Background(fn func()) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if Logger != nil {
					Logger.PrintError(fmt.Errorf("%s", err), nil)
				} else {
					log.Print(err)
				}
			}
		}()

		fn()
	}()
}

// Block until all background tasks are finished
func WaitAll() {
	wg.Wait()
}
