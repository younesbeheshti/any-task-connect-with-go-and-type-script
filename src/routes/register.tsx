import { createFileRoute, Link, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { Mail, Lock, User, ArrowRight, Briefcase, ClipboardList } from "lucide-react";
import { Logo } from "@/components/logo";
import { AuthAside, Field } from "./login";
import { toast } from "sonner";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/register")({
  head: () => ({ meta: [{ title: "Create your account — TaskBridge" }] }),
  component: Register,
});

function Register() {
  const [role, setRole] = useState<"requester" | "agent">("requester");
  const nav = useNavigate();
  return (
    <div className="grid min-h-screen lg:grid-cols-2">
      <div className="flex flex-col px-6 py-8 sm:px-10">
        <Logo />
        <div className="mx-auto flex w-full max-w-md flex-1 flex-col justify-center">
          <h1 className="text-3xl font-bold tracking-tight">Create your account</h1>
          <p className="mt-2 text-sm text-muted-foreground">Choose how you'd like to use TaskBridge.</p>

          <div className="mt-6 grid grid-cols-2 gap-3">
            {([
              { id: "requester", icon: ClipboardList, title: "I need tasks done", desc: "Post tasks and hire agents." },
              { id: "agent", icon: Briefcase, title: "I'm an agent", desc: "Find tasks and earn money." },
            ] as const).map(opt => (
              <button
                key={opt.id}
                onClick={() => setRole(opt.id)}
                className={cn(
                  "rounded-xl border bg-card p-4 text-left transition-all",
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
            onSubmit={(e) => { e.preventDefault(); toast.success("Account created"); nav({ to: "/app" }); }}
            className="mt-6 space-y-4"
          >
            <Field icon={User} label="Full name" placeholder="Jane Doe" />
            <Field icon={Mail} type="email" label="Email" placeholder="you@example.com" />
            <Field icon={Lock} type="password" label="Password" placeholder="At least 8 characters" />
            <p className="text-xs text-muted-foreground">By continuing, you agree to our Terms and Privacy Policy.</p>
            <button className="inline-flex w-full items-center justify-center gap-2 rounded-lg gradient-brand py-2.5 text-sm font-semibold text-white shadow-soft hover:opacity-95">
              Create account <ArrowRight className="h-4 w-4" />
            </button>
          </form>
          <div className="mt-6 text-center text-sm text-muted-foreground">
            Already have an account?{" "}
            <Link to="/login" className="font-medium text-primary hover:underline">Sign in</Link>
          </div>
        </div>
      </div>
      <AuthAside />
    </div>
  );
}
