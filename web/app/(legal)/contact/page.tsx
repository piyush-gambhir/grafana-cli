import type { Metadata } from 'next';
import Link from 'next/link';
import { LegalPage } from '@/components/legal-page';
import { createPageMetadata } from '@/lib/metadata';

export const metadata: Metadata = createPageMetadata({
  title: 'Contact',
  description:
    'Contact and support options for Grafana CLI, an independent, unofficial open-source command-line tool for Grafana.',
  path: '/contact',
});

export default function ContactPage() {
  return (
    <LegalPage
      title="Contact"
      intro={
        <>
          Grafana CLI is a free, open-source project maintained by{' '}
          <strong>Piyush Gambhir</strong>. Support is best-effort. Here are the best
          ways to get in touch.
        </>
      }
    >
      <div className="legal-contact-grid">
        <article className="legal-contact-card">
          <h2>Email</h2>
          <p className="legal-contact-card__link">
            <a href="mailto:developer.piyushgambhir@gmail.com">
              developer.piyushgambhir@gmail.com
            </a>
          </p>
          <p>General questions, privacy, and security reports.</p>
        </article>
        <article className="legal-contact-card">
          <h2>Bugs &amp; features</h2>
          <p className="legal-contact-card__link">
            <a
              href="https://github.com/piyush-gambhir/grafana-cli/issues"
              target="_blank"
              rel="noreferrer"
            >
              GitHub Issues ↗
            </a>
          </p>
          <p>The fastest way to report a bug or request a feature.</p>
        </article>
        <article className="legal-contact-card">
          <h2>Source</h2>
          <p className="legal-contact-card__link">
            <a
              href="https://github.com/piyush-gambhir/grafana-cli"
              target="_blank"
              rel="noreferrer"
            >
              piyush-gambhir/grafana-cli ↗
            </a>
          </p>
          <p>Read the code, open a pull request, or fork it.</p>
        </article>
      </div>

      <section>
        <h2>Security issues</h2>
        <p>
          If you believe you&apos;ve found a security vulnerability, please email{' '}
          <a href="mailto:developer.piyushgambhir@gmail.com">
            developer.piyushgambhir@gmail.com
          </a>{' '}
          with the details rather than opening a public issue. Grafana CLI stores
          credentials only on your own device and operates no servers, but
          responsible disclosure is always appreciated.
        </p>
      </section>

      <section>
        <h2>Response time</h2>
        <p>
          This is an independent side project, not a commercial product. The
          maintainer aims to respond when possible, but no response time or level of
          support is guaranteed. See the <Link href="/terms">Terms of Service</Link>{' '}
          for the full no-warranty terms.
        </p>
      </section>

      <section>
        <h2>Not affiliated with Grafana</h2>
        <p>
          Grafana CLI is an independent, unofficial tool and is not affiliated with,
          endorsed by, or sponsored by Grafana or its vendor. For issues with Grafana
          itself, contact that vendor&apos;s own support channels.
        </p>
      </section>
    </LegalPage>
  );
}
