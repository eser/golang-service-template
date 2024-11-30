/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // output: "standalone",

  experimental: {
    reactCompiler: { compilationMode: "all" },
    // reactCompiler: false,
    // ppr: true, // 'incremental',
  },
};

export default nextConfig;
