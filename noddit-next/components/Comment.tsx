"use client";

import { useState, useEffect } from "react";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import ReplyForm from "./ReplyForm";

interface CommentProps {
  comment: {
    postId: number;
    body: string;
    username: string;
    postScore: number;
    createdDate: string;
    parentPostId?: number;
    subnodditId: number;
    topLevelId?: number;
  };
  subnodditName: string;
  topLevelId: number;
  onRefresh: () => void;
}

interface Vote {
  vote: string;
  username: string;
}

export default function Comment({
  comment,
  subnodditName,
  topLevelId,
  onRefresh,
}: CommentProps) {
  const { user, isSignedIn } = useUser();
  const [showReplyForm, setShowReplyForm] = useState(false);
  const [hasVoted, setHasVoted] = useState(false);
  const [currentVote, setCurrentVote] = useState<string | null>(null);
  const [score, setScore] = useState(comment.postScore);

  useEffect(() => {
    checkVoteStatus();
  }, [comment.postId, user]);

  const checkVoteStatus = async () => {
    if (!isSignedIn || !user) return;

    try {
      const votes = await api.get<Vote[]>(
        `/api/public/post/votes/${comment.postId}`
      );
      const userVote = votes.find((v) => v.username === user?.username);

      if (userVote) {
        setHasVoted(true);
        setCurrentVote(userVote.vote);
      }
    } catch (error) {
      console.error("Failed to check vote status:", error);
    }
  };

  const handleVote = async (voteType: "upvote" | "downvote") => {
    if (!isSignedIn || !user || hasVoted) return;

    try {
      await api.post(
        "/api/post/vote",
        {
          postId: comment.postId,
          username: user?.username,
          vote: voteType,
        },
        true
      );

      setHasVoted(true);
      setCurrentVote(voteType);
      setScore(score + (voteType === "upvote" ? 1 : -1));
    } catch (error) {
      console.error("Failed to vote:", error);
    }
  };

  const handleReplySuccess = () => {
    setShowReplyForm(false);
    onRefresh();
  };

  return (
    <div className="bg-gray-900 rounded-lg p-4">
      <div className="flex gap-4">
        {/* Voting */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={() => handleVote("upvote")}
            disabled={hasVoted}
            className={`text-xl transition ${
              currentVote === "upvote"
                ? "text-orange-500"
                : hasVoted
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-orange-500"
            }`}
          >
            ▲
          </button>
          <span className="font-bold">{score}</span>
          <button
            onClick={() => handleVote("downvote")}
            disabled={hasVoted}
            className={`text-xl transition ${
              currentVote === "downvote"
                ? "text-blue-500"
                : hasVoted
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-blue-500"
            }`}
          >
            ▼
          </button>
        </div>

        {/* Content */}
        <div className="flex-1">
          <div className="text-sm text-gray-400 mb-2">
            <span className="font-semibold">u/{comment.username}</span>
            {" • "}
            <span>{new Date(comment.createdDate).toLocaleDateString()}</span>
          </div>

          <p className="text-gray-300 whitespace-pre-wrap mb-3">
            {comment.body}
          </p>

          {isSignedIn && (
            <button
              onClick={() => setShowReplyForm(!showReplyForm)}
              className="text-sm text-gray-400 hover:text-white transition"
            >
              {showReplyForm ? "Cancel" : "Reply"}
            </button>
          )}

          {showReplyForm && (
            <div className="mt-3">
              <ReplyForm
                subnodditName={subnodditName}
                subnodditId={comment.subnodditId}
                parentPostId={comment.postId}
                topLevelId={topLevelId}
                onSuccess={handleReplySuccess}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
