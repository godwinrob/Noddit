"use client";

import { useAuth, useUser } from "@clerk/nextjs";
import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { setClerkTokenGetter, api } from "@/lib/api";
import UsernameModal from "./UsernameModal";

// --- Noddit user context ---

interface NodditUserCtx {
  username: string | null;
  userId: number | null;
  isReady: boolean;
}

const NodditUserContext = createContext<NodditUserCtx>({
  username: null,
  userId: null,
  isReady: false,
});

export function useNodditUser() {
  return useContext(NodditUserContext);
}

// --- State machine ---

type SyncState = "idle" | "loading" | "needs-username" | "ready";

interface SyncResponse {
  isNew: boolean;
  username?: string;
  userId?: number;
}

export function ClerkTokenProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const { getToken } = useAuth();
  const { user, isSignedIn } = useUser();

  const [syncState, setSyncState] = useState<SyncState>("idle");
  const [dbUsername, setDbUsername] = useState<string | null>(null);
  const [dbUserId, setDbUserId] = useState<number | null>(null);
  const hasSynced = useRef(false);

  // Wire up the Clerk token getter for the API client
  useEffect(() => {
    setClerkTokenGetter(async () => {
      try {
        return await getToken();
      } catch (error) {
        console.error("Failed to get Clerk token:", error);
        return null;
      }
    });
  }, [getToken]);

  // Post-auth sync: check if user exists in DB
  useEffect(() => {
    if (!isSignedIn || !user || hasSynced.current) return;

    const email = user.primaryEmailAddress?.emailAddress;
    if (!email) return;

    const check = async () => {
      setSyncState("loading");
      try {
        const data = await api.post<SyncResponse>(
          "/api/user/sync",
          { email, checkOnly: true },
          true
        );

        if (data.isNew) {
          setSyncState("needs-username");
        } else {
          setDbUsername(data.username ?? null);
          setDbUserId(data.userId ?? null);
          hasSynced.current = true;
          setSyncState("ready");
        }
      } catch (error) {
        console.error("Failed to check user sync:", error);
        // Fall through to ready so the app isn't permanently blocked
        setSyncState("ready");
      }
    };

    check();
  }, [isSignedIn, user]);

  // Reset state when the user signs out
  useEffect(() => {
    if (!isSignedIn) {
      hasSynced.current = false;
      setSyncState("idle");
      setDbUsername(null);
      setDbUserId(null);
    }
  }, [isSignedIn]);

  // Callback for when the user picks a username in the modal
  const handleUsernameChosen = useCallback(
    async (username: string) => {
      const email = user?.primaryEmailAddress?.emailAddress;
      if (!email) throw new Error("No email");

      const data = await api.post<SyncResponse>(
        "/api/user/sync",
        { username, email },
        true
      );

      setDbUsername(data.username ?? username.toLowerCase());
      setDbUserId(data.userId ?? null);
      hasSynced.current = true;
      setSyncState("ready");
    },
    [user]
  );

  const ctxValue: NodditUserCtx = {
    username: dbUsername,
    userId: dbUserId,
    isReady: syncState === "ready",
  };

  return (
    <NodditUserContext.Provider value={ctxValue}>
      {children}
      {syncState === "needs-username" && (
        <UsernameModal onUsernameChosen={handleUsernameChosen} />
      )}
    </NodditUserContext.Provider>
  );
}
