import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Phone, Lock, User, Mail, ArrowLeft, Briefcase, ClipboardList } from "lucide-react";
import { Logo } from "@/components/logo";
import { AuthAside, Field, roleHome } from "./login";
import { useRole } from "@/components/role-context";
import { useAuth } from "@/lib/auth-context";
import { toast } from "sonner";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/register")({
  head: () => ({ meta: [{ title: "ثبت‌نام — تسک‌بریج" }] }),
  component: Register,
});

function Register() {
  const [selectedRole, setSelectedRole] = useState<"requester" | "agent">("requester");
  const [fullName, setFullName] = useState("");
  const [phone, setPhone] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const { register } = useAuth();
  const { setRole } = useRole();
  const nav = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!fullName || !phone || !password) { toast.error("نام، شماره موبایل و رمز عبور الزامی است"); return; }
    setSubmitting(true);
    try {
      const u = await register({
        fullName,
        phone,
        password,
        role: selectedRole,
        ...(email ? { email } : {}),
      });
      setRole(u.role);
      toast.success("حساب شما ساخته شد");
      nav({ to: roleHome(u) as any });
    } catch (err: any) {
      toast.error(err?.message ?? "ثبت‌نام ناموفق بود");
    } finally {
      setSubmitting(false);
    }
  };

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
                type="button"
                onClick={() => setSelectedRole(opt.id)}
                className={cn(
                  "rounded-xl border bg-card p-4 text-right transition-all",
                  selectedRole === opt.id ? "border-primary shadow-glow" : "hover:border-primary/40"
                )}
              >
                <opt.icon className={cn("h-5 w-5", selectedRole === opt.id ? "text-primary" : "text-muted-foreground")} />
                <div className="mt-3 text-sm font-semibold">{opt.title}</div>
                <div className="text-xs text-muted-foreground">{opt.desc}</div>
              </button>
            ))}
          </div>

          <form onSubmit={handleSubmit} className="mt-6 space-y-4">
            <Field
              icon={User}
              label="نام و نام خانوادگی"
              placeholder="مثلاً: علی رضایی"
              value={fullName}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setFullName(e.target.value)}
            />
            <Field
              icon={Phone}
              type="tel"
              label="شماره موبایل"
              placeholder={toFa("09120000000")}
              dir="ltr"
              value={phone}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPhone(e.target.value)}
            />
            <Field
              icon={Mail}
              type="email"
              label={<span>ایمیل <span className="text-muted-foreground font-normal">(اختیاری)</span></span>}
              placeholder="example@email.com"
              dir="ltr"
              value={email}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setEmail(e.target.value)}
            />
            <Field
              icon={Lock}
              type="password"
              label="رمز عبور"
              placeholder={`حداقل ${toFa(8)} کاراکتر`}
              value={password}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => setPassword(e.target.value)}
            />
            <p className="text-xs text-muted-foreground">با ادامه ثبت‌نام، با قوانین و حریم خصوصی موافقت می‌کنید.</p>
            <button
              disabled={submitting}
              className="inline-flex w-full items-center justify-center gap-2 rounded-lg gradient-brand py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95 disabled:opacity-60"
            >
              {submitting ? "در حال ساخت حساب..." : <>ساخت حساب <ArrowLeft className="h-4 w-4" /></>}
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
