import { error } from '@sveltejs/kit';

import type { Code } from './types';

export async function fetchPaste(fetch: typeof globalThis.fetch, id: string): Promise<Code> {
  const res = await fetch(`/api/v1/code/${encodeURIComponent(id)}`);
  if (res.status === 404) {
    throw error(404, 'Paste not found');
  }
  if (!res.ok) {
    throw error(502, 'The paste service is unavailable.');
  }
  const paste = (await res.json()).Data;
  // An id of '.' resolves to the list endpoint
  if (typeof paste?.content !== 'string') {
    throw error(404, 'Paste not found');
  }
  return paste;
}
