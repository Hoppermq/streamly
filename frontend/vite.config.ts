import { defineConfig } from 'vite'
import viteReact from '@vitejs/plugin-react'
import tailwindcss from '@tailwindcss/vite'

import { resolve } from 'node:path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [viteReact(), tailwindcss()],
  resolve: {
    alias: {
      '@': resolve(__dirname, './src'),
    },
  },
  server: {
    host: '0.0.0.0',
    watch: {
      usePolling: true,
    }
  },
  define: {
    'import.meta.env.ZITADEL_API_URL': JSON.stringify(process.env.ZITADEL_API_URL),
    'import.meta.env.PLATFORM_API_URL': JSON.stringify(process.env.PLATFORM_API_URL),
    'import.meta.env.VITE_ZITADEL_CLIENT_ID': JSON.stringify(process.env.WEB_OIDC_CLIENT_ID),
  }
})
