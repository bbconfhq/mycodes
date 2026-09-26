import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [sveltekit()],
  server: {
    proxy: {
      '/api': {
        target: process.env.API_URL ?? 'http://mycodes-server:4000',
        changeOrigin: true,
        secure: false
      }
    }
  }
});
