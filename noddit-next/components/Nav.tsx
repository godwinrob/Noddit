"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import Image from "next/image";
import { useRouter } from "next/navigation";
import { useUser, SignInButton, SignUpButton, UserButton } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

interface Favorite {
  subnodditName: string;
}

export default function Nav() {
  const { isSignedIn } = useUser();
  const { username } = useNodditUser();
  const router = useRouter();

  const [searchTerm, setSearchTerm] = useState("");
  const [favorites, setFavorites] = useState<Favorite[]>([]);

  useEffect(() => {
    if (isSignedIn && username) {
      fetchFavorites();
    }
  }, [isSignedIn, username]);

  const fetchFavorites = async () => {
    if (!username) return;

    try {
      const data = await api.get<Favorite[]>(
        `/api/favorites/${username}`,
        true
      );
      setFavorites(data);
    } catch (error) {
      console.error("Failed to fetch favorites:", error);
    }
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (searchTerm.trim()) {
      router.push(`/search/${encodeURIComponent(searchTerm)}`);
      setSearchTerm("");
    }
  };

  return (
    <nav className="bg-gray-900 border-b border-gray-800">
      <div className="container mx-auto px-4 max-w-6xl">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center gap-8">
            <Link href="/" className="flex items-center">
              <Image
                src="/assets/img/Noddit - Dark.png"
                alt="Noddit Logo"
                width={120}
                height={40}
                className="h-10 w-auto"
                priority
              />
            </Link>

            <div className="hidden md:flex gap-4">
              <Link
                href="/subnoddits"
                className="text-gray-300 hover:text-white transition"
              >
                Communities
              </Link>

              {isSignedIn && favorites.length > 0 && (
                <DropdownMenu>
                  <DropdownMenuTrigger>
                    <Button variant="ghost" className="text-gray-300 hover:text-white">
                      Favorites ▾
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent className="w-48 bg-gray-800 border-gray-700">
                    {favorites.map((fav) => (
                      <DropdownMenuItem key={fav.subnodditName}>
                        <Link
                          href={`/s/${fav.subnodditName}`}
                          className="text-gray-300 hover:text-white cursor-pointer w-full"
                        >
                          n/{fav.subnodditName}
                        </Link>
                      </DropdownMenuItem>
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
            </div>
          </div>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="hidden lg:block flex-1 max-w-md mx-8">
            <Input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search communities..."
              className="w-full bg-gray-800 border-gray-700 focus:border-orange-500"
            />
          </form>

          <div className="flex items-center gap-4">
            {isSignedIn ? (
              <>
                <Link href="/submit">
                  <Button className="bg-orange-600 hover:bg-orange-700">
                    Create Post
                  </Button>
                </Link>

                <div className="flex items-center gap-3">
                  {username && (
                    <Link
                      href="/profile"
                      className="text-gray-300 hover:text-white transition"
                    >
                      u/{username}
                    </Link>
                  )}
                  <UserButton
                    appearance={{
                      elements: {
                        avatarBox: "w-10 h-10"
                      }
                    }}
                  />
                </div>
              </>
            ) : (
              <>
                <SignInButton mode="modal">
                  <Button variant="ghost" className="text-gray-300 hover:text-white">
                    Log In
                  </Button>
                </SignInButton>
                <SignUpButton mode="modal">
                  <Button className="bg-orange-600 hover:bg-orange-700">
                    Sign Up
                  </Button>
                </SignUpButton>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
