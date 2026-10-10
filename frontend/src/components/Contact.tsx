import { usePortfolio } from "../api/PortfolioContext";
import Section, { sectionTitle } from "./Section";
import Profiles from "./Profiles";
import ContactForm from "./ContactForm";

const row = "flex flex-col gap-0.5 sm:flex-row sm:gap-3";
const label = "w-24 shrink-0 pt-px text-sm text-muted";

/**
 * Left: heading, email, location, profile links. Right: form panel whose top
 * is level with the heading and whose bottom meets the last row.
 */
export default function Contact() {
  const { site, profiles } = usePortfolio();
  const withForm = site.show_contact_form;

  return (
    <Section id="contact">
      <div className={`grid items-stretch gap-8 ${withForm ? "md:grid-cols-[3fr_2fr] md:gap-10" : ""}`}>
        <div>
          <h2 className={`${sectionTitle} mb-5`}>{site.contact_title}</h2>
          <dl className="space-y-2.5">
            {site.email && (
              <div className={row}>
                <dt className={label}>{site.contact_email_label}</dt>
                <dd className="text-sm">
                  <a href={`mailto:${site.email}`} className="font-medium text-accent hover:underline">
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

        {withForm && <ContactForm />}
      </div>
    </Section>
  );
}
