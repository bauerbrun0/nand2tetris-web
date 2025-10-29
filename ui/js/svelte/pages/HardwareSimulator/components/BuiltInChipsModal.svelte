<script lang="ts">
  import Modal from "../../../components/modal/Modal.svelte";
  import CloseModalButton from "../../../components/modal/CloseModalButton.svelte";
  import ClipboardIcon from "../../../components/icons/Clipboard.svelte";
  import CheckIcon from "../../../components/icons/Check.svelte";

  let { id }: { id: string } = $props();

  let copiedChip = $state<string | null>(null);

  type BuiltInChip = {
    name: string;
    hdl: string;
    description: string;
  };

  const builtInChips: BuiltInChip[] = [
    {
      name: "And",
      hdl: "And(a = ,b = ,out = );",
      description: "And gate",
    },
    {
      name: "Or",
      hdl: "Or(a = ,b = ,out = );",
      description: "Or gate",
    },
    {
      name: "Not",
      hdl: "Not(in = ,out = );",
      description: "Not gate",
    },
    {
      name: "Xor",
      hdl: "Xor(a = ,b = ,out = );",
      description: "Xor gate",
    },
    {
      name: "Nand",
      hdl: "Nand(a = ,b = ,out = );",
      description: "Nand gate",
    },
    {
      name: "Mux",
      hdl: "Mux(a = ,b = ,sel = ,out = );",
      description: "Selects between two inputs",
    },
    {
      name: "DFF",
      hdl: "DFF(in = ,out = );",
      description: "Data flip-flop gate",
    },
    {
      name: "RAM64",
      hdl: "RAM64(in = ,load = ,address = ,out = );",
      description: "64-word RAM",
    },
  ];
</script>

<Modal {id} classes="max-w-[calc(90%)] h-[calc(90%)]">
  <div
    class="bg-white-500 dark:bg-silver-900 h-full rounded-lg p-4 sm:p-6 md:p-8"
  >
    <div class="space-y-6">
      <div class="flex items-center justify-between">
        <h5 class="w-full text-xl font-medium">Built-In Chips</h5>
        <CloseModalButton modalId={id} />
      </div>
      <div>
        <p>
          These are the built-in chips currently available in the Hardware
          Simulator. You can use these chips in your HDL designs without needing
          to define them yourself.
        </p>
        <div class="mt-4 max-h-[70vh] space-y-6 overflow-y-auto">
          {#each builtInChips as chip (chip.name)}
            <div
              class="dark:border-silver-700 border-white-800 rounded-lg border p-4"
            >
              <h6 class="mb-2 font-semibold">{chip.name}</h6>
              <p class="mb-2 text-sm">
                {chip.description}
              </p>
              <div class="relative">
                <button
                  class={`
                    hover:bg-white-500 dark:bg-silver-700 dark:hover:bg-silver-600 absolute top-1/2 right-2 -translate-y-1/2
                    cursor-pointer rounded bg-white p-2
                  `}
                  onclick={() => {
                    navigator.clipboard.writeText(chip.hdl);
                    copiedChip = chip.name;
                    setTimeout(() => {
                      copiedChip = null;
                    }, 2000);
                  }}
                >
                  {#if copiedChip === chip.name}
                    <CheckIcon classes="h-4 w-4" />
                  {:else}
                    <ClipboardIcon classes="h-4 w-4" />
                  {/if}
                </button>
                <pre
                  class="dark:bg-silver-800 bg-white-100 overflow-x-auto rounded p-4 text-sm select-text">{chip.hdl}</pre>
              </div>
            </div>
          {/each}
        </div>
      </div>
    </div>
  </div>
</Modal>
