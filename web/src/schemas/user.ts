import { z } from "zod";

export const PERM_READ = 1;
export const PERM_WRITE = 2;
export const PERM_ADMIN = 4;

export function hasPermission(userPerms: number, required: number): boolean {
	return (userPerms & required) === required;
}

const passwordSchema = z
	.string()
	.min(8, "Password must be at least 8 characters")
	.regex(/[A-Z]/, "Password must contain at least one uppercase letter")
	.regex(/[a-z]/, "Password must contain at least one lowercase letter")
	.regex(/[0-9]/, "Password must contain at least one number")
	.regex(/[!@#$%^&*]/, "Password must contain at least one special character");

export const UserSchema = z.object({
	id: z.uuid(),
	email: z.email(),
	name: z.string().min(1).max(255),
	permissions: z.number().int().default(3),
	created_at: z.string(),
	updated_at: z.string(),
});

export const LoginInputSchema = z.object({
	email: z.email().min(1, "email is required"),
	password: passwordSchema,
});

export const RegisterInputSchema = z.object({
	email: z.email().min(1, "email is required"),
	name: z.string().min(1, "name is required"),
	password: passwordSchema,
});

export const AuthResponseSchema = z.object({
	user: UserSchema,
	token: z.string(),
});

export type User = z.infer<typeof UserSchema>;
export type LoginInput = z.infer<typeof LoginInputSchema>;
export type RegisterInput = z.infer<typeof RegisterInputSchema>;
export type AuthResponse = z.infer<typeof AuthResponseSchema>;
