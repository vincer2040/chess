import { defineConfig } from "vite";

export default defineConfig({
    plugins: [
    ],
    publicDir: false,
    server: {
        port: 3000,
    },
    build: {
        target: 'esnext',
        manifest: true,
        assetsDir: "assets",
        rollupOptions: {
            input: "src/main.js",
        },
    },
});
