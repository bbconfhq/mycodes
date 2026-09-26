<script lang="ts">
  import type { Editor, EditorConfiguration, EditorFromTextArea } from 'codemirror';
  import { createEventDispatcher, onMount } from 'svelte';

  import 'codemirror/lib/codemirror.css';
  import 'codemirror/addon/scroll/simplescrollbars.css';
  import 'codemirror-github-light/lib/codemirror-github-light-theme.css';
  import 'codemirror-github-dark/lib/codemirror-github-dark-theme.css';

  import { languageMode } from './language';

  export let code = '';
  // Language id from language.ts, already resolved (never 'auto')
  export let language = 'text';

  const dispatch = createEventDispatcher<{ submit: null }>();

  let editor: EditorFromTextArea | null = null;
  let textareaRef: HTMLTextAreaElement;

  export const focus = () => (editor ? editor.focus() : textareaRef.focus());

  const themeFor = (dark: boolean) => (dark ? 'github-dark' : 'github-light');

  onMount(() => {
    let destroyed = false;
    const darkQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const onThemeChange = (e: MediaQueryListEvent) =>
      editor?.setOption('theme', themeFor(e.matches));

    (async () => {
      const { default: CodeMirror } = await import('codemirror');
      await Promise.all([
        import('codemirror/addon/mode/simple'),
        import('codemirror/addon/scroll/simplescrollbars'),
        import('codemirror/addon/edit/closebrackets'),
        import('codemirror/addon/edit/closetag'),
        import('codemirror/addon/edit/continuelist'),
        import('codemirror/addon/edit/matchbrackets'),
        import('codemirror/addon/comment/comment'),
        import('codemirror/mode/clike/clike'),
        import('codemirror/mode/css/css'),
        import('codemirror/mode/dockerfile/dockerfile'),
        import('codemirror/mode/go/go'),
        import('codemirror/mode/htmlmixed/htmlmixed'),
        import('codemirror/mode/javascript/javascript'),
        import('codemirror/mode/markdown/markdown'),
        import('codemirror/mode/php/php'),
        import('codemirror/mode/python/python'),
        import('codemirror/mode/ruby/ruby'),
        import('codemirror/mode/rust/rust'),
        import('codemirror/mode/shell/shell'),
        import('codemirror/mode/sql/sql'),
        import('codemirror/mode/toml/toml'),
        import('codemirror/mode/xml/xml'),
        import('codemirror/mode/yaml/yaml')
      ]);
      if (destroyed) {
        return;
      }

      const submit = () => {
        dispatch('submit');
      };
      const options: EditorConfiguration & Record<string, unknown> = {
        theme: themeFor(darkQuery.matches),
        mode: languageMode(language),
        lineNumbers: true,
        indentWithTabs: true,
        indentUnit: 4,
        tabSize: 4,
        scrollbarStyle: 'overlay',
        autoCloseBrackets: true,
        autoCloseTags: true,
        matchBrackets: true,
        extraKeys: {
          Enter: 'newlineAndIndentContinueMarkdownList',
          'Ctrl-/': 'toggleComment',
          'Cmd-/': 'toggleComment',
          'Ctrl-Enter': submit,
          'Cmd-Enter': submit,
          // Tab is captured by the editor; Esc releases focus for keyboard users.
          Esc: (cm: Editor) => {
            cm.getInputField().blur();
          }
        }
      };
      editor = CodeMirror.fromTextArea(textareaRef, options);
      editor.setValue(code);
      editor.getInputField().setAttribute('aria-label', 'Code');
      editor.on('change', (instance) => {
        code = instance.getValue();
      });
      darkQuery.addEventListener('change', onThemeChange);
    })();

    return () => {
      destroyed = true;
      darkQuery.removeEventListener('change', onThemeChange);
      editor?.toTextArea();
      editor = null;
    };
  });

  $: editor?.setOption('mode', languageMode(language));
</script>

<!-- Stand-in with the editor's height until CodeMirror loads -->
<textarea bind:this={textareaRef} aria-label="Code" readonly value={code} />

<style>
  textarea {
    display: block;
    width: 100%;
    height: 100%;
    border: none;
    resize: none;
    background: transparent;
  }
</style>
