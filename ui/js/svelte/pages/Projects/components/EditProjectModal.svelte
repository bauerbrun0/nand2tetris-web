<script lang="ts">
  let {
    id,
    editProject,
  }: {
    id: string;
    editProject: (
      id: number,
      title: string,
      description: string,
    ) => Promise<boolean>;
  } = $props();

  import Modal from "../../../components/modal/Modal.svelte";
  import CloseModalButton from "../../../components/modal/CloseModalButton.svelte";
  import FormSubmitButton from "../../../components/buttons/FormSubmitButton.svelte";
  import { clickedProject } from "../store";
  import Input from "../../../components/input/Input.svelte";
  import { writable } from "svelte/store";

  let title = writable($clickedProject ? $clickedProject.title : "");
  let description = writable($clickedProject ? $clickedProject.title : "");

  let loading = $state(false);

  clickedProject.subscribe((value) => {
    if (value) {
      title.set(value.title);
      description.set(value.description);
    }
  });

  let closeButtonEl: HTMLButtonElement = $state(
    undefined as unknown as HTMLButtonElement,
  );

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    const shouldCloseModal = await editProject(
      $clickedProject!.id,
      $title,
      $description,
    );
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
          Edit Project: {$clickedProject ? $clickedProject.title : ""}
        </h5>
        <CloseModalButton bind:buttonElement={closeButtonEl} modalId={id} />
      </div>
      <Input
        required={true}
        value={title}
        label="Title"
        name="title"
        placeholder="My First Project"
        type="text"
      />
      <Input
        value={description}
        label="Description"
        name="description"
        placeholder="Description of my project"
        type="textarea"
      />
      <FormSubmitButton text="Edit Project" {loading} />
    </form>
  </div>
</Modal>
