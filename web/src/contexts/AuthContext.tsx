"use client";

import type { User } from "@/schemas/user";
import { useQueryClient } from "@tanstack/react-query";
import { useContext, useEffect, useState, createContext } from "react";

interface AuthContextType {
	user: User | null;
	token: string | null;
	isLoading: boolean;
	login: (token: string, user: User) => void;
	logout: () => void;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
	const [user, setUser] = useState<User | null>(null);
	const [token, setToken] = useState<string | null>(null);
	const [isLoading, setIsLoading] = useState<boolean>(true);
	const queryClient = useQueryClient();

	useEffect(() => {
		// Read the token and user
		const userInStorage = localStorage.getItem("user");
		const tokenInStorage = localStorage.getItem("token");
		if (userInStorage && tokenInStorage) {
			setToken(tokenInStorage);
			setUser(JSON.parse(userInStorage));
		}

		// Hydrate state
		setIsLoading(false);
		//
	}, []);

	function login(token: string, user: User) {
		localStorage.setItem("token", token);
		localStorage.setItem("user", JSON.stringify(user));
		setToken(token);
		setUser(user);
	}

	function logout() {
		localStorage.removeItem("token");
		localStorage.removeItem("user");
		setToken(null);
		setUser(null);
		queryClient.clear();
	}

	return (
		<AuthContext.Provider value={{ user, token, isLoading, login, logout }}>
			{children}
		</AuthContext.Provider>
	);
}

export function useAuth() {
	const context = useContext(AuthContext);
	if (!context) {
		throw new Error("useAuth must be within an AuthProvider");
	}
	return context;
}
