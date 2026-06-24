import { createFileRoute, Link } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ArrowRight, Star, MapPin, Clock, CheckCircle2, X, MessagesSquare, ShieldCheck, Briefcase } from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";

export const Route = createFileRoute("/app/tasks/$id/applications")({
  head: ({ params }) => ({ meta: [{ title: `متقاضیان ${params.id} — تسک‌بریج` }] }),
  component: ApplicationsReviewPage,
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
  eta: string;
  proposedPrice: number | null;
  createdAt: string;
  agent: {
    id: string;
    fullName: string;
    rating: number;
    completedCount: number;
    avatarUrl: string | null;
    isVerified: boolean;
  };
};

type Task = {
  id: string;
  title: string;
  budget: number;
  status: string;
};

const STATUS_COLORS: Record<string, string> = {
  PENDING: "bg-warning/10 text-warning",
  ACCEPTED: "bg-success/10 text-success",
  REJECTED: "bg-destructive/10 text-destructive",
  WITHDRAWN: "bg-muted text-muted-foreground",
};
const STATUS_LABELS: Record<string, string> = {
  PENDING: "در انتظار", ACCEPTED: "پذیرفته‌شده", REJECTED: "رد‌شده", WITHDRAWN: "پس‌گرفته‌شده",
};

function ApplicationsReviewPage() {
  const { id } = Route.useParams();
  const [task, setTask] = useState<Task | null>(null);
  const [apps, setApps] = useState<Application[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [accepting, setAccepting] = useState<string | null>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    const headers = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/tasks/${id}`, { headers }).then(r => r.json()),
      fetch(`${API_BASE}/v1/tasks/${id}/applications`, { headers }).then(r => r.json()),
    ])
      .then(([t, a]) => {
        setTask(t.task ?? t);
        setApps(a.applications ?? a.items ?? []);
      })
      .catch(() => setError("خطا در بارگذاری"))
      .finally(() => setLoading(false));
  }, [id]);

  async function accept(appId: string) {
    const token = getToken();
    if (!token) return;
    setAccepting(appId);
    try {
      const res = await fetch(`${API_BASE}/v1/applications/${appId}/accept`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("مجری پذیرفته شد");
      setApps(prev => prev.map(a =>
        a.id === appId ? { ...a, status: "ACCEPTED" } :
        a.status === "PENDING" ? { ...a, status: "REJECTED" } : a
      ));
    } catch (e: any) {
      toast.error(e.message ?? "خطا در پذیرش مجری");
    } finally {
      setAccepting(null);
    }
  }

  async function reject(appId: string) {
    const token = getToken();
    if (!token) return;
    try {
      const res = await fetch(`${API_BASE}/v1/applications/${appId}/reject`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("درخواست رد شد");
      setApps(prev => prev.map(a => a.id === appId ? { ...a, status: "REJECTED" } : a));
    } catch (e: any) {
      toast.error(e.message ?? "خطا");
    }
  }

  const pending = apps.filter(a => a.status === "PENDING");
  const others = apps.filter(a => a.status !== "PENDING");

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Link to="/app/tasks/$id" params={{ id }}
          className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
          <ArrowRight className="h-4 w-4" /> بازگشت به درخواست
        </Link>
      </div>

      <div>
        <h1 className="text-2xl font-bold tracking-tight">بررسی متقاضیان</h1>
        <p className="text-sm text-muted-foreground">
          {task?.title} · {toFa(apps.length)} متقاضی
        </p>
      </div>

      {pending.length === 0 && others.length === 0 ? (
        <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
          <Briefcase className="mx-auto h-12 w-12 text-muted-foreground/40" />
          <h3 className="mt-4 font-semibold">هنوز متقاضی‌ای ندارید</h3>
          <p className="mt-1 text-sm text-muted-foreground">وقتی مجری‌ها برای درخواست شما اپلای کنند، اینجا نمایش داده می‌شوند.</p>
        </div>
      ) : (
        <div className="space-y-8">
          {pending.length > 0 && (
            <section>
              <h2 className="mb-4 font-display font-semibold">در انتظار بررسی ({toFa(pending.length)})</h2>
              <div className="space-y-4">
                {pending.map(a => <ApplicantCard key={a.id} app={a} onAccept={accept} onReject={reject} accepting={accepting} taskBudget={task?.budget ?? 0} />)}
              </div>
            </section>
          )}
          {others.length > 0 && (
            <section>
              <h2 className="mb-4 font-display font-semibold text-muted-foreground">بررسی‌شده ({toFa(others.length)})</h2>
              <div className="space-y-4">
                {others.map(a => <ApplicantCard key={a.id} app={a} onAccept={accept} onReject={reject} accepting={accepting} taskBudget={task?.budget ?? 0} />)}
              </div>
            </section>
          )}
        </div>
      )}
    </div>
  );
}

function ApplicantCard({ app, onAccept, onReject, accepting, taskBudget }: {
  app: Application;
  onAccept: (id: string) => void;
  onReject: (id: string) => void;
  accepting: string | null;
  taskBudget: number;
}) {
  const agent = app.agent ?? {};
  const isPending = app.status === "PENDING";
  return (
    <div className={`rounded-2xl border bg-card p-5 shadow-soft ${app.status === "ACCEPTED" ? "border-success/30 bg-success/5" : ""}`}>
      <div className="flex flex-wrap items-start gap-4">
        <div className="flex h-14 w-14 items-center justify-center rounded-full bg-gradient-to-br from-primary to-secondary font-bold text-lg text-white">
          {(agent.fullName ?? "؟").slice(0, 2)}
        </div>
        <div className="flex-1 min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <span className="font-semibold">{agent.fullName ?? "مجری"}</span>
            {agent.isVerified && (
              <span className="inline-flex items-center gap-1 rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                <ShieldCheck className="h-3 w-3" /> تایید‌شده
              </span>
            )}
            <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${STATUS_COLORS[app.status]}`}>
              {STATUS_LABELS[app.status]}
            </span>
          </div>
          <div className="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
            <span className="flex items-center gap-1">
              <Star className="h-3 w-3 fill-warning text-warning" />
              {toFa((agent.rating ?? 0).toFixed(1))} امتیاز
            </span>
            <span className="flex items-center gap-1">
              <CheckCircle2 className="h-3 w-3 text-success" />
              {toFa(agent.completedCount ?? 0)} کار تکمیل‌شده
            </span>
            {app.eta && <span className="flex items-center gap-1"><Clock className="h-3 w-3" /> تخمین: {app.eta}</span>}
          </div>
          {app.proposalMessage && (
            <p className="mt-3 rounded-lg bg-muted/50 p-3 text-sm leading-relaxed">{app.proposalMessage}</p>
          )}
        </div>
        <div className="flex flex-col items-end gap-2">
          {app.proposedPrice ? (
            <span className="rounded-lg bg-primary/10 px-3 py-1.5 text-sm font-semibold text-primary">
              {toman(app.proposedPrice)}
            </span>
          ) : (
            <span className="rounded-lg bg-muted px-3 py-1.5 text-sm text-muted-foreground">
              {toman(taskBudget)}
            </span>
          )}
          <span className="text-xs text-muted-foreground">{app.createdAt?.slice(0, 10)}</span>
        </div>
      </div>

      {isPending && (
        <div className="mt-4 flex gap-2 border-t pt-4">
          <button
            onClick={() => onAccept(app.id)}
            disabled={accepting === app.id}
            className="flex flex-1 items-center justify-center gap-1.5 rounded-lg gradient-brand py-2 text-sm font-semibold text-white shadow-soft disabled:opacity-60"
          >
            {accepting === app.id ? "در حال پردازش..." : <><CheckCircle2 className="h-4 w-4" /> پذیرفتن</>}
          </button>
          <Link to="/app/chat" className="flex items-center gap-1.5 rounded-lg border px-4 py-2 text-sm font-medium hover:border-primary hover:text-primary">
            <MessagesSquare className="h-4 w-4" /> گفت‌وگو
          </Link>
          <button
            onClick={() => onReject(app.id)}
            className="flex items-center gap-1.5 rounded-lg border border-destructive/40 px-4 py-2 text-sm font-medium text-destructive hover:bg-destructive/10"
          >
            <X className="h-4 w-4" /> رد
          </button>
        </div>
      )}
    </div>
  );
}

function Skeleton() {
  return (
    <div className="space-y-4">
      {[1, 2, 3].map(i => (
        <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse">
          <div className="flex gap-4">
            <div className="h-14 w-14 rounded-full bg-muted" />
            <div className="flex-1 space-y-2">
              <div className="h-5 w-1/3 rounded bg-muted" />
              <div className="h-4 w-1/2 rounded bg-muted" />
              <div className="h-12 rounded bg-muted" />
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}

function ErrBox({ msg }: { msg: string }) {
  return <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center"><p className="text-sm text-destructive">{msg}</p></div>;
}
