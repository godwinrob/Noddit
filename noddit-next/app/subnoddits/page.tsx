"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { api } from "@/lib/api";

interface Subnoddit {
  subnodditId: number;
  subnodditName: string;
  subnodditDescription: string;
}

export default function AllSubnodditsPage() {
  const [subnoddits, setSubnoddits] = useState<Subnoddit[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSubnoddits();
  }, []);

  const fetchSubnoddits = async () => {
    try {
      const data = await api.get<Subnoddit[]>("/api/public/subnoddits");
      setSubnoddits(data);
    } catch (error) {
      console.error("Failed to fetch subnoddits:", error);
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
    <div className="max-w-4xl mx-auto">
      <h1 className="text-4xl font-bold mb-6">All Communities</h1>

      <div className="grid gap-4">
        {subnoddits.map((sn) => (
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
    </div>
  );
}
