import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useRef, useState } from "react";
import { Send, Search, CheckCheck, MessagesSquare, Paperclip, X, FileText } from "lucide-react";
import { toast } from "sonner";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";
import { useAuth } from "@/lib/auth-context";
import { FileUploader, type UploadedFile } from "@/components/file-uploader";

export const Route = createFileRoute("/app/chat")({
  // `task` (a public task id) opens that conversation directly, e.g. from a task page.
  validateSearch: (search: Record<string, unknown>): { task?: string } => ({
    task: typeof search.task === "string" ? search.task : undefined,
  }),
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

type Attachment = { id: string; url: string; mime: string; name: string; size: number };

type Message = {
  id: string; from: string; text: string; time: string; seen: boolean;
  attachment?: Attachment | null;
};

function Chat() {
  const { user } = useAuth();
  const { task: taskParam } = Route.useSearch();
  const [chats, setChats] = useState<ChatSummary[]>([]);
  const [activeTaskId, setActiveTaskId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [draft, setDraft] = useState("");
  const [pendingAttachment, setPendingAttachment] = useState<Attachment | null>(null);
  const [showUploader, setShowUploader] = useState(false);
  const [sending, setSending] = useState(false);
  const [loadingChats, setLoadingChats] = useState(true);
  const [loadingMsgs, setLoadingMsgs] = useState(false);
  const bottomRef = useRef<HTMLDivElement>(null);
  const scrollRef = useRef<HTMLDivElement>(null);
  // Signature/count of the last-rendered thread, so polling skips redundant
  // re-renders and only auto-scrolls when genuinely new messages arrive.
  const msgSigRef = useRef<string>("");
  const msgCountRef = useRef<number>(0);

  function isAtBottom() {
    const el = scrollRef.current;
    if (!el) return true;
    return el.scrollHeight - el.scrollTop - el.clientHeight < 80;
  }
  function scrollToBottom(smooth = true) {
    bottomRef.current?.scrollIntoView({ behavior: smooth ? "smooth" : "auto" });
  }

  // Maps the API conversation list to the sidebar shape. Keeps a placeholder for
  // a task opened with no messages yet, and zeroes the active chat's unread.
  function applyChats(raw: any[], active: string | null): ChatSummary[] {
    const list: ChatSummary[] = (raw ?? []).map((c: any) => ({
      id: c.id, taskId: c.taskId, name: c.name || c.taskId,
      last: c.last || "", unread: c.taskId === active ? 0 : (c.unread || 0),
      online: c.online ?? false, time: c.time?.slice(11, 16) ?? "",
    }));
    if (taskParam && !list.some(c => c.taskId === taskParam)) {
      list.unshift({ id: taskParam, taskId: taskParam, name: `درخواست ${taskParam}`, last: "", unread: 0, online: false, time: "" });
    }
    return list;
  }

  function loadChats(active: string | null) {
    const token = getToken();
    if (!token) return;
    fetch(`${API_BASE}/v1/chats`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => setChats(applyChats(d.chats ?? [], active)))
      .catch(() => {});
  }

  // Initial load: conversations, then open the requested/first one.
  useEffect(() => {
    const token = getToken();
    if (!token) { setLoadingChats(false); return; }
    fetch(`${API_BASE}/v1/chats`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => {
        const initial = taskParam ?? ((d.chats ?? []).length > 0 ? d.chats[0].taskId : null);
        setChats(applyChats(d.chats ?? [], initial));
        if (initial) openChat(initial);
      })
      .catch(() => {})
      .finally(() => setLoadingChats(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [taskParam]);

  // Poll the active thread + conversation list every 5s.
  useEffect(() => {
    if (!activeTaskId) return;
    const t = setInterval(() => {
      fetchThread(activeTaskId, false);
      loadChats(activeTaskId);
    }, 5000);
    return () => clearInterval(t);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeTaskId]);

  // Fetches the thread. `force` (explicit open) always replaces + scrolls to the
  // bottom; polling only re-renders on real changes and preserves scroll position.
  function fetchThread(taskId: string, force: boolean) {
    const token = getToken();
    if (!token) return;
    const wasAtBottom = isAtBottom();
    fetch(`${API_BASE}/v1/tasks/${taskId}/messages`, { headers: { Authorization: `Bearer ${token}` } })
      .then(r => r.json())
      .then(d => {
        const msgs: Message[] = (d.messages ?? []).map((m: any) => ({
          id: m.id,
          from: m.from === user?.id ? "me" : "them",
          text: m.text || m.message || "",
          time: m.createdAt?.slice(11, 16) ?? "",
          seen: m.readAt != null,
          attachment: m.attachments ?? m.attachment ?? null,
        }));
        const sig = msgs.map(m => `${m.id}:${m.seen ? 1 : 0}`).join("|");
        if (!force && sig === msgSigRef.current) return;
        const grew = msgs.length > msgCountRef.current;
        msgSigRef.current = sig;
        msgCountRef.current = msgs.length;
        setMessages(msgs);
        if (force || (wasAtBottom && grew)) {
          setTimeout(() => scrollToBottom(!force), 50);
        }
        // New incoming messages while the chat is open → mark them read.
        if (msgs.some(m => m.from === "them" && !m.seen)) {
          fetch(`${API_BASE}/v1/tasks/${taskId}/messages/read`, {
            method: "POST", headers: { Authorization: `Bearer ${token}` },
          }).catch(() => {});
          setChats(prev => prev.map(c => c.taskId === taskId ? { ...c, unread: 0 } : c));
        }
      })
      .catch(() => { if (force) setMessages([]); })
      .finally(() => { if (force) setLoadingMsgs(false); });
  }

  function openChat(taskId: string) {
    setActiveTaskId(taskId);
    setLoadingMsgs(true);
    msgSigRef.current = "";
    msgCountRef.current = 0;
    fetchThread(taskId, true);

    const token = getToken();
    if (token) {
      fetch(`${API_BASE}/v1/tasks/${taskId}/messages/read`, {
        method: "POST", headers: { Authorization: `Bearer ${token}` },
      }).catch(() => {});
    }
    setChats(prev => prev.map(c => c.taskId === taskId ? { ...c, unread: 0 } : c));
  }

  async function sendMessage(e: React.FormEvent) {
    e.preventDefault();
    if ((!draft.trim() && !pendingAttachment) || !activeTaskId) return;
    const token = getToken();
    if (!token) return;

    const attachment = pendingAttachment;
    setSending(true);
    try {
      const res = await fetch(`${API_BASE}/v1/tasks/${activeTaskId}/messages`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}`, "Content-Type": "application/json" },
        body: JSON.stringify({
          message: draft,
          ...(attachment ? { attachment } : {}),
        }),
      });
      const data = await res.json();
      if (res.ok) {
        setMessages(prev => {
          const next = [...prev, {
            id: data.id, from: "me", text: draft,
            time: new Date().toLocaleTimeString("fa-IR", { hour: "2-digit", minute: "2-digit" }),
            seen: false,
            attachment,
          }];
          // Keep the poll signature in sync so the next refresh doesn't re-scroll.
          msgCountRef.current = next.length;
          return next;
        });
        setDraft("");
        setPendingAttachment(null);
        setShowUploader(false);
        // Surface the now-real conversation in the sidebar list.
        setChats(prev => prev.map(c => c.taskId === activeTaskId ? { ...c, last: draft || "پیوست" } : c));
        setTimeout(() => bottomRef.current?.scrollIntoView({ behavior: "smooth" }), 50);
      } else {
        toast.error(data?.error?.message ?? "ارسال پیام ناموفق بود");
      }
    } catch {
      toast.error("خطای شبکه");
    } finally {
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

            <div ref={scrollRef} className="flex-1 space-y-3 overflow-y-auto p-5">
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
                        {m.attachment && <AttachmentView attachment={m.attachment} me={me} />}
                        {m.text && <div className="whitespace-pre-wrap">{m.text}</div>}
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
              {showUploader && !pendingAttachment && (
                <div className="mb-2">
                  <FileUploader
                    multiple={false}
                    onChange={(files) => { if (files[0]) { setPendingAttachment(files[0]); setShowUploader(false); } }}
                  />
                </div>
              )}
              {pendingAttachment && (
                <div className="mb-2 flex items-center gap-2 rounded-lg border bg-card px-3 py-2 text-sm">
                  <FileText className="h-4 w-4 shrink-0 text-primary" />
                  <span className="min-w-0 flex-1 truncate">{pendingAttachment.name}</span>
                  <button type="button" onClick={() => setPendingAttachment(null)}
                    className="shrink-0 rounded p-1 text-muted-foreground hover:bg-accent hover:text-foreground" aria-label="حذف پیوست">
                    <X className="h-4 w-4" />
                  </button>
                </div>
              )}
              <div className="flex items-end gap-2 rounded-2xl border bg-background p-2 focus-within:border-primary focus-within:ring-2 focus-within:ring-primary/20">
                <button type="button" onClick={() => setShowUploader(v => !v)}
                  className={cn("grid h-9 w-9 place-items-center rounded-lg text-muted-foreground hover:bg-accent hover:text-foreground", showUploader && "bg-accent text-foreground")}
                  aria-label="افزودن پیوست">
                  <Paperclip className="h-4 w-4" />
                </button>
                <input
                  value={draft} onChange={e => setDraft(e.target.value)}
                  placeholder="پیام خود را بنویسید…"
                  className="flex-1 bg-transparent px-2 py-2 text-sm outline-none"
                />
                <button type="submit" disabled={sending || (!draft.trim() && !pendingAttachment)}
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

/** Renders a chat attachment. Downloads are auth-protected, so we fetch with the
 *  Bearer token and open the result as a blob URL rather than a plain link. */
function AttachmentView({ attachment, me }: { attachment: Attachment; me: boolean }) {
  const [busy, setBusy] = useState(false);

  async function open() {
    const token = getToken();
    if (!token) return;
    setBusy(true);
    try {
      const res = await fetch(`${API_BASE}${attachment.url}`, {
        headers: { Authorization: `Bearer ${token}` },
      });
      if (!res.ok) return;
      const blob = await res.blob();
      const objectUrl = URL.createObjectURL(blob);
      window.open(objectUrl, "_blank");
      setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000);
    } finally {
      setBusy(false);
    }
  }

  return (
    <button
      type="button"
      onClick={open}
      disabled={busy}
      className={cn(
        "mb-1.5 flex w-full items-center gap-2 rounded-lg px-2.5 py-2 text-start text-xs transition-colors disabled:opacity-60",
        me ? "bg-white/15 hover:bg-white/25" : "bg-background hover:bg-accent"
      )}
    >
      <FileText className="h-4 w-4 shrink-0" />
      <span className="min-w-0 flex-1 truncate font-medium">{attachment.name}</span>
    </button>
  );
}
