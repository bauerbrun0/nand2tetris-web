<script lang="ts">
  let {
    id,
    createProject,
  }: {
    id: string;
    createProject: (title: string, description: string) => Promise<void>;
  } = $props();

  import Modal from "../../../components/modal/Modal.svelte";
  import CloseModalButton from "../../../components/modal/CloseModalButton.svelte";
  import Input from "../../../components/input/Input.svelte";
  import FormSubmitButton from "../../../components/buttons/FormSubmitButton.svelte";
  import { writable } from "svelte/store";

  const projectTitle = writable("");
  const projectDescription = writable("");
  let loading = $state(false);

  async function handleSubmit(e: Event) {
    e.preventDefault();
    loading = true;
    await createProject($projectTitle, $projectDescription);
    loading = false;
  }
</script>

<Modal {id} classes="max-w-[500px]">
  <div class="bg-white-500 dark:bg-silver-900 rounded-lg p-4 sm:p-6 md:p-8">
    <form class="space-y-6" onsubmit={handleSubmit}>
      <div class="flex items-center justify-between">
        <h5 class="text-text dark:text-text-dark w-full text-xl font-medium">
          Create New Project
        </h5>
        <CloseModalButton modalId={id} />
      </div>
      <Input
        required={true}
        value={projectTitle}
        label="Title"
        name="title"
        placeholder="My First Project"
        type="text"
      />
      <Input
        value={projectDescription}
        label="Description"
        name="description"
        placeholder="Description of my project"
        type="textarea"
      />
      <FormSubmitButton text="Create Project" {loading} />
    </form>
  </div>
</Modal>
