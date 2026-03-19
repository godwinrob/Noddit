"use client";

import { useState, useEffect, useRef, useCallback } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL;

const USERNAME_REGEX = /^[a-zA-Z0-9][a-zA-Z0-9_]{2,19}$/;
const RESERVED_NAMES = new Set([
  "admin",
  "moderator",
  "noddit",
  "system",
  "deleted",
]);

type Status = "idle" | "checking" | "available" | "error";

interface Props {
  onUsernameChosen: (username: string) => Promise<void>;
}

function validateLocal(value: string): string | null {
  if (value.length < 3) return "Must be at least 3 characters";
  if (value.length > 20) return "Must be 20 characters or fewer";
  if (value.startsWith("_")) return "Cannot start with underscore";
  if (!USERNAME_REGEX.test(value))
    return "Letters, numbers, and underscores only";
  if (RESERVED_NAMES.has(value.toLowerCase())) return "This name is reserved";
  return null;
}

export default function UsernameModal({ onUsernameChosen }: Props) {
  const [value, setValue] = useState("");
  const [status, setStatus] = useState<Status>("idle");
  const [message, setMessage] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  // Auto-focus input on mount
  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  const checkAvailability = useCallback(async (username: string) => {
    try {
      const res = await fetch(
        `${API_URL}/api/public/user/available/${encodeURIComponent(username)}`
      );
      if (!res.ok) {
        setStatus("error");
        setMessage("Server error, try again");
        return;
      }
      const data = await res.json();
      if (data.available) {
        setStatus("available");
        setMessage("Username is available");
      } else {
        setStatus("error");
        setMessage(data.reason || "Username is taken");
      }
    } catch {
      setStatus("error");
      setMessage("Could not check availability");
    }
  }, []);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const raw = e.target.value;
    setValue(raw);

    if (debounceRef.current) clearTimeout(debounceRef.current);

    const localErr = validateLocal(raw);
    if (localErr) {
      setStatus("error");
      setMessage(localErr);
      return;
    }

    setStatus("checking");
    setMessage("Checking availability...");
    debounceRef.current = setTimeout(() => checkAvailability(raw), 400);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (status !== "available" || submitting) return;

    setSubmitting(true);
    try {
      await onUsernameChosen(value);
    } catch {
      setStatus("error");
      setMessage("Failed to create account. Try a different name.");
      setSubmitting(false);
    }
  };

  const statusColor =
    status === "available"
      ? "text-green-400"
      : status === "error"
        ? "text-red-400"
        : "text-gray-400";

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70">
      <form
        onSubmit={handleSubmit}
        className="bg-gray-900 border border-gray-700 rounded-lg p-8 w-full max-w-md mx-4 shadow-xl"
      >
        <h2 className="text-2xl font-bold text-white mb-2">
          Choose your username
        </h2>
        <p className="text-gray-400 mb-6 text-sm">
          This is how other people will see you on Noddit.
        </p>

        <label htmlFor="username-input" className="sr-only">
          Username
        </label>
        <input
          ref={inputRef}
          id="username-input"
          type="text"
          value={value}
          onChange={handleChange}
          placeholder="username"
          maxLength={20}
          autoComplete="off"
          className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-3 text-white placeholder-gray-500 focus:outline-none focus:border-orange-500 text-lg"
        />

        {value.length > 0 && (
          <p className={`mt-2 text-sm ${statusColor}`}>
            {status === "checking" ? "Checking..." : message}
          </p>
        )}

        <p className="mt-3 text-xs text-gray-500">
          3-20 characters. Letters, numbers, and underscores.
        </p>

        <button
          type="submit"
          disabled={status !== "available" || submitting}
          className="mt-6 w-full bg-orange-600 hover:bg-orange-700 disabled:bg-gray-700 disabled:cursor-not-allowed text-white font-semibold py-3 rounded transition"
        >
          {submitting ? "Creating account..." : "Continue"}
        </button>
      </form>
    </div>
  );
}
