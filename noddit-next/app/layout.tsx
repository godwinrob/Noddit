import type { Metadata } from "next";
import "./globals.css";
import { AuthProvider } from "@/lib/auth-context";
import Nav from "@/components/Nav";

export const metadata: Metadata = {
  title: "Noddit - A Reddit-like Message Board",
  description: "Share posts, comment, and vote on your favorite topics",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased min-h-screen">
        <AuthProvider>
          <Nav />
          <main className="container mx-auto px-4 py-8 max-w-6xl">
            {children}
          </main>
        </AuthProvider>
      </body>
    </html>
  );
}
