import hljs from './hljs';

type LanguageInfo = {
  id: string;
  label: string;
  // CodeMirror 5 MIME type used by the editor
  mode: string;
};

export const languages: LanguageInfo[] = [
  { id: 'text', label: 'Plain text', mode: 'text/plain' },
  { id: 'bash', label: 'Bash', mode: 'text/x-sh' },
  { id: 'c', label: 'C', mode: 'text/x-csrc' },
  { id: 'cpp', label: 'C++', mode: 'text/x-c++src' },
  { id: 'csharp', label: 'C#', mode: 'text/x-csharp' },
  { id: 'css', label: 'CSS', mode: 'text/css' },
  { id: 'dockerfile', label: 'Dockerfile', mode: 'text/x-dockerfile' },
  { id: 'go', label: 'Go', mode: 'text/x-go' },
  { id: 'html', label: 'HTML', mode: 'text/html' },
  { id: 'java', label: 'Java', mode: 'text/x-java' },
  { id: 'javascript', label: 'JavaScript', mode: 'text/javascript' },
  { id: 'json', label: 'JSON', mode: 'application/json' },
  { id: 'kotlin', label: 'Kotlin', mode: 'text/x-kotlin' },
  { id: 'less', label: 'Less', mode: 'text/x-less' },
  { id: 'markdown', label: 'Markdown', mode: 'text/x-markdown' },
  { id: 'php', label: 'PHP', mode: 'application/x-httpd-php' },
  { id: 'python', label: 'Python', mode: 'text/x-python' },
  { id: 'ruby', label: 'Ruby', mode: 'text/x-ruby' },
  { id: 'rust', label: 'Rust', mode: 'text/x-rustsrc' },
  { id: 'scala', label: 'Scala', mode: 'text/x-scala' },
  { id: 'scss', label: 'SCSS', mode: 'text/x-scss' },
  { id: 'sql', label: 'SQL', mode: 'text/x-sql' },
  { id: 'toml', label: 'TOML', mode: 'text/x-toml' },
  { id: 'typescript', label: 'TypeScript', mode: 'text/typescript' },
  { id: 'xml', label: 'XML', mode: 'application/xml' },
  { id: 'yaml', label: 'YAML', mode: 'text/x-yaml' }
];

const byId = new Map(languages.map((lang) => [lang.id, lang]));

export const languageLabel = (id: string) => byId.get(id)?.label ?? id;

export const languageMode = (id: string) => byId.get(id)?.mode ?? 'text/plain';

// Languages highlight.js can guess. Plain text is the fallback, not a candidate.
const detectable = languages
  .map((lang) => lang.id)
  .filter((id) => id !== 'text' && hljs.getLanguage(id));

// Only the start of the code is inspected so detection stays cheap on large pastes.
const DETECT_SAMPLE = 5000;
// Below this relevance highlight.js is guessing; plain text is the honest answer.
const MIN_RELEVANCE = 3;

export function detectLanguage(code: string): string {
  const sample = code.slice(0, DETECT_SAMPLE);
  if (sample.trim() === '') {
    return 'text';
  }
  if (/^\s*[[{]/.test(code)) {
    try {
      JSON.parse(code);
      return 'json';
    } catch {
      // Not JSON; let highlight.js guess
    }
  }
  const result = hljs.highlightAuto(sample, detectable);
  if (!result.language || result.relevance < MIN_RELEVANCE) {
    return 'text';
  }
  return result.language;
}
