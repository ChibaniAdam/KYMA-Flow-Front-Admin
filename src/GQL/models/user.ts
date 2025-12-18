export interface User {
  uid: string;
  cn: string;
  sn: string;
  givenName: string;
  mail: string;
  department: string;
  uidNumber: number;
  gidNumber: number;
  homeDirectory: string;
  repositories: string[];
  dn: string;
}

export interface CreateUserInput {
  uid: string;
  cn: string;
  sn: string;
  givenName: string;
  mail: string;
  department: string;
  password: string;
  repositories?: string[];
}

export interface UpdateUserInput {
  uid: string;
  cn?: string;
  sn?: string;
  givenName?: string;
  mail?: string;
  department?: string;
  password?: string;
  repositories?: string[];
}

export interface UserPage {
  items: User[];
  total: number;
  page: number;
  limit: number;
  hasNextPage: boolean;
}

export interface UserFilter {
  department?: string;
  cn?: string;
  mail?: string;
  uid?: string;
}

export interface PaginationInput {
  page: number;
  limit: number;
}
