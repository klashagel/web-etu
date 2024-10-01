/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  // Add this line to explicitly set the app directory as the source
  experimental: {
    appDir: true,
  },
}

module.exports = nextConfig