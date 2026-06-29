import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ArrowUpRight, ClipboardList, CheckCircle2, Wallet, Lock, PlusCircle, MessagesSquare, ListChecks, Bell, Activity } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { useAuth } from "@/lib/auth-context";

export const Route = createFileRoute("/app/")({
  head: () => ({ meta: [{ title: "داشبورد — تسک‌بریج" }] }),
  component: Dashboard,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type UserStats = { postedTasks: number; completedTasks: number; activeTasks: number; totalApplications: number };
type WalletData = { availableBalance: number; lockedBalance: number; currency: string };
type Notification = { id: string; title: string; body: string; isRead: boolean; createdAt: string };
type Task = { id: string; title: string; status: string; budget: number; category: string; city: string };

function Dashboard() {
  const { user } = useAuth();
  const [stats, setStats] = useState<UserStats | null>(null);
  const [wallet, setWallet] = useState<WalletData | null>(null);
  const [tasks, setTasks] = useState<Task[]>([]);
  const [notifs, setNotifs] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/dashboard/stats`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/wallet`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/tasks?mine=true&page=1&pageSize=4`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/notifications?page=1&pageSize=4`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([s, w, t, n]) => {
      if (s) setStats(s);
      if (w && !w.error) setWallet(w);
      if (t) setTasks(t.items ?? []);
      if (n) setNotifs(n.items ?? []);
    }).finally(() => setLoading(false));
  }, []);

  const cards = [
    { label: "درخواست‌های فعال",      value: loading ? "..." : toFa(stats?.activeTasks ?? 0),               icon: ClipboardList, trend: "در حال انجام",     color: "from-primary to-secondary" },
    { label: "درخواست‌های تکمیل‌شده", value: loading ? "..." : toFa(stats?.completedTasks ?? 0),            icon: CheckCircle2,  trend: "مجموع",          color: "from-success to-secondary" },
    { label: "موجودی کیف پول",         value: loading ? "..." : toman(wallet?.availableBalance ?? 0, false), unit: "تومان", icon: Wallet, trend: "قابل برداشت", color: "from-secondary to-primary" },
    { label: "در انتظار آزادسازی",     value: loading ? "..." : toman(wallet?.lockedBalance ?? 0, false),    unit: "تومان", icon: Lock,   trend: "امانت",       color: "from-warning to-destructive" },
  ];

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">سلام {user?.fullName ?? "..."} 👋</h1>
          <p className="text-sm text-muted-foreground">خلاصه‌ای از درخواست‌های امروز شما.</p>
        </div>
        <Link to="/app/tasks/new" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
          <PlusCircle className="h-4 w-4" /> ثبت درخواست جدید
        </Link>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map(c => (
          <div key={c.label} className="relative overflow-hidden rounded-2xl border bg-card p-5 shadow-soft">
            <div className={`absolute -left-6 -top-6 h-24 w-24 rounded-full bg-gradient-to-br ${c.color} opacity-10`} />
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{c.label}</span>
              <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><c.icon className="h-4 w-4" /></span>
            </div>
            <div className="mt-3 font-display text-2xl font-bold tabular-nums">
              {c.value}{c.unit && <span className="ms-1 text-sm text-muted-foreground font-medium">{c.unit}</span>}
            </div>
            <div className="mt-1 text-xs text-muted-foreground">{c.trend}</div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-3">
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold">درخواست‌های اخیر</h2>
            <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline">
              مشاهده همه <ArrowUpRight className="h-3.5 w-3.5 rotate-180" />
            </Link>
          </div>
          {loading ? (
            <div className="grid gap-4 sm:grid-cols-2">
              {[1,2,3,4].map(i => <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse"><div className="space-y-2"><div className="h-5 w-3/4 rounded bg-muted" /><div className="h-4 w-1/2 rounded bg-muted" /></div></div>)}
            </div>
          ) : tasks.length === 0 ? (
            <div className="rounded-2xl border bg-card p-10 text-center shadow-soft">
              <Activity className="mx-auto h-10 w-10 text-muted-foreground/40" />
              <p className="mt-3 text-sm text-muted-foreground">هنوز درخواستی ثبت نشده است.</p>
              <Link to="/app/tasks/new" className="mt-3 inline-block rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white">اولین درخواست را ثبت کنید</Link>
            </div>
          ) : (
            <div className="grid gap-4 sm:grid-cols-2">
              {tasks.map(t => (
                <Link key={t.id} to="/app/tasks/$id" params={{ id: t.id }} className="rounded-2xl border bg-card p-5 shadow-soft hover:border-primary/40 block">
                  <div className="font-semibold line-clamp-1">{t.title}</div>
                  <div className="mt-1 flex items-center gap-2 text-xs text-muted-foreground">
                    <span>{t.category}</span><span>·</span><span>{t.city}</span>
                    <span>·</span><span className="text-primary font-medium">{toman(t.budget)}</span>
                  </div>
                </Link>
              ))}
            </div>
          )}
        </div>

        <div className="space-y-6">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <h3 className="font-display font-semibold">اقدامات سریع</h3>
            <div className="mt-4 space-y-2">
              <QA to="/app/tasks/new" icon={PlusCircle}     label="ثبت درخواست جدید" />
              <QA to="/app/tasks"     icon={ListChecks}     label="مشاهده درخواست‌های فعال" />
              <QA to="/app/chat"      icon={MessagesSquare} label="باز کردن گفت‌وگوها" />
              <QA to="/app/wallet"    icon={Wallet}         label="شارژ کیف پول" />
            </div>
          </div>
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center justify-between">
              <h3 className="font-display font-semibold">اعلان‌ها</h3>
              <Link to="/app/notifications" className="text-xs font-medium text-primary hover:underline">مشاهده همه</Link>
            </div>
            <div className="mt-3 divide-y">
              {notifs.length === 0 ? (
                <p className="py-4 text-xs text-muted-foreground text-center">اعلان جدیدی ندارید</p>
              ) : notifs.slice(0, 4).map(n => (
                <div key={n.id} className="flex items-start gap-3 py-3">
                  <span className="mt-1 grid h-8 w-8 place-items-center rounded-lg bg-primary/10 text-primary"><Bell className="h-4 w-4" /></span>
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-sm font-medium">{n.title}</div>
                    <div className="truncate text-xs text-muted-foreground">{n.body}</div>
                  </div>
                  <div className="text-xs text-muted-foreground whitespace-nowrap">{n.createdAt?.slice(0, 10)}</div>
                </div>
              ))}
            </div>
          </div>
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
