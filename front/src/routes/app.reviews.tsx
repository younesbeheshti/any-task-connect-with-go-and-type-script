import { createFileRoute } from "@tanstack/react-router";
import { useEffect, useState } from "react";
import { Star, MessageSquare } from "lucide-react";
import { toFa } from "@/lib/fa";
import { useAuth } from "@/lib/auth-context";

export const Route = createFileRoute("/app/reviews")({
  head: () => ({ meta: [{ title: "نظرات و امتیازها — تسک‌بریج" }] }),
  component: ReviewsPage,
});

const STORAGE_KEY = "tb-auth";
const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";
function getToken() {
  try { return JSON.parse(localStorage.getItem(STORAGE_KEY) ?? "{}").token ?? null; }
  catch { return null; }
}

type Review = {
  id: string;
  rating: number;
  comment: string;
  createdAt: string;
  taskId: string;
  reviewerId: string;
  reviewedUserId: string;
};

function ReviewsPage() {
  const { user } = useAuth();
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!user) return;
    const token = getToken();
    if (!token) { setError("لطفاً وارد شوید"); setLoading(false); return; }
    fetch(`${API_BASE}/v1/users/${user.id}/reviews`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then(r => r.json())
      .then(d => setReviews(d.reviews ?? []))
      .catch(() => setError("خطا در بارگذاری نظرات"))
      .finally(() => setLoading(false));
  }, [user?.id]);

  const avgRating = reviews.length > 0
    ? reviews.reduce((s, r) => s + r.rating, 0) / reviews.length
    : 0;

  if (loading) return <Skeleton />;
  if (error) return <ErrBox msg={error} />;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight">نظرات و امتیازها</h1>
        <p className="text-sm text-muted-foreground">نظراتی که دیگران درباره شما ثبت کرده‌اند</p>
      </div>

      {reviews.length > 0 && (
        <div className="rounded-2xl border bg-card p-6 shadow-soft">
          <div className="flex items-center gap-6">
            <div className="text-center">
              <div className="font-display text-5xl font-bold tabular-nums">{avgRating.toFixed(1)}</div>
              <div className="mt-1 flex justify-center gap-0.5">
                {[1, 2, 3, 4, 5].map(i => (
                  <Star key={i} className={`h-4 w-4 ${i <= Math.round(avgRating) ? "fill-warning text-warning" : "text-muted-foreground/30"}`} />
                ))}
              </div>
              <div className="mt-1 text-xs text-muted-foreground">{toFa(reviews.length)} نظر</div>
            </div>
            <div className="flex-1 space-y-1.5">
              {[5, 4, 3, 2, 1].map(star => {
                const count = reviews.filter(r => r.rating === star).length;
                const pct = reviews.length > 0 ? (count / reviews.length) * 100 : 0;
                return (
                  <div key={star} className="flex items-center gap-2 text-xs">
                    <span className="w-3 text-muted-foreground">{star}</span>
                    <Star className="h-3 w-3 fill-warning text-warning" />
                    <div className="flex-1 rounded-full bg-muted h-2 overflow-hidden">
                      <div className="h-full rounded-full bg-warning" style={{ width: `${pct}%` }} />
                    </div>
                    <span className="w-6 text-left text-muted-foreground">{toFa(count)}</span>
                  </div>
                );
              })}
            </div>
          </div>
        </div>
      )}

      {reviews.length === 0 ? (
        <div className="rounded-2xl border bg-card p-12 text-center shadow-soft">
          <MessageSquare className="mx-auto h-12 w-12 text-muted-foreground/40" />
          <h3 className="mt-4 font-semibold">هنوز نظری ثبت نشده</h3>
          <p className="mt-1 text-sm text-muted-foreground">پس از تکمیل درخواست‌ها، طرف مقابل می‌تواند نظر ثبت کند.</p>
        </div>
      ) : (
        <div className="space-y-3">
          {reviews.map(r => (
            <ReviewCard key={r.id} review={r} />
          ))}
        </div>
      )}
    </div>
  );
}

function ReviewCard({ review }: { review: Review }) {
  return (
    <div className="rounded-2xl border bg-card p-5 shadow-soft">
      <div className="flex items-start justify-between gap-3">
        <div className="flex gap-0.5">
          {[1, 2, 3, 4, 5].map(i => (
            <Star key={i} className={`h-4 w-4 ${i <= review.rating ? "fill-warning text-warning" : "text-muted-foreground/30"}`} />
          ))}
        </div>
        <span className="text-xs text-muted-foreground tabular-nums">{review.createdAt?.slice(0, 10)}</span>
      </div>
      {review.comment && (
        <p className="mt-3 text-sm leading-relaxed text-foreground/90">{review.comment}</p>
      )}
      <div className="mt-3 text-xs text-muted-foreground">
        برای درخواست <span className="font-mono text-primary">{review.taskId.slice(0, 8)}</span>
      </div>
    </div>
  );
}

function Skeleton() {
  return (
    <div className="space-y-3">
      {[1, 2, 3].map(i => (
        <div key={i} className="rounded-2xl border bg-card p-5 shadow-soft animate-pulse">
          <div className="space-y-2">
            <div className="flex gap-1">{[1, 2, 3, 4, 5].map(j => <div key={j} className="h-4 w-4 rounded bg-muted" />)}</div>
            <div className="h-4 w-3/4 rounded bg-muted" />
          </div>
        </div>
      ))}
    </div>
  );
}

function ErrBox({ msg }: { msg: string }) {
  return <div className="rounded-2xl border border-destructive/30 bg-destructive/5 p-8 text-center"><p className="text-sm text-destructive">{msg}</p></div>;
}
