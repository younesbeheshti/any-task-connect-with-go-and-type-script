import { createFileRoute, Link, Outlet, useRouterState } from "@tanstack/react-router";
import {
  LayoutDashboard, PlusCircle, Wallet, MessagesSquare, Bell, User, Settings, Search, Moon, Sun,
  Shield, Menu, X, Briefcase, ClipboardList, CheckCircle2, Activity, Sparkles, Users, BarChart3, FileBarChart, ArrowLeftRight,
} from "lucide-react";
import { Logo } from "@/components/logo";
import { useTheme } from "@/components/theme-provider";
import { useRole, ROLE_LABELS, type Role } from "@/components/role-context";
import { cn } from "@/lib/utils";
import { useState } from "react";

export const Route = createFileRoute("/app")({
  head: () => ({ meta: [{ title: "داشبورد — تسک‌بریج" }] }),
  component: AppShell,
});

type NavItem = { to: string; label: string; icon: any; end?: boolean; badge?: number };

const navByRole: Record<Role, NavItem[]> = {
  requester: [
    { to: "/app",                label: "داشبورد",                   icon: LayoutDashboard, end: true },
    { to: "/app/tasks/new",      label: "ثبت درخواست جدید",          icon: PlusCircle },
    { to: "/app/tasks",          label: "درخواست‌های من",            icon: ClipboardList },
    { to: "/app/in-progress",    label: "درخواست‌های در حال انجام",  icon: Activity },
    { to: "/app/completed",      label: "درخواست‌های تکمیل شده",     icon: CheckCircle2 },
    { to: "/app/wallet",         label: "کیف پول",                   icon: Wallet },
    { to: "/app/chat",           label: "گفت‌وگوها",                 icon: MessagesSquare, badge: 2 },
    { to: "/app/notifications",  label: "اعلان‌ها",                  icon: Bell, badge: 3 },
    { to: "/app/profile",        label: "پروفایل",                   icon: User },
  ],
  agent: [
    { to: "/app/agent",          label: "داشبورد",                   icon: LayoutDashboard, end: true },
    { to: "/app/tasks",          label: "فرصت‌های جدید",             icon: Sparkles },
    { to: "/app/accepted",       label: "درخواست‌های پذیرفته شده",   icon: Briefcase },
    { to: "/app/completed",      label: "درخواست‌های تکمیل شده",     icon: CheckCircle2 },
    { to: "/app/earnings",       label: "درآمدها",                   icon: BarChart3 },
    { to: "/app/wallet",         label: "کیف پول",                   icon: Wallet },
    { to: "/app/chat",           label: "گفت‌وگوها",                 icon: MessagesSquare, badge: 2 },
    { to: "/app/notifications",  label: "اعلان‌ها",                  icon: Bell, badge: 3 },
    { to: "/app/profile",        label: "پروفایل",                   icon: User },
  ],
  admin: [
    { to: "/app/admin",          label: "داشبورد مدیریت",            icon: LayoutDashboard, end: true },
    { to: "/app/admin/users",    label: "کاربران",                   icon: Users },
    { to: "/app/tasks",          label: "درخواست‌ها",                icon: ClipboardList },
    { to: "/app/admin/finance",  label: "تراکنش‌ها و درآمد",         icon: Wallet },
    { to: "/app/admin/reports",  label: "گزارش‌ها",                  icon: FileBarChart },
    { to: "/app/notifications",  label: "اعلان‌ها",                  icon: Bell },
    { to: "/app/profile",        label: "پروفایل",                   icon: User },
  ],
};

function AppShell() {
  const [open, setOpen] = useState(false);
  return (
    <div className="min-h-screen bg-background">
      <Sidebar onClose={() => setOpen(false)} mobileOpen={open} />
      <div className="lg:pr-72">
        <Topbar onMenu={() => setOpen(true)} />
        <main className="px-4 py-6 sm:px-6 lg:px-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function RoleSwitcher() {
  const { role, setRole } = useRole();
  const roles: Role[] = ["requester", "agent", "admin"];
  return (
    <div className="rounded-xl border bg-card p-1">
      <div className="flex items-center gap-1 text-xs">
        {roles.map(r => (
          <button
            key={r}
            onClick={() => setRole(r)}
            className={cn(
              "flex-1 rounded-lg px-2 py-1.5 font-medium transition-colors",
              role === r ? "bg-primary text-primary-foreground shadow-soft" : "text-muted-foreground hover:bg-accent"
            )}
          >
            {ROLE_LABELS[r]}
          </button>
        ))}
      </div>
    </div>
  );
}

function Sidebar({ mobileOpen, onClose }: { mobileOpen: boolean; onClose: () => void }) {
  const path = useRouterState({ select: s => s.location.pathname });
  const { role } = useRole();
  const nav = navByRole[role];
  return (
    <>
      {mobileOpen && <button className="fixed inset-0 z-40 bg-foreground/30 backdrop-blur-sm lg:hidden" onClick={onClose} />}
      <aside className={cn(
        "fixed inset-y-0 right-0 z-50 flex w-72 flex-col border-l bg-sidebar transition-transform lg:translate-x-0",
        mobileOpen ? "translate-x-0" : "translate-x-full"
      )}>
        <div className="flex items-center justify-between px-5 py-4">
          <Logo />
          <button onClick={onClose} className="lg:hidden"><X className="h-5 w-5" /></button>
        </div>
        <div className="px-3 pb-2">
          <RoleSwitcher />
        </div>
        <nav className="flex-1 space-y-0.5 overflow-y-auto px-3 py-2">
          {nav.map(item => {
            const active = item.end ? path === item.to : path.startsWith(item.to);
            return (
              <Link
                key={item.to}
                to={item.to as any}
                onClick={onClose}
                className={cn(
                  "group flex items-center gap-3 rounded-lg px-3 py-2.5 text-sm font-medium transition-colors",
                  active ? "bg-primary/10 text-primary" : "text-muted-foreground hover:bg-accent hover:text-foreground"
                )}
              >
                <item.icon className={cn("h-4 w-4", active && "text-primary")} />
                <span className="flex-1">{item.label}</span>
                {item.badge && (
                  <span className="rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground">{item.badge}</span>
                )}
              </Link>
            );
          })}
        </nav>
        <div className="border-t p-3">
          <Link to="/app/profile" className="flex items-center gap-3 rounded-lg p-2 hover:bg-accent">
            <span className="grid h-9 w-9 place-items-center rounded-full gradient-brand text-sm font-bold text-white">س‌م</span>
            <div className="min-w-0 flex-1">
              <div className="truncate text-sm font-semibold">سارا محمدی</div>
              <div className="truncate text-xs text-muted-foreground">{ROLE_LABELS[role]} · پرو</div>
            </div>
            <Settings className="h-4 w-4 text-muted-foreground" />
          </Link>
        </div>
      </aside>
    </>
  );
}

function Topbar({ onMenu }: { onMenu: () => void }) {
  const { theme, toggle } = useTheme();
  const { role, setRole } = useRole();
  const cycle = () => setRole(role === "requester" ? "agent" : role === "agent" ? "admin" : "requester");
  return (
    <header className="sticky top-0 z-30 flex items-center gap-3 border-b bg-background/80 px-4 py-3 backdrop-blur sm:px-6 lg:px-8">
      <button onClick={onMenu} className="grid h-9 w-9 place-items-center rounded-lg border bg-card lg:hidden"><Menu className="h-4 w-4" /></button>
      <div className="relative hidden flex-1 max-w-md md:block">
        <Search className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          placeholder="جستجوی درخواست‌ها، انجام‌دهنده‌ها، تراکنش‌ها…"
          className="w-full rounded-lg border bg-card py-2 pr-9 pl-3 text-sm outline-none placeholder:text-muted-foreground focus:border-primary focus:ring-2 focus:ring-primary/20"
        />
      </div>
      <div className="me-auto flex items-center gap-2">
        <button onClick={cycle} title="تغییر نقش" className="hidden sm:inline-flex items-center gap-1.5 rounded-lg border bg-card px-3 py-2 text-xs font-medium hover:bg-accent">
          <ArrowLeftRight className="h-3.5 w-3.5" /> {ROLE_LABELS[role]}
        </button>
        <button onClick={toggle} className="grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent" aria-label="تغییر تم">
          {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </button>
        <Link to="/app/notifications" className="relative grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent">
          <Bell className="h-4 w-4" />
          <span className="absolute left-1.5 top-1.5 h-2 w-2 rounded-full bg-destructive ring-2 ring-background" />
        </Link>
        {role === "requester" && (
          <Link to="/app/tasks/new" className="hidden sm:inline-flex items-center gap-1.5 rounded-lg gradient-brand px-3.5 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            <PlusCircle className="h-4 w-4" /> درخواست جدید
          </Link>
        )}
      </div>
    </header>
  );
}
