import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
	plugins: [react()],
	build: {
		outDir: "static/react", // Output React files into your static folder
		assetsDir: "assets",    // Keep assets organized
		rollupOptions: {
			input: "./components/react/main.jsx",
			output: {
				entryFileNames: "main.js",  // Keep the same name for the main JS file
				chunkFileNames: "chunk.js", // Name for dynamically imported chunks
				assetFileNames: "assets/[name].[ext]", // For other assets like images, etc.
			}
		}
	},
});

