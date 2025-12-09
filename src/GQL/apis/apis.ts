import type { User } from "../models/user";
import type { AuthPayload } from "../models/authPayload";

// ------------------- Queries -------------------

export type MeQuery = {
  me: User | null;
};

export type UserQuery = {
  user: User | null;
};

export type UserQueryVariables = {
  username: string;
};

export type HealthQuery = {
  health: string;
};

// ------------------- Mutations -------------------

export type LoginMutation = {
  login: AuthPayload;
};

export type LoginMutationVariables = {
  username: string;
  password: string;
};

export type RegisterMutation = {
  register: AuthPayload;
};

export type RegisterMutationVariables = {
  username: string;
  password: string;
  email: string;
  firstName: string;
  lastName: string;
  userType: string;
};
