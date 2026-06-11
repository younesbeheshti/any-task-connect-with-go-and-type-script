import { createFileRoute } from "@tanstack/react-router";
import { ArrowDownLeft, ArrowUpRight, Lock, Wallet as WalletIcon, TrendingUp } from "lucide-react";
import { transactions } from "@/lib/mock-data";
import { toman, toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/wallet")({
  head: () => ({ meta: [{ title: "کیف پول — تسک‌بریج" }] }),
  component: Wallet,
});

function Wallet() {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 lg:grid-cols-3">
        <div className="relative overflow-hidden rounded-2xl gradient-brand p-6 text-white shadow-elevated lg:col-span-2">
          <div className="absolute inset-0 grid-bg opacity-15" />
          <div className="relative flex items-start justify-between">
            <div>
              <div className="text-sm opacity-80">موجودی قابل برداشت</div>
              <div className="mt-2 font-display text-4xl font-bold tracking-tight tabular-nums sm:text-5xl">{toman(12400000, false)}</div>
              <div className="mt-1 text-sm opacity-80">تومان · کیف پول اصلی</div>
            </div>
            <WalletIcon className="h-6 w-6 opacity-80" />
          </div>
          <div className="relative mt-8 flex flex-wrap gap-2">
            <button className="rounded-lg bg-white px-4 py-2 text-sm font-semibold text-primary hover:bg-white/90">شارژ کیف پول</button>
            <button className="rounded-lg border border-white/30 bg-white/10 px-4 py-2 text-sm font-semibold backdrop-blur hover:bg-white/20">برداشت وجه</button>
            <button className="rounded-lg border border-white/30 bg-white/10 px-4 py-2 text-sm font-semibold backdrop-blur hover:bg-white/20">صورت‌حساب</button>
          </div>
        </div>
        <div className="space-y-4">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <Lock className="h-4 w-4 text-warning" />
              <span className="text-sm text-muted-foreground">موجودی بلوکه شده</span>
            </div>
            <div className="mt-2 font-display text-2xl font-bold tabular-nums">{toman(3650000)}</div>
            <div className="mt-1 text-xs text-muted-foreground">در {toFa(3)} درخواست فعال</div>
          </div>
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4 text-success" />
              <span className="text-sm text-muted-foreground">گردش این ماه</span>
            </div>
            <div className="mt-2 font-display text-2xl font-bold tabular-nums">+{toman(21400000)}</div>
            <div className="mt-1 text-xs text-muted-foreground">{toFa(12)} تراکنش</div>
          </div>
        </div>
      </div>

      <div className="rounded-2xl border bg-card p-5 shadow-soft">
        <div className="flex items-center justify-between">
          <h2 className="font-display text-lg font-semibold">نمای کلی هزینه‌ها</h2>
          <span className="text-xs text-muted-foreground">{toFa(7)} روز اخیر</span>
        </div>
        <MiniChart />
      </div>

      <div className="rounded-2xl border bg-card shadow-soft">
        <div className="flex items-center justify-between border-b px-5 py-4">
          <h2 className="font-display text-lg font-semibold">تاریخچه تراکنش‌ها</h2>
          <button className="rounded-lg border px-3 py-1.5 text-xs font-medium hover:bg-accent">خروجی CSV</button>
        </div>
        <div className="divide-y">
          {transactions.map(t => {
            const positive = t.amount > 0;
            return (
              <div key={t.id} className="flex items-center gap-4 px-5 py-3.5">
                <span className={cn("grid h-9 w-9 place-items-center rounded-lg", positive ? "bg-success/10 text-success" : "bg-primary/10 text-primary")}>
                  {positive ? <ArrowDownLeft className="h-4 w-4" /> : <ArrowUpRight className="h-4 w-4" />}
                </span>
                <div className="min-w-0 flex-1">
                  <div className="truncate text-sm font-medium">{t.description}</div>
                  <div className="text-xs text-muted-foreground">{t.date} · {t.status}</div>
                </div>
                <div className={cn("text-sm font-semibold tabular-nums whitespace-nowrap", positive ? "text-success" : "text-foreground")}>
                  {positive ? "+" : "-"}{toman(Math.abs(t.amount))}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}

function MiniChart() {
  const data = [40, 65, 35, 80, 55, 95, 70];
  const labels = ["ش","ی","د","س","چ","پ","ج"];
  const max = Math.max(...data);
  return (
    <div className="mt-4 flex h-32 items-end gap-2">
      {data.map((v, i) => (
        <div key={i} className="flex flex-1 flex-col items-center gap-1.5">
          <div className="w-full rounded-md gradient-brand transition-all hover:opacity-80" style={{ height: `${(v / max) * 100}%` }} />
          <span className="text-[10px] text-muted-foreground">{labels[i]}</span>
        </div>
      ))}
    </div>
  );
}
