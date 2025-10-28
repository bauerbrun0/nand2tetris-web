<script lang="ts">
  import Pagination from "./components/Pagination.svelte";
  import ProjectsList from "./components/ProjectsList.svelte";
  import PlusIcon from "../../../svelte/components/icons/Plus.svelte";
  import { createMutation, createQuery } from "@tanstack/svelte-query";
  import {
    fetchProjects,
    createProject,
    deleteProject,
    editProject,
  } from "./requests";
  import { showToast } from "../../../utils/toast";

  import { currentPage, projectsPerPage } from "./store";
  import type { ProjectsResponse } from "../../../types/projects";
  import Loading from "../../components/Loading.svelte";
  import ProjectsError from "./components/ProjectsError.svelte";
  import NewProjectModal from "./components/NewProjectModal.svelte";
  import DeleteProjectModal from "./components/DeleteProjectModal.svelte";
  import EditProjectModal from "./components/EditProjectModal.svelte";

  const query = createQuery<ProjectsResponse>(() => ({
    queryKey: ["projects", $currentPage, $projectsPerPage],
    queryFn: () => fetchProjects($currentPage, $projectsPerPage),
  }));

  const newProjectMutation = createMutation(() => ({
    mutationFn: (data: { title: string; description: string }) =>
      createProject(data.title, data.description),
    onSuccess: () => {
      query.refetch();
    },
  }));

  const deleteProjectMutation = createMutation(() => ({
    mutationFn: (data: { id: number }) => deleteProject(data.id),
    onSuccess: () => {
      query.refetch();
    },
  }));

  const editProjectMutation = createMutation(() => ({
    mutationFn: (data: { id: number; title: string; description: string }) =>
      editProject(data.id, data.title, data.description),
    onSuccess: () => {
      query.refetch();
    },
  }));

  let newProjectModalTrigger: HTMLButtonElement | null = null;

  async function handleDeleteProject(id: number): Promise<boolean> {
    try {
      await deleteProjectMutation.mutateAsync({
        id,
      });
    } catch (error: unknown) {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
      return false;
    }
    showToast({
      duration: 3000,
      message: "Project deleted successfully.",
      variant: "success",
    });
    return true;
  }

  async function handleEditProject(
    id: number,
    title: string,
    description: string,
  ): Promise<boolean> {
    try {
      await editProjectMutation.mutateAsync({
        id,
        title,
        description,
      });
    } catch (error: unknown) {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
      return false;
    }
    showToast({
      duration: 3000,
      message: "Project edited successfully.",
      variant: "success",
    });
    return true;
  }

  async function handleCreateProject(title: string, description: string) {
    try {
      await newProjectMutation.mutateAsync({
        title,
        description,
      });
    } catch (error: unknown) {
      showToast({
        duration: 3000,
        message: (error as Error).message,
        variant: "error",
      });
      return;
    }
    closeNewProjectModal();
  }

  function closeNewProjectModal() {
    if (newProjectModalTrigger) {
      newProjectModalTrigger.click();
    }
  }
</script>

<div class="inner-container-small py-4">
  <div class="flex items-center justify-between">
    <h1 class="text-2xl font-extrabold md:text-3xl">Your Projects</h1>
    <button
      bind:this={newProjectModalTrigger}
      data-modal-target="new-project-modal"
      data-modal-toggle="new-project-modal"
      class="bg-primary hover:bg-primary-600 flex cursor-pointer items-center rounded-md px-3 py-2 font-medium"
    >
      <PlusIcon classes="inline h-5 w-5 mr-2 stroke-2" />
      New Project
    </button>
    <NewProjectModal
      id="new-project-modal"
      createProject={handleCreateProject}
    />
  </div>
  {#if query.isPending}
    <div class="mt-10 flex justify-center">
      <Loading classes="" />
    </div>
  {/if}
  {#if query.error}
    <div class="mx-auto mt-10 w-full md:w-1/2">
      <ProjectsError message={query.error.message} />
      <button
        onclick={() => query.refetch()}
        class="bg-primary hover:bg-primary-600 mt-4 w-full cursor-pointer items-center rounded-md px-3 py-2 text-center font-medium"
      >
        Retry
      </button>
    </div>
  {/if}
  {#if query.isSuccess}
    {#if query.data.projects.length === 0}
      <p class="text-silver-600 mt-4">
        You have no projects yet. Create one to get started!
      </p>
    {:else}
      <ProjectsList projects={query.data.projects} />
      {#if query.data.totalPages > 1}
        <Pagination />
      {/if}
      <DeleteProjectModal
        id="delete-project-modal"
        deleteProject={handleDeleteProject}
      />
      <EditProjectModal
        id="edit-project-modal"
        editProject={handleEditProject}
      />
    {/if}
  {/if}
</div>
