import { defineConfig, loadEnv } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import path from 'path'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // Load env from the project root (one level up from frontend/)
  const env = loadEnv(mode, path.resolve(__dirname, '..'), '')

  return {
    plugins: [svelte()],
    // Make VITE_ prefixed env vars available to the app
    envDir: path.resolve(__dirname, '..'),
    server: {
      proxy: {
        '/api': {
          target: `http://localhost:${env.BACKEND_PORT || 8080}`,
          changeOrigin: true,
          secure: false,
        }
      }
    }
  }
})
