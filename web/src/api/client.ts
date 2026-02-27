const API_BASE = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export async function apiFetch<T>(
	path: string,
	init?: RequestInit,
): Promise<T> {
	const token = getToken();
	const res = await fetch(`${API_BASE}${path}`, {
		...init,
		headers: {
			"Content-Type": "application/json",
			...(token ? { Authorization: `Bearer ${token}` } : {}),
			...init?.headers,
		},
	});
	if (!res.ok) {
		const body = await res.json().catch(() => ({}));
		throw new Error(body.error ?? `HTTP ${res.status}`);
	}
	return res.json();
}

function getToken(): string | null {
	// Read local storage
	if (typeof window === "undefined") return null;
	return localStorage.getItem("token");
}
