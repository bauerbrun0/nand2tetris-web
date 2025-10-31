import type { ProjectsResponse, Project } from "../../../types/projects";

export async function fetchProjects(
  page: number,
  perPage: number,
): Promise<ProjectsResponse> {
  const res = await fetch(`/api/projects?page=${page}&page_size=${perPage}`);
  if (!res.ok) {
    throw new Error("Failed to fetch projects");
  }
  return await res.json();
}

export async function createProject(
  title: string,
  description: string,
): Promise<Project> {
  const res = await fetch(`/api/projects`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ title, description }),
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(
        ("Failed to create project: " + errorData.error) as string,
      );
    }
    throw new Error("Failed to create project");
  }

  return await res.json();
}

export async function deleteProject(id: number): Promise<Project> {
  const res = await fetch(`/api/projects/${id}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(
        ("Failed to delete project: " + errorData.error) as string,
      );
    }
    throw new Error("Failed to delete project");
  }
  return await res.json();
}

export async function editProject(
  id: number,
  title: string,
  description: string,
): Promise<Project> {
  const res = await fetch(`/api/projects/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ title, description }),
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(("Failed to edit project: " + errorData.error) as string);
    }
    throw new Error("Failed to edit project");
  }

  return await res.json();
}
