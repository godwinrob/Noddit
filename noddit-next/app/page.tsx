"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import PostCardCompact from "@/components/PostCardCompact";

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
  const [newPosts, setNewPosts] = useState<Post[]>([]);
  const [activeSubnoddits, setActiveSubnoddits] = useState<Subnoddit[]>([]);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState<"hot" | "new">("hot");

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [posts, newPostsData, subnoddits] = await Promise.all([
        api.get<Post[]>("/api/public/popularposts"),
        api.get<Post[]>("/api/public/recentposts"),
        api.get<Subnoddit[]>("/api/public/subnoddits/active"),
      ]);
      setPopularPosts(posts);
      setNewPosts(newPostsData);
      setActiveSubnoddits(subnoddits);
    } catch (error) {
      console.error("Failed to fetch data:", error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex gap-4">
        <div className="flex-1">
          <div className="bg-gray-800/50 h-10 mb-2 animate-pulse" />
          <div className="space-y-1">
            {[1, 2, 3].map((i) => (
              <div key={i} className="bg-gray-900 border-l-4 border-gray-700 p-2">
                <div className="bg-gray-800/50 h-16 animate-pulse" />
              </div>
            ))}
          </div>
        </div>
        <div className="w-80">
          <div className="bg-gray-900 border border-gray-700 p-3">
            <div className="bg-gray-800/50 h-32 animate-pulse" />
          </div>
        </div>
      </div>
    );
  }

  const currentPosts = activeTab === "hot" ? popularPosts : newPosts;

  return (
    <div className="flex gap-4">
      {/* Main content */}
      <div className="flex-1">
        {/* Tab navigation - old reddit style */}
        <div className="bg-gray-900/50 border-b border-gray-700 mb-2 px-2">
          <div className="flex gap-1 text-sm">
            <button
              onClick={() => setActiveTab("hot")}
              className={`px-3 py-2 ${
                activeTab === "hot"
                  ? "text-orange-500 border-b-2 border-orange-500"
                  : "text-gray-400 hover:text-gray-200"
              }`}
            >
              hot
            </button>
            <button
              onClick={() => setActiveTab("new")}
              className={`px-3 py-2 ${
                activeTab === "new"
                  ? "text-orange-500 border-b-2 border-orange-500"
                  : "text-gray-400 hover:text-gray-200"
              }`}
            >
              new
            </button>
          </div>
        </div>

        {/* Posts */}
        {currentPosts.length === 0 ? (
          <div className="bg-gray-900 border border-gray-700 p-8 text-center">
            <p className="text-gray-400 mb-4">there doesn't seem to be anything here</p>
            <Link href="/submit" className="text-blue-400 hover:underline text-sm">
              create a new post
            </Link>
          </div>
        ) : (
          <div className="space-y-[1px]">
            {currentPosts.map((post) => (
              <PostCardCompact key={post.postId} post={post} onRefresh={fetchData} />
            ))}
          </div>
        )}
      </div>

      {/* Sidebar - old reddit style */}
      <div className="w-80 space-y-3">
        {/* Submit buttons */}
        <div className="bg-gray-900 border border-gray-700">
          <div className="bg-gradient-to-r from-orange-600 to-orange-500 p-3">
            <h3 className="font-bold text-white">Create</h3>
          </div>
          <div className="p-3 space-y-2">
            <Link href="/submit" className="block">
              <button className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 text-sm">
                Submit a new link
              </button>
            </Link>
            <Link href="/submit" className="block">
              <button className="w-full bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-4 text-sm">
                Submit a new text post
              </button>
            </Link>
          </div>
        </div>

        {/* Communities */}
        <div className="bg-gray-900 border border-gray-700">
          <div className="bg-gradient-to-r from-orange-600 to-orange-500 p-2">
            <h3 className="font-bold text-sm text-white">ACTIVE COMMUNITIES</h3>
          </div>
          <div className="p-3">
            {activeSubnoddits.length === 0 ? (
              <p className="text-xs text-gray-400">No communities yet</p>
            ) : (
              <div className="space-y-2">
                {activeSubnoddits.slice(0, 10).map((sn) => (
                  <div key={sn.subnodditId} className="text-xs">
                    <Link
                      href={`/n/${sn.subnodditName}`}
                      className="text-blue-400 hover:underline font-bold"
                    >
                      /n/{sn.subnodditName}
                    </Link>
                    <p className="text-gray-400 mt-0.5">{sn.subnodditDescription}</p>
                  </div>
                ))}
                <Link
                  href="/subnoddits"
                  className="text-blue-400 hover:underline text-xs block mt-3"
                >
                  view more »
                </Link>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
