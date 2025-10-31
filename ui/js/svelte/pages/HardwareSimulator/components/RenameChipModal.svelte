<script lang="ts">
  import Modal from "../../../components/modal/Modal.svelte";
  import CloseModalButton from "../../../components/modal/CloseModalButton.svelte";
  import type { Chip } from "../../../../types/chips";
  import { writable } from "svelte/store";
  import FormSubmitButton from "../../../components/buttons/FormSubmitButton.svelte";
  import Input from "../../../components/input/Input.svelte";

  let {
    id,
    chip,
    renameChip,
  }: {
    id: string;
    chip: Chip | null;
    renameChip: (id: number, name: string) => Promise<boolean>;
  } = $props();

  let name = writable(chip ? chip.name : "");

  $effect(() => {
    if (chip) {
      name.set(chip.name);
    }
  });

  let loading = $state(false);

  let closeButtonEl: HTMLButtonElement = $state(
    undefined as unknown as HTMLButtonElement,
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    const shouldCloseModal = await renameChip((chip as Chip).id, $name);
    loading = false;
    if (shouldCloseModal) {
      closeButtonEl.click();
    }
  }
</script>

<Modal {id} classes="max-w-[300px]">
  <div class="bg-white-500 dark:bg-silver-900 rounded-lg p-4 sm:p-6 md:p-8">
    <form class="space-y-6" onsubmit={handleSubmit}>
      <div class="flex items-center justify-between">
        <h5 class="text-text dark:text-text-dark w-full text-xl font-medium">
          Rename Chip
        </h5>
        <CloseModalButton bind:buttonElement={closeButtonEl} modalId={id} />
      </div>
      <Input
        required={true}
        value={name}
        name="name"
        placeholder="NewChipName"
        type="text"
      />
      <FormSubmitButton text="Rename Chip" {loading} />
    </form>
  </div>
</Modal>
