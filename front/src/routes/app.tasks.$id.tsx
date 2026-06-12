import { createFileRoute, Link, notFound } from "@tanstack/react-router";
import { ArrowRight, MapPin, Clock, Wallet, MessagesSquare, CheckCircle2, Star, ShieldCheck, Paperclip } from "lucide-react";
import { tasks, agents } from "@/lib/mock-data";
import { StatusBadge } from "@/components/status-badge";
import { TaskTimeline } from "@/components/task-timeline";
import { toman, toFa } from "@/lib/fa";
import { toast } from "sonner";

export const Route = createFileRoute("/app/tasks/$id")({
  head: ({ params }) => ({ meta: [{ title: `${params.id} — تسک‌بریج` }] }),
  component: TaskDetails,
  notFoundComponent: () => <div className="p-10 text-center text-muted-foreground">درخواست پیدا نشد</div>,
  loader: ({ params }) => {
    const task = tasks.find(t => t.id === params.id);
    if (!task) throw notFound();
    return { task };
  },
});

function TaskDetails() {
  const { task } = Route.useLoaderData();
  const fee = Math.round(task.budget * 0.08);
  return (
    <div className="space-y-6">
      <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm text-muted-foreground hover:text-foreground">
        <ArrowRight className="h-4 w-4" /> بازگشت به فهرست درخواست‌ها
      </Link>

      <div className="grid gap-6 lg:grid-cols-3">
        <div className="space-y-6 lg:col-span-2">
          <div className="rounded-2xl border bg-card p-6 shadow-soft">
            <div className="flex flex-wrap items-start justify-between gap-3">
              <div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span className="font-mono">{task.id}</span> · <span>ثبت {task.postedAgo}</span> · <span>توسط {task.postedBy}</span>
                </div>
                <h1 className="mt-2 text-2xl font-bold tracking-tight">{task.title}</h1>
              </div>
              <StatusBadge status={task.status} />
            </div>
            <div className="mt-4 flex flex-wrap gap-x-5 gap-y-2 text-sm text-muted-foreground">
              <span className="inline-flex items-center gap-1.5"><MapPin className="h-4 w-4" />{task.city}</span>
              <span className="inline-flex items-center gap-1.5"><Clock className="h-4 w-4" />مهلت: {task.deadline}</span>
              <span className="inline-flex items-center gap-1.5"><Wallet className="h-4 w-4 text-primary" /><span className="font-semibold text-foreground">{toman(task.budget)}</span></span>
              <span className="inline-flex items-center gap-1.5 rounded-full bg-muted px-2 py-0.5 text-xs">{task.category}</span>
            </div>
            <p className="mt-5 text-sm leading-relaxed text-foreground/90">{task.description}</p>
            <div className="mt-5 flex items-center gap-2 text-xs text-muted-foreground">
              <Paperclip className="h-3.5 w-3.5" /> {toFa(2)} پیوست
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 shadow-soft">
            <h2 className="font-display font-semibold">مراحل انجام</h2>
            <div className="mt-6 overflow-x-auto">
              <TaskTimeline status={task.status} />
            </div>
          </div>

          <div className="rounded-2xl border bg-card p-6 shadow-soft">
            <div className="flex items-center justify-between">
              <h2 className="font-display font-semibold">متقاضی‌ها <span className="text-muted-foreground">({toFa(task.applicants)})</span></h2>
            </div>
            <div className="mt-4 space-y-3">
              {agents.map(a => (
                <div key={a.id} className="flex flex-wrap items-center gap-4 rounded-xl border p-4 hover:border-primary/40">
                  <span className="grid h-12 w-12 place-items-center rounded-full gradient-brand font-bold text-white">{a.name[0]}</span>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-center gap-2">
                      <div className="font-semibold">{a.name}</div>
                      <span className="rounded-full bg-success/15 px-2 py-0.5 text-xs font-medium text-success-foreground">{a.badge}</span>
                    </div>
                    <div className="mt-0.5 flex flex-wrap items-center gap-3 text-xs text-muted-foreground">
                      <span className="inline-flex items-center gap-1"><Star className="h-3.5 w-3.5 fill-warning text-warning" />{toFa(a.rating)} · {toFa(a.tasks)} کار</span>
                      <span className="inline-flex items-center gap-1"><MapPin className="h-3.5 w-3.5" />{a.city}</span>
                      <span className="inline-flex items-center gap-1"><Clock className="h-3.5 w-3.5" />آماده {a.eta}</span>
                    </div>
                    <p className="mt-1 text-sm text-muted-foreground">{a.bio}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <div className="rounded-lg bg-primary/10 px-3 py-1.5 text-sm font-semibold text-primary tabular-nums">{toman(a.price)}</div>
                    <Link to="/app/chat" className="grid h-9 w-9 place-items-center rounded-lg border hover:bg-accent"><MessagesSquare className="h-4 w-4" /></Link>
                    <button onClick={() => toast.success(`درخواست به ${a.name} واگذار شد`)} className="rounded-lg gradient-brand px-3.5 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
                      پذیرفتن
                    </button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>

        <aside className="space-y-4">
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <div className="flex items-center gap-2">
              <ShieldCheck className="h-4 w-4 text-success" />
              <h3 className="font-semibold">جزئیات پرداخت</h3>
            </div>
            <dl className="mt-4 space-y-2 text-sm">
              <Row k="مبلغ درخواست" v={toman(task.budget)} />
              <Row k="کارمزد سامانه" v={toman(fee)} />
              <Row k="بلوکه شده در امانت" v={toman(task.budget + fee)} bold />
            </dl>
            <div className="mt-4 rounded-lg bg-success/10 p-3 text-xs text-success-foreground">
              <CheckCircle2 className="me-1 inline h-3.5 w-3.5" /> پرداخت تنها پس از تایید نهایی شما آزاد می‌شود.
            </div>
          </div>
          <div className="rounded-2xl border bg-card p-5 shadow-soft">
            <h3 className="font-semibold">نیاز به کمک دارید؟</h3>
            <p className="mt-1 text-sm text-muted-foreground">تیم پشتیبانی ۲۴/۷ می‌تواند پرداخت را متوقف و میانجی‌گری کند.</p>
            <button className="mt-3 w-full rounded-lg border py-2 text-sm font-medium hover:bg-accent">ثبت اعتراض</button>
          </div>
        </aside>
      </div>
    </div>
  );
}

function Row({ k, v, bold }: any) {
  return (
    <div className={`flex items-center justify-between ${bold ? "font-semibold" : ""}`}>
      <dt className="text-muted-foreground">{k}</dt>
      <dd className="tabular-nums">{v}</dd>
    </div>
  );
}
