import { Link } from "@tanstack/react-router";
import { Clock, MapPin, Users, Wallet } from "lucide-react";
import { StatusBadge } from "./status-badge";
import type { Task } from "@/lib/mock-data";

export function TaskCard({ task }: { task: Task }) {
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
            <span>{task.postedAgo}</span>
          </div>
          <h3 className="mt-1.5 line-clamp-1 text-base font-semibold tracking-tight text-foreground group-hover:text-primary">
            {task.title}
          </h3>
          <p className="mt-1 line-clamp-2 text-sm text-muted-foreground">{task.description}</p>
        </div>
        <StatusBadge status={task.status} />
      </div>

      <div className="mt-4 flex flex-wrap items-center gap-x-4 gap-y-2 text-xs text-muted-foreground">
        <span className="inline-flex items-center gap-1.5"><MapPin className="h-3.5 w-3.5" />{task.city}</span>
        <span className="inline-flex items-center gap-1.5"><Clock className="h-3.5 w-3.5" />{task.deadline}</span>
        <span className="inline-flex items-center gap-1.5"><Users className="h-3.5 w-3.5" />{task.applicants} applicants</span>
        <span className="ml-auto inline-flex items-center gap-1.5 rounded-lg bg-primary/10 px-2 py-1 text-sm font-semibold text-primary">
          <Wallet className="h-3.5 w-3.5" />${task.budget}
        </span>
      </div>
    </Link>
  );
}
