export const appName = 'Grafana CLI';
export const siteUrl = 'https://projects.piyushgambhir.com/grafana-cli';
export const docsRoute = '/docs';
export const docsImageRoute = '/og/docs';
export const docsContentRoute = '/llms.mdx/docs';

// The tracked file behind a docs page: guides live in web/content/docs.
export function sourcePath(pagePath: string): string {
  return `web/content/docs/${pagePath}`;
}

export const gitConfig = {
  user: 'piyush-gambhir',
  repo: 'grafana-cli',
  branch: 'main',
};
