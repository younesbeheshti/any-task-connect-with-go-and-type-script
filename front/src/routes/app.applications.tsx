import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { FileText, MapPin, Clock, Wallet } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";

export const Route = createFileRoute("/app/applications")({
  head: () => ({ meta: [{ title: "درخواست‌های من — تسک‌بریج" }] }),
  component: MyApplicationsPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Application = {
  id: string;
  status: string;
  proposalMessage: string;
  createdAt: string;
  task: { id: string; title: string; budget: number; city: string; deadline: string; status: string };
};

const STATUS_LABELS: Record<string, string> = {
  PENDING: "در انتظار بررسی",
  ACCEPTED: "پذیرفته‌شده",
  REJECTED: "رد‌شده",
  WITHDRAWN: "پس‌گرفته‌شده",
};
const STATUS_COLORS: Record<string, string> = {
  PENDING: "bg-warning/10 text-warning",
  ACCEPTED: "bg-success/10 text-success",
  REJECTED: "bg-destructive/10 text-destructive",
  WITHDRAWN: "bg-muted text-muted-foreground",
};

function MyApplicationsPage() {
  const [apps, setApps] = useState<Application[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<string>("");

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    fetch(`${API_BASE}/v1/me/applications`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => setApps(d.applications ?? d.items ?? []))
      .catch(() => setError("خطا در بارگذاری"))
      .finally(() => setLoading(false));
  }, []);

  async function withdraw(id: string) {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch(`${API_BASE}/v1/applications/${id}/withdraw`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("درخواست پس گرفته شد");
      setApps(prev => prev.map(a => a.id === id ? { ...a, status: "WITHDRAWN" } : a));
    } catch (e: any) { toast.error(e.message); }
  }

  const filtered = filter ? apps.filter(a => a.status === filter) : apps;

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">درخواست‌های من</h1>
          <p className="text-sm text-muted-foreground">{toFa(apps.length)} درخواست ارسال‌شده</p>
        </div>
        <div className="flex gap-2">
          {["", "PENDING", "ACCEPTED", "REJECTED"].map(s => (
            <button key={s} onClick={() => setFilter(s)}
              className={`rounded-lg border px-3 py-1.5 text-xs font-medium transition-colors ${filter === s ? "border-primary bg-primary/10 text-primary" : "hover:border-primary/40"}`}>
              {s ? STATUS_LABELS[s] : "همه"}
            </button>
          ))}
        </div>
      </div>

      {filtered.length === 0 ? (
        <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
          <FileText className="mx-auto h-12 w-12 text-muted-foreground/40" />
          <h3 className="mt-4 font-semibold">درخواستی ثبت نشده</h3>
          <p className="mt-1 text-sm text-muted-foreground">روی فرصت‌های جدید اپلای کنید.</p>
          <Link to="/app/tasks" className="mt-4 inline-block rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white">
            مشاهده فرصت‌ها
          </Link>
        </div>
      ) : (
        <div className="space-y-3">
          {filtered.map(a => (
            <div key={a.id} className="rounded-2xl border bg-card p-5 shadow-soft">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="flex-1 min-w-0">
                  <Link to="/app/tasks/$id" params={{ id: a.task?.id ?? a.id }}
                    className="font-semibold hover:text-primary line-clamp-1">
                    {a.task?.title ?? "بدون عنوان"}
                  </Link>
                  <div className="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
                    <span className="flex items-center gap-1"><MapPin className="h-3 w-3" />{a.task?.city}</span>
                    <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{a.task?.deadline?.slice(0, 10)}</span>
                    <span className="flex items-center gap-1 font-medium text-primary"><Wallet className="h-3 w-3" />{toman(a.task?.budget ?? 0)}</span>
                  </div>
                  {a.proposalMessage && (
                    <p className="mt-2 text-sm text-muted-foreground line-clamp-2">{a.proposalMessage}</p>
                  )}
                </div>
                <span className={`rounded-full px-3 py-1 text-xs font-medium ${STATUS_COLORS[a.status] ?? "bg-muted"}`}>
                  {STATUS_LABELS[a.status] ?? a.status}
                </span>
              </div>
              <div className="mt-3 flex items-center justify-between">
                <span className="text-xs text-muted-foreground">{a.createdAt?.slice(0, 10)}</span>
                {a.status === "PENDING" && (
                  <button onClick={() => withdraw(a.id)}
                    className="text-xs text-muted-foreground hover:text-destructive underline">
                    پس‌گرفتن درخواست
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map(i => (
        <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse">
          <div className="space-y-2">
            <div className="h-5 w-2/3 rounded bg-muted" />
            <div className="h-4 w-1/3 rounded bg-muted" />
          </div>
        </div>
      ))}
    </div>
  );
}

function ErrBox({ msg }: { msg: string }) {
  return <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center"><p className="text-sm text-destructive">{msg}</p></div>;
}
