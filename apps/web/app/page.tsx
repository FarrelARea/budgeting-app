import Link from "next/link";
import { Button, Card, CardContent, CardHeader, CardTitle } from "@budgeting/ui";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

async function getHealth(): Promise<string> {
  try {
    const res = await fetch(`${API_URL}/healthz`, { cache: "no-store" });
    if (!res.ok) return "degraded";
    const body = (await res.json()) as { status: string };
    return body.status ?? "unknown";
  } catch {
    return "unreachable";
  }
}

export default async function HomePage() {
  const health = await getHealth();

  return (
    <main className="mx-auto max-w-3xl px-6 py-16">
      <Card>
        <CardHeader>
          <CardTitle>Budgeting App</CardTitle>
          <p className="text-sm text-zinc-500">Monorepo scaffold: Next.js web + Go API.</p>
        </CardHeader>
        <CardContent className="flex flex-col gap-4">
          <div className="flex items-center justify-between rounded-lg border border-zinc-200 px-4 py-3 text-sm dark:border-zinc-800">
            <span>API status</span>
            <span className="font-mono">{health}</span>
          </div>
          <div className="flex gap-2">
            <Button>Get started</Button>
            <Button variant="secondary">View accounts</Button>
            <Link href="/login">
              <Button variant="ghost">Login</Button>
            </Link>
          </div>
        </CardContent>
      </Card>
    </main>
  );
}
