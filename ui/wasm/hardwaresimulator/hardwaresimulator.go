//go:build js && wasm

package main

import (
	"fmt"
	"math"
	"strings"
	"syscall/js"
	"time"
)

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Get("WASM").Set("startComputing", startComputingWrapper())
	<-make(chan struct{})
}

func sqrt() {
	sum := 0.0
	n := 100_000_000

	for i := 1; i <= n; i++ {
		sum += math.Sqrt(float64(i))
	}
}

func startComputing(n int, delayNS int) {
	start := time.Now()
	setProgress := js.Global().Get("WASM").Get("setProgressWASM")
	setProgress.Invoke("STARTED")

	for i := 1; i <= n; i++ {
		sqrt()
		progress := strings.Repeat("#", i)
		setProgress.Invoke(progress)
		time.Sleep(time.Duration(delayNS) * time.Nanosecond)
	}

	elapsed := time.Since(start).Milliseconds()
	setProgress.Invoke(fmt.Sprintf("Done! Runtime: %d ms", elapsed))
}

func startComputingWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}
		n := args[0].Int()
		delayNS := args[1].Int()
		go startComputing(n, delayNS)
		return nil
	})
	return jsonFunc
}
