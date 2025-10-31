import { init } from "./init.ts";
import HardwareSimulatorEntry from "../svelte/pages/HardwareSimulator/HardwareSimulatorEntry.svelte";
import { mount } from "svelte";

init();

mount(HardwareSimulatorEntry, {
  target: document.getElementById("svelte-app") as HTMLElement,
});
