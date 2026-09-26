import type { PageLoad } from './$types';

import type { Code } from '$lib/types';

export const load: PageLoad = async ({ fetch }) => {
  try {
    const res = await fetch('/api/v1/code');
    if (res.ok) {
      const body = await res.json();
      return { recentItems: (body.Data ?? []) as Code[], recentFailed: false };
    }
  } catch {
    // Fall through: the editor stays usable without the list.
  }
  return { recentItems: [] as Code[], recentFailed: true };
};
