import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { ArrowRight, MapPin, Clock, Wallet, MessagesSquare, CheckCircle2, Star, ShieldCheck, Users, Play, XCircle } from "lucide-react";
import { StatusBadge } from "@/components/status-badge";
import { TaskTimeline } from "@/components/task-timeline";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";
import { useRole } from "@/components/role-context";
import type { ApiTask } from "@/lib/types";

export const Route = createFileRoute("/app/tasks/$id")({
  head: ({ params }) => ({ meta: [{ title: `درخواست — تسک‌بریج` }] }),
  component: TaskDetails,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

function TaskDetails() {
  const { id } = Route.useParams();
  const { role } = useRole();
  const navigate = useNavigate();
  const [task, setTask] = useState<ApiTask | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [acting, setActing] = useState(false);

  useEffect(() => {
    const token = getToken();
    const h: HeadersInit = token ? { Authorization: `Bearer ${token}` } : {};
    fetch(`${API_BASE}/v1/tasks/${id}`, { headers: h })
      .then(r => { if (!r.ok) throw new Error("not found"); return r.json(); })
      .then(d => setTask(d.task ?? d))
      .catch(() => setError("درخواست پیدا نشد"))
      .finally(() => setLoading(false));
  }, [id]);

  async function doAction(action: string, successMsg: string, newStatus?: string) {
    const token = getToken();
    if (!token) return;
    setActing(true);
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${id}/${action}`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success(successMsg);
      if (newStatus) setTask(prev => prev ? { ...prev, status: newStatus } : prev);
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setActing(false);
    }
  }

  async function applyToTask() {
    const token = getToken();
    if (!token) return;
    setActing(true);
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${id}/applications`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ proposalMessage: "درخواست همکاری", eta: "۲۴ ساعت" }),
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("درخواست همکاری ثبت شد");
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setActing(false);
    }
  }

  if (loading) return <LoadingSkeleton />;
  if (error || !task) return (
    <div className="p-10 text-center">
      <p className="text-muted-foreground">{error ?? "درخواست پیدا نشد"}</p>
      <Link to="/app/tasks" className="mt-3 inline-block text-sm text-primary hover:underline">بازگشت به فهرست</Link>
    </div>
  );

  const fee = Math.round(task.budget * 0.08);
  const isOpen = task.status === "OPEN";
  const isAssigned = task.status === "ASSIGNED";
  const isInProgress = task.status === "IN_PROGRESS";
  const isCompleted = task.status === "COMPLETED";
  const isWaiting = task.status === "WAITING_FOR_VERIFICATION";

  return (
    <div className="space-y-6">
      <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowRight className="h-4 w-4" /> بازگشت به فهرست درخواست‌ها
      </Link>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6 lg:col-span-2">
          <div className="rounded-2xl border bg-card p-6 shadow-soft">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="font-mono">{task.id.slice(0, 8)}</span>
                  {task.createdAt && <><span>·</span><span>ثبت {task.createdAt.slice(0, 10)}</span></>}
                  {task.requester && <><span>·</span><span>توسط {task.requester.fullName}</span></>}
                </div>
                <h1 className="mt-2 text-2xl font-bold tracking-tight">{task.title}</h1>
              </div>
              <StatusBadge status={task.status} />
            </div>
            <div className="mt-4 flex flex-wrap gap-x-5 gap-y-2 text-sm text-muted-foreground">
              {task.city && <span className="inline-flex items-center gap-1.5"><MapPin className="h-4 w-4" />{task.city.title}</span>}
              {task.deadline && <span className="inline-flex items-center gap-1.5"><Clock className="h-4 w-4" />مهلت: {task.deadline.slice(0, 10)}</span>}
              <span className="inline-flex items-center gap-1.5"><Wallet className="h-4 w-4 text-primary" /><span className="font-semibold text-foreground">{toman(task.budget)}</span></span>
              {task.category && <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2 py-0.5 text-xs">{task.category.title}</span>}
              <span className="inline-flex items-center gap-1.5"><Users className="h-4 w-4" />{toFa(task.applicantCount ?? 0)} متقاضی</span>
            </div>
            {task.description && <p className="mt-5 text-sm leading-relaxed text-foreground/90">{task.description}</p>}
          </div>

          <div className="rounded-2xl border bg-card p-6 shadow-soft">
            <h2 className="font-display font-semibold">مراحل انجام</h2>
            <div className="mt-6 overflow-x-auto">
              <TaskTimeline status={task.status} />
            </div>
          </div>

          {/* Action buttons per role */}
          <div className="flex flex-wrap gap-3">
            {role === "agent" && isOpen && (
              <button onClick={applyToTask} disabled={acting}
                className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60">
                <Star className="h-4 w-4" /> {acting ? "در حال ارسال..." : "ارسال درخواست همکاری"}
              </button>
            )}
            {role === "agent" && isAssigned && (
              <button onClick={() => doAction("start", "کار شروع شد", "IN_PROGRESS")} disabled={acting}
                className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60">
                <Play className="h-4 w-4" /> {acting ? "..." : "شروع کار"}
              </button>
            )}
            {role === "agent" && isInProgress && (
              <button onClick={() => doAction("complete", "کار تکمیل شد", "COMPLETED")} disabled={acting}
                className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60">
                <CheckCircle2 className="h-4 w-4" /> {acting ? "..." : "اعلام اتمام کار"}
              </button>
            )}
            {role === "requester" && (isCompleted || isWaiting) && (
              <button onClick={() => doAction("verify", "تایید شد — پرداخت انجام می‌شود", "VERIFIED")} disabled={acting}
                className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60">
                <ShieldCheck className="h-4 w-4" /> {acting ? "..." : "تایید و پرداخت"}
              </button>
            )}
            {role === "requester" && isOpen && (
              <Link to="/app/tasks/$id/applications" params={{ id }}
                className="inline-flex items-center gap-2 rounded-lg border px-5 py-2.5 text-sm font-semibold hover:bg-accent">
                <Users className="h-4 w-4" /> بررسی متقاضیان ({toFa(task.applicantCount ?? 0)})
              </Link>
            )}
            {(role === "requester" || role === "admin") && (isOpen || isAssigned) && (
              <button onClick={() => doAction("cancel", "درخواست لغو شد", "CANCELLED")} disabled={acting}
                className="inline-flex items-center gap-2 rounded-lg border border-destructive/40 px-5 py-2.5 text-sm font-semibold text-destructive hover:bg-destructive/10 disabled:opacity-60">
                <XCircle className="h-4 w-4" /> {acting ? "..." : "لغو درخواست"}
              </button>
            )}
            <Link to="/app/chat" className="inline-flex items-center gap-2 rounded-lg border px-4 py-2.5 text-sm font-medium hover:bg-accent">
              <MessagesSquare className="h-4 w-4" /> گفت‌وگو
            </Link>
          </div>
        </div>

        <aside className="space-y-4">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-success" />
              <h3 className="font-semibold">جزئیات پرداخت</h3>
            </div>
            <dl className="mt-4 space-y-2 text-sm">
              <Row k="مبلغ درخواست" v={toman(task.budget)} />
              <Row k="کارمزد سامانه" v={toman(fee)} />
              <Row k="بلوکه شده در امانت" v={toman(task.budget + fee)} bold />
            </dl>
            <div className="mt-4 rounded-lg bg-success/10 p-3 text-xs text-success-foreground">
              <CheckCircle2 className="me-1 inline h-3.5 w-3.5" /> پرداخت تنها پس از تایید نهایی شما آزاد می‌شود.
            </div>
          </div>

          {task.assignedAgent && (
            <div className="rounded-2xl border bg-card p-5 shadow-soft">
              <h3 className="font-semibold">مجری تخصیص‌یافته</h3>
              <div className="mt-3 flex items-center gap-3">
                <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">
                  {task.assignedAgent.fullName.slice(0, 2)}
                </span>
                <div>
                  <div className="font-medium">{task.assignedAgent.fullName}</div>
                  <div className="text-xs text-muted-foreground">مجری</div>
                </div>
              </div>
            </div>
          )}

          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <h3 className="font-semibold">نیاز به کمک دارید؟</h3>
            <p className="mt-1 text-sm text-muted-foreground">تیم پشتیبانی ۲۴/۷ می‌تواند پرداخت را متوقف و میانجی‌گری کند.</p>
            <button className="mt-3 w-full rounded-lg border py-2 text-sm font-medium hover:bg-accent">ثبت اعتراض</button>
          </div>
        </aside>
      </div>
    </div>
  );
}

function Row({ k, v, bold }: { k: string; v: string; bold?: boolean }) {
  return (
    <div className={`flex items-center justify-between ${bold ? "font-semibold" : ""}`}>
      <dt className="text-muted-foreground">{k}</dt>
      <dd className="tabular-nums">{v}</dd>
    </div>
  );
}

function LoadingSkeleton() {
  return (
    <div className="space-y-6 animate-pulse">
      <div className="h-4 w-32 rounded bg-muted" />
      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-4 lg:col-span-2">
          <div className="rounded-2xl border bg-card p-6 space-y-3">
            <div className="h-4 w-1/4 rounded bg-muted" />
            <div className="h-7 w-2/3 rounded bg-muted" />
            <div className="h-4 w-full rounded bg-muted" />
            <div className="h-4 w-3/4 rounded bg-muted" />
          </div>
        </div>
        <div className="rounded-2xl border bg-card p-5 h-48" />
      </div>
    </div>
  );
}
