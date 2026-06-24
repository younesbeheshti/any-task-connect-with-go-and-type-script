import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { Send, Search, CheckCheck, MessagesSquare } from "lucide-react";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";
import { useAuth } from "@/lib/auth-context";

export const Route = createFileRoute("/app/chat")({
  head: () => ({ meta: [{ title: "گفت‌وگوها — تسک‌بریج" }] }),
  component: Chat,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type ChatSummary = {
  id: string; taskId: string; name: string; last: string;
  unread: number; online: boolean; time: string;
};

type Message = {
  id: string; from: string; text: string; time: string; seen: boolean;
};

function Chat() {
  const { user } = useAuth();
  const [chats, setChats] = useState<ChatSummary[]>([]);
  const [activeTaskId, setActiveTaskId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [draft, setDraft] = useState("");
  const [sending, setSending] = useState(false);
  const [loadingChats, setLoadingChats] = useState(true);
  const [loadingMsgs, setLoadingMsgs] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const token = getToken();
    if (!token) { setLoadingChats(false); return; }
    fetch(`${API_BASE}/v1/chats`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => {
        const list: ChatSummary[] = (d.chats ?? []).map((c: any) => ({
          id: c.id, taskId: c.taskId, name: c.name || c.taskId?.slice(0, 8),
          last: c.lastMessage || "", unread: c.unread || 0,
          online: c.online ?? false, time: c.updatedAt?.slice(11, 16) ?? "",
        }));
        setChats(list);
        if (list.length > 0) openChat(list[0].taskId);
      })
      .catch(() => {})
      .finally(() => setLoadingChats(false));
  }, []);

  function openChat(taskId: string) {
    setActiveTaskId(taskId);
    setLoadingMsgs(true);
    const token = getToken();
    if (!token) return;
    fetch(`${API_BASE}/v1/tasks/${taskId}/messages`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => {
        const msgs: Message[] = (d.messages ?? []).map((m: any) => ({
          id: m.id,
          from: m.from === user?.id ? "me" : "them",
          text: m.text || m.message || "",
          time: m.createdAt?.slice(11, 16) ?? "",
          seen: m.readAt != null,
        }));
        setMessages(msgs);
        setTimeout(() => bottomRef.current?.scrollIntoView({ behavior: "smooth" }), 50);
      })
      .catch(() => setMessages([]))
      .finally(() => setLoadingMsgs(false));

    // mark read
    fetch(`${API_BASE}/v1/tasks/${taskId}/messages/read`, {
      method: "POST", headers: { Authorization: `Bearer ${token}` },
    }).catch(() => {});
    setChats(prev => prev.map(c => c.taskId === taskId ? { ...c, unread: 0 } : c));
  }

  async function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    if (!draft.trim() || !activeTaskId) return;
    const token = getToken();
    if (!token) return;

    const activeChat = chats.find(c => c.taskId === activeTaskId);
    setSending(true);
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${activeTaskId}/messages`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({ message: draft, receiverId: "00000000-0000-0000-0000-000000000000" }),
      });
      const data = await res.json();
      if (res.ok) {
        setMessages(prev => [...prev, {
          id: data.id, from: "me", text: draft,
          time: new Date().toLocaleTimeString("fa-IR", { hour: "2-digit", minute: "2-digit" }),
          seen: false,
        }]);
        setDraft("");
        setTimeout(() => bottomRef.current?.scrollIntoView({ behavior: "smooth" }), 50);
      }
    } catch {} finally {
      setSending(false);
    }
  }

  const activeChat = chats.find(c => c.taskId === activeTaskId);

  return (
    <div className="grid h-[calc(100vh-8rem)] gap-0 overflow-hidden rounded-2xl border bg-card shadow-soft md:grid-cols-[300px_1fr]">
      {/* Sidebar */}
      <aside className="flex flex-col border-l">
        <div className="border-b p-4">
          <div className="relative">
            <Search className="pointer-events-none absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <input placeholder="جستجوی گفت‌وگو…"
              className="w-full rounded-lg border bg-background py-2 pr-9 pl-3 text-sm outline-none focus:border-primary focus:ring-2 focus:ring-primary/20" />
          </div>
        </div>
        <div className="flex-1 overflow-y-auto">
          {loadingChats ? (
            Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="flex items-center gap-3 border-b px-4 py-3 animate-pulse">
                <div className="h-10 w-10 rounded-full bg-muted" />
                <div className="flex-1 space-y-1.5">
                  <div className="h-4 w-2/3 rounded bg-muted" />
                  <div className="h-3 w-full rounded bg-muted" />
                </div>
              </div>
            ))
          ) : chats.length === 0 ? (
            <div className="flex flex-col items-center px-4 py-10 text-center">
              <MessagesSquare className="h-8 w-8 text-muted-foreground/40" />
              <p className="mt-2 text-sm text-muted-foreground">گفت‌وگویی ندارید</p>
            </div>
          ) : (
            chats.map(c => (
              <button key={c.id} onClick={() => openChat(c.taskId)}
                className={cn("flex w-full items-center gap-3 border-b px-4 py-3 text-right hover:bg-accent", c.taskId === activeTaskId && "bg-accent")}>
                <div className="relative">
                  <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">
                    {(c.name || "؟")[0]}
                  </span>
                  {c.online && <span className="absolute bottom-0 left-0 h-2.5 w-2.5 rounded-full bg-success ring-2 ring-card" />}
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex items-center justify-between">
                    <div className="truncate text-sm font-semibold">{c.name}</div>
                    <div className="text-xs text-muted-foreground">{c.time}</div>
                  </div>
                  <div className="flex items-center justify-between">
                    <div className="truncate text-xs text-muted-foreground">{c.last}</div>
                    {c.unread > 0 && (
                      <span className="ms-2 rounded-full bg-primary px-1.5 py-0.5 text-[10px] font-semibold text-primary-foreground">
                        {toFa(c.unread)}
                      </span>
                    )}
                  </div>
                </div>
              </button>
            ))
          )}
        </div>
      </aside>

      {/* Message area */}
      <section className="flex min-w-0 flex-col">
        {!activeTaskId ? (
          <div className="flex flex-1 flex-col items-center justify-center gap-3 text-muted-foreground">
            <MessagesSquare className="h-12 w-12 opacity-30" />
            <p className="text-sm">یک گفت‌وگو انتخاب کنید</p>
          </div>
        ) : (
          <>
            <header className="flex items-center gap-3 border-b px-5 py-3">
              <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">
                {(activeChat?.name || "؟")[0]}
              </span>
              <div className="min-w-0 flex-1">
                <div className="font-semibold">{activeChat?.name}</div>
                <div className="text-xs text-muted-foreground">درخواست {activeTaskId?.slice(0, 8)}</div>
              </div>
            </header>

            <div className="flex-1 space-y-3 overflow-y-auto p-5">
              {loadingMsgs ? (
                <div className="flex items-center justify-center py-10 text-muted-foreground text-sm">در حال بارگذاری...</div>
              ) : messages.length === 0 ? (
                <div className="flex flex-col items-center py-10 text-muted-foreground">
                  <MessagesSquare className="h-8 w-8 opacity-30" />
                  <p className="mt-2 text-sm">هنوز پیامی نیست. اولین پیام را بفرستید.</p>
                </div>
              ) : (
                messages.map(m => {
                  const me = m.from === "me";
                  return (
                    <div key={m.id} className={cn("flex", me ? "justify-start" : "justify-end")}>
                      <div className={cn(
                        "max-w-[75%] rounded-2xl px-4 py-2.5 text-sm shadow-soft",
                        me ? "rounded-bl-md gradient-brand text-white" : "rounded-br-md bg-muted"
                      )}>
                        <div className="whitespace-pre-wrap">{m.text}</div>
                        <div className={cn("mt-1 flex items-center justify-start gap-1 text-[10px]",
                          me ? "text-white/80" : "text-muted-foreground")}>
                          {m.time}
                          {me && <CheckCheck className="h-3 w-3" />}
                        </div>
                      </div>
                    </div>
                  );
                })
              )}
              <div ref={bottomRef} />
            </div>

            <form onSubmit={sendMessage} className="border-t p-3">
              <div className="flex items-end gap-2 rounded-2xl border bg-background p-2 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20">
                <input
                  value={draft} onChange={e => setDraft(e.target.value)}
                  placeholder="پیام خود را بنویسید…"
                  className="flex-1 bg-transparent px-2 py-2 text-sm outline-none"
                />
                <button type="submit" disabled={sending || !draft.trim()}
                  className="inline-flex h-9 items-center gap-1.5 rounded-lg gradient-brand px-4 text-sm font-semibold text-white hover:opacity-95 disabled:opacity-50">
                  <Send className="h-4 w-4" /> ارسال
                </button>
              </div>
            </form>
          </>
        )}
      </section>
    </div>
  );
}
