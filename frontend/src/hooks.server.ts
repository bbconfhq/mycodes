import type { HandleFetch } from '@sveltejs/kit';

import { dev } from '$app/environment';
import { env } from '$env/dynamic/private';

// In development the API is not behind the same origin, so server-side
// fetches to /api are sent to it directly.
export const handleFetch: HandleFetch = async ({ event, request, fetch }) => {
  if (dev) {
    const url = new URL(request.url);
    if (url.origin === event.url.origin && url.pathname.startsWith('/api/')) {
      const target = new URL(url.pathname + url.search, env.API_URL ?? 'http://server:4000');
      request = new Request(target, request);
    }
  }
  return fetch(request);
};
