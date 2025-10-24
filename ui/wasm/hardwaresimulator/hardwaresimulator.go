//go:build js && wasm

package main

import (
	"context"
	"syscall/js"
	"time"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/simulator"
)

var jsFuncs map[string]js.Value

var hardwareSimulator *simulator.HardwareSimulator
var cancelSimulationLoop context.CancelFunc

func main() {
	hardwareSimulatorJsObject := js.Global().Get("WASM").Get("HardwareSimulator")

	// exporting go functions to javascript
	hardwareSimulatorJsObject.Set("processHdls", processHdlsWrapper())
	hardwareSimulatorJsObject.Set("evaluate", evaluateWrapper())
	hardwareSimulatorJsObject.Set("tick", tickWrapper())
	hardwareSimulatorJsObject.Set("tock", tockWrapper())
	hardwareSimulatorJsObject.Set("startSimulationLoop", startSimulationLoopWrapper())
	hardwareSimulatorJsObject.Set("stopSimulationLoop", stopSimulationLoopWrapper())

	// getting js functions from javascript
	jsFuncs = make(map[string]js.Value)

	jsFuncs["getInputPins"] = hardwareSimulatorJsObject.Get("getInputPins")
	jsFuncs["setInputPins"] = hardwareSimulatorJsObject.Get("setInputPins")
	jsFuncs["setOutputPins"] = hardwareSimulatorJsObject.Get("setOutputPins")
	jsFuncs["setInternalPins"] = hardwareSimulatorJsObject.Get("setInternalPins")
	<-make(chan struct{})
}

func startSimulationLoop() {
	if cancelSimulationLoop != nil {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancelSimulationLoop = cancel

	js.Global().Get("WASM").Get("HardwareSimulator").Get("setSimulationLoopRunning").Invoke(js.ValueOf(true))
	delayMs := js.Global().Get("WASM").Get("HardwareSimulator").Get("getSimulationDelayMs").Invoke().Int()
	advanceCycle := js.Global().Get("WASM").Get("HardwareSimulator").Get("advanceCycle")
	startingCycleStage := js.Global().Get("WASM").Get("HardwareSimulator").Get("getCycleStage").Invoke().String()

	timer := time.NewTimer(0) // actual duration will be set in waitOrCancel

	defer func() { cancelSimulationLoop = nil }() // clear cancel function when done, allowing to start again

	if startingCycleStage == "tock" {
		tock()
		advanceCycle.Invoke()
		if !waitOrCancel(ctx, timer, delayMs) {
			return
		}
	}

	for {
		tick()
		advanceCycle.Invoke()
		if !waitOrCancel(ctx, timer, delayMs) {
			return
		}

		tock()
		advanceCycle.Invoke()
		if !waitOrCancel(ctx, timer, delayMs) {
			return
		}
	}
}

func waitOrCancel(ctx context.Context, timer *time.Timer, delayMs int) bool {
	timer.Reset(time.Duration(delayMs) * time.Millisecond)

	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func stopSimulationLoop() {
	if cancelSimulationLoop != nil {
		cancelSimulationLoop()
		js.Global().Get("WASM").Get("HardwareSimulator").Get("setSimulationLoopRunning").Invoke(js.ValueOf(false))
	}
}

func tick() {
	inputs := getInputPins()
	outputPins, internalPins := hardwareSimulator.Tick(inputs)
	setOutputAndInternalPins(outputPins, internalPins)
}

func tock() {
	inputs := getInputPins()
	outputPins, internalPins := hardwareSimulator.Tock(inputs)
	setOutputAndInternalPins(outputPins, internalPins)
}

func evaluate() {
	inputs := getInputPins()
	outputPins, internalPins := hardwareSimulator.Evaluate(inputs)
	setOutputAndInternalPins(outputPins, internalPins)
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

func evaluateWrapper() js.Func {
	evaluateFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go evaluate()
		return nil
	})
	return evaluateFunc
}

func tickWrapper() js.Func {
	tickFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go tick()
		return nil
	})
	return tickFunc
}

func tockWrapper() js.Func {
	tockFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go tock()
		return nil
	})
	return tockFunc
}

func stopSimulationLoopWrapper() js.Func {
	stopSimulationLoopFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go stopSimulationLoop()
		return nil
	})
	return stopSimulationLoopFunc
}

func startSimulationLoopWrapper() js.Func {
	startSimulationLoopFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 0 {
			return "Invalid no of arguments passed"
		}
		go startSimulationLoop()
		return nil
	})
	return startSimulationLoopFunc
}

func getInputPins() map[string][]bool {
	inputPinsJS := jsFuncs["getInputPins"].Invoke()
	inputs := make(map[string][]bool)
	length := inputPinsJS.Length()
	for i := 0; i < length; i++ {
		pinObj := inputPinsJS.Index(i)
		pinName := pinObj.Get("name").String()
		bitsJS := pinObj.Get("bits")
		bitsLength := bitsJS.Length()
		bits := make([]bool, bitsLength)
		for j := 0; j < bitsLength; j++ {
			bits[j] = bitsJS.Index(j).Bool()
		}
		inputs[pinName] = bits
	}
	return inputs
}

func setOutputAndInternalPins(outputPins map[string][]bool, internalPins map[string][]bool) {
	outputPinsJS := js.Global().Get("Array").New()
	for outputName, outputBits := range outputPins {
		obj := js.Global().Get("Object").New()
		obj.Set("name", outputName)

		goSlice := make([]any, len(outputBits))
		for i, bit := range outputBits {
			goSlice[i] = bit
		}
		obj.Set("bits", goSlice)

		outputPinsJS.Call("push", obj)
	}
	jsFuncs["setOutputPins"].Invoke(outputPinsJS)

	internalPinsJS := js.Global().Get("Array").New()
	for internalName, internalBits := range internalPins {
		obj := js.Global().Get("Object").New()
		obj.Set("name", internalName)

		goSlice := make([]any, len(internalBits))
		for i, bit := range internalBits {
			goSlice[i] = bit
		}

		obj.Set("bits", goSlice)

		internalPinsJS.Call("push", obj)
	}
	jsFuncs["setInternalPins"].Invoke(internalPinsJS)
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
