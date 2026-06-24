import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import {
  Users, Search, ShieldCheck, ShieldX, X, Phone, Mail,
  Star, CheckCircle2, Calendar, UserCog, ChevronLeft,
} from "lucide-react";
import { toFa } from "@/lib/fa";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/admin/users")({
  head: () => ({ meta: [{ title: "مدیریت کاربران — تسک‌بریج" }] }),
  component: AdminUsersPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type UserItem = {
  id: string; fullName: string; phone: string; email: string | null;
  role: string; isActive: boolean; isVerified: boolean;
  rating: number; completedCount: number; createdAt: string;
};

type PageResult = { items: UserItem[]; total: number; page: number; pageSize: number };

const ROLE_LABELS: Record<string, string> = {
  requester: "کارفرما", agent: "مجری", admin: "مدیر",
  REQUESTER: "کارفرما", AGENT: "مجری", ADMIN: "مدیر",
};
const ROLE_COLORS: Record<string, string> = {
  requester: "bg-secondary/20 text-secondary-foreground",
  agent: "bg-primary/10 text-primary",
  admin: "bg-destructive/10 text-destructive",
  REQUESTER: "bg-secondary/20 text-secondary-foreground",
  AGENT: "bg-primary/10 text-primary",
  ADMIN: "bg-destructive/10 text-destructive",
};

function AdminUsersPage() {
  const [result, setResult] = useState<PageResult | null>(null);
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState("");
  const [roleFilter, setRoleFilter] = useState("");
  const [activeFilter, setActiveFilter] = useState<"" | "true" | "false">("");
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<UserItem | null>(null);
  const [editingRole, setEditingRole] = useState<string>("");
  const [savingRole, setSavingRole] = useState(false);

  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  function load(p = page, q = query, rf = roleFilter, af = activeFilter) {
    const token = getToken();
    if (!token) { setLoading(false); return; }
    setLoading(true);
    const params = new URLSearchParams({ page: String(p), pageSize: "20" });
    if (q) params.set("q", q);
    if (rf) params.set("role", rf);
    if (af !== "") params.set("isActive", af);
    fetch(`${API_BASE}/v1/admin/users?${params}`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => { if (d?.items) setResult(d as PageResult); else toast.error("خطا در بارگذاری کاربران"); })
      .catch(() => toast.error("خطا در بارگذاری کاربران"))
      .finally(() => setLoading(false));
  }

  useEffect(() => { load(); }, [page, roleFilter, activeFilter]);

  useEffect(() => {
    if (debounceRef.current) clearTimeout(debounceRef.current);
    debounceRef.current = setTimeout(() => { setPage(1); load(1, query, roleFilter, activeFilter); }, 350);
    return () => { if (debounceRef.current) clearTimeout(debounceRef.current); };
  }, [query]);

  async function toggleActive(user: UserItem) {
    const token = getToken();
    if (!token) return;
    const action = user.isActive ? "suspend" : "activate";
    try {
      const res = await fetch(`${API_BASE}/v1/admin/users/${user.id}/${action}`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success(user.isActive ? "حساب تعلیق شد" : "حساب فعال شد");
      const updated = { ...user, isActive: !user.isActive };
      setResult(prev => prev ? { ...prev, items: prev.items.map(u => u.id === user.id ? updated : u) } : null);
      if (selected?.id === user.id) setSelected(updated);
    } catch (e: any) { toast.error(e.message); }
  }

  async function saveRole(user: UserItem) {
    if (editingRole === user.role) { setEditingRole(""); return; }
    const token = getToken();
    if (!token) return;
    setSavingRole(true);
    try {
      const res = await fetch(`${API_BASE}/v1/admin/users/${user.id}`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ role: editingRole }),
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("نقش کاربر تغییر کرد");
      const updated = { ...user, role: editingRole };
      setResult(prev => prev ? { ...prev, items: prev.items.map(u => u.id === user.id ? updated : u) } : null);
      setSelected(updated);
      setEditingRole("");
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setSavingRole(false);
    }
  }

  function openDetail(u: UserItem) {
    setSelected(u);
    setEditingRole(u.role);
  }

  const totalPages = result ? Math.ceil(result.total / result.pageSize) : 1;

  return (
    <div className="flex gap-6 min-h-[600px]">
      {/* Table */}
      <div className={cn("flex-1 min-w-0 space-y-4", selected && "hidden lg:block")}>
        <div className="flex flex-wrap items-end justify-between gap-3">
          <div>
            <h1 className="text-2xl font-bold tracking-tight">مدیریت کاربران</h1>
            <p className="text-sm text-muted-foreground">
              {result ? `${toFa(result.total)} کاربر ثبت‌شده` : "در حال بارگذاری..."}
            </p>
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-4 shadow-soft">
          <div className="flex flex-wrap gap-3">
            <div className="relative flex-1 min-w-[200px]">
              <Search className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
              <input value={query} onChange={e => setQuery(e.target.value)}
                placeholder="جستجو نام، تلفن..."
                className="w-full rounded-lg border bg-background py-2 pr-9 pl-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20" />
            </div>
            <select value={roleFilter} onChange={e => { setRoleFilter(e.target.value); setPage(1); }}
              className="rounded-lg border bg-background px-3 py-2 text-sm outline-none focus:border-primary">
              <option value="">همه نقش‌ها</option>
              <option value="requester">کارفرما</option>
              <option value="agent">مجری</option>
              <option value="admin">مدیر</option>
            </select>
            <select value={activeFilter} onChange={e => { setActiveFilter(e.target.value as "" | "true" | "false"); setPage(1); }}
              className="rounded-lg border bg-background px-3 py-2 text-sm outline-none focus:border-primary">
              <option value="">همه وضعیت‌ها</option>
              <option value="true">فعال</option>
              <option value="false">تعلیق‌شده</option>
            </select>
          </div>
        </div>

        <div className="rounded-2xl border bg-card shadow-soft overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead className="border-b bg-muted/30">
                <tr>
                  <th className="px-4 py-3 text-right font-medium text-muted-foreground">کاربر</th>
                  <th className="px-4 py-3 text-right font-medium text-muted-foreground hidden sm:table-cell">نقش</th>
                  <th className="px-4 py-3 text-right font-medium text-muted-foreground hidden md:table-cell">وضعیت</th>
                  <th className="px-4 py-3 text-right font-medium text-muted-foreground hidden lg:table-cell">امتیاز</th>
                  <th className="px-4 py-3 text-right font-medium text-muted-foreground">اقدام</th>
                </tr>
              </thead>
              <tbody className="divide-y">
                {loading ? (
                  Array.from({ length: 8 }).map((_, i) => (
                    <tr key={i}>
                      {[1,2,3,4,5].map((_, j) => (
                        <td key={j} className="px-4 py-3">
                          <div className="h-4 w-full animate-pulse rounded bg-muted" />
                        </td>
                      ))}
                    </tr>
                  ))
                ) : (result?.items?.length ?? 0) === 0 ? (
                  <tr>
                    <td colSpan={5} className="py-12 text-center text-muted-foreground">
                      <Users className="mx-auto mb-2 h-8 w-8 opacity-40" />کاربری یافت نشد
                    </td>
                  </tr>
                ) : (
                  result?.items.map(u => (
                    <tr key={u.id}
                      onClick={() => openDetail(u)}
                      className={cn("hover:bg-muted/30 cursor-pointer", selected?.id === u.id && "bg-primary/5")}>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-3">
                          <span className="grid h-9 w-9 place-items-center rounded-full bg-primary/10 text-primary font-semibold text-sm shrink-0">
                            {u.fullName.slice(0, 2)}
                          </span>
                          <div className="min-w-0">
                            <div className="font-medium truncate">{u.fullName}</div>
                            <div className="text-xs text-muted-foreground">{u.phone}</div>
                          </div>
                        </div>
                      </td>
                      <td className="px-4 py-3 hidden sm:table-cell">
                        <span className={cn("rounded-full px-2 py-0.5 text-xs font-medium", ROLE_COLORS[u.role])}>
                          {ROLE_LABELS[u.role] ?? u.role}
                        </span>
                      </td>
                      <td className="px-4 py-3 hidden md:table-cell">
                        <span className={cn("inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium",
                          u.isActive ? "bg-success/10 text-success" : "bg-destructive/10 text-destructive")}>
                          {u.isActive ? <ShieldCheck className="h-3 w-3" /> : <ShieldX className="h-3 w-3" />}
                          {u.isActive ? "فعال" : "تعلیق"}
                        </span>
                      </td>
                      <td className="px-4 py-3 hidden lg:table-cell tabular-nums text-sm">
                        {u.rating.toFixed(1)} ⭐
                      </td>
                      <td className="px-4 py-3">
                        <div className="flex items-center gap-2" onClick={e => e.stopPropagation()}>
                          <button onClick={() => toggleActive(u)}
                            className={cn("rounded-lg px-3 py-1 text-xs font-medium",
                              u.isActive
                                ? "border border-destructive/40 text-destructive hover:bg-destructive/10"
                                : "border border-success/40 text-success hover:bg-success/10")}>
                            {u.isActive ? "تعلیق" : "فعال‌سازی"}
                          </button>
                          <ChevronLeft className="h-4 w-4 text-muted-foreground" />
                        </div>
                      </td>
                    </tr>
                  ))
                )}
              </tbody>
            </table>
          </div>

          {result && result.total > result.pageSize && (
            <div className="flex items-center justify-between border-t px-4 py-3">
              <span className="text-xs text-muted-foreground">
                صفحه {toFa(page)} از {toFa(totalPages)}
              </span>
              <div className="flex gap-2">
                <button disabled={page === 1} onClick={() => setPage(p => p - 1)}
                  className="rounded border px-3 py-1 text-xs disabled:opacity-40 hover:bg-accent">قبلی</button>
                <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}
                  className="rounded border px-3 py-1 text-xs disabled:opacity-40 hover:bg-accent">بعدی</button>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Detail panel */}
      {selected && (
        <div className="w-full lg:w-80 shrink-0">
          <div className="rounded-2xl border bg-card shadow-soft sticky top-20">
            <div className="flex items-center justify-between border-b px-5 py-4">
              <h2 className="font-display font-semibold">جزئیات کاربر</h2>
              <button onClick={() => setSelected(null)} className="grid h-8 w-8 place-items-center rounded-lg hover:bg-accent">
                <X className="h-4 w-4" />
              </button>
            </div>
            <div className="p-5 space-y-5">
              <div className="flex items-center gap-3">
                <span className="grid h-14 w-14 place-items-center rounded-2xl gradient-brand text-xl font-bold text-white">
                  {selected.fullName.slice(0, 2)}
                </span>
                <div>
                  <div className="font-semibold">{selected.fullName}</div>
                  <span className={cn("rounded-full px-2 py-0.5 text-xs font-medium", ROLE_COLORS[selected.role])}>
                    {ROLE_LABELS[selected.role] ?? selected.role}
                  </span>
                </div>
              </div>

              <dl className="space-y-3 text-sm">
                <DetailRow icon={Phone}    label="تلفن"          value={selected.phone} ltr />
                <DetailRow icon={Mail}     label="ایمیل"         value={selected.email ?? "ثبت نشده"} ltr />
                <DetailRow icon={Star}     label="امتیاز"        value={`${toFa(selected.rating.toFixed(1))} ⭐`} />
                <DetailRow icon={CheckCircle2} label="کارها"     value={`${toFa(selected.completedCount)} تکمیل‌شده`} />
                <DetailRow icon={Calendar} label="تاریخ عضویت"  value={selected.createdAt?.slice(0, 10) ?? "—"} />
              </dl>

              <div>
                <label className="text-xs font-medium text-muted-foreground">تغییر نقش</label>
                <div className="mt-1.5 flex gap-2">
                  <select value={editingRole} onChange={e => setEditingRole(e.target.value)}
                    className="flex-1 rounded-lg border bg-background px-3 py-2 text-sm outline-none focus:border-primary">
                    <option value="requester">کارفرما</option>
                    <option value="agent">مجری</option>
                    <option value="admin">مدیر</option>
                  </select>
                  <button onClick={() => saveRole(selected)} disabled={savingRole || editingRole === selected.role}
                    className="rounded-lg gradient-brand px-3 py-2 text-xs font-semibold text-white disabled:opacity-50">
                    {savingRole ? "..." : "ذخیره"}
                  </button>
                </div>
              </div>

              <div className="flex gap-2 pt-1 border-t">
                <button onClick={() => toggleActive(selected)}
                  className={cn("flex-1 rounded-lg py-2 text-sm font-medium",
                    selected.isActive
                      ? "border border-destructive/40 text-destructive hover:bg-destructive/10"
                      : "border border-success/40 text-success hover:bg-success/10")}>
                  {selected.isActive ? <><ShieldX className="inline h-4 w-4 me-1" />تعلیق حساب</> : <><ShieldCheck className="inline h-4 w-4 me-1" />فعال‌سازی</>}
                </button>
                <button className="flex-1 rounded-lg border py-2 text-sm font-medium hover:bg-accent flex items-center justify-center gap-1">
                  <UserCog className="h-4 w-4" /> ورود به حساب
                </button>
              </div>

              <div className="text-xs text-muted-foreground break-all">
                <span className="font-mono">ID: {selected.id}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

function DetailRow({ icon: Icon, label, value, ltr }: { icon: any; label: string; value: string; ltr?: boolean }) {
  return (
    <div className="flex items-center gap-3">
      <Icon className="h-4 w-4 shrink-0 text-muted-foreground" />
      <span className="w-20 text-muted-foreground">{label}</span>
      <span className={cn("flex-1 font-medium truncate", ltr && "text-left")} dir={ltr ? "ltr" : undefined}>{value}</span>
    </div>
  );
}
