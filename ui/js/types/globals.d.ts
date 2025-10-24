import type { Pin } from "../svelte/pages/HardwareSimulator/types";

export {};

declare global {
  interface Window {
    WASM: {
      HardwareSimulator: {
        // exported JS functions (called *from Go*)
        setProgressWASM: (str: string) => void;
        getHdls: () => Record<string, string>;
        getCurrentHdlFileName: () => string;
        setHardwareSimulatorError: (error: string) => void;
        setInputPins: (pins: Pin[]) => void;
        setOutputPins: (pins: Pin[]) => void;
        setInternalPins: (pins: Pin[]) => void;

        // exported Go functions (called *from JS*)
        startComputing: (n: number, delayNS: number) => void;
        processHdls: () => void;
      };
    };
  }
}
