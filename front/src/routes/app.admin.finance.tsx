import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import {
  Wallet, TrendingUp, ArrowDownLeft, Clock, CheckCircle2, XCircle,
  Users, ClipboardList, Filter,
} from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/admin/finance")({
  head: () => ({ meta: [{ title: "مالی — مدیریت — تسک‌بریج" }] }),
  component: AdminFinancePage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type RevenueStats = { total: number; thisMonth: number; txCount: number };
type FinancialOverview = {
  totalDeposited: number; totalWithdrawn: number; totalLocked: number;
  platformRevenue: number; activeWallets: number;
};
type WithdrawItem = {
  id: string; amount: number; status: string;
  bankAccount: string; shebaNumber: string;
  createdAt: string; userId?: string; userName?: string;
};
type TxItem = {
  id: string; amount: number; type: string; status: string;
  description: string; createdAt: string;
};
type PageResult<T> = { items: T[]; total: number; page: number; pageSize: number };

const WITHDRAW_STATUS_LABELS: Record<string, string> = {
  pending: "در انتظار", approved: "تایید‌شده", rejected: "رد‌شده", paid: "پرداخت‌شده",
  PENDING: "در انتظار", APPROVED: "تایید‌شده", REJECTED: "رد‌شده", PAID: "پرداخت‌شده",
};
const WITHDRAW_STATUS_COLORS: Record<string, string> = {
  pending: "bg-warning/10 text-warning", PENDING: "bg-warning/10 text-warning",
  approved: "bg-success/10 text-success", APPROVED: "bg-success/10 text-success",
  rejected: "bg-destructive/10 text-destructive", REJECTED: "bg-destructive/10 text-destructive",
  paid: "bg-primary/10 text-primary", PAID: "bg-primary/10 text-primary",
};

const TX_TYPE_LABELS: Record<string, string> = {
  ESCROW_LOCK: "بلوکه امانت", ESCROW_RELEASE: "آزادسازی امانت",
  TOP_UP: "شارژ کیف پول", WITHDRAW: "برداشت", PLATFORM_FEE: "کارمزد سامانه",
  REFUND: "بازگشت وجه",
};

function AdminFinancePage() {
  const [stats, setStats] = useState<RevenueStats | null>(null);
  const [overview, setOverview] = useState<FinancialOverview | null>(null);
  const [withdraws, setWithdraws] = useState<WithdrawItem[]>([]);
  const [txs, setTxs] = useState<TxItem[]>([]);
  const [wTotal, setWTotal] = useState(0);
  const [tTotal, setTTotal] = useState(0);
  const [loading, setLoading] = useState(true);
  const [tab, setTab] = useState<"withdraws" | "transactions">("withdraws");
  const [wPage, setWPage] = useState(1);
  const [tPage, setTPage] = useState(1);
  const [wStatusFilter, setWStatusFilter] = useState("");

  function loadWithdraws(p = wPage, sf = wStatusFilter) {
    const token = getToken();
    if (!token) return;
    const params = new URLSearchParams({ page: String(p), pageSize: "15" });
    if (sf) params.set("status", sf);
    fetch(`${API_BASE}/v1/admin/financial/withdraws?${params}`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => { setWithdraws(d.items ?? []); setWTotal(d.total ?? 0); })
      .catch(() => {});
  }

  function loadTxs(p = tPage) {
    const token = getToken();
    if (!token) return;
    fetch(`${API_BASE}/v1/admin/financial/transactions?page=${p}&pageSize=15`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => { setTxs(d.items ?? []); setTTotal(d.total ?? 0); })
      .catch(() => {});
  }

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/admin/revenue/statistics`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/admin/financial/overview`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([s, ov]) => {
      if (s && !s.error) setStats(s);
      if (ov && !ov.error) setOverview(ov);
    }).finally(() => setLoading(false));
    loadWithdraws(1);
    loadTxs(1);
  }, []);

  useEffect(() => { loadWithdraws(); }, [wPage]);
  useEffect(() => { loadTxs(); }, [tPage]);
  useEffect(() => { setWPage(1); loadWithdraws(1, wStatusFilter); }, [wStatusFilter]);

  async function approveWithdraw(id: string) {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch(`${API_BASE}/v1/admin/withdraws/${id}/approve`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("برداشت تایید شد");
      setWithdraws(prev => prev.map(w => w.id === id ? { ...w, status: "approved" } : w));
    } catch (e: any) { toast.error(e.message); }
  }

  async function rejectWithdraw(id: string) {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch(`${API_BASE}/v1/admin/withdraws/${id}/reject`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("برداشت رد شد");
      setWithdraws(prev => prev.map(w => w.id === id ? { ...w, status: "rejected" } : w));
    } catch (e: any) { toast.error(e.message); }
  }

  const wTotalPages = Math.ceil(wTotal / 15);
  const tTotalPages = Math.ceil(tTotal / 15);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">مالی و تراکنش‌ها</h1>
        <p className="text-sm text-muted-foreground">مدیریت درآمد پلتفرم و درخواست‌های برداشت</p>
      </div>

      {/* Revenue stats */}
      <div className="grid gap-4 sm:grid-cols-3">
        <StatCard label="کل درآمد پلتفرم"  value={loading ? "..." : toman(stats?.total ?? 0, false)}     unit="تومان" icon={TrendingUp} />
        <StatCard label="درآمد این ماه"     value={loading ? "..." : toman(stats?.thisMonth ?? 0, false)} unit="تومان" icon={Wallet} />
        <StatCard label="تعداد تراکنش‌ها"   value={loading ? "..." : toFa(stats?.txCount ?? 0)}          icon={ArrowDownLeft} />
      </div>

      {/* Financial overview */}
      {(loading || overview) && (
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <h2 className="font-display font-semibold">نمای کلی مالی</h2>
          <div className="mt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <OverviewStat label="کل واریزی کاربران"   value={loading ? "..." : toman(overview?.totalDeposited ?? 0, false)}  unit="تومان" icon={TrendingUp}  loading={loading} />
            <OverviewStat label="کل برداشت‌ها"         value={loading ? "..." : toman(overview?.totalWithdrawn ?? 0, false)}  unit="تومان" icon={ArrowDownLeft} loading={loading} />
            <OverviewStat label="مبلغ بلوکه (امانت)"  value={loading ? "..." : toman(overview?.totalLocked ?? 0, false)}     unit="تومان" icon={Clock}        loading={loading} />
            <OverviewStat label="کیف پول‌های فعال"    value={loading ? "..." : toFa(overview?.activeWallets ?? 0)}           icon={Users}       loading={loading} />
          </div>
        </div>
      )}

      {/* Tabs */}
      <div className="rounded-2xl border bg-card shadow-soft">
        <div className="flex border-b">
          {(["withdraws", "transactions"] as const).map(t => (
            <button key={t} onClick={() => setTab(t)}
              className={cn("flex-1 py-3.5 text-sm font-medium transition-colors",
                tab === t ? "border-b-2 border-primary text-primary" : "text-muted-foreground hover:text-foreground")}>
              {t === "withdraws" ? `درخواست‌های برداشت (${toFa(wTotal)})` : `تراکنش‌ها (${toFa(tTotal)})`}
            </button>
          ))}
        </div>

        {/* Withdraws */}
        {tab === "withdraws" && (
          <div>
            <div className="flex items-center gap-3 border-b px-5 py-3">
              <Filter className="h-4 w-4 text-muted-foreground" />
              <select value={wStatusFilter} onChange={e => setWStatusFilter(e.target.value)}
                className="rounded-lg border bg-background px-3 py-1.5 text-sm outline-none focus:border-primary">
                <option value="">همه وضعیت‌ها</option>
                <option value="pending">در انتظار</option>
                <option value="approved">تایید‌شده</option>
                <option value="rejected">رد‌شده</option>
                <option value="paid">پرداخت‌شده</option>
              </select>
            </div>
            <div className="divide-y">
              {withdraws.length === 0 ? (
                <div className="p-12 text-center text-muted-foreground">
                  <Clock className="mx-auto mb-2 h-8 w-8 opacity-40" />درخواست برداشتی وجود ندارد
                </div>
              ) : withdraws.map(w => (
                <div key={w.id} className="flex flex-wrap items-center gap-4 px-5 py-4">
                  <div className="flex-1 min-w-0">
                    <div className="font-semibold">{toman(w.amount)}</div>
                    <div className="mt-0.5 text-xs text-muted-foreground">
                      {w.userName && <span className="me-2">{w.userName} ·</span>}
                      <span dir="ltr">{w.shebaNumber || w.bankAccount}</span>
                      <span> · {w.createdAt?.slice(0, 10)}</span>
                    </div>
                  </div>
                  <span className={cn("rounded-full px-2 py-0.5 text-xs font-medium shrink-0",
                    WITHDRAW_STATUS_COLORS[w.status] ?? "bg-muted text-muted-foreground")}>
                    {WITHDRAW_STATUS_LABELS[w.status] ?? w.status}
                  </span>
                  {(w.status === "pending" || w.status === "PENDING") && (
                    <div className="flex gap-2 shrink-0">
                      <button onClick={() => approveWithdraw(w.id)}
                        className="flex items-center gap-1 rounded-lg border border-success/40 px-3 py-1.5 text-xs font-medium text-success hover:bg-success/10">
                        <CheckCircle2 className="h-3.5 w-3.5" /> تایید
                      </button>
                      <button onClick={() => rejectWithdraw(w.id)}
                        className="flex items-center gap-1 rounded-lg border border-destructive/40 px-3 py-1.5 text-xs font-medium text-destructive hover:bg-destructive/10">
                        <XCircle className="h-3.5 w-3.5" /> رد
                      </button>
                    </div>
                  )}
                </div>
              ))}
            </div>
            <Pagination page={wPage} totalPages={wTotalPages} onPrev={() => setWPage(p => p - 1)} onNext={() => setWPage(p => p + 1)} />
          </div>
        )}

        {/* Transactions */}
        {tab === "transactions" && (
          <div>
            <div className="divide-y">
              {txs.length === 0 ? (
                <div className="p-12 text-center text-muted-foreground">
                  <Wallet className="mx-auto mb-2 h-8 w-8 opacity-40" />تراکنشی وجود ندارد
                </div>
              ) : txs.map(t => (
                <div key={t.id} className="flex items-center gap-4 px-5 py-3">
                  <div className="flex-1 min-w-0">
                    <div className="text-sm font-medium">{TX_TYPE_LABELS[t.type] ?? t.description ?? t.type}</div>
                    <div className="text-xs text-muted-foreground">{t.createdAt?.slice(0, 10)}</div>
                  </div>
                  <span className={cn("rounded-full px-2 py-0.5 text-xs font-medium shrink-0",
                    t.type === "PLATFORM_FEE" ? "bg-success/10 text-success" :
                    t.type === "REFUND" ? "bg-destructive/10 text-destructive" :
                    "bg-muted text-muted-foreground")}>
                    {t.type === "PLATFORM_FEE" ? "+" : ""}{toman(t.amount)}
                  </span>
                </div>
              ))}
            </div>
            <Pagination page={tPage} totalPages={tTotalPages} onPrev={() => setTPage(p => p - 1)} onNext={() => setTPage(p => p + 1)} />
          </div>
        )}
      </div>
    </div>
  );
}

function StatCard({ label, value, unit, icon: Icon }: { label: string; value: string; unit?: string; icon: any }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-start justify-between">
        <span className="text-sm text-muted-foreground">{label}</span>
        <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
      </div>
      <div className="mt-3 font-display text-xl font-bold tabular-nums">
        {value}{unit && <span className="ms-1 text-sm font-medium text-muted-foreground">{unit}</span>}
      </div>
    </div>
  );
}

function OverviewStat({ label, value, unit, icon: Icon, loading }: { label: string; value: string; unit?: string; icon: any; loading: boolean }) {
  return (
    <div className="rounded-xl border bg-muted/20 p-4">
      <div className="flex items-center gap-2 text-sm text-muted-foreground">
        <Icon className="h-4 w-4" />{label}
      </div>
      <div className="mt-2 font-display font-bold tabular-nums">
        {loading ? <div className="h-5 w-24 animate-pulse rounded bg-muted" /> : (
          <>{value}{unit && <span className="ms-1 text-sm font-normal text-muted-foreground">{unit}</span>}</>
        )}
      </div>
    </div>
  );
}

function Pagination({ page, totalPages, onPrev, onNext }: { page: number; totalPages: number; onPrev: () => void; onNext: () => void }) {
  if (totalPages <= 1) return null;
  return (
    <div className="flex items-center justify-between border-t px-5 py-3">
      <span className="text-xs text-muted-foreground">صفحه {toFa(page)} از {toFa(totalPages)}</span>
      <div className="flex gap-2">
        <button disabled={page === 1} onClick={onPrev} className="rounded border px-3 py-1 text-xs disabled:opacity-40 hover:bg-accent">قبلی</button>
        <button disabled={page >= totalPages} onClick={onNext} className="rounded border px-3 py-1 text-xs disabled:opacity-40 hover:bg-accent">بعدی</button>
      </div>
    </div>
  );
}
