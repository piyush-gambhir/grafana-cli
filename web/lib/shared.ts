export const appName = 'Grafana CLI';
export const siteUrl = 'https://projects.piyushgambhir.com/grafana-cli';
export const docsRoute = '/docs';
export const docsImageRoute = '/og/docs';
export const docsContentRoute = '/llms.mdx/docs';

// Agents read the llms Markdown outside the site, where root-relative links
// (Markdown links, reference definitions, and <Card href>) resolve against the
// domain root and miss the basePath, so make them absolute. Code is left alone.
export function absoluteLinks(markdown: string): string {
  return markdown
    .split(/(```[\s\S]*?```|`[^`\n]*`)/)
    .map((part, i) =>
      i % 2 === 1
        ? part
        : part
            .replace(/\]\((?=\/(?!\/))/g, `](${siteUrl}`)
            .replace(/^(\s*\[(?!\^)[^\]\n]+\]:\s*)(?=\/(?!\/))/gm, `$1${siteUrl}`)
            .replace(/href="(?=\/(?!\/))/g, `href="${siteUrl}`),
    )
    .join('');
}

// The tracked file behind a docs page: guides live in web/content/docs.
export function sourcePath(pagePath: string): string {
  return `web/content/docs/${pagePath}`;
}

export const gitConfig = {
  user: 'piyush-gambhir',
  repo: 'grafana-cli',
  branch: 'main',
};
