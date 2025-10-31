import { init } from "./init.ts";
import ProjectsEntry from "../svelte/pages/Projects/ProjectsEntry.svelte";
import { mount } from "svelte";

init();

mount(ProjectsEntry, {
  target: document.getElementById("svelte-app") as HTMLElement,
});
