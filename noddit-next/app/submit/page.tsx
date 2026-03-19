"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

export default function SubmitPage() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [imageAddress, setImageAddress] = useState("");
  const [subnodditId, setSubnodditId] = useState("");
  const [subnoddits, setSubnoddits] = useState<Subnoddit[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { isSignedIn } = useUser();
  const { username: nodditUsername } = useNodditUser();
  const router = useRouter();

  useEffect(() => {
    if (!isSignedIn) {
      router.push("/");
      return;
    }

    // Fetch all subnoddits
    const fetchSubnoddits = async () => {
      try {
        const data = await api.get<Subnoddit[]>("/api/public/subnoddits");
        setSubnoddits(data);
        if (data.length > 0) {
          setSubnodditId(data[0].subnodditId.toString());
        }
      } catch (error) {
        console.error("Failed to fetch subnoddits:", error);
      }
    };

    fetchSubnoddits();
  }, [isSignedIn, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    if (!title.trim()) {
      setError("Title is required");
      setLoading(false);
      return;
    }

    if (!body.trim()) {
      setError("Body is required");
      setLoading(false);
      return;
    }

    if (!subnodditId) {
      setError("Please select a community");
      setLoading(false);
      return;
    }

    try {
      const post = await api.post<{ postId: number }>(
        "/api/post/create",
        {
          title,
          body,
          imageAddress: imageAddress || null,
          username: nodditUsername,
          subnodditId: parseInt(subnodditId),
        },
        true
      );

      // Get the subnoddit name for redirect
      const subnoddit = subnoddits.find(
        (s) => s.subnodditId === parseInt(subnodditId)
      );

      if (subnoddit && post.postId) {
        router.push(`/s/${subnoddit.subnodditName}/${post.postId}`);
      } else {
        router.push("/");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create post");
      setLoading(false);
    }
  };

  if (!isSignedIn) {
    return null; // Will redirect
  }

  return (
    <div className="max-w-3xl mx-auto mt-8">
      <div className="bg-gray-900 rounded-lg p-8">
        <h1 className="text-3xl font-bold mb-6">Create a Post</h1>

        {error && (
          <div className="bg-red-500/10 border border-red-500 text-red-500 rounded p-3 mb-4">
            {error}
          </div>
        )}

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label htmlFor="community" className="block text-sm font-medium mb-2">
              Choose a community
            </label>
            <select
              id="community"
              value={subnodditId}
              onChange={(e) => setSubnodditId(e.target.value)}
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
              required
            >
              {subnoddits.map((sn) => (
                <option key={sn.subnodditId} value={sn.subnodditId}>
                  n/{sn.subnodditName}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label htmlFor="title" className="block text-sm font-medium mb-2">
              Title
            </label>
            <input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="An interesting title"
              maxLength={100}
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
              required
            />
            <p className="text-sm text-gray-400 mt-1">{title.length}/100</p>
          </div>

          <div>
            <label htmlFor="body" className="block text-sm font-medium mb-2">
              Text
            </label>
            <textarea
              id="body"
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="What are your thoughts?"
              rows={8}
              maxLength={2000}
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
              required
            />
            <p className="text-sm text-gray-400 mt-1">{body.length}/2000</p>
          </div>

          <div>
            <label htmlFor="image" className="block text-sm font-medium mb-2">
              Image URL (optional)
            </label>
            <input
              id="image"
              type="url"
              value={imageAddress}
              onChange={(e) => setImageAddress(e.target.value)}
              placeholder="https://example.com/image.jpg"
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
            />
          </div>

          <div className="flex gap-3 pt-4">
            <button
              type="submit"
              disabled={loading}
              className="flex-1 bg-orange-600 hover:bg-orange-700 disabled:bg-gray-700 text-white font-semibold py-3 px-4 rounded transition"
            >
              {loading ? "Posting..." : "Post"}
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
        <h2 className="font-semibold mb-2">Posting Guidelines</h2>
        <ul className="text-sm text-gray-400 space-y-1">
          <li>• Be respectful and civil</li>
          <li>• Stay on topic for the community</li>
          <li>• No spam or self-promotion</li>
          <li>• Use descriptive titles</li>
        </ul>
      </div>
    </div>
  );
}
