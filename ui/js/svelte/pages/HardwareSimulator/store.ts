import { writable } from "svelte/store";

export const progressWASM = writable("READY");
export const progressJS = writable("READY");
