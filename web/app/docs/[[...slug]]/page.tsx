import { getPageImage, getPageMarkdownUrl, source } from '@/lib/source';
import {
  DocsBody,
  DocsDescription,
  DocsPage,
  DocsTitle,
  MarkdownCopyButton,
  ViewOptionsPopover,
} from 'fumadocs-ui/layouts/docs/page';
import { notFound } from 'next/navigation';
import { getMDXComponents } from '@/components/mdx';
import type { Metadata } from 'next';
import { createRelativeLink } from 'fumadocs-ui/mdx';
import { serializeJsonLd, siteMetadataDescription } from '@/lib/seo';
import { gitConfig, siteUrl } from '@/lib/shared';
import { site } from '@/lib/site';

function getMetadataDescription(description?: string): string {
  void description;
  return siteMetadataDescription;
}

export default async function Page(props: PageProps<'/docs/[[...slug]]'>) {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();

  const MDX = page.data.body;
  const markdownUrl = `${siteUrl}${getPageMarkdownUrl(page).url}`;
  const canonicalUrl = `${siteUrl}${page.url}`;
  const breadcrumbItems = [
    {
      '@type': 'ListItem',
      position: 1,
      name: 'Home',
      item: siteUrl,
    },
    {
      '@type': 'ListItem',
      position: 2,
      name: 'Documentation',
      item: `${siteUrl}/docs`,
    },
    ...page.slugs.map((slug, index) => ({
      '@type': 'ListItem',
      position: index + 3,
      name:
        index === page.slugs.length - 1
          ? page.data.title
          : slug
              .split('-')
              .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
              .join(' '),
      item: `${siteUrl}/docs/${page.slugs.slice(0, index + 1).join('/')}`,
    })),
  ].filter(
    (item, index, items) => index === 0 || item.item !== items[index - 1]?.item,
  );
  const structuredData = {
    '@context': 'https://schema.org',
    '@graph': [
      {
        '@type': 'BreadcrumbList',
        itemListElement: breadcrumbItems,
      },
      {
        '@type': 'TechArticle',
        '@id': `${canonicalUrl}#tech-article`,
        headline: page.data.title,
        description: getMetadataDescription(page.data.description),
        url: canonicalUrl,
        mainEntityOfPage: canonicalUrl,
        author: {
          '@type': 'Person',
          name: 'Piyush Gambhir',
          url: 'https://github.com/piyush-gambhir',
        },
        publisher: {
          '@type': 'Person',
          name: 'Piyush Gambhir',
          url: 'https://github.com/piyush-gambhir',
        },
        isPartOf: {
          '@type': 'WebSite',
          '@id': `${siteUrl}#website`,
          name: site.name,
          url: siteUrl,
        },
      },
    ],
  };

  return (
    <DocsPage toc={page.data.toc} full={page.data.full}>
      <script
        type="application/ld+json"
        dangerouslySetInnerHTML={{ __html: serializeJsonLd(structuredData) }}
      />
      <DocsTitle>{page.data.title}</DocsTitle>
      <DocsDescription className="mb-0">{page.data.description}</DocsDescription>
      <div className="flex flex-row gap-2 items-center pb-6">
        <MarkdownCopyButton markdownUrl={markdownUrl} />
        <ViewOptionsPopover
          markdownUrl={markdownUrl}
          githubUrl={`https://github.com/${gitConfig.user}/${gitConfig.repo}/blob/${gitConfig.branch}/content/docs/${page.path}`}
        />
      </div>
      <DocsBody>
        <MDX
          components={getMDXComponents({
            // this allows you to link to other pages with relative file paths
            a: createRelativeLink(source, page),
          })}
        />
      </DocsBody>
    </DocsPage>
  );
}

export async function generateStaticParams() {
  return source.generateParams();
}

export async function generateMetadata(props: PageProps<'/docs/[[...slug]]'>): Promise<Metadata> {
  const params = await props.params;
  const page = source.getPage(params.slug);
  if (!page) notFound();
  const description = getMetadataDescription(page.data.description);
  const canonicalUrl = `${siteUrl}${page.url}`;
  const socialImageUrl = `${siteUrl}${getPageImage(page).url}`;

  return {
    title: page.data.title,
    description,
    alternates: { canonical: canonicalUrl },
    openGraph: {
      type: 'article',
      url: canonicalUrl,
      siteName: site.name,
      title: page.data.title,
      description,
      images: [
        {
          url: socialImageUrl,
          width: 1200,
          height: 630,
          alt: `${page.data.title} | ${site.name}`,
        },
      ],
    },
    twitter: {
      card: 'summary_large_image',
      title: page.data.title,
      description,
      images: [socialImageUrl],
    },
  };
}
