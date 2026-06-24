import { Check } from "lucide-react";
import { cn } from "@/lib/utils";

const STEPS: { key: string; label: string }[] = [
  { key: "CREATED",                  label: "ثبت‌شده" },
  { key: "OPEN",                     label: "در انتظار متقاضی" },
  { key: "ASSIGNED",                 label: "پذیرفته‌شده" },
  { key: "IN_PROGRESS",              label: "در حال انجام" },
  { key: "COMPLETED",                label: "تکمیل‌شده" },
  { key: "WAITING_FOR_VERIFICATION", label: "در انتظار تایید" },
  { key: "VERIFIED",                 label: "تایید‌شده" },
  { key: "PAID",                     label: "پرداخت‌شده" },
];

// Map legacy lowercase status values to their index
const STATUS_INDEX: Record<string, number> = {
  posted: 0, awaiting_applicants: 1, accepted: 2, in_progress: 3,
  completed: 4, awaiting_verification: 5, paid: 7, cancelled: -1,
};

export function TaskTimeline({ status }: { status: string }) {
  const idx = STEPS.findIndex(s => s.key === status);
  const currentIdx = idx >= 0 ? idx : (STATUS_INDEX[status] ?? 0);

  return (
    <ol className="flex flex-col gap-0 sm:flex-row sm:items-start sm:gap-0">
      {STEPS.map((step, i) => {
        const done = i < currentIdx;
        const current = i === currentIdx;
        return (
          <li key={step.key} className="flex flex-1 items-start gap-3 sm:flex-col sm:items-center sm:text-center">
            <div className="flex flex-col items-center sm:w-full">
              <div className="flex w-full items-center">
                {i > 0 && <div className={cn("hidden h-0.5 flex-1 sm:block", done || current ? "bg-primary" : "bg-border")} />}
                <div className={cn(
                  "relative z-10 grid h-9 w-9 shrink-0 place-items-center rounded-full border-2 transition-all",
                  done && "border-primary bg-primary text-primary-foreground",
                  current && "border-primary bg-background text-primary shadow-glow",
                  !done && !current && "border-border bg-background text-muted-foreground"
                )}>
                  {done ? <Check className="h-4 w-4" /> : <span className="text-xs font-semibold">{i + 1}</span>}
                  {current && <span className="absolute inset-0 -z-10 animate-ping rounded-full bg-primary/30" />}
                </div>
                {i < STEPS.length - 1 && <div className={cn("hidden h-0.5 flex-1 sm:block", done ? "bg-primary" : "bg-border")} />}
              </div>
              <span className={cn("mt-2 text-xs font-medium", (done || current) ? "text-foreground" : "text-muted-foreground")}>
                {step.label}
              </span>
            </div>
            <div className={cn("ms-1 flex-1 pb-6 sm:hidden", i === STEPS.length - 1 && "pb-0")}>
              <div className={cn("text-sm font-medium", (done || current) ? "text-foreground" : "text-muted-foreground")}>
                {step.label}
              </div>
            </div>
          </li>
        );
      })}
    </ol>
  );
}
