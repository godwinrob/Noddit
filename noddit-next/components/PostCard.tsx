"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import ReplyForm from "./ReplyForm";

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

export default function PostCard({
  post,
  showFullBody = false,
  onRefresh,
}: PostCardProps) {
  const { user, isSignedIn } = useUser();
  const { username: nodditUsername } = useNodditUser();
  const [showReplyForm, setShowReplyForm] = useState(false);
  const [currentVote, setCurrentVote] = useState<string | null>(null);
  const [score, setScore] = useState(post.postScore);
  const [isVoting, setIsVoting] = useState(false);
  const [canDelete, setCanDelete] = useState(false);

  useEffect(() => {
    checkVoteStatus();
    checkDeletePermission();
  }, [post.postId, user]);

  const checkVoteStatus = async () => {
    if (!isSignedIn || !user) return;

    try {
      const votes = await api.get<Vote[]>(`/api/public/post/votes/${post.postId}`);
      const userVote = votes.find((v) => v.username === nodditUsername);
      setCurrentVote(userVote?.vote ?? null);
    } catch (error) {
      console.error("Failed to check vote status:", error);
    }
  };

  const checkDeletePermission = async () => {
    if (!isSignedIn || !user) return;

    // User is author
    if (nodditUsername === post.username) {
      setCanDelete(true);
      return;
    }

    // User is admin/super_admin (check Clerk publicMetadata)
    const role = (user?.publicMetadata as Record<string, unknown>)?.role as string | undefined;
    if (role === "admin" || role === "super_admin") {
      setCanDelete(true);
      return;
    }

    // Check if user is moderator
    try {
      const moderators = await api.get<{ username: string }[]>(
        `/api/public/moderators/${post.subnodditName}`
      );
      if (moderators.some((m) => m.username === nodditUsername)) {
        setCanDelete(true);
      }
    } catch (error) {
      console.error("Failed to check moderator status:", error);
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
    } catch (error) {
      console.error("Failed to vote:", error);
    } finally {
      setIsVoting(false);
    }
  };

  const handleDelete = async () => {
    if (!canDelete || !confirm("Are you sure you want to delete this post?"))
      return;

    try {
      await api.delete(`/api/post/delete/${post.postId}`, true);
      onRefresh?.();
    } catch (error) {
      console.error("Failed to delete post:", error);
    }
  };

  const handleReplySuccess = () => {
    setShowReplyForm(false);
    onRefresh?.();
  };

  return (
    <div className="bg-gray-900 rounded-lg overflow-hidden">
      <div className="flex gap-4 p-4">
        {/* Voting */}
        <div className="flex flex-col items-center gap-1">
          <button
            onClick={() => handleVote("upvote")}
            disabled={!isSignedIn || isVoting}
            className={`text-2xl transition ${
              currentVote === "upvote"
                ? "text-orange-500"
                : !isSignedIn
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-orange-500"
            }`}
          >
            ▲
          </button>
          <span className="font-bold text-lg">{score}</span>
          <button
            onClick={() => handleVote("downvote")}
            disabled={!isSignedIn || isVoting}
            className={`text-2xl transition ${
              currentVote === "downvote"
                ? "text-blue-500"
                : !isSignedIn
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-blue-500"
            }`}
          >
            ▼
          </button>
        </div>

        {/* Content */}
        <div className="flex-1">
          {/* Image (if exists) */}
          {post.imageAddress && (
            <div className="mb-3">
              <Image
                src={post.imageAddress}
                alt={post.title}
                width={600}
                height={400}
                className="rounded max-w-full h-auto"
                unoptimized
              />
            </div>
          )}

          {/* Metadata */}
          <div className="text-sm text-gray-400 mb-2">
            <Link
              href={`/s/${post.subnodditName}`}
              className="hover:underline font-semibold"
            >
              n/{post.subnodditName}
            </Link>
            {" • Posted by "}
            <span className="hover:underline">u/{post.username}</span>
            {" • "}
            <span>{new Date(post.createdDate).toLocaleDateString()}</span>
          </div>

          {/* Title & Body */}
          {showFullBody ? (
            <>
              <h1 className="text-2xl font-bold mb-3">{post.title}</h1>
              <p className="text-gray-300 whitespace-pre-wrap">{post.body}</p>
            </>
          ) : (
            <Link href={`/s/${post.subnodditName}/${post.postId}`}>
              <h2 className="text-xl font-semibold mb-2 hover:text-orange-500 transition">
                {post.title}
              </h2>
            </Link>
          )}

          {/* Actions */}
          <div className="mt-3 flex gap-4 text-sm text-gray-400">
            {showFullBody && isSignedIn && (
              <button
                onClick={() => setShowReplyForm(!showReplyForm)}
                className="hover:text-white transition"
              >
                {showReplyForm ? "Cancel" : "Reply"}
              </button>
            )}
            {canDelete && (
              <button
                onClick={handleDelete}
                className="hover:text-red-500 transition"
              >
                Delete
              </button>
            )}
          </div>

          {/* Reply Form */}
          {showReplyForm && (
            <div className="mt-4">
              <ReplyForm
                subnodditName={post.subnodditName}
                subnodditId={post.subnodditId}
                parentPostId={post.postId}
                topLevelId={post.topLevelId || post.postId}
                onSuccess={handleReplySuccess}
              />
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
