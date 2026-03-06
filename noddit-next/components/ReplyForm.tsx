"use client";

import { useState } from "react";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";

interface ReplyFormProps {
  subnodditName: string;
  subnodditId: number;
  parentPostId: number;
  topLevelId: number;
  onSuccess: () => void;
}

export default function ReplyForm({
  subnodditName,
  subnodditId,
  parentPostId,
  topLevelId,
  onSuccess,
}: ReplyFormProps) {
  const [body, setBody] = useState("");
  const [loading, setLoading] = useState(false);
  const { user } = useUser();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!body.trim() || !user) return;

    setLoading(true);

    try {
      await api.post(
        `/${subnodditName}/${topLevelId}/createreply`,
        {
          body,
          username: user?.username,
          parentPostId,
          topLevelId,
          subnodditId,
        },
        true
      );

      setBody("");
      onSuccess();
    } catch (error) {
      console.error("Failed to post reply:", error);
      alert("Failed to post reply");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        placeholder="What are your thoughts?"
        rows={4}
        maxLength={2000}
        className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
        required
      />
      <div className="flex justify-between items-center">
        <span className="text-sm text-gray-400">{body.length}/2000</span>
        <button
          type="submit"
          disabled={loading || !body.trim()}
          className="bg-orange-600 hover:bg-orange-700 disabled:bg-gray-700 text-white font-semibold py-2 px-6 rounded transition"
        >
          {loading ? "Posting..." : "Reply"}
        </button>
      </div>
    </form>
  );
}
