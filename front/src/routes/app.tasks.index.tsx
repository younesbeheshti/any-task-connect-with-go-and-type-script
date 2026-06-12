import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Search, SlidersHorizontal, MapPin, Briefcase, ArrowUpDown } from "lucide-react";
import { tasks, cities, categories } from "@/lib/mock-data";
import { TaskCard } from "@/components/task-card";
import { useRole } from "@/components/role-context";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/tasks/")({
  head: () => ({ meta: [{ title: "بازار درخواست‌ها — تسک‌بریج" }] }),
  component: Marketplace,
});

function Marketplace() {
  const { role } = useRole();
  const [query, setQuery] = useState("");
  const [city, setCity] = useState<string | null>(null);
  const [cat, setCat] = useState<string | null>(null);
  const [sort, setSort] = useState<"newest" | "budget" | "deadline">("newest");

  const filtered = tasks
    .filter(t => (!city || t.city === city) && (!cat || t.category === cat))
    .filter(t => !query || (t.title + t.description + t.id).toLowerCase().includes(query.toLowerCase()))
    .sort((a, b) => sort === "budget" ? b.budget - a.budget : 0);

  const title = role === "agent" ? "فرصت‌های جدید" : role === "admin" ? "همه درخواست‌ها" : "درخواست‌های من";
  const subtitle = role === "agent"
    ? `${toFa(filtered.length)} فرصت کاری مناسب مهارت شما`
    : `${toFa(filtered.length)} درخواست — مدیریت و پیگیری`;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">{title}</h1>
        <p className="text-sm text-muted-foreground">{subtitle}</p>
      </div>

      <div className="rounded-2xl border bg-card p-4 shadow-soft">
        <div className="grid gap-3 lg:grid-cols-[1fr_auto_auto_auto]">
          <div className="relative">
            <Search className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input
              value={query} onChange={e => setQuery(e.target.value)}
              placeholder="جستجو در درخواست‌ها…"
              className="w-full rounded-lg border bg-background py-2 pr-9 pl-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20"
            />
          </div>
          <Select value={city} onChange={setCity} placeholder="همه شهرها" icon={MapPin} options={cities} />
          <Select value={cat} onChange={setCat} placeholder="همه دسته‌ها" icon={Briefcase} options={categories.map(c => c.label)} />
          <Select value={sort} onChange={(v: any) => setSort(v)} placeholder="مرتب‌سازی" icon={ArrowUpDown} options={["newest", "budget", "deadline"]} labels={{ newest: "جدیدترین", budget: "بیشترین مبلغ", deadline: "نزدیک‌ترین مهلت" }} />
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {filtered.map(t => <TaskCard key={t.id} task={t} showApply={role === "agent"} />)}
        {filtered.length === 0 && (
          <div className="col-span-full rounded-2xl border border-dashed bg-card/50 p-12 text-center">
            <SlidersHorizontal className="mx-auto h-8 w-8 text-muted-foreground" />
            <h3 className="mt-4 font-semibold">درخواستی با این فیلترها پیدا نشد</h3>
            <p className="mt-1 text-sm text-muted-foreground">جستجو را گسترده‌تر کنید یا فیلترها را پاک کنید.</p>
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
        <span className="flex-1 text-right">{value ? (labels?.[value] ?? value) : placeholder}</span>
      </button>
      {open && (
        <>
          <button className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
          <div className="absolute left-0 z-40 mt-1 min-w-[200px] rounded-lg border bg-popover p-1 shadow-elevated">
            <button onClick={() => { onChange(null); setOpen(false); }} className="block w-full rounded-md px-3 py-2 text-right text-sm hover:bg-accent">{placeholder}</button>
            {options.map((o: string) => (
              <button key={o} onClick={() => { onChange(o); setOpen(false); }} className={cn("block w-full rounded-md px-3 py-2 text-right text-sm hover:bg-accent", value === o && "bg-accent font-medium")}>
                {labels?.[o] ?? o}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
