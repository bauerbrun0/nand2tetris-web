//go:build js && wasm

package main

import (
	"syscall/js"

	"github.com/bauerbrun0/nand2tetris-web/internal/hardwaresimulator/simulator"
)

var hardwareSimulator *simulator.HardwareSimulator

func main() {
	js.Global().Get("WASM").Get("HardwareSimulator").Set("processHdls", processHdlsWrapper())
	js.Global().Get("WASM").Get("HardwareSimulator").Set("evaluate", evaluateWrapper())
	js.Global().Get("WASM").Get("HardwareSimulator").Set("tick", tickWrapper())
	js.Global().Get("WASM").Get("HardwareSimulator").Set("tock", tockWrapper())
	<-make(chan struct{})
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

func getInputPins() map[string][]bool {
	getInputPins := js.Global().Get("WASM").Get("HardwareSimulator").Get("getInputPins")
	inputPinsJS := getInputPins.Invoke()
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
	setOutputPins := js.Global().Get("WASM").Get("HardwareSimulator").Get("setOutputPins")
	setInternalPins := js.Global().Get("WASM").Get("HardwareSimulator").Get("setInternalPins")

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
	setOutputPins.Invoke(outputPinsJS)

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
	setInternalPins.Invoke(internalPinsJS)
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
