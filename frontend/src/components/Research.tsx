import { usePortfolio } from "../api/PortfolioContext";
import Section from "./Section";
import PaperEntry from "./PaperEntry";

export default function Research() {
  const { site, papers } = usePortfolio();
  if (papers.length === 0 && !site.research_statement) return null;

  return (
    <Section id="research" title={site.research_title}>
      {site.research_statement && <p className="mb-5 max-w-[720px] text-[15px] leading-relaxed">{site.research_statement}</p>}
      <ul>
        {papers.map((p) => (
          <PaperEntry key={p.id} paper={p} />
        ))}
      </ul>
    </Section>
  );
}
