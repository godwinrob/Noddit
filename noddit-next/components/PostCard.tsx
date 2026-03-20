"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import ReplyForm from "./ReplyForm";
import { Card } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Separator } from "@/components/ui/separator";
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
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);

  useEffect(() => {
    if (isSignedIn && nodditUsername) {
      checkVoteStatus();
      checkDeletePermission();
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
      toast.success(result.vote === null ? "Vote removed" : `${voteType === "upvote" ? "Upvoted" : "Downvoted"}!`);
    } catch (error) {
      console.error("Failed to vote:", error);
      toast.error("Failed to vote. Please try again.");
    } finally {
      setIsVoting(false);
    }
  };

  const handleDelete = async () => {
    if (!canDelete) return;

    try {
      await api.delete(`/api/post/delete/${post.postId}`, true);
      toast.success("Post deleted successfully");
      setShowDeleteDialog(false);
      onRefresh?.();
    } catch (error) {
      console.error("Failed to delete post:", error);
      toast.error("Failed to delete post. Please try again.");
    }
  };

  const handleReplySuccess = () => {
    setShowReplyForm(false);
    onRefresh?.();
  };

  return (
    <Card className="bg-gray-900 border-gray-800 overflow-hidden">
      <div className="flex gap-4 p-4">
        {/* Voting */}
        <div className="flex flex-col items-center gap-1">
          <Button
            onClick={() => handleVote("upvote")}
            disabled={!isSignedIn || isVoting}
            variant="ghost"
            size="sm"
            className={`p-0 h-auto text-2xl transition ${
              currentVote === "upvote"
                ? "text-orange-500"
                : !isSignedIn
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-orange-500"
            }`}
          >
            ▲
          </Button>
          <span className="font-bold text-lg">{score}</span>
          <Button
            onClick={() => handleVote("downvote")}
            disabled={!isSignedIn || isVoting}
            variant="ghost"
            size="sm"
            className={`p-0 h-auto text-2xl transition ${
              currentVote === "downvote"
                ? "text-blue-500"
                : !isSignedIn
                ? "text-gray-600 cursor-not-allowed"
                : "text-gray-400 hover:text-blue-500"
            }`}
          >
            ▼
          </Button>
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
          <div className="flex flex-wrap items-center gap-2 mb-2">
            <Badge variant="secondary" className="bg-orange-900/30 text-orange-400 hover:bg-orange-900/50">
              <Link href={`/s/${post.subnodditName}`}>
                n/{post.subnodditName}
              </Link>
            </Badge>
            <Separator orientation="vertical" className="h-4" />
            <div className="flex items-center gap-2">
              <Avatar className="h-5 w-5">
                <AvatarFallback className="text-xs bg-gray-700">
                  {post.username[0].toUpperCase()}
                </AvatarFallback>
              </Avatar>
              <span className="text-sm text-gray-400">
                <span className="hover:underline">u/{post.username}</span>
              </span>
            </div>
            <Separator orientation="vertical" className="h-4" />
            <span className="text-sm text-gray-500">
              {new Date(post.createdDate).toLocaleDateString()}
            </span>
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
          <div className="mt-3 flex gap-4">
            {showFullBody && isSignedIn && (
              <Button
                onClick={() => setShowReplyForm(!showReplyForm)}
                variant="ghost"
                size="sm"
                className="text-gray-400 hover:text-white"
              >
                {showReplyForm ? "Cancel" : "Reply"}
              </Button>
            )}
            {canDelete && (
              <Dialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
                <DialogTrigger>
                  <Button
                    variant="ghost"
                    size="sm"
                    className="text-gray-400 hover:text-red-500"
                  >
                    Delete
                  </Button>
                </DialogTrigger>
                <DialogContent className="bg-gray-900 border-gray-700">
                  <DialogHeader>
                    <DialogTitle>Delete Post</DialogTitle>
                    <DialogDescription className="text-gray-400">
                      Are you sure you want to delete this post? This action cannot be undone.
                    </DialogDescription>
                  </DialogHeader>
                  <DialogFooter>
                    <Button
                      variant="ghost"
                      onClick={() => setShowDeleteDialog(false)}
                    >
                      Cancel
                    </Button>
                    <Button
                      variant="destructive"
                      onClick={handleDelete}
                    >
                      Delete
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
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
    </Card>
  );
}
