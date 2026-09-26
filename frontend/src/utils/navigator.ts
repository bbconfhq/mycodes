export type LineRange = { start: number; end: number };

export const formatLineHash = ({ start, end }: LineRange) =>
  start === end ? `#L${start}` : `#L${start}-${end}`;

// Replaces the hash without scrolling or adding a history entry. SvelteKit
// keeps router data in history.state, so it has to be preserved.
export const setLineHash = (range: LineRange | null) => {
  const hash = range ? formatLineHash(range) : '';
  if (window.location.hash === hash) {
    return;
  }
  const url = `${window.location.pathname}${window.location.search}${hash}`;
  history.replaceState(history.state, '', url);
};

// Parses '#L12' or '#L12-34' (also '#L12-L34'), clamped to the available lines.
export const parseLineHash = (hash: string, lineCount: number): LineRange | null => {
  const match = /^#L(\d+)(?:-L?(\d+))?$/.exec(hash);
  if (match == null) {
    return null;
  }
  const a = Number(match[1]);
  const b = match[2] == null ? a : Number(match[2]);
  const start = Math.max(1, Math.min(a, b));
  const end = Math.min(lineCount, Math.max(a, b));
  if (start > lineCount || end < 1) {
    return null;
  }
  return { start, end };
};
