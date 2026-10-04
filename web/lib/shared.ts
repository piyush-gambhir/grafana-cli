export const appName = 'Grafana CLI';
export const siteUrl = 'https://projects.piyushgambhir.com/grafana-cli';
export const docsRoute = '/docs';
export const docsImageRoute = '/og/docs';
export const docsContentRoute = '/llms.mdx/docs';

// Agents read the llms Markdown outside the site, where root-relative links
// (Markdown links, reference definitions, and <Card href>) resolve against the
// domain root and miss the basePath, so make them absolute. Code spans and
// fences (a run of N unescaped backticks up to the next run of exactly N) are
// matched first and left alone.
const codeOrRootLink =
  /(?<![\\`])(`+)(?!`)[\s\S]*?(?<!`)\1(?!`)|(?:\]\(<?|href="|^[ \t]*\[(?!\^)[^\]\n]+\]:[ \t]*<?)(?=\/(?!\/))/gm;

export function absoluteLinks(markdown: string): string {
  return markdown.replace(codeOrRootLink, (match, ticks?: string) => (ticks ? match : match + siteUrl));
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
