import type { TaskStatus } from "@/lib/mock-data";
import { cn } from "@/lib/utils";

const map: Record<TaskStatus, { label: string; cls: string; dot: string }> = {
  open:         { label: "Open",         cls: "bg-primary/10 text-primary border-primary/20",      dot: "bg-primary" },
  assigned:     { label: "Assigned",     cls: "bg-secondary/15 text-secondary-foreground border-secondary/30", dot: "bg-secondary" },
  in_progress:  { label: "In Progress",  cls: "bg-warning/15 text-warning-foreground border-warning/30", dot: "bg-warning" },
  completed:    { label: "Completed",    cls: "bg-success/15 text-success-foreground border-success/30", dot: "bg-success" },
  verified:     { label: "Verified",     cls: "bg-success/15 text-success-foreground border-success/30", dot: "bg-success" },
  paid:         { label: "Paid",         cls: "bg-muted text-muted-foreground border-border", dot: "bg-muted-foreground" },
};

export function StatusBadge({ status, className }: { status: TaskStatus; className?: string }) {
  const s = map[status];
  return (
    <span className={cn("inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium", s.cls, className)}>
      <span className={cn("h-1.5 w-1.5 rounded-full", s.dot)} />
      {s.label}
    </span>
  );
}
