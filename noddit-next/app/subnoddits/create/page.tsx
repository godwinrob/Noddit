"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";

export default function CreateSubnodditPage() {
  const [subnodditName, setSubnodditName] = useState("");
  const [subnodditDescription, setSubnodditDescription] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { user, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isAuthenticated) {
      router.push("/login");
    }
  }, [isAuthenticated, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    if (!subnodditName.trim() || !subnodditDescription.trim()) {
      setError("All fields are required");
      setLoading(false);
      return;
    }

    try {
      // Replace spaces with underscores
      const formattedName = subnodditName.replace(/\s+/g, "_");

      await api.post(
        "/api/subnoddits/create",
        {
          subnodditName: formattedName,
          subnodditDescription,
          username: user?.username,
        },
        true
      );

      router.push(`/s/${formattedName}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create community");
      setLoading(false);
    }
  };

  if (!isAuthenticated) {
    return null; // Will redirect
  }

  return (
    <div className="max-w-2xl mx-auto mt-8">
      <div className="bg-gray-900 rounded-lg p-8">
        <h1 className="text-3xl font-bold mb-6">Create a Community</h1>

        {error && (
          <div className="bg-red-500/10 border border-red-500 text-red-500 rounded p-3 mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="name" className="block text-sm font-medium mb-2">
              Community Name
            </label>
            <input
              id="name"
              type="text"
              value={subnodditName}
              onChange={(e) => setSubnodditName(e.target.value)}
              placeholder="my_awesome_community"
              maxLength={30}
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
              required
            />
            <p className="text-sm text-gray-400 mt-1">
              Spaces will be replaced with underscores
            </p>
          </div>

          <div>
            <label
              htmlFor="description"
              className="block text-sm font-medium mb-2"
            >
              Description
            </label>
            <textarea
              id="description"
              value={subnodditDescription}
              onChange={(e) => setSubnodditDescription(e.target.value)}
              placeholder="What is this community about?"
              rows={4}
              maxLength={200}
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
              required
            />
            <p className="text-sm text-gray-400 mt-1">
              {subnodditDescription.length}/200
            </p>
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              disabled={loading}
              className="flex-1 bg-orange-600 hover:bg-orange-700 disabled:bg-gray-700 text-white font-semibold py-3 px-4 rounded transition"
            >
              {loading ? "Creating..." : "Create Community"}
            </button>
            <button
              type="button"
              onClick={() => router.back()}
              className="px-6 bg-gray-700 hover:bg-gray-600 text-white font-semibold py-3 rounded transition"
            >
              Cancel
            </button>
          </div>
        </form>
      </div>

      <div className="bg-gray-900 rounded-lg p-6 mt-4">
        <h2 className="font-semibold mb-2">Community Guidelines</h2>
        <ul className="text-sm text-gray-400 space-y-1">
          <li>• Choose a clear, descriptive name</li>
          <li>• Write a helpful description</li>
          <li>• You'll become the first moderator</li>
          <li>• Keep it respectful and on-topic</li>
        </ul>
      </div>
    </div>
  );
}
