<script lang="ts">
  import { onDestroy, onMount } from 'svelte';

  import { debounce } from '../utils/debounce';

  import CodeMirror from './CodeMirror.svelte';
  import { detectLanguage, languageLabel, languages } from './language';

  import { goto } from '$app/navigation';

  const AUTHOR_KEY = 'mycodes:author';
  const EMPTY_ERROR = 'Paste or type some code first.';

  let code = '';
  let title = '';
  let name = '';
  let language = 'auto';
  let detected = 'text';
  let pending = false;
  let error = '';
  let editor: CodeMirror;
  let modKey = 'Ctrl';

  $: effectiveLanguage = language === 'auto' ? detected : language;

  const scheduleDetect = debounce((value: string) => {
    detected = detectLanguage(value);
  }, 400);
  $: if (language === 'auto') scheduleDetect(code);
  $: if (error === EMPTY_ERROR && code.trim()) error = '';

  onMount(() => {
    if (/Mac|iPhone|iPad/.test(navigator.userAgent)) {
      modKey = '⌘';
    }
    try {
      name = localStorage.getItem(AUTHOR_KEY) ?? '';
    } catch {
      // Storage can be unavailable (private mode); the field just starts empty.
    }
  });
  onDestroy(() => scheduleDetect.cancel());

  const messageFor = (status: number) => {
    if (status === 413) return 'This paste is too large to share.';
    if (status === 400) return 'Check the title (max 191 characters) and code, then try again.';
    return 'Could not share your code. Check your connection and try again.';
  };

  async function submit() {
    if (pending) {
      return;
    }
    if (code.trim() === '') {
      error = EMPTY_ERROR;
      editor.focus();
      return;
    }

    pending = true;
    error = '';
    try {
      const resp = await fetch('/api/v1/code/', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: title.trim() || 'Untitled',
          name: name.trim() || 'anonymous',
          content: code,
          // Detect again so a pending debounce cannot submit a stale guess
          language: language === 'auto' ? detectLanguage(code) : language
        })
      });
      const body = await resp.json().catch(() => null);
      const id = body?.Data?.id;
      if (!resp.ok || typeof id !== 'string') {
        error = messageFor(resp.status);
        return;
      }
      try {
        localStorage.setItem(AUTHOR_KEY, name.trim());
      } catch {
        // Remembering the author is a convenience only.
      }
      await goto(`/${id}`);
    } catch {
      error = messageFor(0);
    } finally {
      pending = false;
    }
  }
</script>

<form class="panel composer" on:submit|preventDefault={submit} aria-busy={pending}>
  <div class="head">
    <label class="sr-only" for="title">Title</label>
    <input
      class="field title"
      id="title"
      name="title"
      type="text"
      placeholder="Untitled…"
      maxlength="191"
      autocomplete="off"
      bind:value={title}
    />
    <label class="sr-only" for="name">Author</label>
    <input
      class="field name"
      id="name"
      name="name"
      type="text"
      placeholder="Author (optional)…"
      maxlength="191"
      autocomplete="off"
      spellcheck="false"
      bind:value={name}
    />
    <label class="sr-only" for="language">Language</label>
    <select class="field language" id="language" name="language" bind:value={language}>
      <option value="auto">Auto · {languageLabel(detected)}</option>
      {#each languages as lang}
        <option value={lang.id}>{lang.label}</option>
      {/each}
    </select>
  </div>

  <div class="editor">
    <CodeMirror bind:this={editor} bind:code language={effectiveLanguage} on:submit={submit} />
  </div>

  <div class="foot">
    <p class="hint">
      Expires in 7 days · <kbd>{modKey}</kbd>&nbsp;<kbd>Enter</kbd> to share
    </p>
    <p class="error" role="alert" aria-live="assertive">{error}</p>
    <button class="btn btn-primary" type="submit" disabled={pending}>
      {pending ? 'Sharing…' : 'Share code'}
    </button>
  </div>
</form>

<style lang="scss">
  .composer {
    overflow: hidden;
  }

  .head {
    display: flex;
    gap: 0.5rem;
    padding: 0.625rem;
    border-bottom: 1px solid var(--border);
    background: var(--surface-2);
  }

  .title {
    flex: 1 1 auto;
    font-weight: 600;
  }

  .name {
    flex: 0 1 13rem;
  }

  .language {
    flex: 0 0 11rem;
  }

  .editor {
    height: clamp(18rem, 55vh, 40rem);
  }

  // Beats the fixed 300px height from codemirror.css
  .editor :global(.CodeMirror) {
    height: 100%;
  }

  .foot {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.625rem 0.625rem 0.625rem 1rem;
    border-top: 1px solid var(--border);
    background: var(--surface-2);
  }

  .hint {
    color: var(--text-muted);
    font-size: 13px;
  }

  .error {
    flex: 1;
    color: var(--danger);
    font-size: 14px;
    text-align: right;
  }

  @media (max-width: 640px) {
    .head {
      flex-wrap: wrap;
    }

    .title {
      flex-basis: 100%;
    }

    .name,
    .language {
      flex: 1 1 0;
    }

    .foot {
      flex-wrap: wrap;
    }

    .hint {
      display: none;
    }

    .error {
      text-align: left;
    }
  }
</style>
