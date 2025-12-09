import type { MeQuery,
              UserQuery,
              UserQueryVariables,
              HealthQuery,
              LoginMutation,
              LoginMutationVariables,
              RegisterMutation,
              RegisterMutationVariables, } from "../GQL/apis/apis";
const GRAPHQL_ENDPOINT = import.meta.env.VITE_GRAPHQL_ENDPOINT!;

async function graphqlRequest<T, V = Record<string, any>>(query: string, variables?: V): Promise<T> {
  const res = await fetch(GRAPHQL_ENDPOINT, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ query, variables }),
  });

  const json = await res.json();
  if (json.errors) {
    throw new Error(json.errors.map((err: any) => err.message).join(", "));
  }
  return json.data as T;
}

// ------------------- Queries -------------------

export const getMe = async (): Promise<MeQuery> => {
  const query = `
    query {
      me {
        id
        username
        email
        firstName
        lastName
        userType
      }
    }
  `;
  return graphqlRequest<MeQuery>(query);
};

export const getUser = async (variables: UserQueryVariables): Promise<UserQuery> => {
  const query = `
    query ($username: String!) {
      user(username: $username) {
        id
        username
        email
        firstName
        lastName
        userType
      }
    }
  `;
  return graphqlRequest<UserQuery, UserQueryVariables>(query, variables);
};

export const getHealth = async (): Promise<HealthQuery> => {
  const query = `
    query {
      health
    }
  `;
  return graphqlRequest<HealthQuery>(query);
};

// ------------------- Mutations -------------------

export const login = async (variables: LoginMutationVariables): Promise<LoginMutation> => {
  const mutation = `
    mutation ($username: String!, $password: String!) {
      login(username: $username, password: $password) {
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
  return graphqlRequest<LoginMutation, LoginMutationVariables>(mutation, variables);
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
