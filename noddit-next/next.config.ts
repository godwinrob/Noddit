import type { NextConfig } from "next";
import { resolve } from "path";
import { loadEnvConfig } from "@next/env";

// Load env from monorepo root .env (for local development; no-op in Docker)
loadEnvConfig(resolve(process.cwd(), ".."));

const nextConfig: NextConfig = {
  output: 'standalone',
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "**",
      },
    ],
  },
};

export default nextConfig;
