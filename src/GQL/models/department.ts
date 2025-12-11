export interface Department {
  ou: string;
  description?: string;
  manager?: string;
  members: string[];
  repositories: string[];
  dn: string;
}

export interface CreateDepartmentInput {
  ou: string;
  description?: string;
  manager?: string;
  repositories?: string[];
}
