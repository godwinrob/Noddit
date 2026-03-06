"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { useUser, useClerk } from "@clerk/nextjs";
import { api } from "@/lib/api";
import PostCard from "@/components/PostCard";

interface Post {
  postId: number;
  title: string;
  body: string;
  username: string;
  subnodditName: string;
  postScore: number;
  createdDate: string;
  imageAddress?: string;
  subnodditId: number;
}

export default function ProfilePage() {
  const { user, isSignedIn } = useUser();
  const { signOut } = useClerk();
  const router = useRouter();

  const [posts, setPosts] = useState<Post[]>([]);
  const [loading, setLoading] = useState(true);
  const [newUsername, setNewUsername] = useState("");
  const [showUsernameModal, setShowUsernameModal] = useState(false);

  useEffect(() => {
    if (!isSignedIn) {
      router.push("/");
      return;
    }

    fetchUserPosts();
  }, [isSignedIn, user, router]);

  const fetchUserPosts = async () => {
    if (!user) return;

    try {
      const username = user.username || user.primaryEmailAddress?.emailAddress?.split('@')[0] || '';
      const data = await api.get<Post[]>(
        `/api/public/post/user/${username}`
      );
      setPosts(data);
    } catch (error) {
      console.error("Failed to fetch user posts:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateUsername = async () => {
    if (!user || !newUsername.trim()) return;

    try {
      const username = user.username || user.primaryEmailAddress?.emailAddress?.split('@')[0] || '';
      await api.put(
        `/api/user/update/username/${username}`,
        { newUsername },
        true
      );

      alert("Username updated! Please sign in again.");
      signOut();
    } catch (error) {
      console.error("Failed to update username:", error);
      alert("Failed to update username");
    }
  };

  if (!isSignedIn || !user) {
    return null; // Will redirect
  }

  return (
    <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
      {/* Posts Section */}
      <div className="lg:col-span-2">
        <h1 className="text-4xl font-bold mb-6">My Posts</h1>

        {loading ? (
          <div className="flex justify-center items-center min-h-[200px]">
            <div className="text-gray-400">Loading...</div>
          </div>
        ) : posts.length === 0 ? (
          <div className="bg-gray-900 rounded-lg p-12 text-center">
            <p className="text-gray-400 mb-4">You haven't posted anything yet</p>
            <button
              onClick={() => router.push("/submit")}
              className="bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition"
            >
              Create your first post
            </button>
          </div>
        ) : (
          <div className="space-y-4">
            {posts.map((post) => (
              <PostCard
                key={post.postId}
                post={post}
                onRefresh={fetchUserPosts}
              />
            ))}
          </div>
        )}
      </div>

      {/* Account Settings Section */}
      <div>
        <div className="bg-gray-900 rounded-lg p-6 sticky top-8">
          <h2 className="text-2xl font-bold mb-4">Account Settings</h2>

          <div className="mb-6">
            <div className="text-sm text-gray-400">Username</div>
            <div className="text-lg font-semibold">{user.username || user.primaryEmailAddress?.emailAddress?.split('@')[0]}</div>
          </div>

          <div className="space-y-3">
            <details className="bg-gray-800 rounded">
              <summary className="cursor-pointer p-3 font-semibold hover:bg-gray-700">
                Change Username
              </summary>
              <div className="p-3 pt-0">
                <input
                  type="text"
                  value={newUsername}
                  onChange={(e) => setNewUsername(e.target.value)}
                  placeholder="New username"
                  className="w-full bg-gray-700 border border-gray-600 rounded px-3 py-2 mb-2 focus:outline-none focus:border-orange-500"
                />
                <button
                  onClick={handleUpdateUsername}
                  disabled={!newUsername.trim()}
                  className="w-full bg-orange-600 hover:bg-orange-700 disabled:bg-gray-700 text-white font-semibold py-2 rounded transition"
                >
                  Update Username
                </button>
                <p className="text-xs text-gray-400 mt-2">
                  You'll be logged out after changing your username
                </p>
              </div>
            </details>

            <details className="bg-gray-800 rounded">
              <summary className="cursor-pointer p-3 font-semibold hover:bg-gray-700">
                Change Password
              </summary>
              <div className="p-3 pt-0">
                <p className="text-sm text-gray-400">
                  Password change feature coming soon
                </p>
              </div>
            </details>

            <details className="bg-gray-800 rounded">
              <summary className="cursor-pointer p-3 font-semibold hover:bg-gray-700">
                Change Email
              </summary>
              <div className="p-3 pt-0">
                <p className="text-sm text-gray-400">
                  Email change feature coming soon
                </p>
              </div>
            </details>
          </div>

          <button
            onClick={() => signOut()}
            className="w-full mt-6 bg-gray-700 hover:bg-gray-600 text-white font-semibold py-2 rounded transition"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  );
}
