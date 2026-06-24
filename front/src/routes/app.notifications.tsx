import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Bell, CheckCheck, ClipboardList, CreditCard, MessagesSquare, UserPlus, BellOff, CheckCircle2 } from "lucide-react";
import { cn } from "@/lib/utils";
import { toast } from "sonner";

export const Route = createFileRoute("/app/notifications")({
  head: () => ({ meta: [{ title: "اعلان‌ها — تسک‌بریج" }] }),
  component: Notifications,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Notification = {
  id: string;
  title: string;
  body: string;
  type: string;
  isRead: boolean;
  deepLink: string | null;
  createdAt: string;
};

type PageResult = { items: Notification[]; total: number; page: number; pageSize: number };

const iconFor = (type: string) => {
  switch (type) {
    case "TASK_ASSIGNED":
    case "TASK_COMPLETED": return ClipboardList;
    case "PAYMENT":        return CreditCard;
    case "APPLICATION":    return UserPlus;
    case "MESSAGE":        return MessagesSquare;
    case "REVIEW":         return CheckCircle2;
    default:               return Bell;
  }
};

type Filter = "all" | "unread";

function Notifications() {
  const [result, setResult] = useState<PageResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [filter, setFilter] = useState<Filter>("all");

  function load(f: Filter = filter) {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    setLoading(true);
    const params = new URLSearchParams({ page: "1", pageSize: "30" });
    if (f === "unread") params.set("unread", "true");
    fetch(`${API_BASE}/v1/notifications?${params}`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => { if (d?.items) setResult(d as PageResult); })
      .catch(() => {})
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, []);

  async function markRead(id: string) {
    const token = getToken();
    if (!token) return;
    await fetch(`${API_BASE}/v1/notifications/${id}/read`, {
      method: "PATCH", headers: { Authorization: `Bearer ${token}` },
    });
    setResult(prev => prev ? {
      ...prev,
      items: prev.items.map(n => n.id === id ? { ...n, isRead: true } : n),
    } : null);
  }

  async function markAllRead() {
    const token = getToken();
    if (!token) return;
    try {
      await fetch(`${API_BASE}/v1/notifications/read-all`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      setResult(prev => prev ? { ...prev, items: prev.items.map(n => ({ ...n, isRead: true })) } : null);
      toast.success("همه اعلان‌ها خوانده‌شده علامت‌گذاری شدند");
    } catch {
      toast.error("خطا در علامت‌گذاری");
    }
  }

  const unreadCount = (result?.items ?? []).filter(n => !n.isRead).length;

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">اعلان‌ها</h1>
          <p className="text-sm text-muted-foreground">
            {unreadCount > 0 ? `${unreadCount} اعلان خوانده‌نشده` : "همه اعلان‌ها خوانده شده‌اند"}
          </p>
        </div>
        {unreadCount > 0 && (
          <button onClick={markAllRead}
            className="inline-flex items-center gap-1.5 rounded-lg border bg-card px-3 py-1.5 text-xs font-medium hover:bg-accent">
            <CheckCheck className="h-3.5 w-3.5" /> علامت‌گذاری همه
          </button>
        )}
      </div>

      <div className="flex gap-2 text-sm">
        {(["all", "unread"] as Filter[]).map(f => (
          <button key={f} onClick={() => { setFilter(f); load(f); }}
            className={cn("rounded-full px-3 py-1.5 font-medium",
              filter === f ? "bg-primary text-primary-foreground" : "border bg-card text-muted-foreground hover:bg-accent"
            )}>
            {f === "all" ? "همه" : "خوانده‌نشده"}
          </button>
        ))}
      </div>

      <div className="overflow-hidden rounded-2xl border bg-card shadow-soft">
        {loading ? (
          <div className="divide-y">
            {[1, 2, 3, 4, 5].map(i => (
              <div key={i} className="flex items-start gap-4 px-5 py-4 animate-pulse">
                <div className="h-10 w-10 rounded-xl bg-muted" />
                <div className="flex-1 space-y-2">
                  <div className="h-4 w-1/3 rounded bg-muted" />
                  <div className="h-3 w-2/3 rounded bg-muted" />
                </div>
              </div>
            ))}
          </div>
        ) : !result?.items.length ? (
          <Empty />
        ) : (
          <ul className="divide-y">
            {result.items.map(n => {
              const Icon = iconFor(n.type);
              return (
                <li
                  key={n.id}
                  onClick={() => !n.isRead && markRead(n.id)}
                  className={cn(
                    "flex cursor-pointer items-start gap-4 px-5 py-4 transition-colors hover:bg-muted/30",
                    !n.isRead && "bg-primary/[0.03]"
                  )}
                >
                  <span className={cn(
                    "grid h-10 w-10 place-items-center rounded-xl",
                    n.isRead ? "bg-muted/50 text-muted-foreground" : "bg-primary/10 text-primary"
                  )}>
                    <Icon className="h-4 w-4" />
                  </span>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <div className={cn("text-sm font-medium", !n.isRead && "font-semibold")}>{n.title}</div>
                      {!n.isRead && <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full bg-primary" />}
                    </div>
                    {n.body && <div className="mt-0.5 text-sm text-muted-foreground line-clamp-2">{n.body}</div>}
                  </div>
                  <div className="shrink-0 text-xs text-muted-foreground whitespace-nowrap">
                    {n.createdAt?.slice(0, 10)}
                  </div>
                </li>
              );
            })}
          </ul>
        )}
      </div>
    </div>
  );
}

function Empty() {
  return (
    <div className="flex flex-col items-center px-6 py-16 text-center">
      <BellOff className="h-10 w-10 text-muted-foreground" />
      <h3 className="mt-4 font-semibold">اعلان جدیدی ندارید</h3>
      <p className="mt-1 text-sm text-muted-foreground">اعلان‌های جدید در اینجا نمایش داده می‌شوند.</p>
    </div>
  );
}
