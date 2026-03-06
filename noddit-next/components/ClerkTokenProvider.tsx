"use client";

import { useAuth } from "@clerk/nextjs";
import { useEffect } from "react";
import { setClerkTokenGetter } from "@/lib/api";

export function ClerkTokenProvider({ children }: { children: React.ReactNode }) {
  const { getToken } = useAuth();

  useEffect(() => {
    // Set the Clerk token getter for the API client
    setClerkTokenGetter(async () => {
      try {
        return await getToken();
      } catch (error) {
        console.error("Failed to get Clerk token:", error);
        return null;
      }
    });
  }, [getToken]);

  return <>{children}</>;
}
