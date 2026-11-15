import { createMDX } from 'fumadocs-mdx/next';

const withMDX = createMDX();

const basePath = '/grafana-cli';

const config = {
  experimental: { agentUpgrade: 'security' },
  basePath,
  env: { NEXT_PUBLIC_BASE_PATH: basePath },
  output: 'export',
  reactStrictMode: true,
};

export default withMDX(config);
