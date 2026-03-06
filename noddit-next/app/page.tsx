"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
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

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

export default function Home() {
  const [popularPosts, setPopularPosts] = useState<Post[]>([]);
  const [activeSubnoddits, setActiveSubnoddits] = useState<Subnoddit[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [posts, subnoddits] = await Promise.all([
        api.get<Post[]>("/api/public/popularposts"),
        api.get<Subnoddit[]>("/api/public/subnoddits/active"),
      ]);
      setPopularPosts(posts);
      setActiveSubnoddits(subnoddits);
    } catch (error) {
      console.error("Failed to fetch data:", error);
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

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
      {/* Main content */}
      <div className="md:col-span-2">
        <h1 className="text-4xl font-bold mb-6">Popular Today</h1>

        {popularPosts.length === 0 ? (
          <div className="bg-gray-900 rounded-lg p-12 text-center">
            <p className="text-gray-400 mb-4">No popular posts today</p>
            <Link
              href="/submit"
              className="inline-block bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition"
            >
              Create the first post
            </Link>
          </div>
        ) : (
          <div className="space-y-4">
            {popularPosts.map((post) => (
              <PostCard key={post.postId} post={post} onRefresh={fetchData} />
            ))}
          </div>
        )}
      </div>

      {/* Sidebar */}
      <div className="space-y-6">
        <div className="bg-gray-900 rounded-lg p-6">
          <h2 className="text-xl font-bold mb-4">Active Communities</h2>
          <div className="space-y-3">
            {activeSubnoddits.map((sn) => (
              <Link
                key={sn.subnodditId}
                href={`/s/${sn.subnodditName}`}
                className="block hover:text-orange-500 transition"
              >
                <div className="font-semibold">n/{sn.subnodditName}</div>
                <div className="text-sm text-gray-400">{sn.subnodditDescription}</div>
              </Link>
            ))}
          </div>

          <Link
            href="/subnoddits"
            className="block mt-4 text-orange-500 hover:underline text-sm"
          >
            View all communities →
          </Link>
        </div>

        <div className="bg-gray-900 rounded-lg p-6">
          <h2 className="text-xl font-bold mb-4">Create</h2>
          <div className="space-y-2">
            <Link
              href="/submit"
              className="block w-full bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded text-center transition"
            >
              Create Post
            </Link>
            <Link
              href="/subnoddits/create"
              className="block w-full bg-gray-700 hover:bg-gray-600 text-white font-semibold py-2 px-4 rounded text-center transition"
            >
              Create Community
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
