export function getProjectSlug(): string {
  let slug = "";
  if (typeof window !== "undefined") {
    const parts = window.location.pathname.split("/");
    slug = parts[2] || ""; // ['', 'projects', 'some-slug']
  }
  return slug;
}
