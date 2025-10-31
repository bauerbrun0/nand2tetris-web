<script lang="ts">
  import { newChipNameInputFocused, newChipNameInputOpen } from "../../store";

  let { createChip }: { createChip: (name: string) => void } = $props();

  let newChipNameInput: HTMLInputElement = $state(
    null as unknown as HTMLInputElement,
  );
  let newChipFileName: string = $state("");

  $effect(() => {
    if (newChipNameInput) {
      newChipNameInputFocused.set(true);
    }
  });

  newChipNameInputFocused.subscribe((focus) => {
    if (newChipNameInput && focus) {
      newChipNameInput.focus();
    }
  });

  function handleFocusOut() {
    if (!newChipNameInputOpen) return;

    newChipNameInputFocused.set(false);
    newChipNameInputOpen.set(false);
    if (newChipFileName.trim() !== "") {
      createChip(newChipFileName);
    }
    newChipFileName = "";
  }

  function handleKeyDown(e: KeyboardEvent) {
    if (e.key === "Enter") {
      newChipNameInputOpen.set(false);
      newChipNameInputFocused.set(false);

      if (newChipFileName.trim() === "") {
        newChipFileName = "";
        return;
      }

      createChip(newChipFileName);
      newChipFileName = "";
      return;
    }
    if (e.key === "Escape") {
      newChipNameInputFocused.set(false);
      newChipNameInputOpen.set(false);
      newChipFileName = "";
      return;
    }
  }
</script>

{#if $newChipNameInputOpen}
  <button
    class="dark:bg-silver-700 bg-white-600 flex w-full cursor-pointer items-center justify-between overflow-auto"
    onclick={() => newChipNameInputFocused.set(true)}
  >
    <div class="flex w-full py-1 pr-4 pl-6 text-left">
      <input
        bind:this={newChipNameInput}
        onkeydown={handleKeyDown}
        onfocusout={handleFocusOut}
        class=" focus:ring-0 focus:outline-none"
        bind:value={newChipFileName}
        style="width: {Math.min(
          Math.max(newChipFileName.length * 9.6, 10),
          180,
        )}px"
      />
      <span>.hdl</span>
    </div>
  </button>
{/if}
