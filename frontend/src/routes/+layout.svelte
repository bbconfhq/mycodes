<script lang="ts">
  import 'pretendard/dist/web/variable/pretendardvariable-dynamic-subset.css';
  import './app.css';

  import { page } from '$app/stores';
</script>

<a class="skip-link" href="#content">Skip to content</a>
<div id="root">
  <header>
    <a class="brand" href="/" translate="no" aria-label="mycod.es home">
      mycod<span>.es</span>
    </a>
    <p class="tagline">Share code for 7 days</p>
    {#if $page.url.pathname !== '/'}
      <a class="btn new" href="/">New paste</a>
    {/if}
  </header>
  <main id="content">
    <slot />
  </main>
</div>

<style lang="scss">
  #root {
    max-width: 1040px;
    margin: 0 auto;
    padding: 1.75rem 1.25rem 3rem;
  }

  header {
    display: flex;
    align-items: baseline;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .brand {
    font-family: var(--font-mono);
    font-size: 1.375rem;
    font-weight: 700;
    letter-spacing: -0.02em;
    color: var(--text);
    text-decoration: none;

    span {
      color: var(--accent);
    }
  }

  .tagline {
    color: var(--text-muted);
    font-size: 14px;
  }

  .new {
    margin-left: auto;
    align-self: center;
  }

  .skip-link {
    position: absolute;
    left: 1rem;
    top: -3rem;
    z-index: 10;
    padding: 0.5rem 0.75rem;
    border-radius: 7px;
    background: var(--accent);
    color: var(--on-accent);

    &:focus {
      top: 1rem;
    }
  }

  @media (max-width: 480px) {
    #root {
      padding: 1.25rem 1rem 2rem;
    }

    .tagline {
      display: none;
    }
  }
</style>
