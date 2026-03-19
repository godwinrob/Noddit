"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useUser, SignInButton, SignUpButton, UserButton } from "@clerk/nextjs";
import { api } from "@/lib/api";
import { useNodditUser } from "@/components/ClerkTokenProvider";

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
            <Link href="/" className="text-2xl font-bold text-orange-500">
              Noddit
            </Link>

            <div className="hidden md:flex gap-4">
              <Link
                href="/subnoddits"
                className="text-gray-300 hover:text-white transition"
              >
                Communities
              </Link>

              {isSignedIn && favorites.length > 0 && (
                <div className="relative group">
                  <button className="text-gray-300 hover:text-white transition">
                    Favorites ▾
                  </button>
                  <div className="absolute left-0 mt-2 w-48 bg-gray-800 rounded shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-50">
                    {favorites.map((fav) => (
                      <Link
                        key={fav.subnodditName}
                        href={`/s/${fav.subnodditName}`}
                        className="block px-4 py-2 text-gray-300 hover:bg-gray-700 hover:text-white"
                      >
                        n/{fav.subnodditName}
                      </Link>
                    ))}
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Search Bar */}
          <form onSubmit={handleSearch} className="hidden lg:block flex-1 max-w-md mx-8">
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder="Search communities..."
              className="w-full bg-gray-800 border border-gray-700 rounded px-4 py-2 focus:outline-none focus:border-orange-500"
            />
          </form>

          <div className="flex items-center gap-4">
            {isSignedIn ? (
              <>
                <Link
                  href="/submit"
                  className="bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition"
                >
                  Create Post
                </Link>

                <div className="flex items-center gap-3">
                  <Link
                    href="/profile"
                    className="text-gray-300 hover:text-white transition"
                  >
                    u/{username}
                  </Link>
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
                  <button className="text-gray-300 hover:text-white transition">
                    Log In
                  </button>
                </SignInButton>
                <SignUpButton mode="modal">
                  <button className="bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition">
                    Sign Up
                  </button>
                </SignUpButton>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
