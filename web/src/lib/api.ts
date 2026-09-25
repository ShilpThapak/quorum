export interface QuestionMeta {
  id: string;
  title: string;
  slug: string;
  difficulty: number;
  tags: string[];
  task: string;
  description?: string;
  timeoutSec?: number;
}

export interface Case {
  name: string;
  passed: boolean;
  elapsed?: number;
}

export interface Report {
  question: string;
  kind: string;
  vetOK: boolean;
  testOK: boolean;
  race: boolean;
  cases: Case[];
  output?: string;
  vetNotes?: string[];
  timeMs: number;
}


export async function listQuestions(): Promise<QuestionMeta[]> {
  const r = await fetch("/api/questions");
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}

export async function loadSource(slug: string, kind: "reference" | "starter"): Promise<string> {
  const r = await fetch(`/api/source/${slug}/${kind}`);
  if (!r.ok) throw new Error(await r.text());
  const d = (await r.json()) as { code: string };
  return d.code;
}

export async function grade(slug: string, code: string): Promise<Report> {
  const r = await fetch("/api/grade", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ question: slug, kind: "user", code }),
  });
  if (!r.ok) throw new Error(await r.text());
  return r.json();
}
