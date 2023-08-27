import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { VitePWA } from 'vite-plugin-pwa'
// import basicSSL from '@vitejs/plugin-basic-ssl'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    https: false,
  },
  plugins: [
    react(),
    VitePWA({
      registerType: 'autoUpdate',
      injectRegister: 'auto',
      devOptions: {
        enabled: true
      },
      includeAssets: ['favicon.ico', 'apple-touch-icon.png', 'mask-icon.svg'],
      manifest: {
        id: 'place.safer.app',
        name: 'SaferPlace',
        short_name: 'SaferPlace',
        description: 'Trying to make the world a little bit safer',
        theme_color: '#1976d2',
        display: 'standalone',
        orientation: 'portrait',
        icons: [
          {
            src: 'pwa-64x64.png',
            sizes: '64x64',
            type: 'image/png',
          },{
            src: 'pwa-192x192.png',
            sizes: '192x192',
            type: 'image/png',
          }, {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
          }, {
            src: 'pwa-512x512.png',
            sizes: '512x512',
            type: 'image/png',
            purpose: 'maskable',
          }
        ]
      }
    }),
    // basicSSL(),
  ],
})
