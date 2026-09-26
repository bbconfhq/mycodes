<script lang="ts">
  import { tick } from 'svelte';

  import { highlightLines } from '../../lib/highlight';
  import { languageLabel } from '../../lib/language';
  import RelativeTime from '../../lib/RelativeTime.svelte';
  import { copyText } from '../../utils/clipboard';
  import {
    formatLineHash,
    parseLineHash,
    setLineHash,
    type LineRange
  } from '../../utils/navigator';

  import type { PageData } from './$types';

  import { afterNavigate } from '$app/navigation';

  export let data: PageData;

  $: paste = data.paste;
  $: lines = highlightLines(paste.content, paste.language);
  $: gutterWidth = `${String(lines.length).length + 2}ch`;
  // Skip rendering off-screen rows of large pastes
  $: large = lines.length > 1000;

  let range: LineRange | null = null;
  let anchor: number | null = null;
  let dragging = false;
  let focusedLine = 1;
  let wrap = false;

  let copied: 'code' | 'link' | null = null;
  let status = '';
  let copiedTimer: ReturnType<typeof setTimeout>;

  async function applyHash(scroll: boolean) {
    range = parseLineHash(window.location.hash, lines.length);
    anchor = range?.start ?? null;
    focusedLine = range?.start ?? 1;
    if (range && scroll) {
      await tick();
      document.getElementById(`L${range.start}`)?.scrollIntoView({ block: 'center' });
    }
  }

  // On Back/Forward SvelteKit restores the scroll position itself
  afterNavigate(({ type }) => applyHash(type !== 'popstate'));

  function select(line: number, extend: boolean) {
    if (extend && anchor != null) {
      range = { start: Math.min(anchor, line), end: Math.max(anchor, line) };
    } else {
      anchor = line;
      range = { start: line, end: line };
    }
    focusedLine = line;
    setLineHash(range);
  }

  const lineOf = (target: EventTarget | null) =>
    Number((target as Element).closest('[data-line]')?.getAttribute('data-line')) || null;

  function onPointerDown(e: PointerEvent) {
    const target = e.target as Element;
    if (e.button !== 0 || !target.classList.contains('ln')) {
      return;
    }
    const line = Number(target.id.slice(1));
    e.preventDefault();
    select(line, e.shiftKey);
    if (!e.shiftKey && e.pointerType === 'mouse') {
      dragging = true;
    }
  }

  function onPointerOver(e: PointerEvent) {
    if (!dragging) {
      return;
    }
    const line = lineOf(e.target);
    if (line != null) {
      select(line, true);
    }
  }

  function onClick(e: MouseEvent) {
    const target = e.target as Element;
    if (!target.classList.contains('ln')) {
      return;
    }
    // Selection already happened on pointerdown; keep the link from scrolling.
    // detail is 0 for keyboard activation (Enter).
    e.preventDefault();
    if (e.detail === 0) {
      select(Number(target.id.slice(1)), e.shiftKey);
    }
  }

  // Line numbers form one Tab stop; arrow keys move between them.
  function onKeyDown(e: KeyboardEvent) {
    const target = e.target as HTMLElement;
    if (!target.classList.contains('ln') || (e.key !== 'ArrowDown' && e.key !== 'ArrowUp')) {
      return;
    }
    e.preventDefault();
    const line = Number(target.id.slice(1)) + (e.key === 'ArrowDown' ? 1 : -1);
    if (line >= 1 && line <= lines.length) {
      focusedLine = line;
      document.getElementById(`L${line}`)?.focus();
    }
  }

  function flash(kind: 'code' | 'link', ok: boolean, message: string) {
    clearTimeout(copiedTimer);
    copied = ok ? kind : null;
    status = ok ? message : 'Copying failed. Select the text and copy it manually.';
    copiedTimer = setTimeout(() => {
      copied = null;
    }, 2000);
  }

  async function copyCode() {
    let text = paste.content;
    if (range) {
      text = text
        .split(/\r\n?|\n/)
        .slice(range.start - 1, range.end)
        .join('\n');
    }
    flash('code', await copyText(text), range ? 'Selected lines copied' : 'Code copied');
  }

  async function copyLink() {
    const url = `${window.location.origin}/${paste.id}${range ? formatLineHash(range) : ''}`;
    flash('link', await copyText(url), 'Link copied');
  }
</script>

<svelte:head>
  <title>{paste.title} by {paste.name} – mycod.es</title>
  <meta
    name="description"
    content={`${languageLabel(paste.language)} snippet shared on mycod.es`}
  />
</svelte:head>

<svelte:window on:hashchange={() => applyHash(true)} on:pointerup={() => (dragging = false)} />

<article>
  <header class="head">
    <h1>{paste.title}</h1>
    <p class="meta">
      <span class="badge">{languageLabel(paste.language)}</span>
      <span>by <strong>{paste.name}</strong></span>
      <span>IP {paste.ip}</span>
      <span>Created <RelativeTime datetime={paste.created_at} /></span>
      <span>Expires <RelativeTime datetime={paste.expired_at} /></span>
    </p>
  </header>

  <div class="panel viewer">
    <div class="toolbar">
      <span class="count">
        {lines.length.toLocaleString()}
        {lines.length === 1 ? 'line' : 'lines'}
        {#if range}<span class="selection">
            · {range.start === range.end
              ? `line ${range.start}`
              : `lines ${range.start}–${range.end}`} selected</span
          >{/if}
      </span>
      <div class="actions">
        <button class="btn" type="button" aria-pressed={wrap} on:click={() => (wrap = !wrap)}>
          Wrap
        </button>
        <a class="btn" href={`/${paste.id}/raw`} data-sveltekit-reload>Raw</a>
        <button class="btn" type="button" on:click={copyLink}>
          {copied === 'link' ? 'Copied!' : 'Copy link'}
        </button>
        <button class="btn btn-accent" type="button" on:click={copyCode}>
          {copied === 'code' ? 'Copied!' : range ? 'Copy lines' : 'Copy code'}
        </button>
      </div>
    </div>

    <div
      class="code"
      class:wrap
      class:large
      class:dragging
      style:--gutter={gutterWidth}
      translate="no"
      on:pointerdown={onPointerDown}
      on:pointerover={onPointerOver}
      on:click={onClick}
      on:keydown={onKeyDown}
    >
      <div class="lines">
        {#each lines as line, i}
          <div
            class="row"
            class:selected={range != null && i + 1 >= range.start && i + 1 <= range.end}
            data-line={i + 1}
          >
            <!-- The number is drawn by CSS so copying code never includes it -->
            <a
              class="ln"
              id={`L${i + 1}`}
              href={`#L${i + 1}`}
              tabindex={i + 1 === focusedLine ? 0 : -1}
              aria-label={`Line ${i + 1}`}
              data-n={i + 1}
            />
            <span class="src">{@html line}</span>
          </div>
        {/each}
      </div>
    </div>
  </div>
  <p class="sr-only" aria-live="polite">{status}</p>
</article>

<style lang="scss">
  .head {
    margin-bottom: 1rem;
  }

  h1 {
    font-size: 1.625rem;
    font-weight: 700;
    letter-spacing: -0.015em;
    line-height: 1.3;
    overflow-wrap: anywhere;
    text-wrap: balance;
  }

  .meta {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    gap: 0.35rem 1rem;
    margin-top: 0.5rem;
    color: var(--text-muted);
    font-size: 14px;

    strong {
      color: var(--text);
      font-weight: 600;
    }
  }

  .viewer {
    /* clip, not hidden: hidden would break the sticky toolbar */
    overflow: clip;
  }

  .toolbar {
    position: sticky;
    top: 0;
    z-index: 2;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.5rem 0.5rem 0.5rem 1rem;
    border-bottom: 1px solid var(--border);
    background: var(--surface-2);
  }

  .count {
    color: var(--text-muted);
    font-size: 13px;
    font-variant-numeric: tabular-nums;
    white-space: nowrap;
  }

  .selection {
    color: var(--accent-text);
  }

  .actions {
    display: flex;
    gap: 0.375rem;
  }

  .btn-accent {
    min-width: 6.5rem;
    border-color: var(--accent);
    color: var(--accent-text);
  }

  .code {
    overflow-x: auto;
    padding: 0.75rem 0;
    font-family: var(--font-mono);
    font-size: 14px;
    line-height: 1.6;
    tab-size: 4;
    color: var(--syn-plain);
    font-variant-ligatures: none;

    &.dragging {
      user-select: none;
      cursor: row-resize;
    }
  }

  .lines {
    width: max-content;
    min-width: 100%;
  }

  .wrap .lines {
    width: auto;
  }

  .row {
    display: grid;
    grid-template-columns: var(--gutter) 1fr;
  }

  .large .row {
    content-visibility: auto;
    contain-intrinsic-size: auto 22.4px;
  }

  .ln {
    position: sticky;
    left: 0;
    padding-right: 1ch;
    border-right: 2px solid transparent;
    background: var(--surface);
    color: var(--gutter-text);
    text-align: right;
    text-decoration: none;
    user-select: none;
    cursor: pointer;
    scroll-margin-top: 4rem;

    &::before {
      content: attr(data-n);
    }

    &:hover {
      color: var(--text);
    }
  }

  .src {
    padding: 0 1.25rem 0 1ch;
    white-space: pre;
  }

  .wrap .src {
    white-space: pre-wrap;
    overflow-wrap: anywhere;
  }

  .selected {
    background: var(--highlight);

    .ln {
      border-right-color: var(--highlight-edge);
      background: linear-gradient(var(--highlight), var(--highlight)), var(--surface);
      color: var(--text);
    }
  }

  @media (max-width: 640px) {
    h1 {
      font-size: 1.3rem;
    }

    .toolbar {
      flex-wrap: wrap;
      gap: 0.5rem;
    }

    .actions {
      width: 100%;

      .btn {
        flex: 1;
        padding: 0 0.5rem;
      }
    }

    .code {
      font-size: 13px;
    }
  }
</style>
