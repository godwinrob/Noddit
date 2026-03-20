"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";
import PostCard from "@/components/PostCard";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Separator } from "@/components/ui/separator";

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
  const [activeTab, setActiveTab] = useState("popular");

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
      console.log("Popular posts:", posts);
      console.log("Recent posts:", newPostsData);
      console.log("Active subnoddits:", subnoddits);
      setPopularPosts(posts);
      setNewPosts(newPostsData);
      setActiveSubnoddits(subnoddits);
    } catch (error) {
      console.error("Failed to fetch data:", error);
      if (error instanceof Error) {
        console.error("Error message:", error.message);
      }
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
        <div className="md:col-span-2 space-y-4">
          <Skeleton className="h-12 w-64 bg-gray-800" />
          {[1, 2, 3].map((i) => (
            <Card key={i} className="bg-gray-900 border-gray-800">
              <CardContent className="p-4">
                <Skeleton className="h-24 w-full bg-gray-800" />
              </CardContent>
            </Card>
          ))}
        </div>
        <div className="space-y-6">
          <Card className="bg-gray-900 border-gray-800">
            <CardContent className="p-6">
              <Skeleton className="h-32 w-full bg-gray-800" />
            </CardContent>
          </Card>
        </div>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
      {/* Main content */}
      <div className="md:col-span-2">
        <h1 className="text-4xl font-bold mb-6">Noddit</h1>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-2 bg-gray-800/50">
            <TabsTrigger value="popular">Popular</TabsTrigger>
            <TabsTrigger value="new">New</TabsTrigger>
          </TabsList>

          <TabsContent value="popular" className="mt-6">
            {popularPosts.length === 0 ? (
              <Card className="bg-gray-900 border-gray-800 p-12 text-center">
                <p className="text-gray-400 mb-4">No popular posts today</p>
                <Link href="/submit">
                  <Button className="bg-orange-600 hover:bg-orange-700">
                    Create the first post
                  </Button>
                </Link>
              </Card>
            ) : (
              <div className="space-y-4">
                {popularPosts.map((post) => (
                  <PostCard key={post.postId} post={post} onRefresh={fetchData} />
                ))}
              </div>
            )}
          </TabsContent>

          <TabsContent value="new" className="mt-6">
            {newPosts.length === 0 ? (
              <Card className="bg-gray-900 border-gray-800 p-12 text-center">
                <p className="text-gray-400 mb-4">No posts yet</p>
                <Link href="/submit">
                  <Button className="bg-orange-600 hover:bg-orange-700">
                    Create the first post
                  </Button>
                </Link>
              </Card>
            ) : (
              <div className="space-y-4">
                {newPosts.map((post) => (
                  <PostCard key={post.postId} post={post} onRefresh={fetchData} />
                ))}
              </div>
            )}
          </TabsContent>
        </Tabs>
      </div>

      {/* Sidebar */}
      <div className="space-y-6">
        <Card className="bg-gray-900 border-gray-800">
          <CardHeader>
            <CardTitle className="text-xl">Active Communities</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            {activeSubnoddits.map((sn, index) => (
              <div key={sn.subnodditId}>
                {index > 0 && <Separator className="my-3" />}
                <Link
                  href={`/s/${sn.subnodditName}`}
                  className="block hover:text-orange-500 transition"
                >
                  <div className="font-semibold">n/{sn.subnodditName}</div>
                  <div className="text-sm text-gray-400">{sn.subnodditDescription}</div>
                </Link>
              </div>
            ))}

            <Separator className="my-3" />
            <Link
              href="/subnoddits"
              className="block text-orange-500 hover:underline text-sm"
            >
              View all communities →
            </Link>
          </CardContent>
        </Card>

        <Card className="bg-gray-900 border-gray-800">
          <CardHeader>
            <CardTitle className="text-xl">Create</CardTitle>
          </CardHeader>
          <CardContent className="space-y-2">
            <Link href="/submit" className="block">
              <Button className="w-full bg-orange-600 hover:bg-orange-700">
                Create Post
              </Button>
            </Link>
            <Link href="/subnoddits/create" className="block">
              <Button variant="secondary" className="w-full">
                Create Community
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
