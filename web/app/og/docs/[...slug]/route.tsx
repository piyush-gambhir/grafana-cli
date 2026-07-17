import { getPageImage, source } from '@/lib/source';
import { notFound } from 'next/navigation';
import { ImageResponse } from 'next/og';
import { site } from '@/lib/site';

import { readFile } from 'node:fs/promises';
import { join } from 'node:path';

const fontBuffer = async (name: string) => {
  const data = await readFile(join(process.cwd(), 'fonts', name));
  return data.buffer.slice(data.byteOffset, data.byteOffset + data.byteLength) as ArrayBuffer;
};

export const revalidate = false;

const hafferXH = fontBuffer('haffer-xh-regular-2.ttf');

export async function GET(_req: Request, { params }: RouteContext<'/og/docs/[...slug]'>) {
  const { slug } = await params;
  const page = source.getPage(slug.slice(0, -1));
  if (!page) notFound();

  return new ImageResponse(
    <div
      style={{
        display: 'flex',
        width: '100%',
        height: '100%',
        flexDirection: 'column',
        justifyContent: 'space-between',
        padding: '70px 76px',
        color: '#f3f4f1',
        background: '#131412',
        fontFamily: 'Haffer XH',
      }}
    >
      <div style={{ display: 'flex', color: '#f2943c', fontSize: 30 }}>
        &gt;_ {site.binary} docs
      </div>
      <div style={{ display: 'flex', flexDirection: 'column', gap: 24 }}>
        <div
          style={{
            display: 'flex',
            maxWidth: 1010,
            fontSize: 78,
            lineHeight: 0.96,
            letterSpacing: '-0.045em',
          }}
        >
          {page.data.title}
        </div>
        <div
          style={{
            display: 'flex',
            maxWidth: 880,
            color: '#b6b8b3',
            fontSize: 30,
            lineHeight: 1.25,
          }}
        >
          {page.data.description}
        </div>
      </div>
      <div
        style={{
          display: 'flex',
          justifyContent: 'space-between',
          color: '#7f827b',
          fontSize: 24,
        }}
      >
        <span>{site.name}</span>
        <span>projects.piyushgambhir.com/grafana-cli</span>
      </div>
    </div>,
    {
      width: 1200,
      height: 630,
      fonts: [
        {
          name: 'Haffer XH',
          data: await hafferXH,
          style: 'normal',
          weight: 400,
        },
      ],
    },
  );
}

export function generateStaticParams() {
  return source.getPages().map((page) => ({
    lang: page.locale,
    slug: getPageImage(page).segments,
  }));
}
