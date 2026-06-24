import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { BarChart3, TrendingUp, Wallet, ArrowDownLeft, Clock } from "lucide-react";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/app/earnings")({
  head: () => ({ meta: [{ title: "درآمدها — تسک‌بریج" }] }),
  component: EarningsPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type WalletData = {
  availableBalance: number;
  lockedBalance: number;
  pendingWithdrawBalance: number;
  totalEarned: number;
  totalSpent: number;
  currency: string;
};

type TxItem = {
  id: string; amount: number; type: string; status: string;
  description: string; createdAt: string;
};

type Stats = { totalEarned: number; totalSpent: number; txCount: number };

function EarningsPage() {
  const [wallet, setWallet] = useState<WalletData | null>(null);
  const [txs, setTxs] = useState<TxItem[]>([]);
  const [stats, setStats] = useState<Stats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    const headers = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/wallet`, { headers }).then(r => r.json()),
      fetch(`${API_BASE}/v1/wallet/statistics`, { headers }).then(r => r.json()),
      fetch(`${API_BASE}/v1/transactions?page=1`, { headers }).then(r => r.json()),
    ])
      .then(([w, s, t]) => {
        setWallet(w);
        setStats(s);
        setTxs(t.items ?? []);
      })
      .catch(() => setError("خطا در بارگذاری اطلاعات درآمد"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  const income = txs.filter(t => t.type === "payment" || t.type === "escrow_out");
  const withdrawable = wallet?.availableBalance ?? 0;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">درآمدها</h1>
          <p className="text-sm text-muted-foreground">خلاصه‌ای از درآمد و تراکنش‌های مالی شما</p>
        </div>
        <Link to="/app/wallet" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white">
          <ArrowDownLeft className="h-4 w-4" /> برداشت وجه
        </Link>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <StatCard label="موجودی قابل برداشت" value={toman(withdrawable)} icon={Wallet} color="from-primary to-secondary" />
        <StatCard label="کل درآمد" value={toman(stats?.totalEarned ?? 0)} icon={TrendingUp} color="from-success to-secondary" />
        <StatCard label="در انتظار تایید" value={toman(wallet?.pendingWithdrawBalance ?? 0)} icon={Clock} color="from-warning to-destructive" />
        <StatCard label="تعداد تراکنش‌ها" value={toFa(stats?.txCount ?? 0)} icon={BarChart3} color="from-secondary to-primary" />
      </div>

      <div className="rounded-2xl border bg-card shadow-soft">
        <div className="flex items-center justify-between border-b px-6 py-4">
          <h2 className="font-display font-semibold">تاریخچه تراکنش‌ها</h2>
          <Link to="/app/wallet" className="text-xs font-medium text-primary hover:underline">مشاهده همه</Link>
        </div>
        {txs.length === 0 ? (
          <div className="p-12 text-center">
            <BarChart3 className="mx-auto h-10 w-10 text-muted-foreground/40" />
            <p className="mt-3 text-sm text-muted-foreground">تراکنشی ثبت نشده است</p>
          </div>
        ) : (
          <div className="divide-y">
            {txs.slice(0, 10).map(tx => (
              <div key={tx.id} className="flex items-center gap-4 px-6 py-3">
                <span className="grid h-9 w-9 place-items-center rounded-lg bg-success/10 text-success">
                  <ArrowDownLeft className="h-4 w-4" />
                </span>
                <div className="flex-1 min-w-0">
                  <div className="truncate text-sm font-medium">{tx.description || txTypeLabel(tx.type)}</div>
                  <div className="text-xs text-muted-foreground">{tx.createdAt?.slice(0, 10)}</div>
                </div>
                <span className="font-semibold text-success">{toman(tx.amount)}</span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function txTypeLabel(t: string) {
  const m: Record<string, string> = {
    payment: "پرداخت از کارفرما", escrow_out: "آزادسازی امانت",
    withdraw: "برداشت", topup: "شارژ", refund: "بازگشت وجه",
  };
  return m[t] ?? t;
}

function StatCard({ label, value, icon: Icon, color }: { label: string; value: string; icon: any; color: string }) {
  return (
    <div className="relative overflow-hidden rounded-2xl border bg-card p-5 shadow-soft">
      <div className={`absolute -left-6 -top-6 h-24 w-24 rounded-full bg-gradient-to-br ${color} opacity-10`} />
      <div className="flex items-start justify-between">
        <span className="text-sm text-muted-foreground">{label}</span>
        <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
      </div>
      <div className="mt-3 font-display text-xl font-bold tabular-nums">{value}</div>
    </div>
  );
}

function Skeleton() {
  return (
    <div className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {[1, 2, 3, 4].map(i => (
          <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse">
            <div className="h-4 w-2/3 rounded bg-muted" />
            <div className="mt-3 h-6 w-1/2 rounded bg-muted" />
          </div>
        ))}
      </div>
    </div>
  );
}

function ErrBox({ msg }: { msg: string }) {
  return <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center"><p className="text-sm text-destructive">{msg}</p></div>;
}
