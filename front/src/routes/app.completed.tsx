import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { CheckCircle2, MapPin, Clock, Wallet, Star } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { useRole } from "@/components/role-context";

export const Route = createFileRoute("/app/completed")({
  head: () => ({ meta: [{ title: "تکمیل‌شده‌ها — تسک‌بریج" }] }),
  component: CompletedPage,
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
  assignedAgentId?: string;
};

function CompletedPage() {
  const { role } = useRole();
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    fetch(`${API_BASE}/v1/tasks?status=VERIFIED,PAID&mine=true`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(r => r.json())
      .then(d => setTasks(d.items ?? []))
      .catch(() => setError("خطا در بارگذاری"))
      .finally(() => setLoading(false));
  }, []);

  const title = role === "agent" ? "درخواست‌های تکمیل‌شده من" : "درخواست‌های تکمیل‌شده";

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
        <p className="text-sm text-muted-foreground">{toFa(tasks.length)} درخواست تکمیل‌شده</p>
      </div>

      {tasks.length === 0 ? (
        <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
          <CheckCircle2 className="mx-auto h-12 w-12 text-muted-foreground/40" />
          <h3 className="mt-4 font-semibold">درخواست تکمیل‌شده‌ای ندارید</h3>
          <p className="mt-1 text-sm text-muted-foreground">درخواست‌های تایید‌شده اینجا نمایش داده می‌شوند.</p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {tasks.map(t => (
            <div key={t.id} className="rounded-2xl border bg-card p-5 shadow-soft space-y-3">
              <div className="flex items-start justify-between gap-2">
                <h3 className="font-semibold line-clamp-2">{t.title}</h3>
                <span className="inline-flex items-center gap-1 rounded-full bg-success/10 px-2 py-0.5 text-xs font-medium text-success">
                  <CheckCircle2 className="h-3 w-3" /> تکمیل‌شده
                </span>
              </div>
              <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
                <span className="flex items-center gap-1"><MapPin className="h-3 w-3" />{t.city}</span>
                <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{t.deadline?.slice(0, 10)}</span>
                <span className="flex items-center gap-1 text-primary font-medium"><Wallet className="h-3 w-3" />{toman(t.budget)}</span>
              </div>
              <div className="flex gap-2 pt-1">
                <Link to="/app/tasks/$id" params={{ id: t.id }}
                  className="flex-1 rounded-lg border py-1.5 text-center text-xs font-medium hover:border-primary hover:text-primary">
                  مشاهده جزئیات
                </Link>
                <Link to="/app/tasks/$id" params={{ id: t.id }}
                  className="flex items-center gap-1 rounded-lg border px-3 py-1.5 text-xs font-medium hover:border-primary hover:text-primary">
                  <Star className="h-3.5 w-3.5" /> ثبت نظر
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
