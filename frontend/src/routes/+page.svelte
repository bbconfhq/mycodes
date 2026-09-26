<script lang="ts">
  import type { PageData } from './$types';

  import Editor from '$lib/Editor.svelte';
  import RecentList from '$lib/RecentList.svelte';

  export let data: PageData;
</script>

<svelte:head>
  <title>mycod.es – Share code for 7 days</title>
  <meta
    name="description"
    content="Paste code, get a link with syntax highlighting and line anchors. Pastes expire after 7 days."
  />
</svelte:head>

<h1 class="sr-only">New paste</h1>
<Editor />
<RecentList items={data.recentItems} failed={data.recentFailed} />
