<script lang="ts">
  import { languageLabel } from './language';
  import RelativeTime from './RelativeTime.svelte';
  import type { Code } from './types';

  export let items: Code[] = [];
  export let failed = false;
</script>

<section aria-labelledby="recent-heading">
  <h2 id="recent-heading">Recent pastes</h2>
  {#if failed}
    <p class="empty">Recent pastes could not be loaded. Refresh to try again.</p>
  {:else if items.length === 0}
    <p class="empty">No pastes yet. Share one above and it will show up here for 7 days.</p>
  {:else}
    <ul class="panel">
      {#each items as item (item.id)}
        <li>
          <a href={`/${item.id}`}>
            <span class="title">{item.title}</span>
            <span class="meta">
              <span class="badge">{languageLabel(item.language)}</span>
              <span class="author">{item.name}</span>
              <RelativeTime datetime={item.created_at} />
            </span>
          </a>
        </li>
      {/each}
    </ul>
  {/if}
</section>

<style lang="scss">
  section {
    margin-top: 2.5rem;
  }

  h2 {
    margin-bottom: 0.625rem;
    color: var(--text-muted);
    font-size: 13px;
    font-weight: 600;
    letter-spacing: 0.04em;
    text-transform: uppercase;
  }

  ul {
    list-style: none;
    overflow: hidden;
  }

  li:not(:last-child) {
    border-bottom: 1px solid var(--border);
  }

  a {
    display: flex;
    align-items: center;
    gap: 1rem;
    min-height: 3rem;
    padding: 0.5rem 1rem;
    color: var(--text);
    text-decoration: none;
    transition: background-color 0.15s;

    &:hover {
      background: var(--surface-2);

      .title {
        color: var(--accent-text);
      }
    }

    &:focus-visible {
      outline-offset: -2px;
    }
  }

  .title {
    flex: 1;
    min-width: 0;
    overflow: hidden;
    font-weight: 500;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .meta {
    display: flex;
    flex-shrink: 0;
    align-items: center;
    gap: 0.75rem;
    color: var(--text-muted);
    font-size: 13px;
    font-variant-numeric: tabular-nums;
  }

  .author {
    max-width: 10rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .empty {
    padding: 1.25rem 1rem;
    border: 1px dashed var(--border-strong);
    border-radius: 10px;
    color: var(--text-muted);
    font-size: 14px;
    text-align: center;
  }

  @media (max-width: 560px) {
    a {
      flex-direction: column;
      align-items: flex-start;
      gap: 0.25rem;
    }

    .title {
      width: 100%;
    }
  }
</style>
