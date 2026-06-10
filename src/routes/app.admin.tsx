import { createFileRoute } from "@tanstack/react-router";
import { Users, ClipboardList, Wallet, TrendingUp, AlertTriangle, MoreHorizontal } from "lucide-react";
import { tasks } from "@/lib/mock-data";
import { StatusBadge } from "@/components/status-badge";

export const Route = createFileRoute("/app/admin")({
  head: () => ({ meta: [{ title: "Admin — TaskBridge" }] }),
  component: Admin,
});

function Admin() {
  const kpis = [
    { label: "Total Users", value: "18,402", trend: "+5.2%", icon: Users },
    { label: "Active Tasks", value: "847", trend: "+12%", icon: ClipboardList },
    { label: "Monthly Revenue", value: "$92,140", trend: "+8.4%", icon: Wallet },
    { label: "Disputes Open", value: "6", trend: "-2", icon: AlertTriangle, danger: true },
  ];
  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Admin overview</h1>
          <p className="text-sm text-muted-foreground">Platform health at a glance.</p>
        </div>
        <div className="flex gap-2">
          <button className="rounded-lg border bg-card px-3.5 py-2 text-sm font-medium hover:bg-accent">Export</button>
          <button className="rounded-lg gradient-brand px-3.5 py-2 text-sm font-semibold text-white hover:opacity-95">Generate report</button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {kpis.map(k => (
          <div key={k.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{k.label}</span>
              <span className={`grid h-9 w-9 place-items-center rounded-lg ${k.danger ? "bg-destructive/10 text-destructive" : "bg-primary/10 text-primary"}`}><k.icon className="h-4 w-4" /></span>
            </div>
            <div className="mt-3 font-display text-3xl font-bold">{k.value}</div>
            <div className={`mt-1 inline-flex items-center gap-1 text-xs ${k.danger ? "text-destructive" : "text-success"}`}>
              <TrendingUp className="h-3 w-3" /> {k.trend}
            </div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="rounded-2xl border bg-card p-5 shadow-soft lg:col-span-2">
          <h2 className="font-display font-semibold">Revenue (last 30 days)</h2>
          <BigChart />
        </div>
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <h2 className="font-display font-semibold">Top cities</h2>
          <div className="mt-4 space-y-3">
            {[
              { city: "Tehran", pct: 92, value: "$24,310" },
              { city: "Istanbul", pct: 74, value: "$18,820" },
              { city: "Dubai", pct: 58, value: "$14,210" },
              { city: "Madrid", pct: 42, value: "$10,940" },
              { city: "Berlin", pct: 30, value: "$8,160" },
            ].map(c => (
              <div key={c.city}>
                <div className="flex items-center justify-between text-sm">
                  <span className="font-medium">{c.city}</span>
                  <span className="text-muted-foreground">{c.value}</span>
                </div>
                <div className="mt-1.5 h-2 overflow-hidden rounded-full bg-muted">
                  <div className="h-full gradient-brand" style={{ width: `${c.pct}%` }} />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      <div className="rounded-2xl border bg-card shadow-soft">
        <div className="flex items-center justify-between border-b px-5 py-4">
          <h2 className="font-display text-lg font-semibold">Recent tasks</h2>
          <button className="rounded-lg border px-3 py-1.5 text-xs font-medium hover:bg-accent">View all</button>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-muted/40 text-left text-xs uppercase tracking-wide text-muted-foreground">
              <tr>
                <th className="px-5 py-3 font-medium">ID</th>
                <th className="px-5 py-3 font-medium">Title</th>
                <th className="px-5 py-3 font-medium">City</th>
                <th className="px-5 py-3 font-medium">Budget</th>
                <th className="px-5 py-3 font-medium">Status</th>
                <th className="px-5 py-3 font-medium" />
              </tr>
            </thead>
            <tbody className="divide-y">
              {tasks.map(t => (
                <tr key={t.id} className="hover:bg-accent/40">
                  <td className="px-5 py-3 font-mono text-xs text-muted-foreground">{t.id}</td>
                  <td className="px-5 py-3 font-medium">{t.title}</td>
                  <td className="px-5 py-3 text-muted-foreground">{t.city}</td>
                  <td className="px-5 py-3 font-semibold tabular-nums">${t.budget}</td>
                  <td className="px-5 py-3"><StatusBadge status={t.status} /></td>
                  <td className="px-5 py-3 text-right"><button className="grid h-8 w-8 place-items-center rounded-lg hover:bg-accent"><MoreHorizontal className="h-4 w-4" /></button></td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function BigChart() {
  const data = [30, 45, 38, 55, 60, 48, 70, 65, 78, 82, 75, 90, 85, 92, 88, 95, 100, 92, 105, 110, 102, 118, 125, 115, 130, 128, 140, 135, 148, 155];
  const max = Math.max(...data);
  const points = data.map((v, i) => `${(i / (data.length - 1)) * 100},${100 - (v / max) * 100}`).join(" ");
  return (
    <div className="mt-4 h-56 w-full">
      <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="h-full w-full">
        <defs>
          <linearGradient id="g" x1="0" x2="0" y1="0" y2="1">
            <stop offset="0%" stopColor="var(--primary)" stopOpacity="0.35" />
            <stop offset="100%" stopColor="var(--primary)" stopOpacity="0" />
          </linearGradient>
        </defs>
        <polyline points={`0,100 ${points} 100,100`} fill="url(#g)" />
        <polyline points={points} fill="none" stroke="var(--primary)" strokeWidth="0.8" vectorEffect="non-scaling-stroke" />
      </svg>
    </div>
  );
}
