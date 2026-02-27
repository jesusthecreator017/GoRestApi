import { getAvatarColor, getInitials } from "@/lib/avatar";
import { cn } from "@/lib/utils";

interface UserAvatarProps {
	name: string;
	className?: string;
}

export function UserAvatar({ name, className }: UserAvatarProps) {
	return (
		<div
			className={cn(
				"flex size-9 items-center justify-center rounded-full text-sm font-medium text-white",
				getAvatarColor(name),
				className,
			)}
		>
			{getInitials(name)}
		</div>
	);
}
