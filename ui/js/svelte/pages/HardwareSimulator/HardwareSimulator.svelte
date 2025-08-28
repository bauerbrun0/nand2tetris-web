<script lang="ts">
  import { onMount } from "svelte";
  import { progressWASM, progressJS } from "./store.ts";

  function sqrt() {
    let sum = 0;
    const n = 100_000_000;

    for (let i = 1; i <= n; i++) {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      sum += Math.sqrt(i);
    }
  }

  async function startComputing(n: number, delayNS: number) {
    const start = performance.now(); // start timer
    progressJS.set("STARTED");

    for (let i = 1; i <= n; i++) {
      sqrt();
      const progress = "#".repeat(i);
      progressJS.set(progress);
      await new Promise((resolve) => setTimeout(resolve, delayNS / 1_000_000));
    }
    const elapsed = performance.now() - start;
    progressJS.set(`Done! Runtime: ${elapsed.toFixed(2)} ms`);
  }

  onMount(() => {
    window.WASM = {} as typeof window.WASM;
    window.WASM.setProgressWASM = (str) => {
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

<div class="rounded-xl p-4 shadow-md">
  <p class="mb-2 font-bold">WASM:</p>
  <p class="mb-4">proress: {$progressWASM}</p>
  <p class="mb-2 font-bold">JS</p>
  <p class="mb-2">progress: {$progressJS}</p>
  <button
    class="rounded-lg bg-red-500 px-3 py-1 text-white"
    on:click={() => {
      startComputing(100, 10);
      window.WASM.startComputing(100, 10);
    }}
  >
    Start JS & WASM
  </button>
  <button
    class="rounded-lg bg-red-500 px-3 py-1 text-white"
    on:click={() => {
      startComputing(100, 10);
    }}
  >
    Start JS
  </button>
  <button
    class="rounded-lg bg-red-500 px-3 py-1 text-white"
    on:click={() => {
      window.WASM.startComputing(100, 10);
    }}
  >
    Start WASM
  </button>
</div>
