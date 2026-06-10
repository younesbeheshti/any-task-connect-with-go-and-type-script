import { createFileRoute, Link, Outlet, useRouterState } from "@tanstack/react-router";
import { LayoutDashboard, ListChecks, PlusCircle, Wallet, MessagesSquare, Bell, User, Settings, Search, Moon, Sun, Shield, Menu, X, Briefcase } from "lucide-react";
import { Logo } from "@/components/logo";
import { useTheme } from "@/components/theme-provider";
import { cn } from "@/lib/utils";
import { useState } from "react";

export const Route = createFileRoute("/app")({
  head: () => ({ meta: [{ title: "Dashboard — TaskBridge" }] }),
  component: AppShell,
});

const nav = [
  { to: "/app", label: "Dashboard", icon: LayoutDashboard, end: true },
  { to: "/app/tasks", label: "Marketplace", icon: ListChecks },
  { to: "/app/tasks/new", label: "Create Task", icon: PlusCircle },
  { to: "/app/agent", label: "Agent", icon: Briefcase },
  { to: "/app/wallet", label: "Wallet", icon: Wallet },
  { to: "/app/chat", label: "Messages", icon: MessagesSquare, badge: 2 },
  { to: "/app/notifications", label: "Notifications", icon: Bell, badge: 3 },
  { to: "/app/profile", label: "Profile", icon: User },
  { to: "/app/admin", label: "Admin", icon: Shield },
] as const;

function AppShell() {
  const [open, setOpen] = useState(false);
  return (
    <div className="min-h-screen bg-background">
      <Sidebar onClose={() => setOpen(false)} mobileOpen={open} />
      <div className="lg:pl-64">
        <Topbar onMenu={() => setOpen(true)} />
        <main className="px-4 py-6 sm:px-6 lg:px-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}

function Sidebar({ mobileOpen, onClose }: { mobileOpen: boolean; onClose: () => void }) {
  const path = useRouterState({ select: s => s.location.pathname });
  return (
    <>
      {mobileOpen && <button className="fixed inset-0 z-40 bg-foreground/30 backdrop-blur-sm lg:hidden" onClick={onClose} />}
      <aside className={cn(
        "fixed inset-y-0 left-0 z-50 flex w-64 flex-col border-r bg-sidebar transition-transform lg:translate-x-0",
        mobileOpen ? "translate-x-0" : "-translate-x-full"
      )}>
        <div className="flex items-center justify-between px-5 py-4">
          <Logo />
          <button onClick={onClose} className="lg:hidden"><X className="h-5 w-5" /></button>
        </div>
        <nav className="flex-1 space-y-0.5 px-3">
          {nav.map(item => {
            const active = item.end ? path === item.to : path.startsWith(item.to);
            return (
              <Link
                key={item.to}
                to={item.to}
                onClick={onClose}
                className={cn(
                  "group flex items-center gap-3 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                  active ? "bg-primary/10 text-primary" : "text-muted-foreground hover:bg-accent hover:text-foreground"
                )}
              >
                <item.icon className={cn("h-4 w-4", active && "text-primary")} />
                <span className="flex-1">{item.label}</span>
                {"badge" in item && item.badge && (
                  <span className="rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground">{item.badge}</span>
                )}
              </Link>
            );
          })}
        </nav>
        <div className="border-t p-3">
          <Link to="/app/profile" className="flex items-center gap-3 rounded-lg p-2 hover:bg-accent">
            <span className="grid h-9 w-9 place-items-center rounded-full gradient-brand text-sm font-bold text-white">SM</span>
            <div className="min-w-0 flex-1">
              <div className="truncate text-sm font-semibold">Sara M.</div>
              <div className="truncate text-xs text-muted-foreground">Requester · Pro</div>
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
  return (
    <header className="sticky top-0 z-30 flex items-center gap-3 border-b bg-background/80 px-4 py-3 backdrop-blur sm:px-6 lg:px-8">
      <button onClick={onMenu} className="grid h-9 w-9 place-items-center rounded-lg border bg-card lg:hidden"><Menu className="h-4 w-4" /></button>
      <div className="relative hidden flex-1 max-w-md md:block">
        <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
        <input
          placeholder="Search tasks, agents, transactions…"
          className="w-full rounded-lg border bg-card py-2 pl-9 pr-3 text-sm outline-none placeholder:text-muted-foreground focus:border-primary focus:ring-2 focus:ring-primary/20"
        />
      </div>
      <div className="ml-auto flex items-center gap-2">
        <button onClick={toggle} className="grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent" aria-label="Toggle theme">
          {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </button>
        <Link to="/app/notifications" className="relative grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent">
          <Bell className="h-4 w-4" />
          <span className="absolute right-1.5 top-1.5 h-2 w-2 rounded-full bg-destructive ring-2 ring-background" />
        </Link>
        <Link to="/app/tasks/new" className="hidden sm:inline-flex items-center gap-1.5 rounded-lg gradient-brand px-3.5 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
          <PlusCircle className="h-4 w-4" /> New task
        </Link>
      </div>
    </header>
  );
}
