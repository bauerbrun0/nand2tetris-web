import { writable, get, type Writable } from "svelte/store";
import type { HardwareSimulatorError, Pin, SimulationSpeed } from "./types";
import { simulationSpeeds } from "./utils/simulation";

export const currentProjectName = writable<string>("");

export const hdls = writable<Record<string, string>>({});

export const currentHdlFileName = writable<string | null>(null);

// Create a writable store that syncs with hdls + currentHdlFileName
export const hdl: Writable<string | null> = (() => {
  const { subscribe, set } = writable<string | null>("");

  // Keep hdl updated when hdls or currentHdlFileName changes
  const unsubscribeHdls = hdls.subscribe((hdls) => {
    const currentName = get(currentHdlFileName);
    if (currentName === null) {
      set(null);
      return;
    }
    set(hdls[currentName] || "");
  });

  const unsubscribeCurrent = currentHdlFileName.subscribe((name) => {
    const currentHdls = get(hdls);
    if (name === null) {
      set(null);
      return;
    }
    set(currentHdls[name] || "");
  });

  return {
    subscribe,

    // When hdl changes, update hdls[currentHdlFileName]
    set(newValue: string) {
      const name = get(currentHdlFileName);
      if (name === null) {
        return;
      }

      if (get(hdls)[name] === undefined) {
        return;
      }

      hdls.update(($hdls) => ({
        ...$hdls,
        [name]: newValue,
      }));
      set(newValue);
    },

    update(fn: (value: string) => string) {
      const currentValue = get(hdl);
      if (currentValue === null) {
        return;
      }
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

export const simulationSpeed = writable<SimulationSpeed>(simulationSpeeds[0]);

export const simulationLoopRunning = writable(false);

export const simulationRunning = writable(false);

export const inputPins = writable<Pin[]>([]);
export const outputPins = writable<Pin[]>([]);
export const internalPins = writable<Pin[]>([]);

export const cycleCount = writable<number>(1);
export const cycleStage = writable<"tick" | "tock">("tick");

export function advanceCycle() {
  const currentStage = get(cycleStage);
  if (currentStage === "tock") {
    cycleStage.set("tick");
    cycleCount.update((n) => n + 1);
  } else {
    cycleStage.set("tock");
  }
}

export function resetCycle() {
  cycleCount.set(1);
  cycleStage.set("tick");
}

export const newChipNameInputFocused = writable<boolean>(false);
export const newChipNameInputOpen = writable<boolean>(false);

export const rightClickedChipName = writable<string | null>(null);
