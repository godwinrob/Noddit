"use client";

import Comment from "./Comment";

interface Post {
  postId: number;
  body: string;
  username: string;
  postScore: number;
  createdDate: string;
  parentPostId?: number;
  subnodditId: number;
  topLevelId?: number;
}

interface CommentSectionProps {
  comments: Post[];
  subnodditName: string;
  topLevelId: number;
  onRefresh: () => void;
}

export default function CommentSection({
  comments,
  subnodditName,
  topLevelId,
  onRefresh,
}: CommentSectionProps) {
  if (comments.length === 0) {
    return (
      <div className="bg-gray-900 rounded-lg p-8 text-center">
        <p className="text-gray-400">No comments yet. Be the first to comment!</p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {comments.map((comment) => (
        <Comment
          key={comment.postId}
          comment={comment}
          subnodditName={subnodditName}
          topLevelId={topLevelId}
          onRefresh={onRefresh}
        />
      ))}
    </div>
  );
}
