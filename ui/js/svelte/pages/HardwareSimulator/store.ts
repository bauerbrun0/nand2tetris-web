import { writable } from "svelte/store";
import type { EditorError } from "../../../types/index.ts";

export const progressWASM = writable("READY");
export const progressJS = writable("READY");

export const editorErrors = writable<EditorError[]>([]);
