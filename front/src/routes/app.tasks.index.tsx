import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { Search, SlidersHorizontal, MapPin, Briefcase, ArrowUpDown, PlusCircle } from "lucide-react";
import { TaskCard } from "@/components/task-card";
import { useRole } from "@/components/role-context";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";
import type { ApiTask, ApiCategory, ApiCity } from "@/lib/types";
import { Link } from "@tanstack/react-router";

export const Route = createFileRoute("/app/tasks/")({
  head: () => ({ meta: [{ title: "بازار درخواست‌ها — تسک‌بریج" }] }),
  component: Marketplace,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

function Marketplace() {
  const { role } = useRole();
  const [tasks, setTasks] = useState<ApiTask[]>([]);
  const [categories, setCategories] = useState<ApiCategory[]>([]);
  const [cities, setCities] = useState<ApiCity[]>([]);
  const [loading, setLoading] = useState(true);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);

  const [query, setQuery] = useState("");
  const [cityId, setCityId] = useState<string | null>(null);
  const [catId, setCatId] = useState<string | null>(null);
  const [sort, setSort] = useState<"newest" | "budget" | "deadline">("newest");

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    const token = getToken();
    const h = token ? { Authorization: `Bearer ${token}` } : {};
    Promise.all([
      fetch(`${API_BASE}/v1/categories`, { headers: h }).then(r => r.json()),
      fetch(`${API_BASE}/v1/cities`, { headers: h }).then(r => r.json()),
    ]).then(([c, ci]) => {
      setCategories(c.categories ?? c ?? []);
      setCities(ci.cities ?? ci ?? []);
    }).catch(() => {});
  }, []);

  function loadTasks(p = page) {
    const token = getToken();
    const h = token ? { Authorization: `Bearer ${token}` } : {};
    setLoading(true);
    const params = new URLSearchParams({ page: String(p), pageSize: "12" });
    if (query) params.set("q", query);
    if (cityId) params.set("cityId", cityId);
    if (catId) params.set("categoryId", catId);
    if (sort === "budget") params.set("sort", "budget");
    if (sort === "deadline") params.set("sort", "deadline");
    if (role === "requester") params.set("mine", "true");
    fetch(`${API_BASE}/v1/tasks?${params}`, { headers: h })
      .then(r => r.json())
      .then(d => {
        setTasks(d.tasks ?? []);
        setTotal(d.total ?? 0);
      })
      .catch(() => {})
      .finally(() => setLoading(false));
  }

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => { setPage(1); loadTasks(1); }, 300);
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current); };
  }, [query, cityId, catId, sort, role]);

  useEffect(() => { loadTasks(); }, [page]);

  const title = role === "agent" ? "فرصت‌های جدید" : role === "admin" ? "همه درخواست‌ها" : "درخواست‌های من";
  const subtitle = role === "agent"
    ? `${toFa(total)} فرصت کاری`
    : `${toFa(total)} درخواست`;

  const totalPages = Math.ceil(total / 12);

  return (
    <div className="space-y-6">
      <div className="flex flex-wrap items-end justify-between gap-3">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">{title}</h1>
          <p className="text-sm text-muted-foreground">{subtitle}</p>
        </div>
        {role === "requester" && (
          <Link to="/app/tasks/new" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            <PlusCircle className="h-4 w-4" /> ثبت درخواست جدید
          </Link>
        )}
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
          <DropSelect
            value={cityId} onChange={v => { setCityId(v); setPage(1); }}
            placeholder="همه شهرها" icon={MapPin}
            options={cities.map(c => ({ value: c.id, label: c.title }))}
          />
          <DropSelect
            value={catId} onChange={v => { setCatId(v); setPage(1); }}
            placeholder="همه دسته‌ها" icon={Briefcase}
            options={categories.map(c => ({ value: c.id, label: c.title }))}
          />
          <DropSelect
            value={sort} onChange={v => { setSort(v as any); setPage(1); }}
            placeholder="مرتب‌سازی" icon={ArrowUpDown}
            options={[
              { value: "newest", label: "جدیدترین" },
              { value: "budget", label: "بیشترین مبلغ" },
              { value: "deadline", label: "نزدیک‌ترین مهلت" },
            ]}
          />
        </div>
      </div>

      {loading ? (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse">
              <div className="space-y-3">
                <div className="h-4 w-1/4 rounded bg-muted" />
                <div className="h-5 w-3/4 rounded bg-muted" />
                <div className="h-4 w-full rounded bg-muted" />
                <div className="h-4 w-2/3 rounded bg-muted" />
              </div>
            </div>
          ))}
        </div>
      ) : tasks.length === 0 ? (
        <div className="rounded-2xl border border-dashed bg-card/50 p-12 text-center">
          <SlidersHorizontal className="mx-auto h-8 w-8 text-muted-foreground" />
          <h3 className="mt-4 font-semibold">درخواستی با این فیلترها پیدا نشد</h3>
          <p className="mt-1 text-sm text-muted-foreground">جستجو را گسترده‌تر کنید یا فیلترها را پاک کنید.</p>
        </div>
      ) : (
        <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
          {tasks.map(t => <TaskCard key={t.id} task={t} showApply={role === "agent"} />)}
        </div>
      )}

      {totalPages > 1 && (
        <div className="flex items-center justify-between">
          <span className="text-xs text-muted-foreground">
            صفحه {toFa(page)} از {toFa(totalPages)}
          </span>
          <div className="flex gap-2">
            <button disabled={page === 1} onClick={() => setPage(p => p - 1)}
              className="rounded border px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-accent">قبلی</button>
            <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}
              className="rounded border px-3 py-1.5 text-sm disabled:opacity-40 hover:bg-accent">بعدی</button>
          </div>
        </div>
      )}
    </div>
  );
}

type DropOption = { value: string; label: string };
function DropSelect({ value, onChange, placeholder, icon: Icon, options }: {
  value: string | null; onChange: (v: string | null) => void;
  placeholder: string; icon: any; options: DropOption[];
}) {
  const [open, setOpen] = useState(false);
  const selected = options.find(o => o.value === value);
  return (
    <div className="relative">
      <button
        onClick={() => setOpen(o => !o)}
        className={cn("flex w-full items-center gap-2 rounded-lg border bg-background px-3 py-2 text-sm hover:bg-accent min-w-[140px]", value && "border-primary/40 text-foreground")}
      >
        <Icon className="h-4 w-4 text-muted-foreground" />
        <span className="flex-1 text-right truncate">{selected?.label ?? placeholder}</span>
      </button>
      {open && (
        <>
          <button className="fixed inset-0 z-30" onClick={() => setOpen(false)} />
          <div className="absolute left-0 z-40 mt-1 min-w-[180px] rounded-lg border bg-popover p-1 shadow-elevated">
            <button onClick={() => { onChange(null); setOpen(false); }} className="block w-full rounded-md px-3 py-2 text-right text-sm hover:bg-accent">{placeholder}</button>
            {options.map(o => (
              <button key={o.value} onClick={() => { onChange(o.value); setOpen(false); }}
                className={cn("block w-full rounded-md px-3 py-2 text-right text-sm hover:bg-accent", value === o.value && "bg-accent font-medium")}>
                {o.label}
              </button>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
