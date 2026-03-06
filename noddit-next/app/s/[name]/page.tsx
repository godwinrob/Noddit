"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";
import { useAuth } from "@/lib/auth-context";
import PostCard from "@/components/PostCard";

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
}

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

interface Favorite {
  subnodditId: number;
}

export default function SubnodditPage() {
  const params = useParams();
  const subnodditName = params.name as string;

  const [posts, setPosts] = useState<Post[]>([]);
  const [subnoddit, setSubnoddit] = useState<Subnoddit | null>(null);
  const [showPopular, setShowPopular] = useState(false);
  const [isFavorited, setIsFavorited] = useState(false);
  const [loading, setLoading] = useState(true);

  const { user, isAuthenticated } = useAuth();

  useEffect(() => {
    fetchData();
  }, [subnodditName, showPopular]);

  const fetchData = async () => {
    try {
      setLoading(true);

      // Fetch subnoddit info
      const snData = await api.get<Subnoddit>(
        `/api/public/subnoddits/${subnodditName}`
      );
      setSubnoddit(snData);

      // Fetch posts (recent or popular)
      const postsData = showPopular
        ? await api.get<Post[]>(`/api/public/allpostspopular/${subnodditName}`)
        : await api.get<Post[]>(`/api/public/allposts/${subnodditName}`);
      setPosts(postsData);

      // Check if favorited (if logged in)
      if (isAuthenticated && user) {
        const favorites = await api.get<Favorite[]>(
          `/api/favorites/${user.username}`,
          true
        );
        setIsFavorited(
          favorites.some((f) => f.subnodditId === snData.subnodditId)
        );
      }
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      setLoading(false);
    }
  };

  const toggleFavorite = async () => {
    if (!isAuthenticated || !user || !subnoddit) return;

    try {
      if (isFavorited) {
        await api.delete(
          `/api/favorites/subnoddit/${subnoddit.subnodditId}`,
          true
        );
      } else {
        await api.post(
          "/api/favorites/create/subnoddit",
          {
            username: user.username,
            subnodditName: subnoddit.subnodditName,
          },
          true
        );
      }
      setIsFavorited(!isFavorited);
    } catch (error) {
      console.error("Failed to toggle favorite:", error);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <div className="text-gray-400">Loading...</div>
      </div>
    );
  }

  if (!subnoddit) {
    return (
      <div className="text-center py-12">
        <h1 className="text-2xl font-bold mb-4">Community Not Found</h1>
        <Link href="/subnoddits" className="text-orange-500 hover:underline">
          Browse all communities
        </Link>
      </div>
    );
  }

  return (
    <div>
      {/* Subnoddit Header */}
      <div className="bg-gray-900 rounded-lg p-6 mb-6">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-3xl font-bold mb-2">n/{subnoddit.subnodditName}</h1>
            <p className="text-gray-400">{subnoddit.subnodditDescription}</p>
          </div>

          {isAuthenticated && (
            <button
              onClick={toggleFavorite}
              className={`px-4 py-2 rounded font-semibold transition ${
                isFavorited
                  ? "bg-gray-700 hover:bg-gray-600 text-white"
                  : "bg-orange-600 hover:bg-orange-700 text-white"
              }`}
            >
              {isFavorited ? "Unfavorite" : "Favorite"}
            </button>
          )}
        </div>

        {/* Sort Toggle */}
        <div className="mt-4 flex gap-2">
          <button
            onClick={() => setShowPopular(false)}
            className={`px-4 py-2 rounded transition ${
              !showPopular
                ? "bg-orange-600 text-white"
                : "bg-gray-800 text-gray-300 hover:bg-gray-700"
            }`}
          >
            Recent
          </button>
          <button
            onClick={() => setShowPopular(true)}
            className={`px-4 py-2 rounded transition ${
              showPopular
                ? "bg-orange-600 text-white"
                : "bg-gray-800 text-gray-300 hover:bg-gray-700"
            }`}
          >
            Popular
          </button>
        </div>
      </div>

      {/* Posts List */}
      {posts.length === 0 ? (
        <div className="bg-gray-900 rounded-lg p-12 text-center">
          <p className="text-gray-400 mb-4">No posts yet in this community</p>
          <Link
            href="/submit"
            className="inline-block bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition"
          >
            Create the first post
          </Link>
        </div>
      ) : (
        <div className="space-y-4">
          {posts.map((post) => (
            <PostCard key={post.postId} post={post} onRefresh={fetchData} />
          ))}
        </div>
      )}
    </div>
  );
}
