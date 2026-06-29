import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Briefcase, MapPin, Clock, Wallet, MessagesSquare, Play } from "lucide-react";
import { StatusBadge } from "@/components/status-badge";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";

export const Route = createFileRoute("/app/accepted")({
  head: () => ({ meta: [{ title: "درخواست‌های پذیرفته — تسک‌بریج" }] }),
  component: AcceptedPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Task = {
  id: string; title: string; status: string; budget: number;
  category: string; city: string; deadline: string;
  requesterId: string;
};

function AcceptedPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    fetch(`${API_BASE}/v1/tasks?status=ASSIGNED,IN_PROGRESS&assigned=true`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(r => r.json())
      .then(d => setTasks(d.items ?? []))
      .catch(() => setError("خطا در بارگذاری"))
      .finally(() => setLoading(false));
  }, []);

  async function handleStart(taskId: string) {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${taskId}/start`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) {
        const d = await res.json();
        throw new Error(d?.error?.message ?? "خطا");
      }
      toast.success("شروع کار ثبت شد");
      setTasks(prev => prev.map(t => t.id === taskId ? { ...t, status: "IN_PROGRESS" } : t));
    } catch (e: any) {
      toast.error(e.message ?? "خطا در شروع کار");
    }
  }

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">درخواست‌های پذیرفته‌شده</h1>
        <p className="text-sm text-muted-foreground">{toFa(tasks.length)} درخواست منتظر اقدام شما</p>
      </div>

      {tasks.length === 0 ? (
        <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
          <Briefcase className="mx-auto h-12 w-12 text-muted-foreground/40" />
          <h3 className="mt-4 font-semibold">درخواست پذیرفته‌شده‌ای ندارید</h3>
          <p className="mt-1 text-sm text-muted-foreground">پس از قبول شدن درخواستی، اینجا نمایش داده می‌شود.</p>
          <Link to="/app/tasks" className="mt-4 inline-block rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white">
            مشاهده فرصت‌ها
          </Link>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {tasks.map(t => (
            <div key={t.id} className="rounded-2xl border bg-card p-5 shadow-soft space-y-3">
              <div className="flex items-start justify-between gap-2">
                <h3 className="font-semibold line-clamp-2">{t.title}</h3>
                <StatusBadge status={t.status} />
              </div>
              <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
                <span className="flex items-center gap-1"><MapPin className="h-3 w-3" />{t.city}</span>
                <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{t.deadline?.slice(0, 10)}</span>
                <span className="flex items-center gap-1 text-primary font-medium"><Wallet className="h-3 w-3" />{toman(t.budget)}</span>
              </div>
              <div className="flex gap-2 pt-1">
                {t.status === "ASSIGNED" && (
                  <button onClick={() => handleStart(t.id)}
                    className="flex items-center gap-1 rounded-lg gradient-brand px-3 py-1.5 text-xs font-semibold text-white">
                    <Play className="h-3.5 w-3.5" /> شروع کار
                  </button>
                )}
                <Link to="/app/tasks/$id" params={{ id: t.id }}
                  className="flex-1 rounded-lg border py-1.5 text-center text-xs font-medium hover:border-primary hover:text-primary">
                  جزئیات
                </Link>
                <Link to="/app/chat" className="flex items-center gap-1 rounded-lg border px-3 py-1.5 text-xs font-medium hover:border-primary">
                  <MessagesSquare className="h-3.5 w-3.5" />
                </Link>
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
    <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
      {[1, 2, 3].map(i => (
        <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft">
          <div className="animate-pulse space-y-3">
            <div className="h-5 w-3/4 rounded bg-muted" />
            <div className="h-4 w-1/2 rounded bg-muted" />
            <div className="h-8 rounded bg-muted" />
          </div>
        </div>
      ))}
    </div>
  );
}

function ErrBox({ msg }: { msg: string }) {
  return <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center"><p className="text-sm text-destructive">{msg}</p></div>;
}
