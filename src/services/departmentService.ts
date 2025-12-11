import { graphqlRequest } from "./graphqlRequest";
import type { Department, CreateDepartmentInput } from "../GQL/models/department";

export async function getDepartment(ou: string): Promise<Department> {
  const query = `
    query ($ou: String!) {
      department(ou: $ou) {
        ou description manager members repositories dn
      }
    }
  `;
  return graphqlRequest<{ department: Department }, { ou: string }>(query, { ou }).then(res => res.department);
}

export async function listDepartments(): Promise<Department[]> {
  const query = `
    query {
      departments {
        ou description manager members repositories dn
      }
    }
  `;
  return graphqlRequest<{ departments: Department[] }>(query).then(res => res.departments);
}

export async function createDepartment(input: CreateDepartmentInput): Promise<Department> {
  const mutation = `
    mutation ($input: CreateDepartmentInput!) {
      createDepartment(input: $input) {
        ou description manager members repositories dn
      }
    }
  `;
  return graphqlRequest<{ createDepartment: Department }, { input: CreateDepartmentInput }>(mutation, { input }).then(res => res.createDepartment);
}

export async function deleteDepartment(ou: string): Promise<boolean> {
  const mutation = `
    mutation ($ou: String!) {
      deleteDepartment(ou: $ou)
    }
  `;
  return graphqlRequest<{ deleteDepartment: boolean }, { ou: string }>(mutation, { ou }).then(res => res.deleteDepartment);
}
