export interface Repository {
  id: string;
  name: string;
  fullName: string;
  description?: string;
  private: boolean;
  fork: boolean;
  stars: number;
  forks: number;
  language?: string;
  size: number;
  cloneUrl: string;
  sshUrl: string;
  htmlUrl: string;
  defaultBranch: string;
  createdAt: string;
  updatedAt: string;
  owner: {
    id: string;
    login: string;
    fullName?: string;
    email?: string;
    avatarUrl?: string;
  };
}

// Pagination
export interface RepositoryPage {
  items: Repository[];
  total: number;
  limit: number;
  offset: number;
  hasMore: boolean;
}

export interface RepositoryStats {
  totalCount: number;
  publicCount: number;
  privateCount: number;
  languages: {
    language: string;
    count: number;
  }[];
}
export interface PaginationInput {
  offset?: number;
  limit?: number;
}

