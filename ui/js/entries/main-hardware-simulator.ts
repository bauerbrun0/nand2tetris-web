import { init } from "./init.ts";
import HardwareSimulator from "../svelte/pages/HardwareSimulator/HardwareSimulator.svelte";
import { mount } from "svelte";

init();

mount(HardwareSimulator, {
  target: document.getElementById("svelte-app") as HTMLElement,
});
