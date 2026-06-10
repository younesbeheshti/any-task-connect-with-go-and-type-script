import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Search, SlidersHorizontal, MapPin, Briefcase, Clock, ArrowUpDown } from "lucide-react";
import { tasks, cities, categories } from "@/lib/mock-data";
import { TaskCard } from "@/components/task-card";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/tasks/")({
  head: () => ({ meta: [{ title: "Marketplace — TaskBridge" }] }),
  component: Marketplace,
});

function Marketplace() {
  const [query, setQuery] = useState("");
  const [city, setCity] = useState<string | null>(null);
  const [cat, setCat] = useState<string | null>(null);
  const [sort, setSort] = useState<"newest" | "budget" | "deadline">("newest");

  const filtered = tasks
    .filter(t => (!city || t.city === city) && (!cat || t.category === cat))
    .filter(t => !query || (t.title + t.description + t.id).toLowerCase().includes(query.toLowerCase()))
    .sort((a, b) => sort === "budget" ? b.budget - a.budget : 0);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Task marketplace</h1>
        <p className="text-sm text-muted-foreground">{filtered.length} tasks available — find one that matches your skills.</p>
      </div>

      <div className="rounded-2xl border bg-card p-4 shadow-soft">
        <div className="grid gap-3 lg:grid-cols-[1fr_auto_auto_auto]">
          <div className="relative">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input
              value={query} onChange={e => setQuery(e.target.value)}
              placeholder="Search tasks…"
              className="w-full rounded-lg border bg-background py-2 pl-9 pr-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20"
            />
          </div>
          <Select value={city} onChange={setCity} placeholder="All cities" icon={MapPin} options={cities} />
          <Select value={cat} onChange={setCat} placeholder="All categories" icon={Briefcase} options={categories.map(c => c.label)} />
          <Select value={sort} onChange={(v: any) => setSort(v)} placeholder="Sort" icon={ArrowUpDown} options={["newest", "budget", "deadline"]} labels={{ newest: "Newest", budget: "Highest budget", deadline: "Closest deadline" }} />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {filtered.map(t => <TaskCard key={t.id} task={t} />)}
        {filtered.length === 0 && (
          <div className="col-span-full rounded-2xl border border-dashed bg-card/50 p-12 text-center">
            <SlidersHorizontal className="mx-auto h-8 w-8 text-muted-foreground" />
            <h3 className="mt-4 font-semibold">No tasks match your filters</h3>
            <p className="mt-1 text-sm text-muted-foreground">Try widening your search or clearing filters.</p>
          </div>
        )}
      </div>
    </div>
  );
}

function Select({ value, onChange, placeholder, icon: Icon, options, labels }: any) {
  const [open, setOpen] = useState(false);
  return (
    <div className="relative">
      <button
        onClick={() => setOpen(o => !o)}
        className={cn("flex w-full items-center gap-2 rounded-lg border bg-background px-3 py-2 text-sm hover:bg-accent", value && "border-primary/40 text-foreground")}
      >
        <Icon className="h-4 w-4 text-muted-foreground" />
        <span className="flex-1 text-left">{value ? (labels?.[value] ?? value) : placeholder}</span>
      </button>
      {open && (
        <>
          <button className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
          <div className="absolute right-0 z-40 mt-1 min-w-[180px] rounded-lg border bg-popover p-1 shadow-elevated">
            <button onClick={() => { onChange(null); setOpen(false); }} className="block w-full rounded-md px-3 py-2 text-left text-sm hover:bg-accent">{placeholder}</button>
            {options.map((o: string) => (
              <button key={o} onClick={() => { onChange(o); setOpen(false); }} className={cn("block w-full rounded-md px-3 py-2 text-left text-sm hover:bg-accent", value === o && "bg-accent font-medium")}>
                {labels?.[o] ?? o}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
