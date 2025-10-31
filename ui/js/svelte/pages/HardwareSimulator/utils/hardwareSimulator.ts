import { get } from "svelte/store";
import {
  currentHdlFileName,
  hdls,
  hardwareSimulatorError,
  inputPins,
  outputPins,
  internalPins,
  simulationSpeed,
  simulationLoopRunning,
  advanceCycle,
  cycleStage,
} from "../store";
import type { Pin } from "../types";

export async function loadHardwareSimulator() {
  window.WASM = {} as typeof window.WASM;
  window.WASM.HardwareSimulator = {} as typeof window.WASM.HardwareSimulator;
  window.WASM.HardwareSimulator.getHdls = () => {
    return get(hdls);
  };
  window.WASM.HardwareSimulator.getCurrentHdlFileName = () => {
    return get(currentHdlFileName) || "";
  };
  window.WASM.HardwareSimulator.setHardwareSimulatorError = (error: string) => {
    hardwareSimulatorError.set({ message: error });
  };
  window.WASM.HardwareSimulator.setInputPins = (pins: Pin[]) => {
    inputPins.set(pins);
  };
  window.WASM.HardwareSimulator.setOutputPins = (pins: Pin[]) => {
    outputPins.set(pins);
  };
  window.WASM.HardwareSimulator.setInternalPins = (pins: Pin[]) => {
    internalPins.set(pins);
  };
  window.WASM.HardwareSimulator.getInputPins = (): Pin[] => {
    return get(inputPins);
  };
  window.WASM.HardwareSimulator.getSimulationDelayMs = (): number => {
    return get(simulationSpeed).delayMs;
  };
  window.WASM.HardwareSimulator.setSimulationLoopRunning = (
    value: boolean,
  ): void => {
    simulationLoopRunning.set(value);
  };
  window.WASM.HardwareSimulator.advanceCycle = (): void => {
    advanceCycle();
  };
  window.WASM.HardwareSimulator.getCycleStage = (): "tick" | "tock" => {
    return get(cycleStage);
  };

  const go = new Go();
  return WebAssembly.instantiateStreaming(
    fetch("/static/wasm/hardware_simulator.wasm"),
    go.importObject,
  ).then((result) => {
    go.run(result.instance);
  });
}
