import type { Metadata } from "next";
import "./globals.css";
import { ClerkProvider } from "@clerk/nextjs";
import { ClerkTokenProvider } from "@/components/ClerkTokenProvider";
import Nav from "@/components/Nav";
import { Geist } from "next/font/google";
import { cn } from "@/lib/utils";
import { Toaster } from "sonner";

const geist = Geist({subsets:['latin'],variable:'--font-sans'});

export const metadata: Metadata = {
  title: "Noddit - An opensource message board",
  description: "Share posts, comment, and vote on your favorite topics",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" className={cn("font-sans", geist.variable)}>
      <body className="antialiased min-h-screen">
        <ClerkProvider>
          <ClerkTokenProvider>
            <Nav />
            <main className="container mx-auto px-4 py-8 max-w-6xl">
              {children}
            </main>
            <Toaster position="bottom-right" theme="dark" richColors />
          </ClerkTokenProvider>
        </ClerkProvider>
      </body>
    </html>
  );
}
