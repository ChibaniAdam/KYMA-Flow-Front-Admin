import { graphqlRequest } from "./graphqlRequest";
import type { Group } from "../GQL/models/group";

export async function getGroup(cn: string): Promise<Group> {
  const query = `
    query ($cn: String!) {
      group(cn: $cn) {
        cn gidNumber members dn
      }
    }
  `;
  return graphqlRequest<{ group: Group }, { cn: string }>(query, { cn }).then(res => res.group);
}

export async function createGroup(cn: string, description?: string): Promise<Group> {
  const mutation = `
    mutation ($cn: String!, $description: String) {
      createGroup(cn: $cn, description: $description) {
        cn gidNumber members dn
      }
    }
  `;
  return graphqlRequest<{ createGroup: Group }, { cn: string; description?: string }>(mutation, { cn, description }).then(res => res.createGroup);
}

export async function addUserToGroup(uid: string, groupCn: string): Promise<boolean> {
  const mutation = `
    mutation ($uid: String!, $groupCn: String!) {
      addUserToGroup(uid: $uid, groupCn: $groupCn)
    }
  `;
  return graphqlRequest<{ addUserToGroup: boolean }, { uid: string; groupCn: string }>(mutation, { uid, groupCn }).then(res => res.addUserToGroup);
}
