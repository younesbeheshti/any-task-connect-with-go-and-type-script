// Persian (Farsi) localization helpers

const FA_DIGITS = ["۰","۱","۲","۳","۴","۵","۶","۷","۸","۹"];

export function toFa(input: string | number): string {
  return String(input).replace(/\d/g, d => FA_DIGITS[+d]);
}

/** Format toman with Persian digits and Persian thousand separator. */
export function toman(amount: number, withUnit = true): string {
  const abs = Math.abs(amount);
  const formatted = abs.toLocaleString("en-US").replace(/,/g, "٬");
  const persian = toFa(formatted);
  const sign = amount < 0 ? "-" : "";
  return withUnit ? `${sign}${persian} تومان` : `${sign}${persian}`;
}

/** Convert Gregorian Date to Jalali (yyyy/mm/dd) string with Persian digits. */
export function toJalali(date: Date = new Date()): string {
  // Algorithm from Behdad Esfahbod's libjalali (compact form)
  const gy = date.getFullYear();
  const gm = date.getMonth() + 1;
  const gd = date.getDate();
  const g_d_m = [0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334];
  let jy = gy <= 1600 ? 0 : 979;
  let gy2 = gy <= 1600 ? gy - 621 : gy - 1600;
  const gy3 = gm > 2 ? gy2 + 1 : gy2;
  let days = 365 * gy2 + Math.floor((gy3 + 3) / 4) - Math.floor((gy3 + 99) / 100) + Math.floor((gy3 + 399) / 400) - 80 + gd + g_d_m[gm - 1];
  jy += 33 * Math.floor(days / 12053);
  days %= 12053;
  jy += 4 * Math.floor(days / 1461);
  days %= 1461;
  if (days > 365) {
    jy += Math.floor((days - 1) / 365);
    days = (days - 1) % 365;
  }
  const jm = days < 186 ? 1 + Math.floor(days / 31) : 7 + Math.floor((days - 186) / 30);
  const jd = 1 + (days < 186 ? days % 31 : (days - 186) % 30);
  return toFa(`${jy}/${String(jm).padStart(2, "0")}/${String(jd).padStart(2, "0")}`);
}

export function relativeTimeFa(minutesAgo: number): string {
  if (minutesAgo < 1) return "همین الان";
  if (minutesAgo < 60) return `${toFa(minutesAgo)} دقیقه پیش`;
  const h = Math.floor(minutesAgo / 60);
  if (h < 24) return `${toFa(h)} ساعت پیش`;
  const d = Math.floor(h / 24);
  if (d === 1) return "دیروز";
  if (d < 30) return `${toFa(d)} روز پیش`;
  const m = Math.floor(d / 30);
  return `${toFa(m)} ماه پیش`;
}

export function deadlineFa(daysFromNow: number): string {
  if (daysFromNow === 0) return "امروز";
  if (daysFromNow === 1) return "فردا";
  if (daysFromNow < 0) return `${toFa(Math.abs(daysFromNow))} روز گذشته`;
  return `${toFa(daysFromNow)} روز دیگر`;
}
