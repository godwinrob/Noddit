"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { api } from "@/lib/api";

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

export default function SearchPage() {
  const params = useParams();
  const searchTerm = params.term ? (params.term as string[])[0] : "";

  const [results, setResults] = useState<Subnoddit[]>([]);
  const [loading, setLoading] = useState(false);
  const [searched, setSearched] = useState(false);

  useEffect(() => {
    if (searchTerm) {
      performSearch(searchTerm);
    }
  }, [searchTerm]);

  const performSearch = async (term: string) => {
    if (!term.trim()) return;

    setLoading(true);
    setSearched(true);

    try {
      const data = await api.get<Subnoddit[]>(
        `/api/public/subnoddits/search/${encodeURIComponent(term)}`
      );
      setResults(data);
    } catch (error) {
      console.error("Search failed:", error);
      setResults([]);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <h1 className="text-4xl font-bold mb-6">
        {searchTerm ? `Search Results for "${searchTerm}"` : "Search"}
      </h1>

      {loading && (
        <div className="flex justify-center items-center min-h-[200px]">
          <div className="text-gray-400">Searching...</div>
        </div>
      )}

      {!loading && searched && results.length === 0 && (
        <div className="bg-gray-900 rounded-lg p-12 text-center">
          <p className="text-gray-400 text-lg mb-4">No results found</p>
          <p className="text-gray-500">
            Try searching for a different term or{" "}
            <Link href="/subnoddits" className="text-orange-500 hover:underline">
              browse all communities
            </Link>
          </p>
        </div>
      )}

      {results.length > 0 && (
        <div className="grid gap-4">
          {results.map((sn) => (
            <Link
              key={sn.subnodditId}
              href={`/s/${sn.subnodditName}`}
              className="bg-gray-900 rounded-lg p-6 hover:bg-gray-800 transition"
            >
              <h2 className="text-2xl font-bold mb-2 text-orange-500">
                n/{sn.subnodditName}
              </h2>
              <p className="text-gray-400">{sn.subnodditDescription}</p>
            </Link>
          ))}
        </div>
      )}

      {!searchTerm && !searched && (
        <div className="bg-gray-900 rounded-lg p-12 text-center">
          <p className="text-gray-400">
            Use the search bar above to find communities
          </p>
        </div>
      )}
    </div>
  );
}
