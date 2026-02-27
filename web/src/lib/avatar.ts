const COLORS = [
	"bg-red-600",
	"bg-orange-600",
	"bg-amber-600",
	"bg-yellow-600",
	"bg-lime-600",
	"bg-green-600",
	"bg-emerald-600",
	"bg-teal-600",
	"bg-cyan-600",
	"bg-blue-600",
	"bg-indigo-600",
	"bg-purple-600",
];

export function getAvatarColor(name: string): string {
	let hash = 5381;
	for (let i = 0; i < name.length; i++) {
		hash = (hash * 33) ^ name.charCodeAt(i);
	}
	return COLORS[Math.abs(hash) % COLORS.length];
}

export function getInitials(name: string): string {
	const parts = name.trim().split(/\s+/);
	if (parts.length === 1) {
		return parts[0][0].toUpperCase();
	}
	return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}
