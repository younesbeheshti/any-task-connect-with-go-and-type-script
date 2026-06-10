import { createFileRoute } from "@tanstack/react-router";
import { ClipboardList, Briefcase, Wallet, Star } from "lucide-react";
import { tasks } from "@/lib/mock-data";
import { TaskCard } from "@/components/task-card";

export const Route = createFileRoute("/app/agent")({
  head: () => ({ meta: [{ title: "Agent Dashboard — TaskBridge" }] }),
  component: AgentDashboard,
});

function AgentDashboard() {
  const open = tasks.filter(t => t.status === "open");
  const mine = tasks.filter(t => ["assigned", "in_progress"].includes(t.status));
  const cards = [
    { label: "Available Tasks", value: open.length.toString(), icon: ClipboardList },
    { label: "Assigned Tasks", value: mine.length.toString(), icon: Briefcase },
    { label: "Earnings (mo)", value: "$1,820", icon: Wallet },
    { label: "Rating", value: "4.9", icon: Star },
  ];
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Agent dashboard</h1>
        <p className="text-sm text-muted-foreground">Tasks tailored to your skills and location.</p>
      </div>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {cards.map(c => (
          <div key={c.label} className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-start justify-between">
              <span className="text-sm text-muted-foreground">{c.label}</span>
              <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><c.icon className="h-4 w-4" /></span>
            </div>
            <div className="mt-3 font-display text-3xl font-bold">{c.value}</div>
          </div>
        ))}
      </div>

      <section>
        <h2 className="mb-3 font-display text-lg font-semibold">Recommended for you</h2>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {open.map(t => <TaskCard key={t.id} task={t} />)}
        </div>
      </section>

      <section>
        <h2 className="mb-3 font-display text-lg font-semibold">Your active tasks</h2>
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {mine.map(t => <TaskCard key={t.id} task={t} />)}
        </div>
      </section>
    </div>
  );
}
