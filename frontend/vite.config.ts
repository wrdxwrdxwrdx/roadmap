import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    host: '0.0.0.0', // Для работы в Docker
    port: 3000,
    watch: {
      usePolling: true, // Необходимо для hot reload в Docker
    },
    proxy: {
      '/api': {
        target: process.env.VITE_API_URL || 'http://api:8080',
        changeOrigin: true,
        secure: false,
        rewrite: (path) => path.replace(/^\/api/, ''),
      },
    },
  },
})

