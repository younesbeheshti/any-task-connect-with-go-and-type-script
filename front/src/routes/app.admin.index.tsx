import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Users, ClipboardList, Wallet, TrendingUp, AlertTriangle, MoreHorizontal } from "lucide-react";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/app/admin/")({
  head: () => ({ meta: [{ title: "مدیریت — تسک‌بریج" }] }),
  component: Admin,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Metrics = {
  totalUsers: number; totalTasks: number; totalRevenue: number; openDisputes: number;
  activeAgents?: number; requesters?: number;
};
type Task = {
  id: string; title: string; status: string; budget: number;
  category?: { title: string }; city?: { title: string };
};

const STATUS_LABELS: Record<string, string> = {
  OPEN: "باز", ASSIGNED: "تخصیص‌یافته", IN_PROGRESS: "در حال انجام",
  COMPLETED: "تکمیل‌شده", CANCELLED: "لغو‌شده", PAID: "پرداخت‌شده",
  CREATED: "ثبت‌شده", VERIFIED: "تایید‌شده", WAITING_FOR_VERIFICATION: "در انتظار تایید",
};

function Admin() {
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/admin/metrics`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/tasks?page=1&pageSize=8`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([m, t]) => {
      if (m && !m.error) setMetrics(m);
      if (t) setTasks(t.tasks ?? []);
    }).finally(() => setLoading(false));
  }, []);

  const kpis = [
    { label: "تعداد کاربران",        value: loading ? "..." : toFa(metrics?.totalUsers ?? 0),     icon: Users,         danger: false },
    { label: "درخواست‌های فعال",     value: loading ? "..." : toFa(metrics?.totalTasks ?? 0),     icon: ClipboardList, danger: false },
    { label: "درآمد کل",             value: loading ? "..." : toman(metrics?.totalRevenue ?? 0, false), unit: "تومان", icon: Wallet, danger: false },
    { label: "اعتراضات باز",         value: loading ? "..." : toFa(metrics?.openDisputes ?? 0),   icon: AlertTriangle, danger: true },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-2">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">پنل مدیریت</h1>
          <p className="text-sm text-muted-foreground">وضعیت سلامت پلتفرم در یک نگاه.</p>
        </div>
        <div className="flex gap-2">
          <Link to="/app/admin/reports" className="rounded-lg border bg-card px-3.5 py-2 text-sm font-medium hover:bg-accent">مشاهده گزارش‌ها</Link>
          <Link to="/app/admin/users" className="rounded-lg gradient-brand px-3.5 py-2 text-sm font-semibold text-white hover:opacity-95">مدیریت کاربران</Link>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {kpis.map(k => (
          <div key={k.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{k.label}</span>
              <span className={`grid h-9 w-9 place-items-center rounded-lg ${k.danger ? "bg-destructive/10 text-destructive" : "bg-primary/10 text-primary"}`}>
                <k.icon className="h-4 w-4" />
              </span>
            </div>
            <div className="mt-3 font-display text-2xl font-bold tabular-nums">
              {k.value}{(k as any).unit && <span className="ms-1 text-sm text-muted-foreground font-medium">{(k as any).unit}</span>}
            </div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="rounded-2xl border bg-card p-5 shadow-soft lg:col-span-2">
          <h2 className="font-display font-semibold">روند درآمد (تخمینی)</h2>
          <BigChart />
        </div>
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <h2 className="font-display font-semibold">دسترسی سریع</h2>
          <div className="mt-4 space-y-2">
            <QA to="/app/admin/users"   label="مدیریت کاربران" />
            <QA to="/app/admin/finance" label="مالی و تراکنش‌ها" />
            <QA to="/app/admin/reports" label="گزارش‌ها و آمار" />
            <QA to="/app/tasks"         label="همه درخواست‌ها" />
          </div>
        </div>
      </div>

      <div className="rounded-2xl border bg-card shadow-soft">
        <div className="flex items-center justify-between border-b px-5 py-4">
          <h2 className="font-display text-lg font-semibold">آخرین درخواست‌ها</h2>
          <Link to="/app/tasks" className="rounded-lg border px-3 py-1.5 text-xs font-medium hover:bg-accent">مشاهده همه</Link>
        </div>
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-muted/40 text-right text-xs uppercase tracking-wide text-muted-foreground">
              <tr>
                <th className="px-5 py-3 font-medium">شناسه</th>
                <th className="px-5 py-3 font-medium">عنوان</th>
                <th className="px-5 py-3 font-medium">دسته‌بندی</th>
                <th className="px-5 py-3 font-medium">مبلغ</th>
                <th className="px-5 py-3 font-medium">وضعیت</th>
                <th className="px-5 py-3 font-medium" />
              </tr>
            </thead>
            <tbody className="divide-y">
              {loading ? (
                Array.from({ length: 5 }).map((_, i) => (
                  <tr key={i} className="animate-pulse">
                    <td className="px-5 py-3"><div className="h-4 w-24 rounded bg-muted" /></td>
                    <td className="px-5 py-3"><div className="h-4 w-40 rounded bg-muted" /></td>
                    <td className="px-5 py-3"><div className="h-4 w-20 rounded bg-muted" /></td>
                    <td className="px-5 py-3"><div className="h-4 w-24 rounded bg-muted" /></td>
                    <td className="px-5 py-3"><div className="h-5 w-20 rounded-full bg-muted" /></td>
                    <td className="px-5 py-3" />
                  </tr>
                ))
              ) : tasks.length === 0 ? (
                <tr><td colSpan={6} className="px-5 py-10 text-center text-sm text-muted-foreground">درخواستی یافت نشد</td></tr>
              ) : (
                tasks.map(t => (
                  <tr key={t.id} className="hover:bg-accent/40">
                    <td className="px-5 py-3 font-mono text-xs text-muted-foreground">{t.id?.slice(0, 8)}</td>
                    <td className="px-5 py-3 font-medium max-w-[200px] truncate">{t.title}</td>
                    <td className="px-5 py-3 text-muted-foreground">{t.category?.title ?? "—"}</td>
                    <td className="px-5 py-3 font-semibold tabular-nums whitespace-nowrap">{toman(t.budget)}</td>
                    <td className="px-5 py-3">
                      <span className="rounded-full bg-muted px-2 py-0.5 text-xs font-medium">
                        {STATUS_LABELS[t.status] ?? t.status}
                      </span>
                    </td>
                    <td className="px-5 py-3 text-left">
                      <Link to="/app/tasks/$id" params={{ id: t.id }} className="grid h-8 w-8 place-items-center rounded-lg hover:bg-accent">
                        <MoreHorizontal className="h-4 w-4" />
                      </Link>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}

function QA({ to, label }: { to: string; label: string }) {
  return (
    <Link to={to as any} className="block rounded-lg border p-3 text-sm font-medium hover:border-primary/40 hover:bg-accent">
      {label}
    </Link>
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
