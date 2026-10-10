import { useState, type FormEvent } from "react";
import { ArrowRight, Check } from "lucide-react";
import { usePortfolio } from "../api/PortfolioContext";
import { ApiError, sendMessage } from "../api/client";

type Status = "idle" | "sending" | "sent" | "error";

const field =
  "w-full rounded-lg border border-line bg-bg px-3 py-2 text-sm outline-none placeholder:text-muted focus:border-accent focus:ring-2 focus:ring-accent/15";

/**
 * Posts to the backend (POST /api/v1/messages); submissions show up in the
 * admin inbox. All labels come from settings.
 */
export default function ContactForm() {
  const { site } = usePortfolio();
  const [status, setStatus] = useState<Status>("idle");
  const [error, setError] = useState("");

  async function onSubmit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    const form = e.currentTarget;
    const data = Object.fromEntries(new FormData(form)) as Record<string, string>;

    setStatus("sending");
    try {
      await sendMessage({ name: data.name, email: data.email, message: data.message, website: data.website });
      form.reset();
      setStatus("sent");
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "");
      setStatus("error");
    }
  }

  if (status === "sent") {
    return (
      <div className="panel flex h-full flex-col items-center justify-center gap-2 rounded-2xl border border-line bg-card p-6 text-center">
        <span className="flex h-9 w-9 items-center justify-center rounded-full bg-accent/10 text-accent">
          <Check size={18} />
        </span>
        <p className="font-medium">{site.form_success_title}</p>
        <p className="text-sm text-muted">{site.form_success_text}</p>
      </div>
    );
  }

  return (
    // Light panel on the band (.panel restores the normal palette): top level
    // with the section heading; the message box grows so the panel ends level
    // with the last profile row.
    <form
      onSubmit={onSubmit}
      className="panel flex h-full flex-col gap-2.5 rounded-2xl border border-line bg-card p-4 shadow-[0_24px_60px_-20px_rgba(0,0,0,0.5)]"
    >
      <div className="mb-1">
        <h3 className="text-[15px] font-semibold">{site.form_title}</h3>
        {site.form_subtitle && <p className="text-xs text-muted">{site.form_subtitle}</p>}
      </div>
      <div className="grid gap-2.5 sm:grid-cols-2">
        <label htmlFor="cf-name" className="sr-only">{site.form_name_placeholder}</label>
        <input id="cf-name" name="name" required maxLength={120} placeholder={site.form_name_placeholder} className={field} />
        <label htmlFor="cf-email" className="sr-only">{site.form_email_placeholder}</label>
        <input id="cf-email" name="email" type="email" required maxLength={254} placeholder={site.form_email_placeholder} className={field} />
      </div>
      <label htmlFor="cf-message" className="sr-only">{site.form_message_placeholder}</label>
      <textarea
        id="cf-message"
        name="message"
        required
        maxLength={5000}
        placeholder={site.form_message_placeholder}
        className={`${field} min-h-[72px] flex-1 resize-none`}
      />
      {/* Honeypot: hidden from people, filled by bots, dropped by the server. */}
      <input type="text" name="website" tabIndex={-1} autoComplete="off" aria-hidden className="hidden" />

      <button
        type="submit"
        disabled={status === "sending"}
        className="inline-flex w-full items-center justify-center gap-2 rounded-lg bg-accent px-4 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-50"
      >
        {status === "sending" ? site.form_sending_label : site.form_submit_label}
        {status !== "sending" && <ArrowRight size={15} />}
      </button>

      {status === "error" && (
        <p className="mt-1 text-center text-xs text-red-500" role="alert">
          {error || site.form_error_text}
          {site.email && (
            <>
              {" "}
              <a href={`mailto:${site.email}`} className="underline">
                {site.email}
              </a>
            </>
          )}
        </p>
      )}
    </form>
  );
}
