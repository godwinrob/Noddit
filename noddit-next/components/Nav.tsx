"use client";

import { useState, useEffect } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAuth } from "@/lib/auth-context";
import { api } from "@/lib/api";

interface Favorite {
  subnodditName: string;
}

export default function Nav() {
  const { user, logout, isAuthenticated } = useAuth();
  const router = useRouter();

  const [searchTerm, setSearchTerm] = useState("");
  const [favorites, setFavorites] = useState<Favorite[]>([]);

  useEffect(() => {
    if (isAuthenticated && user) {
      fetchFavorites();
    }
  }, [isAuthenticated, user]);

  const fetchFavorites = async () => {
    if (!user) return;

    try {
      const data = await api.get<Favorite[]>(
        `/api/favorites/${user.username}`,
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

              {isAuthenticated && favorites.length > 0 && (
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
            {isAuthenticated ? (
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
                    u/{user?.username}
                  </Link>
                  <button
                    onClick={logout}
                    className="text-gray-400 hover:text-white transition text-sm"
                  >
                    Logout
                  </button>
                </div>
              </>
            ) : (
              <>
                <Link
                  href="/login"
                  className="text-gray-300 hover:text-white transition"
                >
                  Log In
                </Link>
                <Link
                  href="/register"
                  className="bg-orange-600 hover:bg-orange-700 text-white font-semibold py-2 px-4 rounded transition"
                >
                  Sign Up
                </Link>
              </>
            )}
          </div>
        </div>
      </div>
    </nav>
  );
}
