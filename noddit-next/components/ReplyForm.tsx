"use client";

import { useState } from "react";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { toast } from "sonner";

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
  const { username: nodditUsername } = useNodditUser();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!body.trim() || !user || !nodditUsername) return;

    setLoading(true);

    try {
      await api.post(
        `/${subnodditName}/${topLevelId}/createreply`,
        {
          body,
          username: nodditUsername,
          parentPostId,
          topLevelId,
          subnodditId,
        },
        true
      );

      setBody("");
      toast.success("Reply posted successfully!");
      onSuccess();
    } catch (error) {
      console.error("Failed to post reply:", error);
      toast.error("Failed to post reply. Please try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-3">
      <Textarea
        value={body}
        onChange={(e) => setBody(e.target.value)}
        placeholder="What are your thoughts?"
        rows={4}
        maxLength={2000}
        className="w-full bg-gray-800 border-gray-700 focus:border-orange-500 resize-none"
        required
      />
      <div className="flex justify-between items-center">
        <span className="text-sm text-gray-400">{body.length}/2000</span>
        <Button
          type="submit"
          disabled={loading || !body.trim()}
          className="bg-orange-600 hover:bg-orange-700"
        >
          {loading ? "Posting..." : "Reply"}
        </Button>
      </div>
    </form>
  );
}
