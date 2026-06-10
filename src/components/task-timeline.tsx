import { Check } from "lucide-react";
import { cn } from "@/lib/utils";
import type { TaskStatus } from "@/lib/mock-data";

const STEPS: { key: TaskStatus; label: string }[] = [
  { key: "open",        label: "Created" },
  { key: "assigned",    label: "Assigned" },
  { key: "in_progress", label: "In Progress" },
  { key: "completed",   label: "Completed" },
  { key: "verified",    label: "Verified" },
  { key: "paid",        label: "Paid" },
];

export function TaskTimeline({ status }: { status: TaskStatus }) {
  const currentIdx = STEPS.findIndex(s => s.key === status);
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
                <div
                  className={cn(
                    "relative z-10 grid h-9 w-9 shrink-0 place-items-center rounded-full border-2 transition-all",
                    done && "border-primary bg-primary text-primary-foreground",
                    current && "border-primary bg-background text-primary shadow-glow",
                    !done && !current && "border-border bg-background text-muted-foreground"
                  )}
                >
                  {done ? <Check className="h-4 w-4" /> : <span className="text-xs font-semibold">{i + 1}</span>}
                  {current && <span className="absolute inset-0 -z-10 animate-ping rounded-full bg-primary/30" />}
                </div>
                {i < STEPS.length - 1 && <div className={cn("hidden h-0.5 flex-1 sm:block", done ? "bg-primary" : "bg-border")} />}
              </div>
              <span className={cn("mt-2 text-xs font-medium", (done || current) ? "text-foreground" : "text-muted-foreground")}>{step.label}</span>
            </div>
            <div className={cn("ml-1 flex-1 pb-6 sm:hidden", i === STEPS.length - 1 && "pb-0")}>
              <div className={cn("text-sm font-medium", (done || current) ? "text-foreground" : "text-muted-foreground")}>{step.label}</div>
            </div>
          </li>
        );
      })}
    </ol>
  );
}
