import { request } from "./client";
import type { AuthResponse, SignupResponse, User } from "./types";

export interface LoginInput {
  email: string;
  password: string;
}

export interface SignupInput {
  name: string;
  email: string;
  password: string;
}

export function login(input: LoginInput): Promise<AuthResponse> {
  return request<AuthResponse>("/api/login", { method: "POST", body: input });
}

export async function signup(input: SignupInput): Promise<User> {
  const response = await request<SignupResponse>("/api/signup", {
    method: "POST",
    body: input,
  });
  return response.user;
}
