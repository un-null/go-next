import { request } from "./client";
import type { User } from "./types";

export function getUser(token: string, id: string): Promise<User> {
  return request<User>(`/api/users/${id}`, { token });
}
