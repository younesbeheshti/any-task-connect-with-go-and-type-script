import { createFileRoute } from "@tanstack/react-router";
import { Sparkles, Briefcase, Wallet, Star, TrendingUp } from "lucide-react";
import { tasks } from "@/lib/mock-data";
import { TaskCard } from "@/components/task-card";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/app/agent")({
  head: () => ({ meta: [{ title: "داشبورد انجام‌دهنده — تسک‌بریج" }] }),
  component: AgentDashboard,
});

function AgentDashboard() {
  const open = tasks.filter(t => t.status === "awaiting_applicants" || t.status === "posted");
  const mine = tasks.filter(t => t.status === "accepted" || t.status === "in_progress");
  const cards = [
    { label: "فرصت‌های جدید",     value: toFa(open.length),  icon: Sparkles },
    { label: "درخواست‌های فعال",  value: toFa(mine.length),  icon: Briefcase },
    { label: "درآمد این ماه",     value: toman(18200000),    icon: Wallet, big: true },
    { label: "امتیاز کاربری",     value: toFa("4.9"),        icon: Star },
  ];
  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">داشبورد انجام‌دهنده</h1>
          <p className="text-sm text-muted-foreground">فرصت‌های همخوان با مهارت‌ها و شهر شما.</p>
        </div>
        <div className="inline-flex items-center gap-1.5 rounded-lg border bg-card px-3 py-1.5 text-xs">
          <TrendingUp className="h-3.5 w-3.5 text-success" />
          <span className="text-success">۱۲٪</span>
          <span className="text-muted-foreground">رشد درآمد نسبت به ماه قبل</span>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map(c => (
          <div key={c.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{c.label}</span>
              <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><c.icon className="h-4 w-4" /></span>
            </div>
            <div className={`mt-3 font-display font-bold tabular-nums ${c.big ? "text-xl" : "text-3xl"}`}>{c.value}</div>
          </div>
        ))}
      </div>

      <section>
        <h2 className="mb-3 font-display text-lg font-semibold">پیشنهادهای ویژه شما</h2>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {open.map(t => <TaskCard key={t.id} task={t} showApply />)}
        </div>
      </section>

      <section>
        <h2 className="mb-3 font-display text-lg font-semibold">درخواست‌های در دست انجام</h2>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {mine.map(t => <TaskCard key={t.id} task={t} />)}
        </div>
      </section>
    </div>
  );
}
