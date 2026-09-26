<script lang="ts">
  import { page } from '$app/stores';
</script>

<svelte:head>
  <title>{$page.status === 404 ? 'Paste not found' : 'Error'} – mycod.es</title>
</svelte:head>

<section class="panel">
  <p class="status">{$page.status}</p>
  {#if $page.status === 404}
    <h1>This paste doesn’t exist or has expired</h1>
    <p>Pastes are kept for 7 days. Ask the author to share it again.</p>
  {:else}
    <h1>Something went wrong</h1>
    <p>{$page.error?.message ?? 'The paste could not be loaded.'} Try again in a moment.</p>
  {/if}
  <a class="btn btn-primary" href="/">Create a new paste</a>
</section>

<style lang="scss">
  section {
    max-width: 34rem;
    margin: 4rem auto;
    padding: 2.5rem 2rem;
    text-align: center;
  }

  .status {
    font-family: var(--font-mono);
    font-size: 3rem;
    font-weight: 700;
    line-height: 1;
    color: var(--accent);
  }

  h1 {
    margin: 1rem 0 0.5rem;
    font-size: 1.25rem;
    text-wrap: balance;
  }

  p:not(.status) {
    margin-bottom: 1.5rem;
    color: var(--text-muted);
  }
</style>
