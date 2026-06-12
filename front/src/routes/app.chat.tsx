import { createFileRoute } from "@tanstack/react-router";
import { useState } from "react";
import { Send, Paperclip, Search, Smile, CheckCheck, Image as ImageIcon, Mic } from "lucide-react";
import { chats, messages } from "@/lib/mock-data";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

export const Route = createFileRoute("/app/chat")({
  head: () => ({ meta: [{ title: "گفت‌وگوها — تسک‌بریج" }] }),
  component: Chat,
});

function Chat() {
  const [active, setActive] = useState(chats[0].id);
  const [draft, setDraft] = useState("");
  const current = chats.find(c => c.id === active)!;
  return (
    <div className="grid h-[calc(100vh-8rem)] gap-0 overflow-hidden rounded-2xl border bg-card shadow-soft md:grid-cols-[320px_1fr]">
      <aside className="flex flex-col border-l">
        <div className="border-b p-4">
          <div className="relative">
            <Search className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input placeholder="جستجوی گفت‌وگو…" className="w-full rounded-lg border bg-background py-2 pr-9 pl-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20" />
          </div>
        </div>
        <div className="flex-1 overflow-y-auto">
          {chats.map(c => (
            <button key={c.id} onClick={() => setActive(c.id)} className={cn("flex w-full items-center gap-3 border-b px-4 py-3 text-right hover:bg-accent", c.id === active && "bg-accent")}>
              <div className="relative">
                <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">{c.name[0]}</span>
                {c.online && <span className="absolute bottom-0 left-0 h-2.5 w-2.5 rounded-full bg-success ring-2 ring-card" />}
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex items-center justify-between">
                  <div className="truncate text-sm font-semibold">{c.name}</div>
                  <div className="text-xs text-muted-foreground">{c.time}</div>
                </div>
                <div className="flex items-center justify-between">
                  <div className="truncate text-xs text-muted-foreground">{c.last}</div>
                  {c.unread > 0 && <span className="ms-2 rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground">{toFa(c.unread)}</span>}
                </div>
              </div>
            </button>
          ))}
        </div>
      </aside>

      <section className="flex min-w-0 flex-col">
        <header className="flex items-center gap-3 border-b px-5 py-3">
          <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">{current.name[0]}</span>
          <div className="min-w-0 flex-1">
            <div className="font-semibold">{current.name}</div>
            <div className="text-xs text-muted-foreground">{current.online ? "آنلاین" : "آفلاین"} · درخواست {current.taskId}</div>
          </div>
        </header>

        <div className="flex-1 space-y-3 overflow-y-auto p-5">
          <DayDivider label="امروز" />
          {messages.map(m => {
            const me = m.from === "me";
            return (
              <div key={m.id} className={cn("flex", me ? "justify-start" : "justify-end")}>
                <div className={cn(
                  "max-w-[75%] rounded-2xl px-4 py-2.5 text-sm shadow-soft",
                  me ? "rounded-bl-md gradient-brand text-white" : "rounded-br-md bg-muted"
                )}>
                  <div className="whitespace-pre-wrap">{m.text}</div>
                  <div className={cn("mt-1 flex items-center justify-start gap-1 text-[10px]", me ? "text-white/80" : "text-muted-foreground")}>
                    {m.time}
                    {me && <CheckCheck className="h-3 w-3" />}
                  </div>
                </div>
              </div>
            );
          })}
        </div>

        <form onSubmit={(e) => { e.preventDefault(); setDraft(""); }} className="border-t p-3">
          <div className="flex items-end gap-2 rounded-2xl border bg-background p-2 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20">
            <button type="button" className="grid h-9 w-9 place-items-center rounded-lg text-muted-foreground hover:bg-accent"><Paperclip className="h-4 w-4" /></button>
            <button type="button" className="grid h-9 w-9 place-items-center rounded-lg text-muted-foreground hover:bg-accent"><ImageIcon className="h-4 w-4" /></button>
            <button type="button" className="grid h-9 w-9 place-items-center rounded-lg text-muted-foreground hover:bg-accent"><Mic className="h-4 w-4" /></button>
            <input
              value={draft} onChange={e => setDraft(e.target.value)}
              placeholder="پیام خود را بنویسید…"
              className="flex-1 bg-transparent px-2 py-2 text-sm outline-none"
            />
            <button type="button" className="grid h-9 w-9 place-items-center rounded-lg text-muted-foreground hover:bg-accent"><Smile className="h-4 w-4" /></button>
            <button className="inline-flex h-9 items-center gap-1.5 rounded-lg gradient-brand px-4 text-sm font-semibold text-white hover:opacity-95">
              <Send className="h-4 w-4" /> ارسال
            </button>
          </div>
        </form>
      </section>
    </div>
  );
}

function DayDivider({ label }: { label: string }) {
  return (
    <div className="my-2 flex items-center gap-3">
      <div className="h-px flex-1 bg-border" />
      <span className="text-xs font-medium text-muted-foreground">{label}</span>
      <div className="h-px flex-1 bg-border" />
    </div>
  );
}
