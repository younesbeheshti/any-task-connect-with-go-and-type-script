import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Sparkles, Briefcase, Wallet, Star, BarChart3, ArrowUpRight, MapPin, Clock } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { useAuth } from "@/lib/auth-context";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/agent")({
  head: () => ({ meta: [{ title: "داشبورد انجام‌دهنده — تسک‌بریج" }] }),
  component: AgentDashboard,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Task = {
  id: string; title: string; budget: number; status: string;
  category: string; city: string; deadline: string;
};
type Application = {
  id: string; status: string;
  task: { id: string; title: string; budget: number; status: string };
};
type WalletData = { availableBalance: number; totalEarned?: number };

function AgentDashboard() {
  const { user } = useAuth();
  const [wallet, setWallet] = useState<WalletData | null>(null);
  const [openTasks, setOpenTasks] = useState<Task[]>([]);
  const [apps, setApps] = useState<Application[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/wallet`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/tasks?status=OPEN&page=1&pageSize=6`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/me/applications`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([w, t, a]) => {
      if (w && !w.error) setWallet(w);
      if (t) setOpenTasks(t.items ?? []);
      if (a) setApps(a.applications ?? a.items ?? []);
    }).finally(() => setLoading(false));
  }, []);

  const activeApps = apps.filter(a => a.status === "ACCEPTED");
  const avgRating = (user as any)?.rating ?? 0;

  const cards = [
    { label: "فرصت‌های جدید",     value: loading ? "..." : toFa(openTasks.length),              icon: Sparkles },
    { label: "درخواست‌های فعال",  value: loading ? "..." : toFa(activeApps.length),              icon: Briefcase },
    { label: "موجودی کیف پول",     value: loading ? "..." : toman(wallet?.availableBalance ?? 0, false), unit: "تومان", icon: Wallet },
    { label: "امتیاز کاربری",     value: loading ? "..." : toFa(avgRating ? avgRating.toFixed(1) : "—"),   icon: Star },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            سلام {user?.fullName ?? "..."} 👋
          </h1>
          <p className="text-sm text-muted-foreground">فرصت‌های همخوان با مهارت‌ها و شهر شما.</p>
        </div>
        <Link to="/app/tasks" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
          <Sparkles className="h-4 w-4" /> مشاهده فرصت‌ها
        </Link>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map(c => (
          <div key={c.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{c.label}</span>
              <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><c.icon className="h-4 w-4" /></span>
            </div>
            <div className="mt-3 font-display text-2xl font-bold tabular-nums">
              {c.value}{(c as any).unit && <span className="ms-1 text-sm text-muted-foreground font-medium">{(c as any).unit}</span>}
            </div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-3">
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold">فرصت‌های جدید</h2>
            <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline">
              مشاهده همه <ArrowUpRight className="h-3.5 w-3.5 rotate-180" />
            </Link>
          </div>
          {loading ? (
            <div className="grid gap-4 sm:grid-cols-2">
              {[1,2,3,4].map(i => <div key={i} className="rounded-2xl border bg-card p-5 animate-pulse"><div className="space-y-2"><div className="h-5 w-3/4 bg-muted rounded" /><div className="h-4 w-1/2 bg-muted rounded" /></div></div>)}
            </div>
          ) : openTasks.length === 0 ? (
            <div className="rounded-2xl border bg-card p-10 text-center">
              <Sparkles className="mx-auto h-10 w-10 text-muted-foreground/40" />
              <p className="mt-3 text-sm text-muted-foreground">در حال حاضر فرصت جدیدی وجود ندارد.</p>
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              {openTasks.map(t => (
                <Link key={t.id} to="/app/tasks/$id" params={{ id: t.id }}
                  className="rounded-2xl border bg-card p-5 shadow-soft hover:border-primary/40 block">
                  <div className="font-semibold line-clamp-1">{t.title}</div>
                  <div className="mt-1 flex flex-wrap gap-x-3 gap-y-1 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1"><MapPin className="h-3 w-3" />{t.city}</span>
                    {t.deadline && <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{t.deadline?.slice(0, 10)}</span>}
                  </div>
                  <div className="mt-2 font-semibold text-primary text-sm">{toman(t.budget)}</div>
                </Link>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-4">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <h3 className="font-display font-semibold">اقدامات سریع</h3>
            <div className="mt-4 space-y-2">
              <QA to="/app/tasks"       icon={Sparkles}  label="مشاهده فرصت‌های جدید" />
              <QA to="/app/applications" icon={Briefcase} label="درخواست‌های من" />
              <QA to="/app/earnings"    icon={BarChart3}  label="گزارش درآمد" />
              <QA to="/app/wallet"      icon={Wallet}     label="کیف پول" />
            </div>
          </div>

          {apps.length > 0 && (
            <div className="rounded-2xl border bg-card p-5 shadow-soft">
              <div className="flex items-center justify-between">
                <h3 className="font-display font-semibold">آخرین درخواست‌ها</h3>
                <Link to="/app/applications" className="text-xs font-medium text-primary hover:underline">همه</Link>
              </div>
              <div className="mt-3 divide-y">
                {apps.slice(0, 4).map(a => (
                  <div key={a.id} className="flex items-center justify-between py-2.5 text-sm">
                    <span className="truncate font-medium">{a.task?.title ?? "..."}</span>
                    <span className={cn("rounded-full px-2 py-0.5 text-xs font-medium", statusColor(a.status))}>
                      {statusLabel(a.status)}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function QA({ to, icon: Icon, label }: any) {
  return (
    <Link to={to} className="flex items-center gap-3 rounded-lg border p-3 text-sm font-medium hover:border-primary/40 hover:bg-accent">
      <span className="grid h-8 w-8 place-items-center rounded-lg bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
      <span className="flex-1">{label}</span>
      <ArrowUpRight className="h-4 w-4 text-muted-foreground rotate-180" />
    </Link>
  );
}

function statusColor(s: string) {
  switch (s) {
    case "PENDING": return "bg-warning/10 text-warning";
    case "ACCEPTED": return "bg-success/10 text-success";
    case "REJECTED": return "bg-destructive/10 text-destructive";
    default: return "bg-muted text-muted-foreground";
  }
}
function statusLabel(s: string) {
  switch (s) {
    case "PENDING": return "در انتظار";
    case "ACCEPTED": return "پذیرفته‌شده";
    case "REJECTED": return "رد‌شده";
    case "WITHDRAWN": return "پس‌گرفته‌شده";
    default: return s;
  }
}
