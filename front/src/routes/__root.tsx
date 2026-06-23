import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import {
  Outlet,
  Link,
  createRootRouteWithContext,
  useRouter,
  HeadContent,
  Scripts,
} from "@tanstack/react-router";
import { useEffect, type ReactNode } from "react";

import appCss from "../styles.css?url";
import { reportLovableError } from "../lib/lovable-error-reporting";
import { Toaster } from "@/components/ui/sonner";
import { ThemeProvider } from "@/components/theme-provider";
import { RoleProvider } from "@/components/role-context";
import { AuthProvider } from "@/lib/auth-context";

function NotFoundComponent() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="max-w-md text-center">
        <div className="mx-auto mb-6 inline-flex h-16 w-16 items-center justify-center rounded-2xl gradient-brand shadow-glow">
          <span className="text-3xl font-bold text-white">۴۰۴</span>
        </div>
        <h1 className="text-3xl font-bold text-foreground">صفحه پیدا نشد</h1>
        <p className="mt-2 text-sm text-muted-foreground">
          آدرسی که دنبالش هستید موجود نیست یا منتقل شده است.
        </p>
        <div className="mt-6">
          <Link to="/" className="inline-flex items-center justify-center rounded-lg gradient-brand px-5 py-2.5 text-sm font-medium text-white shadow-soft hover:opacity-95">
            بازگشت به خانه
          </Link>
        </div>
      </div>
    </div>
  );
}

function ErrorComponent({ error, reset }: { error: Error; reset: () => void }) {
  console.error(error);
  const router = useRouter();
  useEffect(() => { reportLovableError(error, { boundary: "tanstack_root_error_component" }); }, [error]);
  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4">
      <div className="max-w-md text-center">
        <h1 className="text-xl font-semibold text-foreground">خطایی رخ داد</h1>
        <p className="mt-2 text-sm text-muted-foreground">دوباره تلاش کنید یا به صفحه اصلی بازگردید.</p>
        <div className="mt-6 flex justify-center gap-2">
          <button onClick={() => { router.invalidate(); reset(); }} className="rounded-lg bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:opacity-90">
            تلاش مجدد
          </button>
          <a href="/" className="rounded-lg border border-input bg-background px-4 py-2 text-sm font-medium hover:bg-accent">
            صفحه اصلی
          </a>
        </div>
      </div>
    </div>
  );
}

export const Route = createRootRouteWithContext<{ queryClient: QueryClient }>()({
  head: () => ({
    meta: [
      { charSet: "utf-8" },
      { name: "viewport", content: "width=device-width, initial-scale=1" },
      { title: "تسک‌بریج — انجام کارهای شما از هر کجا" },
      { name: "description", content: "تسک‌بریج شما را به انجام‌دهنده‌های مورد اعتماد محلی برای کارهای اداری، پزشکی، دانشگاهی، حقوقی و دولتی در هر شهر متصل می‌کند." },
      { name: "author", content: "TaskBridge" },
      { name: "theme-color", content: "#2563EB" },
      { property: "og:title", content: "تسک‌بریج — انجام کارهای شما از هر کجا" },
      { property: "og:description", content: "پلتفرم مطمئن انجام درخواست‌های اداری و خدمات شهری در سراسر ایران، با پرداخت امانی." },
      { property: "og:type", content: "website" },
      { name: "twitter:card", content: "summary_large_image" },
    ],
    links: [
      { rel: "stylesheet", href: appCss },
      { rel: "preconnect", href: "https://fonts.googleapis.com" },
      { rel: "preconnect", href: "https://fonts.gstatic.com", crossOrigin: "anonymous" },
      { rel: "stylesheet", href: "https://cdn.jsdelivr.net/gh/rastikerdar/vazirmatn@v33.003/Vazirmatn-font-face.css" },
    ],
  }),
  shellComponent: RootShell,
  component: RootComponent,
  notFoundComponent: NotFoundComponent,
  errorComponent: ErrorComponent,
});

function RootShell({ children }: { children: ReactNode }) {
  return (
    <html lang="fa" dir="rtl" suppressHydrationWarning>
      <head><HeadContent /></head>
      <body>
        {children}
        <Scripts />
      </body>
    </html>
  );
}

function RootComponent() {
  const { queryClient } = Route.useRouteContext();
  return (
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <ThemeProvider>
          <RoleProvider>
            <Outlet />
            <Toaster richColors position="top-left" />
          </RoleProvider>
        </ThemeProvider>
      </AuthProvider>
    </QueryClientProvider>
  );
}
