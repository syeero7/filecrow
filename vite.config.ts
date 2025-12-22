import { defineConfig } from "vite";
import tailwindcss from "@tailwindcss/vite";
import { ViteMinifyPlugin } from "vite-plugin-minify";

export default defineConfig({
  plugins: [tailwindcss(), ViteMinifyPlugin()],
});
