import { cn } from "@/lib/utils";

export type TaskStatus = string;

const map: Record<string, { label: string; cls: string; dot: string }> = {
  // Backend status values (uppercase)
  CREATED:                  { label: "ثبت‌شده",        cls: "bg-muted text-muted-foreground border-border",               dot: "bg-muted-foreground" },
  OPEN:                     { label: "در انتظار متقاضی", cls: "bg-primary/10 text-primary border-primary/20",             dot: "bg-primary" },
  ASSIGNED:                 { label: "پذیرفته‌شده",     cls: "bg-secondary/15 text-secondary-foreground border-secondary/30", dot: "bg-secondary" },
  IN_PROGRESS:              { label: "در حال انجام",    cls: "bg-warning/15 text-warning-foreground border-warning/30",   dot: "bg-warning" },
  COMPLETED:                { label: "تکمیل‌شده",       cls: "bg-success/15 text-success-foreground border-success/30",   dot: "bg-success" },
  WAITING_FOR_VERIFICATION: { label: "در انتظار تایید", cls: "bg-warning/15 text-warning-foreground border-warning/30",   dot: "bg-warning" },
  VERIFIED:                 { label: "تایید‌شده",       cls: "bg-success/15 text-success-foreground border-success/30",   dot: "bg-success" },
  PAID:                     { label: "پرداخت‌شده",      cls: "bg-success/15 text-success-foreground border-success/30",   dot: "bg-success" },
  CANCELLED:                { label: "لغو‌شده",         cls: "bg-destructive/10 text-destructive border-destructive/20",  dot: "bg-destructive" },
  // Legacy lowercase values (kept for backward compat)
  posted:                   { label: "ثبت‌شده",         cls: "bg-muted text-muted-foreground border-border",               dot: "bg-muted-foreground" },
  awaiting_applicants:      { label: "در انتظار متقاضی", cls: "bg-primary/10 text-primary border-primary/20",             dot: "bg-primary" },
  accepted:                 { label: "پذیرفته‌شده",     cls: "bg-secondary/15 text-secondary-foreground border-secondary/30", dot: "bg-secondary" },
  in_progress:              { label: "در حال انجام",    cls: "bg-warning/15 text-warning-foreground border-warning/30",   dot: "bg-warning" },
  completed:                { label: "تکمیل‌شده",       cls: "bg-success/15 text-success-foreground border-success/30",   dot: "bg-success" },
  awaiting_verification:    { label: "در انتظار تایید", cls: "bg-warning/15 text-warning-foreground border-warning/30",   dot: "bg-warning" },
  verified:                 { label: "تایید‌شده — در انتظار تسویه", cls: "bg-success/15 text-success-foreground border-success/30", dot: "bg-success" },
  paid:                     { label: "پرداخت‌شده",      cls: "bg-success/15 text-success-foreground border-success/30",   dot: "bg-success" },
  cancelled:                { label: "لغو‌شده",         cls: "bg-destructive/10 text-destructive border-destructive/20",  dot: "bg-destructive" },
};

const fallback = { label: "نامشخص", cls: "bg-muted text-muted-foreground border-border", dot: "bg-muted-foreground" };

export function statusLabel(s: string): string { return (map[s] ?? fallback).label; }

export function StatusBadge({ status, className }: { status: string; className?: string }) {
  const s = map[status] ?? fallback;
  return (
    <span className={cn("inline-flex items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium", s.cls, className)}>
      <span className={cn("h-1.5 w-1.5 rounded-full", s.dot)} />
      {s.label}
    </span>
  );
}
