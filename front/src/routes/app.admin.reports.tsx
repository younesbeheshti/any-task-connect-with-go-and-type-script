import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { FileBarChart, Users, ClipboardList, TrendingUp, Download, Wallet, CheckCircle2 } from "lucide-react";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/app/admin/reports")({
  head: () => ({ meta: [{ title: "گزارش‌ها — مدیریت — تسک‌بریج" }] }),
  component: AdminReportsPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Metrics = {
  totalUsers: number;
  activeUsers?: number;
  totalTasks: number;
  activeTasks?: number;
  completedTasks?: number;
  totalRevenue: number;
  openDisputes: number;
};

type DailyPoint = { date: string; revenue: number; count?: number };

function AdminReportsPage() {
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [daily, setDaily] = useState<DailyPoint[]>([]);
  const [monthly, setMonthly] = useState<DailyPoint[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/admin/metrics`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/admin/revenue/daily`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/admin/revenue/monthly`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([m, d, mo]) => {
      if (m && !m.error) setMetrics(m);
      if (d) setDaily(d.items ?? d.data ?? []);
      if (mo) setMonthly(mo.items ?? mo.data ?? []);
    }).finally(() => setLoading(false));
  }, []);

  const kpis = [
    { label: "کل کاربران",            value: loading ? "..." : toFa(metrics?.totalUsers ?? 0),        sub: `${toFa(metrics?.activeUsers ?? 0)} فعال`,        icon: Users,         color: "text-primary bg-primary/10" },
    { label: "کل درخواست‌ها",         value: loading ? "..." : toFa(metrics?.totalTasks ?? 0),         sub: `${toFa(metrics?.activeTasks ?? 0)} فعال`,         icon: ClipboardList, color: "text-secondary bg-secondary/10" },
    { label: "درخواست‌های تکمیل‌شده", value: loading ? "..." : toFa(metrics?.completedTasks ?? 0),     sub: "تا به امروز",                                     icon: CheckCircle2,  color: "text-success bg-success/10" },
    { label: "کل درآمد پلتفرم",       value: loading ? "..." : toman(metrics?.totalRevenue ?? 0, false), sub: "کارمزد جمع‌شده",                                icon: FileBarChart,  color: "text-warning bg-warning/10", unit: "تومان" },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">گزارش‌های پلتفرم</h1>
          <p className="text-sm text-muted-foreground">شاخص‌های کلیدی عملکرد تسک‌بریج</p>
        </div>
        <button className="flex items-center gap-1.5 rounded-lg border bg-card px-4 py-2 text-sm font-medium hover:bg-accent">
          <Download className="h-4 w-4" /> خروجی Excel
        </button>
      </div>

      {/* KPI cards */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {kpis.map(k => (
          <div key={k.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{k.label}</span>
              <span className={`grid h-9 w-9 place-items-center rounded-lg ${k.color}`}><k.icon className="h-4 w-4" /></span>
            </div>
            <div className="mt-3 font-display text-xl font-bold tabular-nums">
              {k.value}{(k as any).unit && <span className="ms-1 text-sm font-medium text-muted-foreground">{(k as any).unit}</span>}
            </div>
            <div className="mt-1 text-xs text-muted-foreground">{k.sub}</div>
          </div>
        ))}
      </div>

      {/* Charts */}
      <div className="grid gap-6 lg:grid-cols-2">
        <ChartCard
          title="درآمد روزانه (۳۰ روز اخیر)"
          data={daily.slice(-30)}
          loading={loading}
          valueKey="revenue"
          labelKey="date"
          formatValue={v => toman(v, false)}
          emptyMsg="داده‌ای برای نمایش وجود ندارد"
        />
        <ChartCard
          title="درآمد ماهانه"
          data={monthly.slice(-12)}
          loading={loading}
          valueKey="revenue"
          labelKey="date"
          formatValue={v => toman(v, false)}
          emptyMsg="داده‌ای برای نمایش وجود ندارد"
        />
      </div>

      {/* Stats tables */}
      <div className="grid gap-6 lg:grid-cols-2">
        <ReportCard title="وضعیت کاربران" rows={loading || !metrics ? [] : [
          { label: "کل کاربران",      value: toFa(metrics.totalUsers) },
          { label: "کاربران فعال",    value: toFa(metrics.activeUsers ?? 0) },
          { label: "کاربران غیرفعال", value: toFa((metrics.totalUsers) - (metrics.activeUsers ?? 0)) },
        ]} loading={loading} />
        <ReportCard title="وضعیت درخواست‌ها" rows={loading || !metrics ? [] : [
          { label: "کل درخواست‌ها",          value: toFa(metrics.totalTasks) },
          { label: "درخواست‌های فعال",        value: toFa(metrics.activeTasks ?? 0) },
          { label: "درخواست‌های تکمیل‌شده",  value: toFa(metrics.completedTasks ?? 0) },
          { label: "اعتراضات باز",            value: toFa(metrics.openDisputes) },
        ]} loading={loading} />
      </div>

      {/* Revenue breakdown */}
      <div className="rounded-2xl border bg-card p-5 shadow-soft">
        <div className="flex items-center justify-between">
          <h2 className="font-display font-semibold">خلاصه مالی</h2>
          <Wallet className="h-5 w-5 text-muted-foreground" />
        </div>
        <div className="mt-4 grid gap-4 sm:grid-cols-3">
          <FinanceStat label="درآمد کل پلتفرم" value={toman(metrics?.totalRevenue ?? 0)} loading={loading} color="text-success" />
          <FinanceStat label="اعتراضات باز" value={toFa(metrics?.openDisputes ?? 0)} loading={loading} color="text-destructive" />
          <FinanceStat label="کل کاربران ثبت‌شده" value={toFa(metrics?.totalUsers ?? 0)} loading={loading} color="text-primary" />
        </div>
      </div>
    </div>
  );
}

function ChartCard({ title, data, loading, valueKey, labelKey, formatValue, emptyMsg }: {
  title: string;
  data: any[];
  loading: boolean;
  valueKey: string;
  labelKey: string;
  formatValue: (v: number) => string;
  emptyMsg: string;
}) {
  if (loading) {
    return (
      <div className="rounded-2xl border bg-card p-5 shadow-soft">
        <div className="h-5 w-1/2 rounded bg-muted animate-pulse" />
        <div className="mt-4 h-48 rounded-xl bg-muted animate-pulse" />
      </div>
    );
  }

  const hasData = data.length > 0 && data.some(d => (d[valueKey] ?? 0) > 0);

  if (!hasData) {
    return (
      <div className="rounded-2xl border bg-card p-5 shadow-soft">
        <h3 className="font-display font-semibold">{title}</h3>
        <div className="mt-4 flex h-48 items-center justify-center rounded-xl bg-muted/30">
          <div className="text-center">
            <TrendingUp className="mx-auto h-8 w-8 text-muted-foreground/40" />
            <p className="mt-2 text-sm text-muted-foreground">{emptyMsg}</p>
          </div>
        </div>
      </div>
    );
  }

  const values = data.map(d => d[valueKey] as number);
  const max = Math.max(...values, 1);
  const pts = values.map((v, i) => `${(i / Math.max(values.length - 1, 1)) * 100},${100 - (v / max) * 100}`).join(" ");
  const peak = values.indexOf(Math.max(...values));

  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-center justify-between">
        <h3 className="font-display font-semibold">{title}</h3>
        <span className="text-xs text-muted-foreground">{toFa(data.length)} نقطه</span>
      </div>
      <div className="mt-4 h-48 w-full">
        <svg viewBox="0 0 100 100" preserveAspectRatio="none" className="h-full w-full">
          <defs>
            <linearGradient id={`grad-${title}`} x1="0" x2="0" y1="0" y2="1">
              <stop offset="0%" stopColor="var(--primary)" stopOpacity="0.3" />
              <stop offset="100%" stopColor="var(--primary)" stopOpacity="0" />
            </linearGradient>
          </defs>
          <polyline points={`0,100 ${pts} 100,100`} fill={`url(#grad-${title})`} />
          <polyline points={pts} fill="none" stroke="var(--primary)" strokeWidth="1.2" vectorEffect="non-scaling-stroke" />
          {values.length > 0 && (
            <circle
              cx={(peak / Math.max(values.length - 1, 1)) * 100}
              cy={100 - (values[peak] / max) * 100}
              r="1.5" fill="var(--primary)" vectorEffect="non-scaling-stroke"
            />
          )}
        </svg>
      </div>
      <div className="mt-2 flex items-center justify-between text-xs text-muted-foreground">
        <span>{data[0]?.[labelKey]?.slice(0, 10)}</span>
        <span className="text-primary font-medium">اوج: {formatValue(values[peak])}</span>
        <span>{data[data.length - 1]?.[labelKey]?.slice(0, 10)}</span>
      </div>
    </div>
  );
}

function ReportCard({ title, rows, loading }: { title: string; rows: { label: string; value: string }[]; loading: boolean }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <h3 className="font-display font-semibold">{title}</h3>
      <div className="mt-4 space-y-3">
        {loading ? (
          Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="flex items-center justify-between animate-pulse">
              <div className="h-4 w-1/3 rounded bg-muted" />
              <div className="h-4 w-1/4 rounded bg-muted" />
            </div>
          ))
        ) : rows.map(r => (
          <div key={r.label} className="flex items-center justify-between">
            <span className="text-sm text-muted-foreground">{r.label}</span>
            <span className="font-semibold tabular-nums">{r.value}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

function FinanceStat({ label, value, loading, color }: { label: string; value: string; loading: boolean; color: string }) {
  return (
    <div className="rounded-xl border bg-muted/20 p-4">
      <div className="text-sm text-muted-foreground">{label}</div>
      <div className={`mt-2 font-display text-xl font-bold tabular-nums ${color}`}>
        {loading ? <span className="inline-block h-6 w-24 animate-pulse rounded bg-muted" /> : value}
      </div>
    </div>
  );
}
