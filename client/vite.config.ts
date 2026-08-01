import { defineConfig, loadEnv, type ProxyOptions } from 'vite'

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd(), '')
  const backend = (env.VITE_API_BASE_URL || 'http://localhost:8080').replace(/\/+$/, '')

  const authProxy: ProxyOptions = {
    target: backend,
    changeOrigin: true,
    // Only forward POST requests to the backend; let Vite serve everything
    // else (e.g. /register.html) so the proxy prefix never shadows pages.
    bypass: (req) => (req.method === 'POST' ? undefined : req.url),
  }

  return {
    server: {
      proxy: {
        '/register': authProxy,
        '/login': authProxy,
      },
    },
  }
})
