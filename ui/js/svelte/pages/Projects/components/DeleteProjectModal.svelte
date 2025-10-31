<script lang="ts">
  let {
    id,
    deleteProject,
  }: {
    id: string;
    deleteProject: (id: number) => Promise<boolean>;
  } = $props();

  import Modal from "../../../components/modal/Modal.svelte";
  import CloseModalButton from "../../../components/modal/CloseModalButton.svelte";
  import FormSubmitButton from "../../../components/buttons/FormSubmitButton.svelte";
  import { clickedProject } from "../store";

  let loading = $state(false);

  let closeButtonEl: HTMLButtonElement = $state(
    undefined as unknown as HTMLButtonElement,
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    const shouldCloseModal = await deleteProject($clickedProject!.id);
    loading = false;
    if (shouldCloseModal) {
      closeButtonEl.click();
    }
  }
</script>

<Modal {id} classes="max-w-[500px]">
  <div class="bg-white-500 dark:bg-silver-900 rounded-lg p-4 sm:p-6 md:p-8">
    <form class="space-y-6" onsubmit={handleSubmit}>
      <div class="flex items-center justify-between">
        <h5 class="text-text dark:text-text-dark w-full text-xl font-medium">
          Delete Project: {$clickedProject ? $clickedProject.title : ""}
        </h5>
        <CloseModalButton bind:buttonElement={closeButtonEl} modalId={id} />
      </div>
      <div class="text-center">
        Are you sure you want to delete this project? This action cannot be
        undone.
      </div>
      <FormSubmitButton
        text="Delete Project"
        {loading}
        classes={`bg-red-700
          hover:bg-red-600 focus:ring-red-300 dark:hover:bg-red-800 dark:focus:ring-red-800
          disabled:bg-red-400
	    `}
      />
    </form>
  </div>
</Modal>
