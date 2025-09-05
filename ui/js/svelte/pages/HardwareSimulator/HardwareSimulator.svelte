<script lang="ts">
  import { onMount } from "svelte";
  import { progressWASM } from "./store.ts";
  import Editor from "./components/Editor/Editor.svelte";

  onMount(() => {
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
  });
</script>

<Editor />
