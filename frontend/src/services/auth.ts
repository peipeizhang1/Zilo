import { httpRequest } from "./http";

export type RegisterReq = {
  email: string;
  password: string;
  nickname?: string;
};

export type LoginReq = {
  email: string;
  password: string;
};

export type LoginResp = {
  accessToken: string;
  refreshToken: string;
  expireAt: number;
};

export function register(payload: RegisterReq) {
  return httpRequest<{ userId: number; email: string }>("/api/v1/auth/register", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export function login(payload: LoginReq) {
  return httpRequest<LoginResp>("/api/v1/auth/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

