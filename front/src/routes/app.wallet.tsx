import { createFileRoute } from "@tanstack/react-router";
import { ArrowDownLeft, ArrowUpRight, Lock, Wallet as WalletIcon, TrendingUp, X } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";
import { useEffect, useState } from "react";

export const Route = createFileRoute("/app/wallet")({
  head: () => ({ meta: [{ title: "کیف پول — تسک‌بریج" }] }),
  component: WalletPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

function getToken(): string | null {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return null;
    const parsed = JSON.parse(stored);
    return parsed?.token ?? null;
  } catch {
    return null;
  }
}

async function apiFetch<T>(path: string, token: string): Promise<T> {
  const res = await fetch(`${API_BASE}/v1${path}`, {
    headers: { Authorization: `Bearer ${token}` },
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error?.message ?? "خطای سرور");
  return data as T;
}

async function apiPost<T>(path: string, token: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}/v1${path}`, {
    method: "POST",
    headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  if (!res.ok) throw new Error(data?.error?.message ?? "خطای سرور");
  return data as T;
}

type WalletData = {
  id: string;
  availableBalance: number;
  lockedBalance: number;
  pendingWithdrawBalance: number;
  totalEarned: number;
  totalSpent: number;
  currency: string;
};

type TxItem = {
  id: string;
  amount: number;
  currency: string;
  type: string;
  status: string;
  description: string;
  date: string;
};

function WalletPage() {
  const [wallet, setWallet] = useState<WalletData | null>(null);
  const [txs, setTxs] = useState<TxItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showTopUp, setShowTopUp] = useState(false);
  const [showWithdraw, setShowWithdraw] = useState(false);

  useEffect(() => {
    const token = getToken();
    if (!token) {
      setError("لطفاً وارد حساب کاربری شوید");
      setLoading(false);
      return;
    }
    Promise.all([
      apiFetch<WalletData>("/wallet", token).catch(() => null),
      apiFetch<{ items: TxItem[] }>("/transactions", token).catch(() => ({ items: [] })),
    ]).then(([walletData, txData]) => {
      setWallet(walletData);
      setTxs(txData?.items ?? []);
    }).catch((e) => setError(e.message)).finally(() => setLoading(false));
  }, []);

  if (loading) {
    return <div className="flex h-48 items-center justify-center text-muted-foreground">در حال بارگذاری...</div>;
  }

  if (error) {
    return <div className="rounded-xl border border-destructive/30 bg-destructive/5 p-6 text-destructive">{error}</div>;
  }

  const available = wallet?.availableBalance ?? 0;
  const locked = wallet?.lockedBalance ?? 0;

  return (
    <div className="space-y-6">
      {showTopUp && <TopUpModal onClose={() => setShowTopUp(false)} />}
      {showWithdraw && <WithdrawModal onClose={() => setShowWithdraw(false)} onSuccess={() => {
        setShowWithdraw(false);
        // re-fetch wallet after withdraw
        const token = getToken();
        if (token) apiFetch<WalletData>("/wallet", token).then(setWallet).catch(() => {});
      }} />}

      <div className="grid gap-4 lg:grid-cols-3">
        <div className="relative overflow-hidden rounded-2xl gradient-brand p-6 text-white shadow-elevated lg:col-span-2">
          <div className="absolute inset-0 grid-bg opacity-15" />
          <div className="relative flex items-start justify-between">
            <div>
              <div className="text-sm opacity-80">موجودی قابل برداشت</div>
              <div className="mt-2 font-display text-4xl font-bold tracking-tight tabular-nums sm:text-5xl">
                {toman(available, false)}
              </div>
              <div className="mt-1 text-sm opacity-80">تومان · کیف پول اصلی</div>
            </div>
            <WalletIcon className="h-6 w-6 opacity-80" />
          </div>
          <div className="relative mt-8 flex flex-wrap gap-2">
            <button
              onClick={() => setShowTopUp(true)}
              className="rounded-lg bg-white px-4 py-2 text-sm font-semibold text-primary hover:bg-white/90"
            >
              شارژ کیف پول
            </button>
            <button
              onClick={() => setShowWithdraw(true)}
              className="rounded-lg border border-white/30 bg-white/10 px-4 py-2 text-sm font-semibold backdrop-blur hover:bg-white/20"
            >
              برداشت وجه
            </button>
          </div>
        </div>

        <div className="space-y-4">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <Lock className="h-4 w-4 text-warning" />
              <span className="text-sm text-muted-foreground">موجودی بلوکه شده</span>
            </div>
            <div className="mt-2 font-display text-2xl font-bold tabular-nums">{toman(locked)}</div>
          </div>
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <TrendingUp className="h-4 w-4 text-success" />
              <span className="text-sm text-muted-foreground">کل درآمد</span>
            </div>
            <div className="mt-2 font-display text-2xl font-bold tabular-nums">
              +{toman(wallet?.totalEarned ?? 0)}
            </div>
            <div className="mt-1 text-xs text-muted-foreground">{toFa(txs.length)} تراکنش</div>
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
        </div>
        <div className="divide-y">
          {txs.length === 0 ? (
            <div className="px-5 py-8 text-center text-sm text-muted-foreground">تراکنشی یافت نشد</div>
          ) : (
            txs.map((t) => {
              const positive = t.type === "topup" || t.type === "release";
              return (
                <div key={t.id} className="flex items-center gap-4 px-5 py-3.5">
                  <span className={cn("grid h-9 w-9 place-items-center rounded-lg", positive ? "bg-success/10 text-success" : "bg-primary/10 text-primary")}>
                    {positive ? <ArrowDownLeft className="h-4 w-4" /> : <ArrowUpRight className="h-4 w-4" />}
                  </span>
                  <div className="min-w-0 flex-1">
                    <div className="truncate text-sm font-medium">{t.description || txTypeLabel(t.type)}</div>
                    <div className="text-xs text-muted-foreground">{t.date?.slice(0, 10)} · {txStatusLabel(t.status)}</div>
                  </div>
                  <div className={cn("text-sm font-semibold tabular-nums whitespace-nowrap", positive ? "text-success" : "text-foreground")}>
                    {positive ? "+" : "-"}{toman(Math.abs(t.amount))}
                  </div>
                </div>
              );
            })
          )}
        </div>
      </div>
    </div>
  );
}

function txTypeLabel(type: string): string {
  const map: Record<string, string> = {
    topup: "شارژ کیف پول",
    escrow_in: "بلوکه اسکرو",
    release: "آزادسازی اسکرو",
    refund: "استرداد وجه",
    withdraw: "برداشت",
    fee: "کارمزد",
  };
  return map[type] ?? type;
}

function txStatusLabel(status: string): string {
  const map: Record<string, string> = {
    pending: "در انتظار",
    completed: "تکمیل شده",
    failed: "ناموفق",
    locked: "بلوکه",
    released: "آزاد شده",
    reversed: "برگشت خورده",
  };
  return map[status] ?? status;
}

function TopUpModal({ onClose }: { onClose: () => void }) {
  const [amount, setAmount] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = getToken();
    if (!token) { setError("احراز هویت لازم است"); return; }
    setLoading(true);
    setError(null);
    try {
      await apiPost("/wallet/topup", token, { amount: parseInt(amount) * 10 });
      onClose();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "خطا");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="w-full max-w-md rounded-2xl border bg-card p-6 shadow-elevated">
        <div className="flex items-center justify-between">
          <h3 className="font-display text-lg font-semibold">شارژ کیف پول</h3>
          <button onClick={onClose}><X className="h-5 w-5" /></button>
        </div>
        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="text-sm text-muted-foreground">مبلغ (تومان)</label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="mt-1 w-full rounded-lg border bg-background px-3 py-2 text-sm"
              placeholder="مثال: ۵۰۰۰۰۰"
              required
              min={1000}
            />
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg gradient-brand py-2.5 text-sm font-semibold text-white disabled:opacity-50"
          >
            {loading ? "در حال پردازش..." : "ادامه"}
          </button>
        </form>
      </div>
    </div>
  );
}

function WithdrawModal({ onClose, onSuccess }: { onClose: () => void; onSuccess: () => void }) {
  const [amount, setAmount] = useState("");
  const [sheba, setSheba] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = getToken();
    if (!token) { setError("احراز هویت لازم است"); return; }
    setLoading(true);
    setError(null);
    try {
      await apiPost("/wallet/withdraw", token, {
        amount: parseInt(amount) * 10,
        shebaNumber: sheba,
      });
      onSuccess();
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "خطا");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="w-full max-w-md rounded-2xl border bg-card p-6 shadow-elevated">
        <div className="flex items-center justify-between">
          <h3 className="font-display text-lg font-semibold">برداشت وجه</h3>
          <button onClick={onClose}><X className="h-5 w-5" /></button>
        </div>
        <form onSubmit={handleSubmit} className="mt-4 space-y-4">
          <div>
            <label className="text-sm text-muted-foreground">مبلغ (تومان)</label>
            <input
              type="number"
              value={amount}
              onChange={(e) => setAmount(e.target.value)}
              className="mt-1 w-full rounded-lg border bg-background px-3 py-2 text-sm"
              placeholder="مثال: ۵۰۰۰۰۰"
              required
              min={1000}
            />
          </div>
          <div>
            <label className="text-sm text-muted-foreground">شماره شبا</label>
            <input
              type="text"
              value={sheba}
              onChange={(e) => setSheba(e.target.value)}
              className="mt-1 w-full rounded-lg border bg-background px-3 py-2 text-sm"
              placeholder="IR..."
              dir="ltr"
            />
          </div>
          {error && <p className="text-sm text-destructive">{error}</p>}
          <button
            type="submit"
            disabled={loading}
            className="w-full rounded-lg bg-primary py-2.5 text-sm font-semibold text-primary-foreground disabled:opacity-50"
          >
            {loading ? "در حال پردازش..." : "ثبت درخواست"}
          </button>
        </form>
      </div>
    </div>
  );
}

function MiniChart() {
  const data = [40, 65, 35, 80, 55, 95, 70];
  const labels = ["ش", "ی", "د", "س", "چ", "پ", "ج"];
  const max = Math.max(...data);
  return (
    <div className="mt-4 flex h-32 items-end gap-2">
      {data.map((v, i) => (
        <div key={i} className="flex flex-1 flex-col items-center gap-1.5">
          <div
            className="w-full rounded-md gradient-brand transition-all hover:opacity-80"
            style={{ height: `${(v / max) * 100}%` }}
          />
          <span className="text-[10px] text-muted-foreground">{labels[i]}</span>
        </div>
      ))}
    </div>
  );
}
