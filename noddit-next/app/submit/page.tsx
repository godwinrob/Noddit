"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useUser } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import { Alert, AlertDescription } from "@/components/ui/alert";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { toast } from "sonner";

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

export default function SubmitPage() {
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [imageAddress, setImageAddress] = useState("");
  const [subnodditId, setSubnodditId] = useState("");
  const [subnoddits, setSubnoddits] = useState<Subnoddit[]>([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const { isSignedIn } = useUser();
  const { username: nodditUsername } = useNodditUser();
  const router = useRouter();

  useEffect(() => {
    if (!isSignedIn) {
      router.push("/");
      return;
    }

    // Fetch all subnoddits
    const fetchSubnoddits = async () => {
      try {
        const data = await api.get<Subnoddit[]>("/api/public/subnoddits");
        setSubnoddits(data);
        if (data.length > 0) {
          setSubnodditId(data[0].subnodditId.toString());
        }
      } catch (error) {
        console.error("Failed to fetch subnoddits:", error);
      }
    };

    fetchSubnoddits();
  }, [isSignedIn, router]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    if (!title.trim()) {
      setError("Title is required");
      setLoading(false);
      return;
    }

    if (!body.trim()) {
      setError("Body is required");
      setLoading(false);
      return;
    }

    if (!subnodditId) {
      setError("Please select a community");
      setLoading(false);
      return;
    }

    try {
      const post = await api.post<{ postId: number }>(
        "/api/post/create",
        {
          title,
          body,
          imageAddress: imageAddress || null,
          username: nodditUsername,
          subnodditId: parseInt(subnodditId),
        },
        true
      );

      // Get the subnoddit name for redirect
      const subnoddit = subnoddits.find(
        (s) => s.subnodditId === parseInt(subnodditId)
      );

      toast.success("Post created successfully!");

      if (subnoddit && post.postId) {
        router.push(`/n/${subnoddit.subnodditName}/${post.postId}`);
      } else {
        router.push("/");
      }
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : "Failed to create post";
      setError(errorMessage);
      toast.error(errorMessage);
      setLoading(false);
    }
  };

  if (!isSignedIn) {
    return null; // Will redirect
  }

  return (
    <div className="max-w-3xl mx-auto mt-8">
      <Card className="bg-gray-900 border-gray-800">
        <CardHeader>
          <CardTitle className="text-3xl">Create a Post</CardTitle>
        </CardHeader>
        <CardContent>
          {error && (
            <Alert variant="destructive" className="mb-4">
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <form onSubmit={handleSubmit} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="community">Choose a community</Label>
            <Select
              value={subnodditId}
              onValueChange={(value) => value && setSubnodditId(value)}
              required
            >
              <SelectTrigger className="bg-gray-800 border-gray-700">
                <SelectValue placeholder="Select a community" />
              </SelectTrigger>
              <SelectContent className="bg-gray-800 border-gray-700">
                {subnoddits.map((sn) => (
                  <SelectItem key={sn.subnodditId} value={sn.subnodditId.toString()}>
                    n/{sn.subnodditName}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="title">Title</Label>
            <Input
              id="title"
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              placeholder="An interesting title"
              maxLength={100}
              className="bg-gray-800 border-gray-700 focus:border-orange-500"
              required
            />
            <p className="text-sm text-gray-400">{title.length}/100</p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="body">Text</Label>
            <Textarea
              id="body"
              value={body}
              onChange={(e) => setBody(e.target.value)}
              placeholder="What are your thoughts?"
              rows={8}
              maxLength={2000}
              className="bg-gray-800 border-gray-700 focus:border-orange-500 resize-none"
              required
            />
            <p className="text-sm text-gray-400">{body.length}/2000</p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="image">Image URL (optional)</Label>
            <Input
              id="image"
              type="url"
              value={imageAddress}
              onChange={(e) => setImageAddress(e.target.value)}
              placeholder="https://example.com/image.jpg"
              className="bg-gray-800 border-gray-700 focus:border-orange-500"
            />
          </div>

          <div className="flex gap-3 pt-4">
            <Button
              type="submit"
              disabled={loading}
              className="flex-1 bg-orange-600 hover:bg-orange-700"
            >
              {loading ? "Posting..." : "Post"}
            </Button>
            <Button
              type="button"
              variant="secondary"
              onClick={() => router.back()}
            >
              Cancel
            </Button>
          </div>
        </form>
        </CardContent>
      </Card>

      <Card className="bg-gray-900 border-gray-800 mt-4">
        <CardContent className="pt-6">
          <h2 className="font-semibold mb-3">Posting Guidelines</h2>
          <Separator className="mb-3" />
          <ul className="text-sm text-gray-400 space-y-2">
            <li className="flex items-start gap-2">
              <span className="text-orange-500 mt-0.5">•</span>
              <span>Be respectful and civil</span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-orange-500 mt-0.5">•</span>
              <span>Stay on topic for the community</span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-orange-500 mt-0.5">•</span>
              <span>No spam or self-promotion</span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-orange-500 mt-0.5">•</span>
              <span>Use descriptive titles</span>
            </li>
          </ul>
        </CardContent>
      </Card>
    </div>
  );
}
