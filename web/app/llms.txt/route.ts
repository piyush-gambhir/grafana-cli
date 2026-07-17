import { source } from '@/lib/source';
import { llms } from 'fumadocs-core/source';
import { siteUrl } from '@/lib/shared';
import { site } from '@/lib/site';
import { getOtherSuiteProjects } from '@/lib/suite';

export const revalidate = false;

export function GET() {
  const intro =
    'Grafana CLI is an independent, unofficial command-line interface built for coding agents and any shell-capable agent harness, including Claude Code, OpenAI Codex, and Cursor. It provides structured JSON/YAML output, read-only safety, and no-input automation for dashboards, datasources, and alerting.';
  const docsIndex = llms(source)
    .index()
    .replaceAll('](/', `](${siteUrl}/`);
  const relatedProjects = getOtherSuiteProjects(site.repo)
    .map(({ name, website }) => `- [${name}](${website})`)
    .join('\n');

  return new Response(
    `${docsIndex}\n\n${intro}\n\n## Related independent CLI projects\n\n${relatedProjects}\n`,
  );
}
