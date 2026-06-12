import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Phone, Lock, User, ArrowLeft, Briefcase, ClipboardList } from "lucide-react";
import { Logo } from "@/components/logo";
import { AuthAside, Field } from "./login";
import { useRole } from "@/components/role-context";
import { toast } from "sonner";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/register")({
  head: () => ({ meta: [{ title: "ثبت‌نام — تسک‌بریج" }] }),
  component: Register,
});

function Register() {
  const [role, setRoleLocal] = useState<"requester" | "agent">("requester");
  const { setRole } = useRole();
  const nav = useNavigate();
  return (
    <div className="grid min-h-screen lg:grid-cols-2">
      <div className="flex flex-col px-6 py-8 sm:px-10">
        <Logo />
        <div className="mx-auto flex w-full max-w-md flex-1 flex-col justify-center">
          <h1 className="text-3xl font-bold tracking-tight">ساخت حساب کاربری</h1>
          <p className="mt-2 text-sm text-muted-foreground">انتخاب کنید چگونه می‌خواهید از تسک‌بریج استفاده کنید.</p>

          <div className="mt-6 grid grid-cols-2 gap-3">
            {([
              { id: "requester", icon: ClipboardList, title: "می‌خوام کارم انجام بشه", desc: "ثبت درخواست و استخدام انجام‌دهنده." },
              { id: "agent",     icon: Briefcase,     title: "می‌خوام انجام‌دهنده باشم", desc: "پیدا کردن درخواست و کسب درآمد." },
            ] as const).map(opt => (
              <button
                key={opt.id}
                onClick={() => setRoleLocal(opt.id)}
                className={cn(
                  "rounded-xl border bg-card p-4 text-right transition-all",
                  role === opt.id ? "border-primary shadow-glow" : "hover:border-primary/40"
                )}
              >
                <opt.icon className={cn("h-5 w-5", role === opt.id ? "text-primary" : "text-muted-foreground")} />
                <div className="mt-3 text-sm font-semibold">{opt.title}</div>
                <div className="text-xs text-muted-foreground">{opt.desc}</div>
              </button>
            ))}
          </div>

          <form
            onSubmit={(e) => { e.preventDefault(); setRole(role); toast.success("حساب شما ساخته شد"); nav({ to: "/app" }); }}
            className="mt-6 space-y-4"
          >
            <Field icon={User}  label="نام و نام خانوادگی" placeholder="مثلاً: سارا محمدی" />
            <Field icon={Phone} type="tel" label="شماره موبایل" placeholder={toFa("09120000000")} dir="ltr" />
            <Field icon={Lock}  type="password" label="رمز عبور" placeholder={`حداقل ${toFa(8)} کاراکتر`} />
            <p className="text-xs text-muted-foreground">با ادامه ثبت‌نام، با قوانین و حریم خصوصی موافقت می‌کنید.</p>
            <button className="inline-flex w-full items-center justify-center gap-2 rounded-lg gradient-brand py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95">
              ساخت حساب <ArrowLeft className="h-4 w-4" />
            </button>
          </form>
          <div className="mt-6 text-center text-sm text-muted-foreground">
            قبلاً حساب ساخته‌اید؟{" "}
            <Link to="/login" className="font-medium text-primary hover:underline">ورود</Link>
          </div>
        </div>
      </div>
      <AuthAside />
    </div>
  );
}
