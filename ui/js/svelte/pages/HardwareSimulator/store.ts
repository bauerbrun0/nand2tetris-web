import { writable, get, type Writable } from "svelte/store";
import type { HardwareSimulatorError } from "./types";

export const progressWASM = writable("READY");
export const progressJS = writable("READY");

export const currentProjectName = writable<string>("First Project");

export const hdls = writable<Record<string, string>>({
  NotChip: `CHIP NotChip {
    IN in;
    OUT out;

    PARTS:
    Nand(a = in, b = in, out = out);
}`,
  AndChip: `CHIP AndChip {
    IN a, b;
    OUT out;

    PARTS:
    Nand(a = a, b = b, out = aNandB);
    NotChip(in = aNandB, out = out);
}`,
});

const initialCurrentHdlFileName = Object.keys(get(hdls)).sort()[0] || "";
export const currentHdlFileName = writable<string>(initialCurrentHdlFileName);

// Create a writable store that syncs with hdls + currentHdlFileName
export const hdl: Writable<string> = (() => {
  const { subscribe, set } = writable("");

  // Keep hdl updated when hdls or currentHdlFileName changes
  const unsubscribeHdls = hdls.subscribe(($hdls) => {
    const currentName = get(currentHdlFileName);
    set($hdls[currentName] || "");
  });

  const unsubscribeCurrent = currentHdlFileName.subscribe(($name) => {
    const currentHdls = get(hdls);
    set(currentHdls[$name] || "");
  });

  return {
    subscribe,

    // When hdl changes, update hdls[currentHdlFileName]
    set(newValue: string) {
      const name = get(currentHdlFileName);
      hdls.update(($hdls) => ({
        ...$hdls,
        [name]: newValue,
      }));
      set(newValue);
    },

    update(fn: (value: string) => string) {
      const currentValue = get(hdl);
      const newValue = fn(currentValue);
      hdl.set(newValue);
    },

    // Optional cleanup if needed
    destroy() {
      unsubscribeHdls();
      unsubscribeCurrent();
    },
  };
})();

export const hardwareSimulatorError = writable<HardwareSimulatorError | null>(
  null,
);
currentHdlFileName.subscribe(() => {
  hardwareSimulatorError.set(null);
});
