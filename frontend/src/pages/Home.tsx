import { usePortfolio } from "../api/PortfolioContext";
import Hero from "../components/Hero";
import About from "../components/About";
import Work from "../components/Work";
import Research from "../components/Research";

/** The single page. Sections are toggled from the admin's Settings → Sections. */
export default function Home() {
  const { site } = usePortfolio();
  return (
    <>
      <Hero />
      {site.show_about && <About />}
      {site.show_work && <Work />}
      {site.show_research && <Research />}
    </>
  );
}
