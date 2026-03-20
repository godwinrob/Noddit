"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import { toast } from "sonner";

interface Post {
  postId: number;
  title: string;
  body: string;
  username: string;
  subnodditName: string;
  postScore: number;
  createdDate: string;
  imageAddress?: string;
  parentPostId?: number;
  topLevelId?: number;
  subnodditId: number;
}

interface Vote {
  vote: string;
  username: string;
}

interface PostCardProps {
  post: Post;
  showFullBody?: boolean;
  onRefresh?: () => void;
}

export default function PostCardCompact({
  post,
  showFullBody = false,
  onRefresh,
}: PostCardProps) {
  const { user, isSignedIn } = useUser();
  const { username: nodditUsername } = useNodditUser();
  const [currentVote, setCurrentVote] = useState<string | null>(null);
  const [score, setScore] = useState(post.postScore);
  const [isVoting, setIsVoting] = useState(false);

  useEffect(() => {
    if (isSignedIn && nodditUsername) {
      checkVoteStatus();
    }
  }, [post.postId, isSignedIn, nodditUsername]);

  const checkVoteStatus = async () => {
    if (!isSignedIn || !nodditUsername) return;

    try {
      const votes = await api.get<Vote[]>(`/api/public/post/votes/${post.postId}`);
      const userVote = votes.find((v) => v.username === nodditUsername);
      if (userVote) {
        setCurrentVote(userVote.vote);
      }
    } catch (error) {
      console.error("Failed to check vote status:", error);
    }
  };

  const handleVote = async (voteType: "upvote" | "downvote") => {
    if (!isSignedIn || !user || isVoting) return;
    setIsVoting(true);

    try {
      const result = await api.post<{ vote: string | null; score: number }>(
        "/api/post/vote",
        {
          postId: post.postId,
          username: nodditUsername,
          vote: voteType,
        },
        true
      );

      setCurrentVote(result.vote);
      setScore(result.score);
      toast.success(result.vote === null ? "Vote removed" : `${voteType === "upvote" ? "Upvoted" : "Downvoted"}!`);
    } catch (error) {
      console.error("Failed to vote:", error);
      toast.error("Failed to vote. Please try again.");
    } finally {
      setIsVoting(false);
    }
  };

  const formatTimeAgo = (dateString: string) => {
    const date = new Date(dateString);
    const now = new Date();
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 60) return "just now";
    if (seconds < 3600) return `${Math.floor(seconds / 60)} minutes ago`;
    if (seconds < 86400) return `${Math.floor(seconds / 3600)} hours ago`;
    return `${Math.floor(seconds / 86400)} days ago`;
  };

  return (
    <div className="bg-gray-900 border-l-4 border-gray-700 hover:border-orange-500 transition-colors">
      <div className="flex gap-2 p-2">
        {/* Vote column */}
        <div className="flex flex-col items-center w-12 text-xs">
          <button
            onClick={() => handleVote("upvote")}
            disabled={!isSignedIn || isVoting}
            className={`${
              currentVote === "upvote"
                ? "text-orange-500"
                : "text-gray-500 hover:text-orange-500"
            } disabled:text-gray-700 disabled:cursor-not-allowed`}
            title="upvote"
          >
            ▲
          </button>
          <span className={`font-bold ${score > 0 ? "text-orange-500" : score < 0 ? "text-blue-500" : "text-gray-400"}`}>
            {score}
          </span>
          <button
            onClick={() => handleVote("downvote")}
            disabled={!isSignedIn || isVoting}
            className={`${
              currentVote === "downvote"
                ? "text-blue-500"
                : "text-gray-500 hover:text-blue-500"
            } disabled:text-gray-700 disabled:cursor-not-allowed`}
            title="downvote"
          >
            ▼
          </button>
        </div>

        {/* Thumbnail */}
        <Link href={`/n/${post.subnodditName}/${post.postId}`} className="flex-shrink-0">
          {post.imageAddress ? (
            <div className="w-[70px] h-[70px] bg-gray-800 border border-gray-700 overflow-hidden">
              <Image
                src={post.imageAddress}
                alt={post.title}
                width={70}
                height={70}
                className="w-full h-full object-cover"
                unoptimized
              />
            </div>
          ) : (
            <div className="w-[70px] h-[70px] bg-gray-800 border border-gray-700 flex items-center justify-center">
              <svg className="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
              </svg>
            </div>
          )}
        </Link>

        {/* Content column */}
        <div className="flex-1 min-w-0">
          <div className="text-xs text-gray-400 mb-1">
            <Link href={`/n/${post.subnodditName}`} className="hover:underline font-bold">
              /n/{post.subnodditName}
            </Link>
            <span className="mx-1">•</span>
            <span>submitted {formatTimeAgo(post.createdDate)} by </span>
            <Link href={`/u/${post.username}`} className="hover:underline">
              {post.username}
            </Link>
          </div>

          <Link href={`/n/${post.subnodditName}/${post.postId}`}>
            <h3 className="text-base font-normal text-blue-400 hover:underline mb-1">
              {post.title}
            </h3>
          </Link>

          <div className="flex gap-3 text-xs text-gray-500">
            <Link href={`/n/${post.subnodditName}/${post.postId}`} className="hover:underline">
              comments
            </Link>
            <button className="hover:underline">share</button>
            <button className="hover:underline">save</button>
            <button className="hover:underline">hide</button>
          </div>
        </div>
      </div>
    </div>
  );
}
