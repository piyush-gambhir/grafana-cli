import type { Metadata } from 'next';
import { siteUrl } from './shared';
import { site } from './site';

export const repoUrl = `https://github.com/${site.repo}`;
export const socialImage = `${siteUrl}/og/docs/image.png`;
export const siteMetadataDescription =
  'Independent, unofficial Grafana CLI for Claude Code, OpenAI Codex, Cursor, or shell agents: JSON/YAML, read-only, no-input dashboard/datasource/alert automation.';

interface PageMetadataOptions {
  title: string;
  description: string;
  path: string;
}

export function createPageMetadata({
  title,
  description,
  path,
}: PageMetadataOptions): Metadata {
  const canonicalUrl = `${siteUrl}${path}`;

  return {
    title,
    description,
    alternates: { canonical: canonicalUrl },
    openGraph: {
      type: 'website',
      url: canonicalUrl,
      siteName: site.name,
      title,
      description,
      images: [
        {
          url: socialImage,
          width: 1200,
          height: 630,
          alt: `${site.name} documentation`,
        },
      ],
    },
    twitter: {
      card: 'summary_large_image',
      title,
      description,
      images: [socialImage],
    },
  };
}

export function serializeJsonLd(data: unknown): string {
  return JSON.stringify(data).replace(/</g, '\\u003c');
}
