import { createFileRoute } from "@tanstack/react-router";
import { Bell, CheckCheck, ClipboardList, CreditCard, MessagesSquare, UserPlus, BellOff } from "lucide-react";
import { notifications } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/notifications")({
  head: () => ({ meta: [{ title: "Notifications — TaskBridge" }] }),
  component: Notifications,
});

const iconFor = (type: string) => {
  switch (type) {
    case "task_assigned": return ClipboardList;
    case "task_completed": return CheckCheck;
    case "payment": return CreditCard;
    case "application": return UserPlus;
    case "message": return MessagesSquare;
    default: return Bell;
  }
};

function Notifications() {
  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <div className="flex items-end justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Notifications</h1>
          <p className="text-sm text-muted-foreground">Stay on top of every task event.</p>
        </div>
        <button className="rounded-lg border bg-card px-3 py-1.5 text-xs font-medium hover:bg-accent">Mark all as read</button>
      </div>

      <div className="flex gap-2 text-sm">
        {["All", "Unread", "Tasks", "Payments", "Messages"].map((t, i) => (
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
                  <div className="shrink-0 text-xs text-muted-foreground">{n.time}</div>
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
      <h3 className="mt-4 font-semibold">You're all caught up</h3>
      <p className="mt-1 text-sm text-muted-foreground">New notifications will appear here.</p>
    </div>
  );
}
