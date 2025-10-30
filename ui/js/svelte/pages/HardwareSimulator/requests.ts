import type { Chip } from "../../../types/chips";
import type { Project } from "../../../types/projects";

export async function fetchProjectBySlug(slug: string): Promise<Project> {
  const res = await fetch(`/api/projects/${slug}/by-slug`);
  if (!res.ok) {
    throw new Error("Failed to fetch projects");
  }
  return await res.json();
}

export async function fetchProjectChips(id: number): Promise<Chip[]> {
  const res = await fetch(`/api/projects/${id}/chips`);
  if (!res.ok) {
    throw new Error("Failed to fetch project chips");
  }
  return await res.json();
}

export async function createChip(
  name: string,
  projectId: number,
): Promise<Chip> {
  const res = await fetch(`/api/projects/${projectId}/chips`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(("Failed to create chip: " + errorData.error) as string);
    }
    throw new Error("Failed to create chip");
  }

  return await res.json();
}

export async function updateChipHdl(
  projectId: number,
  id: number,
  hdl: string,
): Promise<Chip> {
  const res = await fetch(`/api/projects/${projectId}/chips/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ hdl }),
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(("Failed to update chip: " + errorData.error) as string);
    }
    throw new Error("Failed to update chip");
  }

  return await res.json();
}

export async function updateChipName(
  projectId: number,
  id: number,
  name: string,
): Promise<Chip> {
  const res = await fetch(`/api/projects/${projectId}/chips/${id}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ name }),
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(("Failed to rename chip: " + errorData.error) as string);
    }
    throw new Error("Failed to rename chip");
  }

  return await res.json();
}

export async function deleteChipRequest(
  projectId: number,
  id: number,
): Promise<Chip> {
  const res = await fetch(`/api/projects/${projectId}/chips/${id}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    const errorData = await res.json();
    if (errorData && errorData.error && typeof errorData.error === "string") {
      throw new Error(("Failed to delete chip: " + errorData.error) as string);
    }
    throw new Error("Failed to delete chip");
  }

  return await res.json();
}
