import { apiFetch } from "./client";
import type { LoginInput, RegisterInput, AuthResponse } from "@/schemas/user";

export const authApi = {
    login: (data: LoginInput) => 
        apiFetch<AuthResponse>("/v1/users/login", { method: "POST", body: JSON.stringify(data) }),
    register: (data: RegisterInput) => 
        apiFetch<AuthResponse>("/v1/users/register", { method: "POST", body: JSON.stringify(data) }),
};