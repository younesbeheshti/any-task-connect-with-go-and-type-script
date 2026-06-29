import { useState } from "react";
import { Link } from "@tanstack/react-router";
import { Clock, MapPin, Users, Wallet, Star, Check, MessageSquare } from "lucide-react";
import { toast } from "sonner";
import { StatusBadge } from "./status-badge";
import type { ApiTask } from "@/lib/types";
import { toman, toFa } from "@/lib/fa";

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken(): string | null {
  try {
    return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null;
  } catch {
    return null;
  }
}

export function TaskCard({
  task,
  showApply = false,
  alreadyApplied = false,
  unread = 0,
}: {
  task: ApiTask;
  showApply?: boolean;
  /** True when the current agent has already applied to this task. */
  alreadyApplied?: boolean;
  /** Unread chat messages for this task (rendered as a badge). */
  unread?: number;
}) {
  const postedDate = task.createdAt?.slice(0, 10) ?? "";
  const isOpen = task.status === "awaiting_applicants" || task.status === "posted";
  const [applying, setApplying] = useState(false);
  const [applied, setApplied] = useState(alreadyApplied);

  async function apply(e: React.MouseEvent) {
    // The whole card is a link; don't navigate when applying from it.
    e.preventDefault();
    e.stopPropagation();
    const token = getToken();
    if (!token) {
      toast.error("لطفاً ابتدا وارد شوید");
      return;
    }
    setApplying(true);
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${task.id}/applications`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ proposalMessage: "درخواست همکاری", eta: "۲۴ ساعت" }),
      });
      const data = await res.json().catch(() => ({}));
      if (!res.ok) throw new Error(data?.error?.message ?? "خطا در ثبت درخواست");
      setApplied(true);
      toast.success("درخواست همکاری ثبت شد");
    } catch (err: unknown) {
      toast.error(err instanceof Error ? err.message : "خطا");
    } finally {
      setApplying(false);
    }
  }

  return (
    <Link
      to="/app/tasks/$id"
      params={{ id: task.id }}
      className="group block rounded-2xl border bg-card p-5 shadow-soft transition-all hover:-translate-y-0.5 hover:shadow-elevated"
    >
      <div className="flex items-start justify-between gap-3">
        <div className="min-w-0">
          <div className="flex items-center gap-2 text-xs text-muted-foreground">
            <span className="font-mono">{task.id}</span>
            <span>·</span>
            <span>{postedDate}</span>
          </div>
          <h3 className="mt-1.5 line-clamp-1 text-base font-semibold tracking-tight text-foreground group-hover:text-primary">
            {task.title}
          </h3>
          <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{task.description}</p>
        </div>
        <div className="flex shrink-0 flex-col items-end gap-1.5">
          <StatusBadge status={task.status} />
          {unread > 0 && (
            <span className="inline-flex items-center gap-1 rounded-full bg-primary px-2 py-0.5 text-[10px] font-semibold text-primary-foreground">
              <MessageSquare className="h-3 w-3" /> {toFa(unread)} پیام جدید
            </span>
          )}
        </div>
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-muted-foreground">
        {task.city && (
          <span className="inline-flex items-center gap-1.5"><MapPin className="h-3.5 w-3.5" />{task.city}</span>
        )}
        {task.deadline && (
          <span className="inline-flex items-center gap-1.5"><Clock className="h-3.5 w-3.5" />مهلت: {task.deadline.slice(0, 10)}</span>
        )}
        {task.category && (
          <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2 py-0.5">{task.category}</span>
        )}
        <span className="inline-flex items-center gap-1.5"><Users className="h-3.5 w-3.5" />{toFa(task.applicantsCount ?? 0)} متقاضی</span>
        <span className="ms-auto inline-flex items-center gap-1.5 rounded-lg bg-primary/10 px-2 py-1 text-sm font-semibold text-primary">
          <Wallet className="h-3.5 w-3.5" />{toman(task.budget)}
        </span>
      </div>

      {showApply && isOpen && (
        <div className="mt-4 flex justify-end">
          <button
            type="button"
            onClick={apply}
            disabled={applying || applied}
            className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-3.5 py-1.5 text-xs font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60"
          >
            {applied ? (
              <><Check className="h-3.5 w-3.5" /> درخواست ثبت شد</>
            ) : applying ? (
              "در حال ارسال..."
            ) : (
              <><Star className="h-3.5 w-3.5" /> ارسال درخواست همکاری</>
            )}
          </button>
        </div>
      )}
    </Link>
  );
}
