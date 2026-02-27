"use client";

import { useMutation } from "@tanstack/react-query";
import { authApi } from "@/api/auth";
import { useAuth } from "@/contexts/AuthContext";
import { useRouter } from "next/navigation";
import type { LoginInput, RegisterInput, AuthResponse } from "@/schemas/user";

export function useLogin() {
  const auth = useAuth();
  const router = useRouter();
  return useMutation({
    mutationFn: (data: LoginInput) => authApi.login(data),
    onSuccess: (data: AuthResponse) => {
      auth.login(data.token, data.user);
      router.push("/");
    },
  });
}

export function useRegister() {
  const auth = useAuth();
  const router = useRouter();
  return useMutation({
    mutationFn: (data: RegisterInput) => authApi.register(data),
    onSuccess: (data: AuthResponse) => {
      auth.login(data.token, data.user);
      router.push("/");
    },
  });
}
