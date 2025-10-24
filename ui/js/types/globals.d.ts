import type { Pin } from "../svelte/pages/HardwareSimulator/types";

export {};

declare global {
  interface Window {
    WASM: {
      HardwareSimulator: {
        // exported JS functions (called *from Go*)
        getHdls: () => Record<string, string>;
        getCurrentHdlFileName: () => string;
        setHardwareSimulatorError: (error: string) => void;
        setInputPins: (pins: Pin[]) => void;
        setOutputPins: (pins: Pin[]) => void;
        setInternalPins: (pins: Pin[]) => void;
        getInputPins: () => Pin[];
        getSimulationDelayMs: () => number;
        setSimulationLoopRunning: (running: boolean) => void;
        advanceCycle: () => void;
        getCycleStage: () => "tick" | "tock";

        // exported Go functions (called *from JS*)
        startComputing: (n: number, delayNS: number) => void;
        processHdls: () => void;
        evaluate: () => void;
        tick: () => void;
        tock: () => void;
        startSimulationLoop: () => void;
        stopSimulationLoop: () => void;
      };
    };
  }
}
