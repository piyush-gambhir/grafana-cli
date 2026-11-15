export const appName = 'Grafana CLI';
export const siteUrl = 'https://projects.piyushgambhir.com/grafana-cli';
export const docsRoute = '/docs';
export const docsImageRoute = '/og/docs';
export const docsContentRoute = '/llms.mdx/docs';

const codeOrRootLink =
  /(?<![\\`])(`+)(?!`)[\s\S]*?(?<!`)\1(?!`)|(?:\]\(<?|href="|^[ \t]*\[(?!\^)[^\]\n]+\]:[ \t]*<?)(?=\/(?!\/))/gm;

export function absoluteLinks(markdown: string): string {
  return markdown.replace(codeOrRootLink, (match, ticks?: string) => (ticks ? match : match + siteUrl));
}

export function sourcePath(pagePath: string): string {
  return `web/content/docs/${pagePath}`;
}

export const gitConfig = {
  user: 'piyush-gambhir',
  repo: 'grafana-cli',
  branch: 'main',
};
