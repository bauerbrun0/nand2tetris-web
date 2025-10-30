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
      name: "DMux",
      hdl: "DMux(in = ,sel = ,a = ,b = );",
      description: "Routes the input to one of two outputs",
    },
    {
      name: "DMux4Way",
      hdl: "DMux4Way(in = ,sel = ,a = ,b = ,c = ,d = );",
      description: "Routes the input to one of four outputs",
    },
    {
      name: "DMux8Way",
      hdl: "DMux8Way(in = ,sel = ,a = ,b = ,c = ,d = ,e = ,f = ,g = ,h = );",
      description: "Routes the input to one of eight outputs",
    },
    {
      name: "And16",
      hdl: "And16(a = ,b = ,out = );",
      description: "16-bit And",
    },
    {
      name: "Or16",
      hdl: "Or16(a = ,b = ,out = );",
      description: "16-bit Or",
    },
    {
      name: "Or8Way",
      hdl: "Or8Way(in = ,out = );",
      description: "8-way Or",
    },
    {
      name: "Mux16",
      hdl: "Mux16(a = ,b = ,sel = ,out = );",
      description: "Selects between two 16-bit inputs",
    },
    {
      name: "Mux4Way16",
      hdl: "Mux4Way16(a = ,b = ,c = ,d = ,sel = ,out = );",
      description: "Selects between four 16-bit inputs",
    },
    {
      name: "Mux8Way16",
      hdl: "Mux8Way16(a = ,b = ,c = ,d = ,e = ,f = ,g = ,h = ,sel = ,out = );",
      description: "Selects between eight 16-bit inputs",
    },
    {
      name: "HalfAdder",
      hdl: "HalfAdder(a = ,b = ,sum = ,carry = );",
      description: "Adds up two bits",
    },
    {
      name: "FullAdder",
      hdl: "FullAdder(a = ,b = ,c = ,sum = ,carry = );",
      description: "Adds up three bits",
    },
    {
      name: "Inc16",
      hdl: "Inc16(in = ,out = );",
      description: "Sets out to in + 1",
    },
    {
      name: "Add16",
      hdl: "Add16(a = ,b = ,out = );",
      description: "Adds up two 16-bit two's complement values",
    },
    {
      name: "Not16",
      hdl: "Not16(in = ,out = );",
      description: "16-bit Not",
    },
    {
      name: "DFF",
      hdl: "DFF(in = ,out = );",
      description: "Data flip-flop gate",
    },
    {
      name: "Bit",
      hdl: "Bit(in = ,load = ,out = );",
      description: "1-bit register",
    },
    {
      name: "Register",
      hdl: "Register(in = ,load = ,out = );",
      description: "16-bit register",
    },
    {
      name: "PC",
      hdl: "PC(in = ,load = ,inc = ,reset = ,out = );",
      description: "Program Counter",
    },
    {
      name: "RAM8",
      hdl: "RAM8(in = ,load = ,address = ,out = );",
      description: "8-word RAM",
    },
    {
      name: "RAM64",
      hdl: "RAM64(in = ,load = ,address = ,out = );",
      description: "64-word RAM",
    },
    {
      name: "RAM512",
      hdl: "RAM512(in = ,load = ,address = ,out = );",
      description: "512-word RAM",
    },
    {
      name: "RAM4K",
      hdl: "RAM4K(in = ,load = ,address = ,out = );",
      description: "4K RAM",
    },
    {
      name: "RAM16K",
      hdl: "RAM16K(in = ,load = ,address = ,out = );",
      description: "16K RAM",
    },
    {
      name: "RAM16K",
      hdl: "RAM16K(in = ,load = ,address = ,out = );",
      description: "16K-word RAM",
    },
    {
      name: "PC",
      hdl: "PC(in = ,load = ,inc = ,reset = ,out = );",
      description: "Program Counter",
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
