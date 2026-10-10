import { usePortfolio } from "../api/PortfolioContext";
import Section from "./Section";
import Profiles from "./Profiles";

/** Profile links, grouped by row. Email and location live in the footer. */
export default function Contact() {
  const { site, profiles } = usePortfolio();
  if (profiles.length === 0) return null;

  return (
    <Section id="contact" title={site.contact_title}>
      <Profiles />
    </Section>
  );
}
