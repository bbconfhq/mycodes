import type { PageLoad } from './$types';

import { fetchPaste } from '$lib/paste';

export const load: PageLoad = async ({ params, fetch }) => ({
  paste: await fetchPaste(fetch, params.id)
});
