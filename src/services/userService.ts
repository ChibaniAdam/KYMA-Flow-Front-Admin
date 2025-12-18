import { graphqlRequest } from "./graphqlRequest";
import type { User, CreateUserInput, UpdateUserInput, UserPage, UserFilter, PaginationInput } from "../GQL/models/user";
import type { LoginMutation, LoginMutationVariables, MeQuery, RegisterMutation, RegisterMutationVariables } from "../GQL/apis/apis";



export const login = async (variables: LoginMutationVariables): Promise<LoginMutation> => {
  const mutation = `
    mutation ($uid: String!, $password: String!) {
      login(uid: $uid, password: $password) {
        token
        user {
          uid
          cn
          sn
          givenName
          mail
          department
          uidNumber
          gidNumber
          homeDirectory
          repositories
          dn
        }
      }
    }
  `;
  return graphqlRequest<LoginMutation, LoginMutationVariables>(mutation, variables);
};

export const logout = (navigate: (path: string) => void): void => {
  localStorage.removeItem("token");
  localStorage.removeItem("user");
  navigate("/login");
};

export const register = async (variables: RegisterMutationVariables): Promise<RegisterMutation> => {
  const mutation = `
    mutation ($username: String!, $password: String!, $email: String!, $firstName: String!, $lastName: String!, $userType: String!) {
      register(
        username: $username
        password: $password
        email: $email
        firstName: $firstName
        lastName: $lastName
        userType: $userType
      ) {
        token
        user {
          id
          username
          email
          firstName
          lastName
          userType
        }
      }
    }
  `;
  return graphqlRequest<RegisterMutation, RegisterMutationVariables>(mutation, variables);
};

export async function getUser(uid: string): Promise<User> {
  const query = `
    query ($uid: String!) {
      user(uid: $uid) {
        uid cn sn givenName mail department uidNumber gidNumber homeDirectory repositories dn
      }
    }
  `;
  return graphqlRequest<{ user: User }, { uid: string }>(query, { uid }).then(res => res.user);
}

export const getMe = async (): Promise<MeQuery> => {
  const query = `
    query {
      me {
        uid
        cn
        sn
        givenName
        mail
        department
        uidNumber
        gidNumber
        homeDirectory
        repositories
        dn
      }
    }
  `;
  return graphqlRequest<MeQuery>(query);
};

export async function listUsers(
  filter?: UserFilter,
  pagination?: PaginationInput
): Promise<UserPage> {
  const query = `
    query ($filter: SearchFilterInput, $pagination: PaginationInput) {
      users(filter: $filter, pagination: $pagination) {
        items {
          uid cn sn givenName mail department
          uidNumber gidNumber homeDirectory repositories dn
        }
        total
        page
        limit
        hasNextPage
      }
    }
  `;

  return graphqlRequest<
    { users: UserPage },
    { filter?: UserFilter; pagination?: PaginationInput }
  >(query, { filter, pagination }).then(res => res.users);
}

export async function createUser(input: CreateUserInput): Promise<User> {
  const mutation = `
    mutation ($input: CreateUserInput!) {
      createUser(input: $input) {
        uid cn sn givenName mail department uidNumber gidNumber homeDirectory repositories dn
      }
    }
  `;
  return graphqlRequest<{ createUser: User }, { input: CreateUserInput }>(mutation, { input }).then(res => res.createUser);
}

export async function updateUser(input: UpdateUserInput): Promise<User> {
  const mutation = `
    mutation ($input: UpdateUserInput!) {
      updateUser(input: $input) {
        uid cn sn givenName mail department uidNumber gidNumber homeDirectory repositories dn
      }
    }
  `;
  return graphqlRequest<{ updateUser: User }, { input: UpdateUserInput }>(mutation, { input }).then(res => res.updateUser);
}

export async function deleteUser(uid: string): Promise<boolean> {
  const mutation = `
    mutation ($uid: String!) {
      deleteUser(uid: $uid)
    }
  `;
  return graphqlRequest<{ deleteUser: boolean }, { uid: string }>(mutation, { uid }).then(res => res.deleteUser);
}
