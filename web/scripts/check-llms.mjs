// Verify the built agent Markdown (llms.txt, llms-full.txt, per-page content.md)
// links only to absolute URLs: agents read these files outside the site, where
// root-relative links (/docs/...) resolve against the domain root and miss the
// basePath, and unevaluated JSX (href={...}) is not a link at all.
import assert from 'node:assert/strict';
import { readdir, readFile } from 'node:fs/promises';

const out = new URL('../out/', import.meta.url);
const pages = (await readdir(new URL('llms.mdx/', out), { recursive: true }))
  .filter((f) => f.endsWith('.md'))
  .map((f) => `llms.mdx/${f}`);
assert.ok(pages.length > 0, 'no per-page Markdown under out/llms.mdx');
// Code examples are not links; skip fenced blocks and inline code.
const code = /```[\s\S]*?```|`[^`\n]*`/g;
// Inline links, <a>/<Card> hrefs, and reference definitions ([x]: /path).
const link = /\]\(([^)\s]*)|href="([^"]*)"|^\s*\[(?!\^)[^\]\n]+\]:\s*(\S*)/gm;
for (const file of ['llms.txt', 'llms-full.txt', ...pages]) {
  const text = (await readFile(new URL(file, out), 'utf8')).replace(code, '');
  for (const [match, md, href, ref] of text.matchAll(link)) {
    const target = md ?? href ?? ref;
    assert.ok(/^(https?:\/\/|mailto:|#)/.test(target), `${file} has a non-absolute link: ${match.trim()}`);
  }
}
const full = await readFile(new URL('llms-full.txt', out), 'utf8');
assert.ok(full.includes('https://projects.piyushgambhir.com/grafana-cli/docs/'), 'llms-full.txt has no absolute docs links');
console.log(`Agent Markdown links are absolute in ${pages.length + 2} files`);
