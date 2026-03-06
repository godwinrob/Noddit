"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { api } from "@/lib/api";
import PostCard from "@/components/PostCard";
import CommentSection from "@/components/CommentSection";

interface Post {
  postId: number;
  title: string;
  body: string;
  username: string;
  subnodditName: string;
  subnodditId: number;
  postScore: number;
  createdDate: string;
  imageAddress?: string;
  parentPostId?: number;
  topLevelId?: number;
}

export default function PostDetailPage() {
  const params = useParams();
  const subnodditName = params.name as string;
  const postId = params.id as string;

  const [post, setPost] = useState<Post | null>(null);
  const [replies, setReplies] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Validate postId is numeric
    if (isNaN(Number(postId))) {
      window.location.href = "/";
      return;
    }

    fetchPostData();
  }, [subnodditName, postId]);

  const fetchPostData = async () => {
    try {
      setLoading(true);

      const [postData, repliesData] = await Promise.all([
        api.get<Post>(`/api/public/${subnodditName}/${postId}`),
        api.get<Post[]>(`/api/public/${subnodditName}/${postId}/replies`),
      ]);

      setPost(postData);
      setReplies(repliesData);
    } catch (error) {
      console.error("Failed to fetch post:", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="text-gray-400">Loading...</div>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="text-center py-12">
        <h1 className="text-2xl font-bold mb-4">Post Not Found</h1>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      {/* Main Post */}
      <PostCard post={post} showFullBody={true} onRefresh={fetchPostData} />

      {/* Comments Section */}
      <div className="mt-6">
        <CommentSection
          comments={replies}
          subnodditName={subnodditName}
          topLevelId={parseInt(postId)}
          onRefresh={fetchPostData}
        />
      </div>
    </div>
  );
}
