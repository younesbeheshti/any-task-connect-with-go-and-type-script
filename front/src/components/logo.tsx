import { Link } from "@tanstack/react-router";
import { Bridge } from "./brand";

export function Logo({ withText = true, size = 28 }: { withText?: boolean; size?: number }) {
  return (
    <Link to="/" className="inline-flex items-center gap-2 font-display font-bold tracking-tight">
      <Bridge size={size} />
      {withText && <span className="text-lg">تسک‌بریج</span>}
    </Link>
  );
}
