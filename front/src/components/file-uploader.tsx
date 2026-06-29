import { useEffect, useRef, useState } from "react";
import { Upload, X, FileText, Loader2 } from "lucide-react";
import { toFa } from "@/lib/fa";
import { cn } from "@/lib/utils";

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

function getToken(): string | null {
  try {
    return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null;
  } catch {
    return null;
  }
}

/** Attachment metadata returned by POST /v1/files (matches backend chat.Attachment). */
export type UploadedFile = {
  id: string;
  url: string;
  name: string;
  mime: string;
  size: number;
};

const DEFAULT_MAX_MB = 50;
const ACCEPTED_EXT = [
  "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "csv",
  "zip", "rar", "7z",
  "png", "jpg", "jpeg", "gif", "webp", "svg",
  "mp4", "avi", "mov", "mkv", "mp3", "wav",
];

type Pending = {
  key: string;
  name: string;
  size: number;
  progress: number;
  error?: string;
  done?: UploadedFile;
};

type Props = {
  /** Called with the cumulative list of successfully uploaded files. */
  onChange: (files: UploadedFile[]) => void;
  multiple?: boolean;
  maxSizeMb?: number;
  className?: string;
};

function fmtSize(bytes: number): string {
  if (bytes < 1024) return `${toFa(bytes)} بایت`;
  if (bytes < 1024 * 1024) return `${toFa(Math.round(bytes / 1024))} کیلوبایت`;
  return `${toFa((bytes / (1024 * 1024)).toFixed(1))} مگابایت`;
}

function uploadOne(file: File, onProgress: (pct: number) => void): Promise<UploadedFile> {
  return new Promise((resolve, reject) => {
    const token = getToken();
    if (!token) {
      reject(new Error("احراز هویت لازم است"));
      return;
    }
    const form = new FormData();
    form.append("files", file);
    const xhr = new XMLHttpRequest();
    xhr.open("POST", `${API_BASE}/v1/files`);
    xhr.setRequestHeader("Authorization", `Bearer ${token}`);
    xhr.upload.onprogress = (e) => {
      if (e.lengthComputable) onProgress(Math.round((e.loaded / e.total) * 100));
    };
    xhr.onload = () => {
      try {
        const data = JSON.parse(xhr.responseText);
        if (xhr.status >= 200 && xhr.status < 300) {
          const uploaded: UploadedFile | undefined = data?.files?.[0];
          if (uploaded) resolve(uploaded);
          else reject(new Error("پاسخ نامعتبر سرور"));
        } else {
          reject(new Error(data?.error?.message ?? "بارگذاری ناموفق بود"));
        }
      } catch {
        reject(new Error("بارگذاری ناموفق بود"));
      }
    };
    xhr.onerror = () => reject(new Error("خطای شبکه"));
    xhr.send(form);
  });
}

export function FileUploader({ onChange, multiple = true, maxSizeMb = DEFAULT_MAX_MB, className }: Props) {
  const [items, setItems] = useState<Pending[]>([]);
  const inputRef = useRef<HTMLInputElement>(null);

  // Notify the parent of completed uploads after each render — never inside a
  // setState updater (that runs during render and triggers a parent setState).
  const onChangeRef = useRef(onChange);
  onChangeRef.current = onChange;
  useEffect(() => {
    onChangeRef.current(items.filter((i) => i.done).map((i) => i.done!));
  }, [items]);

  function handleFiles(fileList: FileList | null) {
    if (!fileList || fileList.length === 0) return;
    const files = Array.from(fileList);

    files.forEach((file) => {
      const key = `${file.name}-${file.size}-${Date.now()}-${Math.random()}`;
      const ext = file.name.split(".").pop()?.toLowerCase() ?? "";

      // Client-side validation: type + size.
      if (!ACCEPTED_EXT.includes(ext)) {
        setItems((prev) => [...prev, { key, name: file.name, size: file.size, progress: 0, error: "نوع فایل مجاز نیست" }]);
        return;
      }
      if (file.size > maxSizeMb * 1024 * 1024) {
        setItems((prev) => [...prev, { key, name: file.name, size: file.size, progress: 0, error: `حجم فایل بیش از ${toFa(maxSizeMb)} مگابایت است` }]);
        return;
      }

      setItems((prev) => [...prev, { key, name: file.name, size: file.size, progress: 0 }]);
      uploadOne(file, (pct) => {
        setItems((prev) => prev.map((i) => (i.key === key ? { ...i, progress: pct } : i)));
      })
        .then((done) => {
          setItems((prev) => prev.map((i) => (i.key === key ? { ...i, progress: 100, done } : i)));
        })
        .catch((err: Error) => {
          setItems((prev) => prev.map((i) => (i.key === key ? { ...i, error: err.message } : i)));
        });
    });

    if (inputRef.current) inputRef.current.value = "";
  }

  function remove(key: string) {
    setItems((prev) => prev.filter((i) => i.key !== key));
  }

  return (
    <div className={cn("space-y-2", className)}>
      <label
        onDragOver={(e) => e.preventDefault()}
        onDrop={(e) => {
          e.preventDefault();
          handleFiles(e.dataTransfer.files);
        }}
        className="flex cursor-pointer flex-col items-center justify-center gap-2 rounded-lg border border-dashed bg-background py-8 text-sm text-muted-foreground hover:border-primary/40 hover:bg-accent"
      >
        <Upload className="h-5 w-5" />
        <span>فایل‌ها را اینجا رها کنید یا کلیک کنید</span>
        <span className="text-xs">حداکثر {toFa(maxSizeMb)} مگابایت برای هر فایل</span>
        <input
          ref={inputRef}
          type="file"
          multiple={multiple}
          className="hidden"
          onChange={(e) => handleFiles(e.target.files)}
        />
      </label>

      {items.length > 0 && (
        <ul className="space-y-2">
          {items.map((it) => (
            <li
              key={it.key}
              className={cn(
                "flex items-center gap-3 rounded-lg border bg-card px-3 py-2 text-sm",
                it.error && "border-destructive/40 bg-destructive/5"
              )}
            >
              {it.done ? (
                <FileText className="h-4 w-4 shrink-0 text-success" />
              ) : it.error ? (
                <X className="h-4 w-4 shrink-0 text-destructive" />
              ) : (
                <Loader2 className="h-4 w-4 shrink-0 animate-spin text-primary" />
              )}
              <div className="min-w-0 flex-1">
                <div className="truncate font-medium">{it.name}</div>
                {it.error ? (
                  <div className="text-xs text-destructive">{it.error}</div>
                ) : it.done ? (
                  <div className="text-xs text-muted-foreground">{fmtSize(it.size)}</div>
                ) : (
                  <div className="mt-1 h-1.5 w-full overflow-hidden rounded-full bg-muted">
                    <div className="h-full gradient-brand transition-all" style={{ width: `${it.progress}%` }} />
                  </div>
                )}
              </div>
              <button
                type="button"
                onClick={() => remove(it.key)}
                className="shrink-0 rounded p-1 text-muted-foreground hover:bg-accent hover:text-foreground"
                aria-label="حذف"
              >
                <X className="h-4 w-4" />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
