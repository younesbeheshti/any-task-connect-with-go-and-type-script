import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Phone, Lock, Eye, EyeOff, ArrowLeft } from "lucide-react";
import { Logo } from "@/components/logo";
import { toast } from "sonner";
import { toFa } from "@/lib/fa";

export const Route = createFileRoute("/login")({
  head: () => ({ meta: [{ title: "ورود — تسک‌بریج" }] }),
  component: Login,
});

function Login() {
  const [show, setShow] = useState(false);
  const nav = useNavigate();
  return (
    <div className="grid min-h-screen lg:grid-cols-2">
      <div className="flex flex-col px-6 py-8 sm:px-10">
        <Logo />
        <div className="mx-auto flex w-full max-w-sm flex-1 flex-col justify-center">
          <h1 className="text-3xl font-bold tracking-tight">خوش آمدید</h1>
          <p className="mt-2 text-sm text-muted-foreground">برای مدیریت درخواست‌ها و انجام‌دهنده‌های خود وارد شوید.</p>
          <form
            onSubmit={(e) => { e.preventDefault(); toast.success("با موفقیت وارد شدید"); nav({ to: "/app" }); }}
            className="mt-8 space-y-4"
          >
            <Field icon={Phone} type="tel" label="شماره موبایل" placeholder={toFa("09120000000")} dir="ltr" />
            <Field
              icon={Lock}
              type={show ? "text" : "password"}
              label="رمز عبور"
              placeholder="••••••••"
              trailing={
                <button type="button" onClick={() => setShow(s => !s)} className="text-muted-foreground hover:text-foreground">
                  {show ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                </button>
              }
            />
            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center gap-2 text-muted-foreground">
                <input type="checkbox" className="h-4 w-4 rounded border-input text-primary focus:ring-primary" />
                مرا به خاطر بسپار
              </label>
              <a href="#" className="font-medium text-primary hover:underline">فراموشی رمز عبور؟</a>
            </div>
            <button className="inline-flex w-full items-center justify-center gap-2 rounded-lg gradient-brand py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95">
              ورود <ArrowLeft className="h-4 w-4" />
            </button>
          </form>
          <div className="mt-6 text-center text-sm text-muted-foreground">
            تازه به تسک‌بریج پیوسته‌اید؟{" "}
            <Link to="/register" className="font-medium text-primary hover:underline">ساخت حساب کاربری</Link>
          </div>
        </div>
      </div>
      <AuthAside />
    </div>
  );
}

export function Field({ icon: Icon, label, trailing, ...props }: any) {
  return (
    <label className="block">
      <span className="text-sm font-medium">{label}</span>
      <div className="mt-1.5 flex items-center gap-2 rounded-lg border bg-card px-3 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20">
        {Icon && <Icon className="h-4 w-4 shrink-0 text-muted-foreground" />}
        <input {...props} className="w-full bg-transparent py-2.5 text-sm outline-none placeholder:text-muted-foreground" />
        {trailing}
      </div>
    </label>
  );
}

export function AuthAside() {
  return (
    <div className="relative hidden overflow-hidden bg-muted lg:block">
      <div className="absolute inset-0 gradient-brand opacity-95" />
      <div className="absolute inset-0 grid-bg opacity-20" />
      <div className="relative flex h-full flex-col justify-between p-12 text-white">
        <div className="flex items-center gap-2 text-sm font-medium opacity-80">
          <span className="h-1.5 w-1.5 rounded-full bg-white" />
          فعال در {toFa(10)} شهر بزرگ ایران
        </div>
        <div>
          <p className="font-display text-3xl font-bold leading-relaxed">«تسک‌بریج مثل داشتن یه آشنای مطمئن توی هر شهره.»</p>
          <p className="mt-6 text-sm opacity-80">— لیلا هاشمی، طراح محصول</p>
        </div>
        <div className="grid grid-cols-3 gap-6 border-t border-white/15 pt-8 text-sm">
          <div><div className="font-display text-2xl font-bold tabular-nums">{toFa("12k+")}</div><div className="opacity-70">درخواست انجام‌شده</div></div>
          <div><div className="font-display text-2xl font-bold tabular-nums">{toFa("4.9")}</div><div className="opacity-70">میانگین امتیاز</div></div>
          <div><div className="font-display text-2xl font-bold tabular-nums">۱۰۰٪</div><div className="opacity-70">تضمین بازگشت وجه</div></div>
        </div>
      </div>
    </div>
  );
}
