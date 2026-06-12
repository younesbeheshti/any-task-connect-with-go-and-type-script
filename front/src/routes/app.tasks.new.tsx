import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Calendar, MapPin, Tag, Wallet, FileText, Upload, Info, ShieldCheck } from "lucide-react";
import { categories, cities } from "@/lib/mock-data";
import { toast } from "sonner";
import { toman, toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/tasks/new")({
  head: () => ({ meta: [{ title: "ثبت درخواست جدید — تسک‌بریج" }] }),
  component: NewTask,
});

function NewTask() {
  const [title, setTitle] = useState("");
  const [desc, setDesc] = useState("");
  const [cat, setCat] = useState(categories[0].label);
  const [city, setCity] = useState(cities[0]);
  const [budget, setBudget] = useState(750000);
  const [deadline, setDeadline] = useState("");
  const fee = Math.round(budget * 0.08);
  const escrow = budget + fee;
  const nav = useNavigate();

  return (
    <div className="grid gap-6 lg:grid-cols-[1fr_360px]">
      <form
        onSubmit={(e) => { e.preventDefault(); toast.success("درخواست شما با موفقیت ثبت شد"); nav({ to: "/app/tasks" }); }}
        className="space-y-6 rounded-2xl border bg-card p-6 shadow-soft"
      >
        <header>
          <h1 className="text-2xl font-bold tracking-tight">ثبت درخواست جدید</h1>
          <p className="text-sm text-muted-foreground">جزئیات کامل بدهید تا انجام‌دهنده‌ها بتوانند بهترین پیشنهاد را ارائه دهند.</p>
        </header>

        <Field label="عنوان درخواست" icon={FileText}>
          <input value={title} onChange={e => setTitle(e.target.value)} required placeholder="مثلاً: دریافت ریزنمرات از دانشگاه تهران" className="w-full bg-transparent py-2.5 text-sm outline-none" />
        </Field>

        <div>
          <label className="text-sm font-medium">شرح درخواست</label>
          <textarea
            value={desc} onChange={e => setDesc(e.target.value)}
            placeholder="توضیح دهید چه کاری باید انجام شود، در کجا، و چه مدارکی لازم است…"
            rows={5}
            className="mt-1.5 w-full rounded-lg border bg-background p-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20"
          />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <SelectField label="دسته‌بندی" icon={Tag} value={cat} onChange={setCat} options={categories.map(c => c.label)} />
          <SelectField label="شهر" icon={MapPin} value={city} onChange={setCity} options={cities} />
        </div>

        <div className="grid gap-4 sm:grid-cols-2">
          <div>
            <label className="text-sm font-medium">مبلغ پیشنهادی (تومان)</label>
            <div className="mt-1.5 flex items-center gap-3 rounded-lg border bg-background px-3">
              <Wallet className="h-4 w-4 text-muted-foreground" />
              <input type="range" min={100000} max={5000000} step={50000} value={budget} onChange={e => setBudget(+e.target.value)} className="flex-1 accent-[var(--color-primary)]" />
              <span className="w-32 text-left text-xs font-semibold tabular-nums">{toman(budget)}</span>
            </div>
          </div>
          <Field label="مهلت انجام" icon={Calendar}>
            <input type="date" value={deadline} onChange={e => setDeadline(e.target.value)} className="w-full bg-transparent py-2.5 text-sm outline-none" />
          </Field>
        </div>

        <div>
          <label className="text-sm font-medium">پیوست‌ها</label>
          <label className="mt-1.5 flex cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border border-dashed bg-background py-8 text-sm text-muted-foreground hover:border-primary/40 hover:bg-accent">
            <Upload className="h-5 w-5" />
            <span>فایل‌ها را اینجا رها کنید یا کلیک کنید</span>
            <span className="text-xs">PDF، JPG، PNG — حداکثر {toFa(10)} مگابایت</span>
            <input type="file" multiple className="hidden" />
          </label>
        </div>

        <div className="flex justify-end gap-2">
          <button type="button" className="rounded-lg border bg-background px-4 py-2 text-sm font-medium hover:bg-accent">ذخیره پیش‌نویس</button>
          <button className="inline-flex items-center gap-2 rounded-lg gradient-brand px-5 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            ثبت درخواست — {toman(escrow)}
          </button>
        </div>
      </form>

      <aside className="space-y-4">
        <div className="rounded-2xl border bg-card p-5 shadow-soft">
          <div className="flex items-center gap-2">
            <ShieldCheck className="h-4 w-4 text-success" />
            <h3 className="font-semibold">پیش‌نمایش پرداخت امانی</h3>
          </div>
          <dl className="mt-4 space-y-2 text-sm">
            <Row k="مبلغ درخواست" v={toman(budget)} />
            <Row k="کارمزد سامانه (۸٪)" v={toman(fee)} />
            <div className="my-2 h-px bg-border" />
            <Row k="مبلغ بلوکه‌شده در امانت" v={toman(escrow)} bold />
          </dl>
          <p className="mt-3 text-xs text-muted-foreground">مبلغ تنها پس از تایید نهایی شما به انجام‌دهنده پرداخت می‌شود. در صورت لغو پیش از پذیرش، بازگشت کامل وجه.</p>
        </div>

        <div className="rounded-2xl border bg-primary/5 p-5 text-sm">
          <div className="flex items-center gap-2 font-semibold text-primary"><Info className="h-4 w-4" /> زمان تخمینی انجام</div>
          <div className="mt-2 font-display text-2xl font-bold">~ {toFa(24)} ساعت</div>
          <p className="mt-1 text-xs text-muted-foreground">بر اساس درخواست‌های مشابه در {city}.</p>
        </div>

        <div className="rounded-2xl border bg-card p-5 text-sm">
          <h4 className="font-semibold">نکات افزایش شانس پذیرش</h4>
          <ul className="mt-3 space-y-2 text-xs text-muted-foreground list-disc pr-4">
            <li>عنوان روشن و دقیق بنویسید</li>
            <li>تمام مدارک لازم را پیوست کنید</li>
            <li>مبلغ منصفانه‌ای پیشنهاد دهید</li>
            <li>مهلت معقول تعیین کنید</li>
          </ul>
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
      <dd className="tabular-nums">{v}</dd>
    </div>
  );
}
