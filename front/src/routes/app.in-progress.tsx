import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Activity, MapPin, Clock, Wallet, MessagesSquare } from "lucide-react";
import { StatusBadge } from "@/components/status-badge";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/app/in-progress")({
  head: () => ({ meta: [{ title: "در حال انجام — تسک‌بریج" }] }),
  component: InProgressPage,
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
  applicationCount: number;
};

function InProgressPage() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    // Includes COMPLETED/WAITING_FOR_VERIFICATION/VERIFIED so the agent keeps the
    // task until they confirm receipt of payment (VERIFIED → PAID).
    fetch(`${API_BASE}/v1/tasks?status=ASSIGNED,IN_PROGRESS,COMPLETED,WAITING_FOR_VERIFICATION,VERIFIED&mine=true`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(r => r.json())
      .then(d => setTasks(d.items ?? []))
      .catch(() => setError("خطا در بارگذاری درخواست‌ها"))
      .finally(() => setLoading(false));
  }, []);

  if (loading) return <LoadingState />;
  if (error) return <ErrorState msg={error} />;

  return (
    <div className="space-y-6">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">درخواست‌های در حال انجام</h1>
          <p className="text-sm text-muted-foreground">{toFa(tasks.length)} درخواست فعال</p>
        </div>
      </div>

      {tasks.length === 0 ? (
        <EmptyState />
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {tasks.map(t => <TaskCard key={t.id} task={t} />)}
        </div>
      )}
    </div>
  );
}

function TaskCard({ task }: { task: Task }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft space-y-3">
      <div className="flex items-start justify-between gap-2">
        <h3 className="font-semibold line-clamp-2">{task.title}</h3>
        <StatusBadge status={task.status} />
      </div>
      <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
        <span className="flex items-center gap-1"><MapPin className="h-3 w-3" />{task.city}</span>
        <span className="flex items-center gap-1"><Clock className="h-3 w-3" />{task.deadline?.slice(0, 10)}</span>
        <span className="flex items-center gap-1 text-primary font-medium"><Wallet className="h-3 w-3" />{toman(task.budget)}</span>
      </div>
      <div className="flex gap-2 pt-1">
        <Link to="/app/tasks/$id" params={{ id: task.id }}
          className="flex-1 rounded-lg border py-1.5 text-center text-xs font-medium hover:border-primary hover:text-primary">
          مشاهده جزئیات
        </Link>
        <Link to="/app/chat" className="flex items-center gap-1 rounded-lg border px-3 py-1.5 text-xs font-medium hover:border-primary hover:text-primary">
          <MessagesSquare className="h-3.5 w-3.5" /> گفت‌وگو
        </Link>
      </div>
    </div>
  );
}

function LoadingState() {
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

function EmptyState() {
  return (
    <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
      <Activity className="mx-auto h-12 w-12 text-muted-foreground/40" />
      <h3 className="mt-4 font-semibold">درخواست فعالی ندارید</h3>
      <p className="mt-1 text-sm text-muted-foreground">درخواست‌هایی که در حال انجام هستند اینجا نمایش داده می‌شوند.</p>
      <Link to="/app/tasks/new" className="mt-4 inline-block rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white">
        ثبت درخواست جدید
      </Link>
    </div>
  );
}

function ErrorState({ msg }: { msg: string }) {
  return (
    <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center">
      <p className="text-sm text-destructive">{msg}</p>
    </div>
  );
}
