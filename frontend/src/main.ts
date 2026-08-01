import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

const app = mount(App, {
  target: document.getElementById('app')!,
})

// RNF-03: cache rules content offline. Only in production builds (a SW in dev
// would break Vite HMR). The SW precaches /index.html + caches /api/v1/content.
if ('serviceWorker' in navigator && import.meta.env.PROD) {
  window.addEventListener('load', () => {
    navigator.serviceWorker
      .register('/sw.js')
      .catch((err) => console.error('service worker registration failed', err))
  })
}

export default app
