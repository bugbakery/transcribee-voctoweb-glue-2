import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

// https://vite.dev/config/
export default defineConfig({
  base: "",
  plugins: [react()],
  server: {
    proxy: {
      '/api/': {
        target: 'http://localhost:8090',
        changeOrigin: true,
        secure: false,
      },
      '/_/': {
        target: 'http://localhost:8090',
        changeOrigin: true,
        secure: false,
      },
    },
  }
})
