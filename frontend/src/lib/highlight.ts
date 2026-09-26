import { getParseRule, Language, tokenize } from '@calor/core';

import hljs from './hljs';

// Both highlighters slow down quadratically on long lines (a 50k-char line
// takes ~10s in highlight.js), so such pastes are shown as plain text.
const MAX_HIGHLIGHT_CHARS = 200_000;
const MAX_LINE_LENGTH = 1_000;
// calor is quadratic in total length, so it is only used for small pastes.
const MAX_CALOR_CHARS = 10_000;

const calorLanguages = new Set<string>(Object.values(Language));
const calorRemap: Record<string, string> = { go: 'golang' };

const escapeHTML = (text: string) =>
  text.replace(/[&<>]/g, (match) => (match === '&' ? '&amp;' : match === '<' ? '&lt;' : '&gt;'));

// Returns one well-formed HTML string per line of code.
export function highlightLines(code: string, language: string): string[] {
  const calorLanguage = calorRemap[language] ?? language;
  const tooLong =
    code.length > MAX_HIGHLIGHT_CHARS ||
    code.split(/\r\n?|\n/).some((line) => line.length > MAX_LINE_LENGTH);
  let html = '';
  if (!tooLong && code.length <= MAX_CALOR_CHARS && calorLanguages.has(calorLanguage)) {
    try {
      html = tokenize(code, getParseRule(calorLanguage as Language))
        .map((token) => `<span class="calor-${token.kind}">${escapeHTML(token.value)}</span>`)
        .join('');
    } catch {
      // calor 0.0.5 throws for some languages it lists (e.g. c); use highlight.js
    }
  }
  if (!html && !tooLong && language !== 'text' && hljs.getLanguage(language)) {
    html = hljs.highlight(code, { language, ignoreIllegals: true }).value;
  }
  if (!html) {
    html = escapeHTML(code);
  }

  // Spans that cross a newline (multi-line comments, strings) are closed at the
  // end of the line and reopened on the next one.
  const lines: string[] = [];
  const open: string[] = [];
  const pattern = /<span[^>]*>|<\/span>|\r\n?|\n/g;
  let line = '';
  let last = 0;
  let match: RegExpExecArray | null;
  while ((match = pattern.exec(html)) !== null) {
    line += html.slice(last, match.index);
    last = pattern.lastIndex;
    const token = match[0];
    if (token === '</span>') {
      open.pop();
      line += token;
    } else if (token.startsWith('<span')) {
      open.push(token);
      line += token;
    } else {
      lines.push(line + '</span>'.repeat(open.length));
      line = open.join('');
    }
  }
  lines.push(line + html.slice(last));

  // A trailing newline does not start a new visible line.
  if (lines.length > 1 && lines[lines.length - 1].replace(/<[^>]*>/g, '') === '') {
    lines.pop();
  }
  return lines;
}
