import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowUpRight, ClipboardList, CheckCircle2, Wallet, Lock, PlusCircle, MessagesSquare, ListChecks, Bell } from "lucide-react";
import { tasks, notifications } from "@/lib/mock-data";
import { TaskCard } from "@/components/task-card";
import { toman, toFa } from "@/lib/fa";
import { useAuth } from "@/lib/auth-context";

export const Route = createFileRoute("/app/")({
  head: () => ({ meta: [{ title: "داشبورد — تسک‌بریج" }] }),
  component: Dashboard,
});

function Dashboard() {
  const { user } = useAuth();
  const cards = [
    { label: "درخواست‌های فعال",      value: toFa(4),                icon: ClipboardList, trend: "+۲ این هفته",            color: "from-primary to-secondary" },
    { label: "درخواست‌های تکمیل شده",  value: toFa(27),               icon: CheckCircle2,  trend: "+۵ این ماه",             color: "from-success to-secondary" },
    { label: "موجودی کیف پول",         value: toman(12400000, false), unit: "تومان",       icon: Wallet, trend: "قابل برداشت", color: "from-secondary to-primary" },
    { label: "در انتظار آزادسازی",     value: toman(3650000, false),  unit: "تومان",       icon: Lock,   trend: "در ۳ درخواست", color: "from-warning to-destructive" },
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
            <div className="mt-3 font-display text-2xl font-bold tabular-nums">{c.value}{c.unit && <span className="ms-1 text-sm text-muted-foreground font-medium">{c.unit}</span>}</div>
            <div className="mt-1 text-xs text-muted-foreground">{c.trend}</div>
          </div>
        ))}
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-3">
          <div className="flex items-center justify-between">
            <h2 className="font-display text-lg font-semibold">درخواست‌های اخیر</h2>
            <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline">مشاهده همه <ArrowUpRight className="h-3.5 w-3.5 rotate-180" /></Link>
          </div>
          <div className="grid gap-4 sm:grid-cols-2">
            {tasks.slice(0, 4).map(t => <TaskCard key={t.id} task={t} />)}
          </div>
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
              {notifications.slice(0, 4).map(n => (
                <div key={n.id} className="flex items-start gap-3 py-3">
                  <span className="mt-1 grid h-8 w-8 place-items-center rounded-lg bg-primary/10 text-primary"><Bell className="h-4 w-4" /></span>
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-sm font-medium">{n.title}</div>
                    <div className="truncate text-xs text-muted-foreground">{n.desc}</div>
                  </div>
                  <div className="text-xs text-muted-foreground whitespace-nowrap">{n.time}</div>
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
