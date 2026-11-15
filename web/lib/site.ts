import {
  BellRing,
  Database,
  KeyRound,
  LayoutDashboard,
  ShieldCheck,
  Users,
  type LucideIcon,
} from 'lucide-react';

export interface Feature {
  icon: LucideIcon;
  title: string;
  body: string;
}

export interface SiteConfig {
  name: string;
  binary: string;
  repo: string;
  tagline: string;
  description: string;
  badge: string;
  installCommand: string;
  features: Feature[];
  exampleTitle: string;
  example: string;
  compatible?: string[];
  featuresTitle?: string;
  featuresSubtitle?: string;
  ctaBody?: string;
  accent?: string;
  accentName?: string;
  accentHex?: string;
}

export const site: SiteConfig = {
  name: 'Grafana CLI',
  binary: 'grafana',
  repo: 'piyush-gambhir/grafana-cli',
  tagline: 'Grafana from your terminal',
  description:
    'An independent, unofficial open-source CLI for managing Grafana instances, dashboards, data sources, alerts, and access with scriptable, agent-safe workflows.',
  badge: 'Open-source · Cloud & self-hosted',
  accent: 'oklch(0.75 0.15 60)',
  accentName: 'amber',
  accentHex: '#f2943c',
  installCommand:
    'curl -sSfL https://raw.githubusercontent.com/piyush-gambhir/grafana-cli/main/install.sh | sh',
  features: [
    {
      icon: LayoutDashboard,
      title: 'Dashboards & folders',
      body: 'Search, inspect, import, export, version, restore, and organize dashboards, folders, permissions, and library elements.',
    },
    {
      icon: Database,
      title: 'Datasources & queries',
      body: 'Manage datasources and run LogQL or PromQL through Grafana\'s proxy with relative or absolute time ranges.',
    },
    {
      icon: BellRing,
      title: 'Unified alerting',
      body: 'Work with alert rules, contact points, notification policies, mute timings, templates, and silences.',
    },
    {
      icon: Users,
      title: 'Access & automation',
      body: 'Manage organizations, users, teams, service accounts, and service account tokens from one CLI.',
    },
    {
      icon: KeyRound,
      title: 'Profiles & authentication',
      body: 'Connect with a service account token or basic auth, save named profiles, and override credentials per command or environment.',
    },
    {
      icon: ShieldCheck,
      title: 'Agent-safe workflows',
      body: 'Use JSON or YAML output with --read-only, --no-input, and --quiet for predictable automation and safer exploration.',
    },
  ],
  exampleTitle: 'An eight-line tour',
  example: `# Save a named connection profile
grafana login
# Confirm the authenticated user
grafana user current -o json
# Find production dashboards
grafana dashboard list --query "production" -o json
# Query Prometheus through Grafana
grafana datasource query <prometheus-uid> --expr 'up' \\
  --query-type instant -o json`,
  compatible: [
    "HTTP API",
    "Dashboards",
    "Alerting",
    "Data sources",
    "Service accounts",
    "Folders",
  ],
};
