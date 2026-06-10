import { createFileRoute } from "@tanstack/react-router";
import { BadgeCheck, MapPin, Star, Wallet, ShieldCheck, Mail, Phone, Pencil } from "lucide-react";
import { tasks } from "@/lib/mock-data";
import { StatusBadge } from "@/components/status-badge";

export const Route = createFileRoute("/app/profile")({
  head: () => ({ meta: [{ title: "Profile — TaskBridge" }] }),
  component: Profile,
});

const reviews = [
  { id: "r1", from: "Karim N.", rating: 5, text: "Crystal clear instructions and quick to verify. A dream client.", time: "3d ago" },
  { id: "r2", from: "Aylin O.", rating: 5, text: "Communication was excellent throughout the task.", time: "1w ago" },
  { id: "r3", from: "Mahdi T.", rating: 4, text: "Great task, would happily work with again.", time: "2w ago" },
];

function Profile() {
  return (
    <div className="space-y-6">
      <div className="relative overflow-hidden rounded-2xl border bg-card shadow-soft">
        <div className="h-32 gradient-brand" />
        <div className="flex flex-wrap items-end gap-4 px-6 pb-6">
          <span className="-mt-12 grid h-24 w-24 shrink-0 place-items-center rounded-2xl border-4 border-card bg-card font-display text-3xl font-bold gradient-text shadow-soft">SM</span>
          <div className="min-w-0 flex-1">
            <div className="flex items-center gap-2">
              <h1 className="text-2xl font-bold tracking-tight">Sara Moradi</h1>
              <span className="inline-flex items-center gap-1 rounded-full bg-success/15 px-2 py-0.5 text-xs font-medium text-success-foreground">
                <BadgeCheck className="h-3.5 w-3.5" /> Verified
              </span>
            </div>
            <div className="mt-1 flex flex-wrap items-center gap-x-4 gap-y-1 text-sm text-muted-foreground">
              <span className="inline-flex items-center gap-1.5"><MapPin className="h-3.5 w-3.5" />Berlin, Germany</span>
              <span className="inline-flex items-center gap-1.5"><Mail className="h-3.5 w-3.5" />sara@example.com</span>
              <span className="inline-flex items-center gap-1.5"><Phone className="h-3.5 w-3.5" />+49 30 1234 5678</span>
            </div>
          </div>
          <button className="inline-flex items-center gap-1.5 rounded-lg border bg-card px-3.5 py-2 text-sm font-medium hover:bg-accent">
            <Pencil className="h-4 w-4" /> Edit profile
          </button>
        </div>
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-4">
        <Stat label="Rating" value="4.9" icon={Star} sub="From 47 reviews" />
        <Stat label="Completed Tasks" value="27" icon={BadgeCheck} sub="98% success rate" />
        <Stat label="Wallet" value="$1,240" icon={Wallet} sub="$365 in escrow" />
        <Stat label="Verification" value="Full" icon={ShieldCheck} sub="ID + Phone + Email" />
      </div>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="rounded-2xl border bg-card p-5 shadow-soft lg:col-span-2">
          <h2 className="font-display font-semibold">Activity history</h2>
          <div className="mt-4 divide-y">
            {tasks.slice(0, 5).map(t => (
              <div key={t.id} className="flex items-center gap-4 py-3">
                <div className="min-w-0 flex-1">
                  <div className="truncate text-sm font-medium">{t.title}</div>
                  <div className="text-xs text-muted-foreground">{t.id} · {t.city} · {t.postedAgo}</div>
                </div>
                <StatusBadge status={t.status} />
                <div className="hidden text-sm font-semibold text-foreground sm:block">${t.budget}</div>
              </div>
            ))}
          </div>
        </div>

        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <div className="flex items-baseline gap-2">
            <h2 className="font-display font-semibold">Reviews</h2>
            <span className="text-sm text-muted-foreground">47 total</span>
          </div>
          <div className="mt-2 flex items-center gap-2">
            <div className="font-display text-3xl font-bold">4.9</div>
            <div className="flex">{Array.from({length:5}).map((_,i) => <Star key={i} className="h-4 w-4 fill-warning text-warning" />)}</div>
          </div>
          <div className="mt-4 space-y-4">
            {reviews.map(r => (
              <div key={r.id} className="rounded-xl border p-3">
                <div className="flex items-center justify-between">
                  <div className="text-sm font-semibold">{r.from}</div>
                  <div className="text-xs text-muted-foreground">{r.time}</div>
                </div>
                <div className="mt-1 flex">{Array.from({length:r.rating}).map((_,i) => <Star key={i} className="h-3.5 w-3.5 fill-warning text-warning" />)}</div>
                <p className="mt-1.5 text-sm text-muted-foreground">"{r.text}"</p>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}

function Stat({ label, value, icon: Icon, sub }: any) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-start justify-between">
        <span className="text-sm text-muted-foreground">{label}</span>
        <span className="grid h-9 w-9 place-items-center rounded-lg bg-primary/10 text-primary"><Icon className="h-4 w-4" /></span>
      </div>
      <div className="mt-3 font-display text-2xl font-bold">{value}</div>
      <div className="mt-0.5 text-xs text-muted-foreground">{sub}</div>
    </div>
  );
}
