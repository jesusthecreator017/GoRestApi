import type { Metadata } from "next";
import { Providers } from "./providers";
import { Navbar } from "@/components/Navbar";
import "./globals.css";

export const metadata: Metadata = {
	title: "FSWithGo",
	description: "Issue tracker",
};

export default function RootLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	return (
		<html lang="en">
			<body suppressHydrationWarning>
				<Providers>
					<div className="min-h-screen bg-background text-foreground">
						<Navbar />
						{children}
					</div>
				</Providers>
			</body>
		</html>
	);
}
