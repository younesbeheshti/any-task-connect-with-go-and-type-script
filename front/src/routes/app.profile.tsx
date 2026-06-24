import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import {
  BadgeCheck, MapPin, Star, Wallet, ShieldCheck, Mail, Phone,
  Pencil, Check, X, ClipboardList, CheckCircle2, TrendingUp, User, Lock,
} from "lucide-react";
import { toman, toFa } from "@/lib/fa";
import { useAuth } from "@/lib/auth-context";
import { useRole } from "@/components/role-context";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/profile")({
  head: () => ({ meta: [{ title: "پروفایل — تسک‌بریج" }] }),
  component: ProfileRouter,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Review = { id: string; rating: number; comment: string; createdAt: string; taskId: string };
type WalletData = { availableBalance: number; lockedBalance: number };
type DashStats = { postedTasks?: number; completedTasks?: number; activeTasks?: number; totalApplications?: number };

function ProfileRouter() {
  const { role } = useRole();
  if (role === "admin") return <AdminProfile />;
  if (role === "agent") return <AgentProfile />;
  return <RequesterProfile />;
}

/* ─── SHARED HOOKS ─────────────────────────────────────────────────────── */

function useProfileData() {
  const { user, logout } = useAuth();
  const [reviews, setReviews] = useState<Review[]>([]);
  const [wallet, setWallet] = useState<WalletData | null>(null);
  const [stats, setStats] = useState<DashStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!user) return;
    const token = getToken();
    if (!token) { setLoading(false); return; }
    const h = { Authorization: `Bearer ${token}` };
    Promise.all([
      fetch(`${API_BASE}/v1/users/${user.id}/reviews`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/wallet`, { headers: h }).then(r => r.json()).catch(() => null),
      fetch(`${API_BASE}/v1/dashboard/stats`, { headers: h }).then(r => r.json()).catch(() => null),
    ]).then(([rv, w, s]) => {
      if (rv) setReviews(rv.reviews ?? []);
      if (w && !w.error) setWallet(w);
      if (s) setStats(s);
    }).finally(() => setLoading(false));
  }, [user?.id]);

  const avgRating = reviews.length > 0
    ? reviews.reduce((s, r) => s + r.rating, 0) / reviews.length : 0;

  return { user, logout, reviews, wallet, stats, loading, avgRating };
}

/* ─── EDIT PROFILE FORM ─────────────────────────────────────────────────── */

function useEditProfile() {
  const { user } = useAuth();
  const [editing, setEditing] = useState(false);
  const [fullName, setFullName] = useState(user?.fullName ?? "");
  const [email, setEmail] = useState(user?.email ?? "");
  const [saving, setSaving] = useState(false);

  const save = async () => {
    const token = getToken();
    if (!token) return;
    setSaving(true);
    try {
      const res = await fetch(`${API_BASE}/v1/me`, {
        method: "PATCH",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ fullName, email: email || undefined }),
      });
      if (!res.ok) { const d = await res.json(); throw new Error(d?.error?.message ?? "خطا"); }
      toast.success("پروفایل به‌روزرسانی شد");
      setEditing(false);
    } catch (e: any) {
      toast.error(e.message);
    } finally {
      setSaving(false);
    }
  };

  return { editing, setEditing, fullName, setFullName, email, setEmail, saving, save };
}

/* ─── SHARED COMPONENTS ─────────────────────────────────────────────────── */

function ProfileHeader({ user, editing, setEditing, fullName, setFullName, email, setEmail, saving, save }: any) {
  const initials = user?.fullName?.slice(0, 2) ?? "؟";
  return (
    <div className="relative overflow-hidden rounded-2xl border bg-card shadow-soft">
      <div className="h-32 gradient-brand" />
      <div className="flex flex-wrap items-end gap-4 px-6 pb-6">
        <span className="-mt-12 grid h-24 w-24 shrink-0 place-items-center rounded-2xl border-4 border-card bg-primary/10 font-display text-3xl font-bold text-primary shadow-soft">
          {initials}
        </span>
        <div className="min-w-0 flex-1">
          {editing ? (
            <div className="space-y-2">
              <input value={fullName} onChange={e => setFullName(e.target.value)}
                className="rounded-lg border bg-background px-3 py-1.5 text-sm font-semibold outline-none focus:border-primary w-full max-w-xs" />
              <input value={email} onChange={e => setEmail(e.target.value)} placeholder="ایمیل (اختیاری)"
                className="rounded-lg border bg-background px-3 py-1.5 text-sm outline-none focus:border-primary w-full max-w-xs" dir="ltr" />
            </div>
          ) : (
            <>
              <div className="flex items-center gap-2 flex-wrap">
                <h1 className="text-2xl font-bold tracking-tight">{user?.fullName}</h1>
                <span className="inline-flex items-center gap-1 rounded-full bg-success/15 px-2 py-0.5 text-xs font-medium text-success-foreground">
                  <BadgeCheck className="h-3.5 w-3.5" /> تایید شده
                </span>
              </div>
              <div className="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted-foreground">
                {user?.email && <span className="inline-flex items-center gap-1.5"><Mail className="h-3.5 w-3.5" />{user.email}</span>}
                <span className="inline-flex items-center gap-1.5" dir="ltr"><Phone className="h-3.5 w-3.5" />{user?.phone}</span>
              </div>
            </>
          )}
        </div>
        <div className="flex gap-2">
          {editing ? (
            <>
              <button onClick={save} disabled={saving}
                className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-3.5 py-2 text-sm font-medium text-white hover:opacity-90 disabled:opacity-60">
                <Check className="h-4 w-4" /> {saving ? "ذخیره..." : "ذخیره"}
              </button>
              <button onClick={() => setEditing(false)}
                className="inline-flex items-center gap-1.5 rounded-lg border bg-card px-3.5 py-2 text-sm font-medium hover:bg-accent">
                <X className="h-4 w-4" /> انصراف
              </button>
            </>
          ) : (
            <button onClick={() => setEditing(true)}
              className="inline-flex items-center gap-1.5 rounded-lg border bg-card px-3.5 py-2 text-sm font-medium hover:bg-accent">
              <Pencil className="h-4 w-4" /> ویرایش
            </button>
          )}
        </div>
      </div>
    </div>
  );
}

function ReviewsPanel({ reviews, avgRating, loading }: { reviews: Review[]; avgRating: number; loading: boolean }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-baseline gap-2">
        <h2 className="font-display font-semibold">نظرات دریافتی</h2>
        <span className="text-sm text-muted-foreground">{toFa(reviews.length)} نظر</span>
      </div>
      {reviews.length > 0 && (
        <div className="mt-2 flex items-center gap-2">
          <div className="font-display text-3xl font-bold tabular-nums">{avgRating.toFixed(1)}</div>
          <div className="flex">{[1,2,3,4,5].map(i => (
            <Star key={i} className={cn("h-4 w-4", i <= Math.round(avgRating) ? "fill-warning text-warning" : "text-muted-foreground/30")} />
          ))}</div>
        </div>
      )}
      <div className="mt-4 space-y-3">
        {loading ? (
          [1,2].map(i => <div key={i} className="h-16 rounded-xl bg-muted animate-pulse" />)
        ) : reviews.length === 0 ? (
          <p className="text-sm text-muted-foreground text-center py-6">هنوز نظری ثبت نشده</p>
        ) : (
          reviews.slice(0, 4).map(r => (
            <div key={r.id} className="rounded-xl border p-3">
              <div className="flex items-center justify-between">
                <div className="flex">{[1,2,3,4,5].map(i => (
                  <Star key={i} className={cn("h-3.5 w-3.5", i <= r.rating ? "fill-warning text-warning" : "text-muted-foreground/30")} />
                ))}</div>
                <div className="text-xs text-muted-foreground">{r.createdAt?.slice(0, 10)}</div>
              </div>
              {r.comment && <p className="mt-1.5 text-sm text-muted-foreground">«{r.comment}»</p>}
            </div>
          ))
        )}
      </div>
    </div>
  );
}

/* ─── REQUESTER PROFILE ─────────────────────────────────────────────────── */

function RequesterProfile() {
  const { user, reviews, wallet, stats, loading, avgRating } = useProfileData();
  const edit = useEditProfile();

  const statCards = [
    { label: "درخواست‌های ثبت‌شده",  value: loading ? "..." : toFa(stats?.postedTasks ?? 0),    icon: ClipboardList },
    { label: "تکمیل‌شده‌ها",          value: loading ? "..." : toFa(stats?.completedTasks ?? 0), icon: CheckCircle2  },
    { label: "موجودی کیف پول",         value: loading ? "..." : toman(wallet?.availableBalance ?? 0, false), unit: "تومان", icon: Wallet },
    { label: "امتیاز دریافتی",        value: loading ? "..." : (avgRating ? avgRating.toFixed(1) : "—"),      icon: Star },
  ];

  return (
    <div className="space-y-6">
      <ProfileHeader user={user} {...edit} />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map(c => <StatCard key={c.label} {...c} />)}
      </div>
      <div className="grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2">
          <AccountInfoCard user={user} role="کارفرما" wallet={wallet} loading={loading} />
        </div>
        <ReviewsPanel reviews={reviews} avgRating={avgRating} loading={loading} />
      </div>
    </div>
  );
}

/* ─── AGENT PROFILE ─────────────────────────────────────────────────────── */

function AgentProfile() {
  const { user, reviews, wallet, stats, loading, avgRating } = useProfileData();
  const edit = useEditProfile();

  const statCards = [
    { label: "کارهای تکمیل‌شده",     value: loading ? "..." : toFa(user?.completedCount ?? 0),   icon: CheckCircle2  },
    { label: "درخواست‌های ارسالی",   value: loading ? "..." : toFa(stats?.totalApplications ?? 0), icon: ClipboardList },
    { label: "موجودی کیف پول",        value: loading ? "..." : toman(wallet?.availableBalance ?? 0, false), unit: "تومان", icon: Wallet },
    { label: "امتیاز میانگین",       value: loading ? "..." : (avgRating ? avgRating.toFixed(1) : "—"),       icon: Star },
  ];

  return (
    <div className="space-y-6">
      <ProfileHeader user={user} {...edit} />
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map(c => <StatCard key={c.label} {...c} />)}
      </div>
      <div className="grid gap-6 lg:grid-cols-3">
        <div className="lg:col-span-2 space-y-4">
          <AccountInfoCard user={user} role="مجری" wallet={wallet} loading={loading} />
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <h3 className="font-display font-semibold">آمار کارها</h3>
            <div className="mt-4 grid grid-cols-2 gap-4 text-sm">
              <div className="rounded-xl bg-muted/40 p-4">
                <div className="text-muted-foreground">کارهای تکمیل‌شده</div>
                <div className="mt-1 font-display text-2xl font-bold">{toFa(user?.completedCount ?? 0)}</div>
              </div>
              <div className="rounded-xl bg-muted/40 p-4">
                <div className="text-muted-foreground">امتیاز کاربری</div>
                <div className="mt-1 font-display text-2xl font-bold flex items-center gap-1">
                  {avgRating ? avgRating.toFixed(1) : "—"}
                  <Star className="h-5 w-5 fill-warning text-warning" />
                </div>
              </div>
              <div className="rounded-xl bg-muted/40 p-4">
                <div className="text-muted-foreground">موجودی قابل برداشت</div>
                <div className="mt-1 font-display text-lg font-bold">{toman(wallet?.availableBalance ?? 0)}</div>
              </div>
              <div className="rounded-xl bg-muted/40 p-4">
                <div className="text-muted-foreground">در انتظار آزادسازی</div>
                <div className="mt-1 font-display text-lg font-bold">{toman(wallet?.lockedBalance ?? 0)}</div>
              </div>
            </div>
          </div>
        </div>
        <ReviewsPanel reviews={reviews} avgRating={avgRating} loading={loading} />
      </div>
    </div>
  );
}

/* ─── ADMIN PROFILE ─────────────────────────────────────────────────────── */

function AdminProfile() {
  const { user } = useAuth();
  const edit = useEditProfile();

  return (
    <div className="space-y-6">
      <ProfileHeader user={user} {...edit} />
      <div className="grid gap-6 lg:grid-cols-2">
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <h3 className="font-display font-semibold">اطلاعات حساب</h3>
          <dl className="mt-4 divide-y text-sm">
            <InfoRow icon={User}       label="نام کامل"    value={user?.fullName ?? "—"} />
            <InfoRow icon={Phone}      label="تلفن"        value={user?.phone ?? "—"} ltr />
            <InfoRow icon={Mail}       label="ایمیل"       value={user?.email ?? "ثبت نشده"} ltr />
            <InfoRow icon={ShieldCheck} label="نقش"        value="مدیر سیستم" />
            <InfoRow icon={BadgeCheck} label="وضعیت"       value="فعال" />
          </dl>
        </div>
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <h3 className="font-display font-semibold">امنیت حساب</h3>
          <div className="mt-4 space-y-3">
            <SecurityItem label="رمز عبور" status="تنظیم شده" icon={Lock} ok />
            <SecurityItem label="تایید تلفن" status="تایید شده" icon={Phone} ok />
            <SecurityItem label="تایید ایمیل" status={user?.email ? "تایید شده" : "تنظیم نشده"} icon={Mail} ok={!!user?.email} />
          </div>
          <button className="mt-5 w-full rounded-lg border py-2.5 text-sm font-medium hover:bg-accent">
            تغییر رمز عبور
          </button>
        </div>
      </div>
    </div>
  );
}

/* ─── SMALL REUSABLE PIECES ─────────────────────────────────────────────── */

function AccountInfoCard({ user, role, wallet, loading }: any) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <h3 className="font-display font-semibold">اطلاعات حساب</h3>
      <dl className="mt-4 divide-y text-sm">
        <InfoRow icon={User}  label="نام کامل"  value={user?.fullName ?? "—"} />
        <InfoRow icon={Phone} label="تلفن"       value={user?.phone ?? "—"} ltr />
        <InfoRow icon={Mail}  label="ایمیل"      value={user?.email ?? "ثبت نشده"} ltr />
        <InfoRow icon={ShieldCheck} label="نقش"  value={role} />
        <InfoRow icon={Wallet} label="موجودی"
          value={loading ? "در حال بارگذاری..." : toman(wallet?.availableBalance ?? 0)} />
      </dl>
    </div>
  );
}

function StatCard({ label, value, unit, icon: Icon }: { label: string; value: string; unit?: string; icon: any }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-start justify-between">
        <span className="text-sm text-muted-foreground">{label}</span>
        <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
      </div>
      <div className="mt-3 font-display text-xl font-bold tabular-nums">
        {value}{unit && <span className="ms-1 text-sm font-medium text-muted-foreground">{unit}</span>}
      </div>
    </div>
  );
}

function InfoRow({ icon: Icon, label, value, ltr }: { icon: any; label: string; value: string; ltr?: boolean }) {
  return (
    <div className="flex items-center gap-3 py-3">
      <Icon className="h-4 w-4 shrink-0 text-muted-foreground" />
      <span className="w-28 text-muted-foreground">{label}</span>
      <span className={cn("flex-1 font-medium", ltr && "text-left")} dir={ltr ? "ltr" : undefined}>{value}</span>
    </div>
  );
}

function SecurityItem({ label, status, icon: Icon, ok }: { label: string; status: string; icon: any; ok: boolean }) {
  return (
    <div className="flex items-center justify-between rounded-lg border p-3">
      <div className="flex items-center gap-3">
        <Icon className="h-4 w-4 text-muted-foreground" />
        <span className="text-sm font-medium">{label}</span>
      </div>
      <span className={cn("inline-flex items-center gap-1 rounded-full px-2 py-0.5 text-xs font-medium",
        ok ? "bg-success/10 text-success" : "bg-muted text-muted-foreground")}>
        {ok ? <CheckCircle2 className="h-3 w-3" /> : <X className="h-3 w-3" />}
        {status}
      </span>
    </div>
  );
}
