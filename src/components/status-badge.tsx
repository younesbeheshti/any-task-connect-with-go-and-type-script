import type { TaskStatus } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

const map: Record<TaskStatus, { label: string; cls: string; dot: string }> = {
  posted:                { label: "ثبت شده",         cls: "bg-muted text-muted-foreground border-border",                       dot: "bg-muted-foreground" },
  awaiting_applicants:   { label: "در انتظار متقاضی", cls: "bg-primary/10 text-primary border-primary/20",                       dot: "bg-primary" },
  accepted:              { label: "پذیرفته شده",     cls: "bg-secondary/15 text-secondary-foreground border-secondary/30",      dot: "bg-secondary" },
  in_progress:           { label: "در حال انجام",    cls: "bg-warning/15 text-warning-foreground border-warning/30",            dot: "bg-warning" },
  completed:             { label: "تکمیل شده",       cls: "bg-success/15 text-success-foreground border-success/30",            dot: "bg-success" },
  awaiting_verification: { label: "در انتظار تایید",  cls: "bg-warning/15 text-warning-foreground border-warning/30",            dot: "bg-warning" },
  paid:                  { label: "پرداخت شده",      cls: "bg-success/15 text-success-foreground border-success/30",            dot: "bg-success" },
  cancelled:             { label: "لغو شده",         cls: "bg-destructive/10 text-destructive border-destructive/20",            dot: "bg-destructive" },
};

export function statusLabel(s: TaskStatus) { return map[s].label; }

export function StatusBadge({ status, className }: { status: TaskStatus; className?: string }) {
  const s = map[status];
  return (
    <span className={cn("inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium", s.cls, className)}>
      <span className={cn("h-1.5 w-1.5 rounded-full", s.dot)} />
      {s.label}
    </span>
  );
}
