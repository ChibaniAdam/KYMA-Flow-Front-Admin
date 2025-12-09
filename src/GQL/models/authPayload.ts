import type { User } from "./user"
export type AuthPayload = {
  token: String
  user: User
}