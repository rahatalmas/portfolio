import { usePortfolio } from "../api/PortfolioContext";
import Section, { sectionTitle } from "./Section";
import Profiles from "./Profiles";

const row = "flex flex-col gap-0.5 sm:flex-row sm:gap-3";
const label = "w-24 shrink-0 pt-px text-sm text-muted";

/**
 * The inverted band at the foot of the page: heading, email, location and
 * profile links. The footer continues the same band.
 */
export default function Contact() {
  const { site, profiles } = usePortfolio();

  return (
    <Section id="contact" tone="band">
      <div className="max-w-[720px]">
        <div>
          <h2 className={`${sectionTitle} mb-5 md:text-4xl`}>{site.contact_title}</h2>
          <dl className="space-y-2.5">
            {site.email && (
              <div className={row}>
                <dt className={label}>{site.contact_email_label}</dt>
                <dd className="text-sm">
                  <a href={`mailto:${site.email}`} className="font-medium text-fg underline decoration-line underline-offset-4 hover:decoration-fg">
                    {site.email}
                  </a>
                </dd>
              </div>
            )}
            {site.location && (
              <div className={row}>
                <dt className={label}>{site.contact_location_label}</dt>
                <dd className="text-sm">{site.location}</dd>
              </div>
            )}
          </dl>

          {profiles.length > 0 && (
            <>
              <div className="my-2.5 border-t border-line" />
              <Profiles />
            </>
          )}
        </div>
      </div>
    </Section>
  );
}
