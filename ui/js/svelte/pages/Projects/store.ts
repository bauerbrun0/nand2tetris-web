import { readable, writable } from "svelte/store";
import type { Project } from "../../../types/projects";

export const currentPage = writable(1);
export const pageCount = writable(3);
export const projectsPerPage = readable(5);

export const clickedProject = writable<Project | null>(null);
