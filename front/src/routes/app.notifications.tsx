import { createFileRoute } from "@tanstack/react-router";
import { Bell, CheckCheck, ClipboardList, CreditCard, MessagesSquare, UserPlus, BellOff } from "lucide-react";
import { notifications } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/notifications")({
  head: () => ({ meta: [{ title: "اعلان‌ها — تسک‌بریج" }] }),
  component: Notifications,
});

const iconFor = (type: string) => {
  switch (type) {
    case "task_assigned":  return ClipboardList;
    case "task_completed": return CheckCheck;
    case "payment":        return CreditCard;
    case "application":    return UserPlus;
    case "message":        return MessagesSquare;
    default:               return Bell;
  }
};

function Notifications() {
  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">اعلان‌ها</h1>
          <p className="text-sm text-muted-foreground">از تمامی رویدادهای درخواست‌های خود مطلع باشید.</p>
        </div>
        <button className="rounded-lg border bg-card px-3 py-1.5 text-xs font-medium hover:bg-accent">علامت‌گذاری همه به‌عنوان خوانده‌شده</button>
      </div>

      <div className="flex flex-wrap gap-2 text-sm">
        {["همه", "خوانده‌نشده", "درخواست‌ها", "پرداخت‌ها", "پیام‌ها"].map((t, i) => (
          <button key={t} className={cn("rounded-full px-3 py-1.5 font-medium", i === 0 ? "bg-primary text-primary-foreground" : "border bg-card text-muted-foreground hover:bg-accent")}>{t}</button>
        ))}
      </div>

      <div className="overflow-hidden rounded-2xl border bg-card shadow-soft">
        {notifications.length === 0 ? (
          <Empty />
        ) : (
          <ul className="divide-y">
            {notifications.map(n => {
              const Icon = iconFor(n.type);
              return (
                <li key={n.id} className={cn("flex items-start gap-4 px-5 py-4", n.unread && "bg-primary/[0.03]")}>
                  <span className="grid h-10 w-10 place-items-center rounded-xl bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <div className="text-sm font-semibold">{n.title}</div>
                      {n.unread && <span className="h-1.5 w-1.5 rounded-full bg-primary" />}
                    </div>
                    <div className="mt-0.5 text-sm text-muted-foreground">{n.desc}</div>
                  </div>
                  <div className="shrink-0 text-xs text-muted-foreground whitespace-nowrap">{n.time}</div>
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
