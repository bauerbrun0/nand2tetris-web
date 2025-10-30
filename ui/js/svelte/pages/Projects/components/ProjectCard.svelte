<script lang="ts">
  import type { Project } from "../../../../types/projects";
  import TrashIcon from "../../../../svelte/components/icons/Trash.svelte";
  import EditIcon from "../../../../svelte/components/icons/Edit.svelte";
  import { clickedProject } from "../store";
  import SaveIcon from "../../../components/icons/Save.svelte";

  let {
    project,
  }: {
    project: Project;
  } = $props();

  function handleClick(event: MouseEvent) {
    clickedProject.set(project);
    event.stopPropagation();
    event.preventDefault();
  }
</script>

<a
  href={`/projects/${project.slug}`}
  class={`
       dark:bg-silver-900 dark:hover:bg-silver-800 hover:bg-white-500 flex w-full items-center
       justify-between rounded-lg bg-white p-6
    `}
>
  <div>
    <h5 class="mb-2 text-2xl font-bold tracking-tight">
      {project.title}
    </h5>
    <p class="font-normal">
      {project.description}
    </p>
    <div class="mt-4 flex gap-4 text-sm">
      <div
        data-tooltip-target={`tooltip-last-updated-${project.id}`}
        class="flex items-center gap-2"
      >
        <SaveIcon classes="h-4 w-4 " />
        <span>{new Date(project.updated).toLocaleString()}</span>
      </div>
      <div
        id={`tooltip-last-updated-${project.id}`}
        role="tooltip"
        class="tooltip dark:bg-silver-900 bg-silver-100 invisible absolute z-10 inline-block rounded-lg px-3 py-2 text-sm opacity-0 shadow-xs"
      >
        Last Updated on {new Date(project.updated).toLocaleString()}
        <div class="tooltip-arrow" data-popper-arrow></div>
      </div>
    </div>
  </div>
  <div class="flex items-center gap-3">
    <button
      onclick={handleClick}
      data-modal-target="edit-project-modal"
      data-modal-toggle="edit-project-modal"
      class="hover:dark:bg-silver-700 hover:bg-white-900 cursor-pointer rounded-md px-3 py-2"
    >
      <EditIcon classes="h-5 w-5" />
    </button>
    <button
      onclick={handleClick}
      data-modal-target="delete-project-modal"
      data-modal-toggle="delete-project-modal"
      class="hover:dark:bg-silver-700 hover:bg-white-900 stroke-text dark:stroke-text-dark cursor-pointer rounded-md px-3 py-2 hover:stroke-red-500"
    >
      <TrashIcon classes="h-5 w-5 stroke-[1px]" />
    </button>
  </div>
</a>
