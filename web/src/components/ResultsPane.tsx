import type { Report, Case } from "../lib/api";

export function StatusStrip({ r }: { r: Report }) {
  const ok = r.testOK;
  const bits: string[] = [];
  bits.push(r.vetOK ? "vet ok" : "vet FAIL");
  bits.push(r.race ? "DATA RACE" : "no race");
  bits.push(`${r.timeMs}ms`);
  return (
    <div className={`strip ${ok ? "ok" : "bad"}`}>
      <b>{ok ? "TEST PASS" : "TEST FAIL"}</b>
      <span>{bits.join(" · ")}</span>
    </div>
  );
}

export function CaseList({ cases }: { cases: Case[] }) {
  const list = cases ?? [];
  return (
    <ul className="cases">
      {list.map((c) => (
        <li key={c.name} className={c.passed ? "ok" : "bad"}>
          <span className="mark">{c.passed ? "✓" : "✗"}</span>
          {c.name}
        </li>
      ))}
    </ul>
  );
}

export function ReportPanel({ r }: { r: Report }) {
  return (
    <section className="report">
      <StatusStrip r={r} />
      {r.vetNotes && r.vetNotes.length > 0 && (
        <pre className="vetnote">{r.vetNotes.join("\n")}</pre>
      )}
      <CaseList cases={r.cases} />
      {r.output && <pre className="out">{r.output}</pre>}
    </section>
  );
}
