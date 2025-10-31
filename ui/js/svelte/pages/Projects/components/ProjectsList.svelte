<script lang="ts">
  import type { Project } from "../../../../types/projects";
  import ProjectCard from "./ProjectCard.svelte";
  import { initModals, initTooltips } from "flowbite";

  $effect(() => {
    if (projects) {
      initModals();
      initTooltips();
    }
  });

  let {
    projects,
  }: {
    projects: Project[];
  } = $props();
</script>

<div class="my-4 flex w-full flex-col gap-2">
  {#each projects as project (project.id)}
    <ProjectCard {project} />
  {/each}
</div>
