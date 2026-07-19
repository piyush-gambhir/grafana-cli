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
  /** Display name, e.g. "Acme CLI" */
  name: string;
  /** The binary invoked in examples, e.g. "acme" */
  binary: string;
  /** GitHub "owner/repo" */
  repo: string;
  /** One-line hero heading */
  tagline: string;
  /** Hero sub-paragraph */
  description: string;
  /** Small pill above the heading */
  badge: string;
  /** One-line install command shown in the hero */
  installCommand: string;
  /** Feature cards */
  features: Feature[];
  /** Title above the code block */
  exampleTitle: string;
  /** Shell example rendered in the terminal card */
  example: string;
  /** Optional: tech / query languages this CLI speaks (logo strip) */
  compatible?: string[];
  /** Optional: features section heading (default: "Everything, from one binary") */
  featuresTitle?: string;
  /** Optional: features section subheading */
  featuresSubtitle?: string;
  /** Optional: CTA band body (default mentions installing the binary) */
  ctaBody?: string;
  /** Optional: per-site accent expressed as an OKLCH color */
  accent?: string;
  /** Optional: human-readable accent name */
  accentName?: string;
  /** Optional: hex twin of the accent, for surfaces without oklch() support (OG images) */
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
