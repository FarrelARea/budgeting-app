import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@budgeting/ui", "@budgeting/types"],
  typedRoutes: true
};

export default nextConfig;
