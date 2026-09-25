import { useEffect, useMemo, useState } from "react";
import Editor from "@monaco-editor/react";
import { EditorPane } from "./components/EditorPane";
import {
  grade,
  listQuestions,
  loadSource,
  type QuestionMeta,
  type Report,
  type Case,
} from "./lib/api";

const KINDS = ["reference", "starter", "user"] as const;

// Drafts are autosaved per question so a refresh never wipes your work.
const DRAFT_KEY = (slug: string) => "quorum:draft:" + slug;
const loadedPlaceholder = "// loading…";

function fmtDuration(ms: number): string {
  if (ms >= 1000) return (ms / 1000).toFixed(1) + "s";
  return ms + "ms";
}

export default function App() {
  const [qs, setQs] = useState<QuestionMeta[]>([]);
  const [slug, setSlug] = useState<string>("");
  const [code, setCode] = useState<string>(loadedPlaceholder);
  const [report, setReport] = useState<Report | null>(null);
  const [busy, setBusy] = useState(false);
  const [err, setErr] = useState<string>("");

  useEffect(() => {
    listQuestions()
      .then((l) => {
        setQs(l);
        if (l.length > 0 && !slug) setSlug(l[0].slug);
      })
      .catch((e) => setErr(String(e)));
  }, []);

  const meta = useMemo(
    () => qs.find((q) => q.slug === slug) ?? null,
    [qs, slug]
  );

  useEffect(() => {
    if (!slug) return;
    let alive = true;
    setCode("// loading…");
    loadSource(slug, "starter")
      .then((c) => {
        if (!alive) return;
        let saved: string | null = null;
        try {
          saved = localStorage.getItem(DRAFT_KEY(slug));
        } catch {}
        setCode(saved !== null ? saved : c);
      })
      .catch((e) => alive && setErr(String(e)));
    return () => {
      alive = false;
    };
  }, [slug]);

  useEffect(() => {
    if (!slug || !code || code === loadedPlaceholder) return;
    try {
      localStorage.setItem(DRAFT_KEY(slug), code);
    } catch {}
  }, [code, slug]);

  async function load(kind: "reference" | "starter") {
    if (!slug) return;
    setBusy(true);
    try {
      if (kind === "starter") {
        try {
          localStorage.removeItem(DRAFT_KEY(slug));
        } catch {}
      }
      setCode(await loadSource(slug, kind));
    } catch (e) {
      setErr(String(e));
    } finally {
      setBusy(false);
    }
  }

  async function run() {
    if (!slug) return;
    setBusy(true);
    setErr("");
    try {
      setReport(await grade(slug, code));
    } catch (e) {
      setErr(String(e));
    } finally {
      setBusy(false);
    }
  }

  return (
    <div className="app">
      <aside className="panel side">
        <h1>Quorum</h1>
        <nav>
          {qs.map((q) => (
            <button
              key={q.slug}
              className={q.slug === slug ? "q active" : "q"}
              onClick={() => setSlug(q.slug)}
            >
              <span className="dots">
                {"●".repeat(q.difficulty ?? 1)}
              </span>
              <span className="qtitle">{q.title ?? q.slug}</span>
            </button>
          ))}
        </nav>
        <div className="howto">
          {KINDS.join(" / ")} are graded with{" "}
          <code>go vet</code> + <code>go test -race</code>.
        </div>
      </aside>

      <main className="panel">
        <header className="panehead">
          <div>
            <h2>{meta?.title ?? slug}</h2>
            <p className="task">{meta?.task ?? ""}</p>
          </div>
          <div className="toolbar">
            <button onClick={() => load("reference")} disabled={busy}>
              Load reference
            </button>
            <button onClick={() => load("starter")} disabled={busy}>
              Reset starter
            </button>
            <button className="run" onClick={run} disabled={busy}>
              {busy ? "Grading…" : "Run tests"}
            </button>
          </div>
        </header>
        <EditorPane code={code} onChange={(v) => setCode(v ?? "")} />
        {err && <div className="banner err">{err}</div>}
      </main>

      <ResultsPane report={report} busy={busy} />
    </div>
  );
}

function ResultsPane({ report, busy }: { report: Report | null; busy: boolean }) {
  return (
    <aside className="panel results">
      <h2>Results</h2>
      {busy && !report && <p className="hint">assembling scratch module…</p>}
      {!report && !busy && (
        <p className="hint">
          Write your solution in the editor, then hit <b>Run tests</b>.
        </p>
      )}
      {report && (
        <>
          <StatusLine r={report} />
          {report.vetNotes && report.vetNotes.length > 0 && (
            <pre className="vet">{report.vetNotes.join("\n")}</pre>
          )}
          <ul className="cases">
            {(report.cases ?? []).map((c: Case) => (
              <li key={c.name} className={c.passed ? "ok" : "bad"}>
                <span className="mark">{c.passed ? "✓" : "✗"}</span>
                <span className="cname">{c.name}</span>
                {c.elapsed !== undefined && (
                  <span className="ctime">{fmtDuration(c.elapsed)}</span>
                )}
              </li>
            ))}
          </ul>
          <pre className="output">{report.output}</pre>
        </>
      )}
    </aside>
  );
}

function StatusLine({ r }: { r: Report }) {
  const cls = r.testOK ? "ok" : "bad";
  const verdict = r.testOK ? "ALL TESTS PASS" : "FAILED";
  const notes = [
    r.vetOK ? "vet clean" : "vet failed",
    r.race ? "data race" : "no race",
    fmtDuration(r.timeMs),
  ].join(" · ");
  return (
    <div className={`status ${cls}`}>
      <b>{verdict}</b>
      <span className="sub">{notes}</span>
    </div>
  );
}
