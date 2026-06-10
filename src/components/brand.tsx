export function Bridge({ size = 28 }: { size?: number }) {
  return (
    <span
      className="grid place-items-center rounded-xl gradient-brand shadow-glow"
      style={{ width: size, height: size }}
      aria-hidden
    >
      <svg width={size * 0.6} height={size * 0.6} viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.4" strokeLinecap="round" strokeLinejoin="round">
        <path d="M3 12c3-3 6-3 9 0s6 3 9 0" />
        <path d="M5 18V12" />
        <path d="M19 18V12" />
        <path d="M3 18h18" />
      </svg>
    </span>
  );
}
