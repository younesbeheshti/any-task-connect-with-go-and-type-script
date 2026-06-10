import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Calendar, MapPin, Tag, Wallet, FileText, Upload, Info, ShieldCheck } from "lucide-react";
import { categories, cities } from "@/lib/mock-data";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/tasks/new")({
  head: () => ({ meta: [{ title: "Create Task — TaskBridge" }] }),
  component: NewTask,
});

function NewTask() {
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [cat, setCat] = useState(categories[0].label);
  const [city, setCity] = useState(cities[0]);
  const [budget, setBudget] = useState(75);
  const [deadline, setDeadline] = useState("");
  const fee = Math.round(budget * 0.08);
  const escrow = budget + fee;
  const nav = useNavigate();

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_360px]">
      <form
        onSubmit={(e) => { e.preventDefault(); toast.success("Task posted to marketplace"); nav({ to: "/app/tasks" }); }}
        className="space-y-6 rounded-2xl border bg-card p-6 shadow-soft"
      >
        <header>
          <h1 className="text-2xl font-bold tracking-tight">Create a new task</h1>
          <p className="text-sm text-muted-foreground">Give agents the details they need to deliver perfectly.</p>
        </header>

        <Field label="Title" icon={FileText}>
          <input value={title} onChange={e => setTitle(e.target.value)} required placeholder="e.g. Pick up university transcripts" className="w-full bg-transparent py-2.5 text-sm outline-none" />
        </Field>

        <div>
          <label className="text-sm font-medium">Description</label>
          <textarea
            value={desc} onChange={e => setDesc(e.target.value)}
            placeholder="Describe what needs to happen, where, and any documents involved…"
            rows={5}
            className="mt-1.5 w-full rounded-lg border bg-background p-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20"
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <SelectField label="Category" icon={Tag} value={cat} onChange={setCat} options={categories.map(c => c.label)} />
          <SelectField label="City" icon={MapPin} value={city} onChange={setCity} options={cities} />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label className="text-sm font-medium">Budget (USD)</label>
            <div className="mt-1.5 flex items-center gap-3 rounded-lg border bg-background px-3">
              <Wallet className="h-4 w-4 text-muted-foreground" />
              <input type="range" min={10} max={500} step={5} value={budget} onChange={e => setBudget(+e.target.value)} className="flex-1 accent-[var(--color-primary)]" />
              <span className="w-16 text-right text-sm font-semibold tabular-nums">${budget}</span>
            </div>
          </div>
          <Field label="Deadline" icon={Calendar}>
            <input type="date" value={deadline} onChange={e => setDeadline(e.target.value)} className="w-full bg-transparent py-2.5 text-sm outline-none" />
          </Field>
        </div>

        <div>
          <label className="text-sm font-medium">Attachments</label>
          <label className="mt-1.5 flex cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border border-dashed bg-background py-8 text-sm text-muted-foreground hover:border-primary/40 hover:bg-accent">
            <Upload className="h-5 w-5" />
            <span>Drop files here or click to upload</span>
            <span className="text-xs">PDF, JPG, PNG — up to 10MB each</span>
            <input type="file" multiple className="hidden" />
          </label>
        </div>

        <div className="flex justify-end gap-2">
          <button type="button" className="rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent">Save draft</button>
          <button className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            Post task — ${escrow}
          </button>
        </div>
      </form>

      <aside className="space-y-4">
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <div className="flex items-center gap-2">
            <ShieldCheck className="h-4 w-4 text-success" />
            <h3 className="font-semibold">Escrow preview</h3>
          </div>
          <dl className="mt-4 space-y-2 text-sm">
            <Row k="Task budget" v={`$${budget}`} />
            <Row k="Service fee (8%)" v={`$${fee}`} />
            <div className="my-2 h-px bg-border" />
            <Row k="Held in escrow" v={`$${escrow}`} bold />
          </dl>
          <p className="mt-3 text-xs text-muted-foreground">Funds release to your agent only after you verify completion. Refundable if cancelled before assignment.</p>
        </div>

        <div className="rounded-2xl border bg-primary/5 p-5 text-sm">
          <div className="flex items-center gap-2 font-semibold text-primary"><Info className="h-4 w-4" /> Estimated completion</div>
          <div className="mt-2 font-display text-2xl font-bold">~24 hours</div>
          <p className="mt-1 text-xs text-muted-foreground">Based on similar tasks in {city}.</p>
        </div>
      </aside>
    </div>
  );
}

function Field({ label, icon: Icon, children }: any) {
  return (
    <div>
      <label className="text-sm font-medium">{label}</label>
      <div className="mt-1.5 flex items-center gap-2 rounded-lg border bg-background px-3 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20">
        <Icon className="h-4 w-4 text-muted-foreground" />
        {children}
      </div>
    </div>
  );
}

function SelectField({ label, icon: Icon, value, onChange, options }: any) {
  return (
    <Field label={label} icon={Icon}>
      <select value={value} onChange={e => onChange(e.target.value)} className="w-full bg-transparent py-2.5 text-sm outline-none">
        {options.map((o: string) => <option key={o}>{o}</option>)}
      </select>
    </Field>
  );
}

function Row({ k, v, bold }: { k: string; v: string; bold?: boolean }) {
  return (
    <div className={cn("flex items-center justify-between", bold && "font-semibold")}>
      <dt className="text-muted-foreground">{k}</dt>
      <dd>{v}</dd>
    </div>
  );
}
