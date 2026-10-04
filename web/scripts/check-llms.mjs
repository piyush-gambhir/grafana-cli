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
const link = /\]\(([^)\s]*)|href="([^"]*)"/g;
for (const file of ['llms.txt', 'llms-full.txt', ...pages]) {
  const text = await readFile(new URL(file, out), 'utf8');
  for (const [match, md, href] of text.matchAll(link)) {
    const target = md ?? href;
    assert.ok(/^(https?:\/\/|mailto:|#)/.test(target), `${file} has a non-absolute link: ${match}`);
  }
}
const full = await readFile(new URL('llms-full.txt', out), 'utf8');
assert.ok(full.includes('https://projects.piyushgambhir.com/grafana-cli/docs/'), 'llms-full.txt has no absolute docs links');
console.log(`Agent Markdown links are absolute in ${pages.length + 2} files`);
