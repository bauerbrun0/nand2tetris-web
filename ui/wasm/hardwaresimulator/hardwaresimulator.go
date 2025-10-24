//go:build js && wasm

package main

import (
	"fmt"
	"math"
	"strings"
	"syscall/js"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/simulator"
)

var hardwareSimulator *simulator.HardwareSimulator

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Get("WASM").Get("HardwareSimulator").Set("startComputing", startComputingWrapper())
	js.Global().Get("WASM").Get("HardwareSimulator").Set("processHdls", processHdlsWrapper())
	<-make(chan struct{})
}

func sqrt() {
	sum := 0.0
	n := 100_000_000

	for i := 1; i <= n; i++ {
		sum += math.Sqrt(float64(i))
	}
}

func processHdls() {
	hardwareSimulator = simulator.New()
	hardwareSimulatorJSFuncs := js.Global().Get("WASM").Get("HardwareSimulator")
	getHdls := hardwareSimulatorJSFuncs.Get("getHdls")
	getCurrentHdlFileName := hardwareSimulatorJSFuncs.Get("getCurrentHdlFileName")
	setInputPins := hardwareSimulatorJSFuncs.Get("setInputPins")
	setOutputPins := hardwareSimulatorJSFuncs.Get("setOutputPins")
	setInternalPins := hardwareSimulatorJSFuncs.Get("setInternalPins")
	object := js.Global().Get("Object")

	hdls := JSValueToMap(getHdls.Invoke())
	currentHdlFileName := getCurrentHdlFileName.Invoke().String()

	hardwareSimulator.SetChipHDLs(hdls)
	inputs, outputs, internals, err := hardwareSimulator.Process(currentHdlFileName)
	if err != nil {
		hardwareSimulatorJSFuncs.Get("setHardwareSimulatorError").Invoke(err.Error())
		return
	}

	inputPins := js.Global().Get("Array").New()
	for inputName, inputWidth := range inputs {
		obj := object.New()
		obj.Set("name", inputName)

		goSlice := make([]any, inputWidth)
		for i := range goSlice {
			goSlice[i] = false
		}
		obj.Set("bits", goSlice)

		inputPins.Call("push", obj)
	}
	setInputPins.Invoke(inputPins)

	outputPins := js.Global().Get("Array").New()
	for outputName, outputWidth := range outputs {
		obj := object.New()
		obj.Set("name", outputName)

		goSlice := make([]any, outputWidth)
		for i := range goSlice {
			goSlice[i] = false
		}
		obj.Set("bits", goSlice)

		outputPins.Call("push", obj)
	}
	setOutputPins.Invoke(outputPins)

	internalPins := js.Global().Get("Array").New()
	for internalName, internalWidth := range internals {
		obj := object.New()
		obj.Set("name", internalName)

		goSlice := make([]any, internalWidth)
		for i := range goSlice {
			goSlice[i] = false
		}
		obj.Set("bits", goSlice)

		internalPins.Call("push", obj)
	}
	setInternalPins.Invoke(internalPins)

}

func startComputing(n int, delayNS int) {
	start := time.Now()
	setProgress := js.Global().Get("WASM").Get("HardwareSimulator").Get("setProgressWASM")
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

func processHdlsWrapper() js.Func {
	processHdlsFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go processHdls()
		return nil
	})
	return processHdlsFunc
}

func JSValueToMap(v js.Value) map[string]string {
	result := make(map[string]string)

	// Get all keys of the JS object
	keys := js.Global().Get("Object").Call("keys", v)
	length := keys.Length()

	for i := 0; i < length; i++ {
		key := keys.Index(i).String()
		value := v.Get(key).String()
		result[key] = value
	}

	return result
}
