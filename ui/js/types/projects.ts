export type Project = {
  id: number;
  slug: string;
  userId: number;
  title: string;
  description: string;
  created: string;
  updated: string;
};

export type ProjectsResponse = {
  projects: Project[];
  totalCount: number;
  page: number;
  pageSize: number;
  totalPages: number;
};
