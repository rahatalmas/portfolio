import { useState } from "react";
import { usePortfolio } from "../api/PortfolioContext";

/** Round avatar with a hairline ring. Falls back to initials if the image is missing. */
export default function Portrait({ className = "" }: { className?: string }) {
  const { site } = usePortfolio();
  const [failed, setFailed] = useState(false);
  const showImage = site.photo && !failed;

  return (
    <div className={`overflow-hidden rounded-full border border-line bg-card ring-4 ring-bg ${className}`}>
      {showImage ? (
        <img
          src={site.photo}
          alt={`Portrait of ${site.name}`}
          onError={() => setFailed(true)}
          className="h-full w-full object-cover object-[50%_20%]"
        />
      ) : (
        <div className="flex h-full w-full items-center justify-center font-mono text-xl text-muted">
          {site.initials}
        </div>
      )}
    </div>
  );
}
