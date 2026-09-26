import type { RequestHandler } from './$types';

import { fetchPaste } from '$lib/paste';

export const GET: RequestHandler = async ({ params, fetch }) => {
  const paste = await fetchPaste(fetch, params.id);
  return new Response(paste.content, {
    headers: {
      'Content-Type': 'text/plain; charset=utf-8',
      'X-Content-Type-Options': 'nosniff',
      // Pastes expire, so never serve a cached copy
      'Cache-Control': 'no-store'
    }
  });
};
