import { progressWASM } from "../store";

export function loadHardwareSimulator() {
  window.WASM = {} as typeof window.WASM;
  window.WASM.HardwareSimulator = {} as typeof window.WASM.HardwareSimulator;
  window.WASM.HardwareSimulator.setProgressWASM = (str) => {
    progressWASM.set(str);
  };

  const go = new Go();
  WebAssembly.instantiateStreaming(
    fetch("/static/wasm/hardware_simulator.wasm"),
    go.importObject,
  ).then((result) => {
    go.run(result.instance);
  });
}
