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
  const token = localStorage.getItem('authToken');

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(GRAPHQL_ENDPOINT, {
    method: "POST",
    headers,
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

export const getUser = async (variables: UserQueryVariables): Promise<UserQuery> => {
  const query = `
    query ($uid: String!) {
      user(uid: $uid) {
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
  return graphqlRequest<UserQuery, UserQueryVariables>(query, variables);
};

export const getHealth = async (): Promise<HealthQuery> => {
  const query = `
    query {
      health {
        status
        timestamp
        ldap
      }
    }
  `;
  return graphqlRequest<HealthQuery>(query);
};

// ------------------- Mutations -------------------

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
