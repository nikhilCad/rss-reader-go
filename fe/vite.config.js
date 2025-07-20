import { defineConfig } from "vite";
import preact from "@preact/preset-vite";

export default defineConfig({
  plugins: [preact()],
  server: {
    proxy: {
      // Proxy API requests to the Go backend
      "/feeds": "http://localhost:8080",
      "/posts": "http://localhost:8080",
      "/read": "http://localhost:8080",
      "/unread": "http://localhost:8080",
      "/refresh": "http://localhost:8080",
      // Add any other API endpoints you use
    },
  },
});
