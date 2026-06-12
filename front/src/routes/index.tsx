import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowLeft, Shield, Wallet, MessagesSquare, MapPin, Star, FileText, Stethoscope, GraduationCap, Scale, Landmark, Package, Zap, Lock, Sparkles, ChevronDown, Sun, Moon } from "lucide-react";
import { Logo } from "@/components/logo";
import { useTheme } from "@/components/theme-provider";
import { useState } from "react";
import { toman, toFa } from "@/lib/fa";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "تسک‌بریج — انجام کارهای روزمره از هر کجا" },
      { name: "description", content: "پلتفرم مطمئن انجام درخواست‌های اداری، پزشکی، دانشگاهی، حقوقی و دولتی در سراسر ایران، با پرداخت امانی و انجام‌دهنده‌های تایید شده." },
      { property: "og:title", content: "تسک‌بریج — انجام کارهای روزمره از هر کجا" },
      { property: "og:description", content: "اتصال درخواست‌دهنده‌ها به انجام‌دهنده‌های محلی مورد اعتماد. پرداخت امانی، امتیازدهی، و تایید هویت." },
    ],
  }),
  component: Landing,
});

const cats = [
  { icon: FileText,       label: "اداری",         desc: "فرم‌ها، مراجعات، تشکیل پرونده" },
  { icon: Stethoscope,    label: "پزشکی",         desc: "تحویل دارو، جواب آزمایش" },
  { icon: GraduationCap,  label: "دانشگاهی",      desc: "ریزنمرات، آموزش، امتحانات" },
  { icon: Scale,          label: "حقوقی",         desc: "دفترخانه، تایید مدارک" },
  { icon: Landmark,       label: "دولتی",         desc: "مجوز، کارت ملی، تمدید" },
  { icon: Package,        label: "پیک و تحویل",   desc: "حمل بسته در شهر" },
];

const faqs = [
  { q: "پرداخت امانی چطور کار می‌کند؟",          a: "به محض ثبت درخواست، مبلغ پرداختی شما در امانت پلتفرم بلوکه می‌شود و فقط پس از تایید نهایی شما به انجام‌دهنده پرداخت می‌گردد." },
  { q: "انجام‌دهنده‌ها چطور تایید می‌شوند؟",     a: "هر انجام‌دهنده مراحل احراز هویت با کارت ملی، تایید شماره موبایل و تایید آدرس را طی می‌کند و بر اساس نظرات کاربران رتبه‌بندی می‌شود." },
  { q: "اگر مشکلی پیش بیاید چه می‌شود؟",          a: "تیم پشتیبانی ۲۴ ساعته ما می‌تواند پرداخت را متوقف کند، میانجی‌گری انجام دهد و حداکثر ظرف ۷۲ ساعت وجه شما را بازگرداند." },
  { q: "تسک‌بریج در کدام شهرها فعال است؟",       a: "هم‌اکنون در ۱۰ شهر بزرگ ایران فعال هستیم و هر ماه شهرهای بیشتری به ما اضافه می‌شوند." },
];

function Landing() {
  return (
    <div className="min-h-screen bg-background">
      <Nav />
      <Hero />
      <TrustStrip />
      <HowItWorks />
      <Categories />
      <Benefits />
      <Testimonials />
      <FAQSection />
      <CTA />
      <Footer />
    </div>
  );
}

function Nav() {
  const { theme, toggle } = useTheme();
  return (
    <header className="sticky top-0 z-40 glass">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-4 py-3 sm:px-6">
        <Logo />
        <nav className="hidden items-center gap-8 text-sm font-medium text-muted-foreground md:flex">
          <a href="#how" className="hover:text-foreground">نحوه کار</a>
          <a href="#categories" className="hover:text-foreground">دسته‌بندی‌ها</a>
          <a href="#benefits" className="hover:text-foreground">ویژگی‌ها</a>
          <a href="#faq" className="hover:text-foreground">پرسش‌های متداول</a>
        </nav>
        <div className="flex items-center gap-2">
          <button onClick={toggle} className="grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent" aria-label="تغییر تم">
            {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </button>
          <Link to="/login" className="hidden text-sm font-medium text-muted-foreground hover:text-foreground sm:inline">ورود</Link>
          <Link to="/register" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            شروع کنید <ArrowLeft className="h-4 w-4" />
          </Link>
        </div>
      </div>
    </header>
  );
}

function Hero() {
  return (
    <section className="relative overflow-hidden">
      <div className="absolute inset-0 grid-bg [mask-image:radial-gradient(ellipse_at_center,black_30%,transparent_70%)]" />
      <div className="absolute -right-32 top-20 h-72 w-72 rounded-full bg-primary/20 blur-3xl" />
      <div className="absolute -left-20 top-40 h-72 w-72 rounded-full bg-secondary/20 blur-3xl" />

      <div className="relative mx-auto max-w-7xl px-4 pb-24 pt-16 sm:px-6 lg:pt-24">
        <div className="mx-auto max-w-3xl text-center">
          <div className="inline-flex items-center gap-2 rounded-full border bg-card/60 px-3 py-1 text-xs font-medium text-muted-foreground backdrop-blur">
            <span className="inline-block h-1.5 w-1.5 rounded-full bg-success animate-pulse" />
            فعال در {toFa(10)} شهر ایران
          </div>
          <h1 className="mt-6 text-4xl font-bold leading-[1.4] tracking-tight sm:text-5xl lg:text-6xl">
            کارهای روزمره‌تان را <br className="hidden sm:block" />
            <span className="gradient-text">از هر کجا انجام دهید</span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-base text-muted-foreground sm:text-lg">
            تسک‌بریج شما را به انجام‌دهنده‌های مورد اعتماد محلی متصل می‌کند تا کارهای اداری، پزشکی، دانشگاهی، حقوقی و دولتی شما در هر شهر انجام شود.
          </p>
          <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link to="/register" className="inline-flex w-full items-center justify-center gap-2 rounded-xl gradient-brand px-6 py-3 text-sm font-semibold text-white shadow-glow hover:opacity-95 sm:w-auto">
              ثبت درخواست <ArrowLeft className="h-4 w-4" />
            </Link>
            <Link to="/register" className="inline-flex w-full items-center justify-center gap-2 rounded-xl border bg-card px-6 py-3 text-sm font-semibold hover:bg-accent sm:w-auto">
              انجام‌دهنده شوید
            </Link>
          </div>
          <div className="mt-6 flex flex-wrap items-center justify-center gap-x-6 gap-y-2 text-xs text-muted-foreground">
            <span className="inline-flex items-center gap-1.5"><Shield className="h-3.5 w-3.5 text-primary" /> پرداخت امانی</span>
            <span className="inline-flex items-center gap-1.5"><Star className="h-3.5 w-3.5 text-warning" /> امتیاز {toFa("4.9")}</span>
            <span className="inline-flex items-center gap-1.5"><Lock className="h-3.5 w-3.5 text-success" /> احراز هویت شده</span>
          </div>
        </div>

        <HeroMockup />
      </div>
    </section>
  );
}

function HeroMockup() {
  return (
    <div className="relative mx-auto mt-16 max-w-5xl">
      <div className="absolute -inset-4 rounded-3xl gradient-brand opacity-30 blur-2xl" />
      <div className="relative rounded-2xl border bg-card shadow-elevated">
        <div className="flex items-center gap-2 border-b px-4 py-3" dir="ltr">
          <div className="flex gap-1.5">
            <span className="h-2.5 w-2.5 rounded-full bg-destructive/70" />
            <span className="h-2.5 w-2.5 rounded-full bg-warning/70" />
            <span className="h-2.5 w-2.5 rounded-full bg-success/70" />
          </div>
          <div className="mx-auto rounded-md bg-muted px-3 py-1 text-xs text-muted-foreground">app.taskbridge.ir/tasks</div>
        </div>
        <div className="grid gap-4 p-4 sm:grid-cols-3">
          {[
            { city: "تهران",   title: "دریافت ریزنمرات دانشگاه",  price: 850000,  badge: "در انتظار متقاضی", color: "bg-primary/10 text-primary" },
            { city: "مشهد",    title: "تمدید کارت اقامت",          price: 620000,  badge: "پذیرفته شده",      color: "bg-secondary/15 text-secondary-foreground" },
            { city: "اصفهان",  title: "تایید سه سند دفترخانه",     price: 2200000, badge: "در حال انجام",     color: "bg-warning/15 text-warning-foreground" },
          ].map((t, i) => (
            <div key={i} className="rounded-xl border bg-background p-4 shadow-soft">
              <div className="flex items-center justify-between text-xs">
                <span className="inline-flex items-center gap-1 text-muted-foreground"><MapPin className="h-3 w-3" />{t.city}</span>
                <span className={`rounded-full px-2 py-0.5 font-medium ${t.color}`}>{t.badge}</span>
              </div>
              <div className="mt-2 font-semibold">{t.title}</div>
              <div className="mt-3 flex items-center justify-between">
                <span className="text-xs text-muted-foreground">{toFa(3)} متقاضی</span>
                <span className="rounded-md bg-primary/10 px-2 py-0.5 text-xs font-bold text-primary tabular-nums">{toman(t.price)}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function TrustStrip() {
  return (
    <div className="border-y bg-muted/40">
      <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-center gap-x-10 gap-y-3 px-4 py-6 text-sm font-medium text-muted-foreground sm:px-6">
        <span>قابل اعتماد در شهرهای بزرگ ایران</span>
        {["تهران","مشهد","اصفهان","شیراز","تبریز","کرج"].map(n => (
          <span key={n} className="font-display text-base font-bold opacity-60">{n}</span>
        ))}
      </div>
    </div>
  );
}

function HowItWorks() {
  const steps = [
    { n: "۰۱", icon: FileText, title: "ثبت درخواست",       desc: "بنویسید چه کاری نیاز دارید، مبلغ پیشنهادی و شهر را مشخص کنید." },
    { n: "۰۲", icon: Sparkles, title: "انتخاب انجام‌دهنده", desc: "پیشنهادها، امتیازات و قیمت‌ها را ببینید. پیش از واگذاری گفتگو کنید." },
    { n: "۰۳", icon: Shield,   title: "پرداخت پس از تایید", desc: "وجه شما در امانت می‌ماند تا کار با موفقیت تایید شود." },
  ];
  return (
    <section id="how" className="mx-auto max-w-7xl px-4 py-24 sm:px-6">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">چطور کار می‌کند؟</h2>
        <p className="mt-3 text-muted-foreground">از ثبت درخواست تا تحویل — فقط در سه مرحله.</p>
      </div>
      <div className="mt-14 grid gap-6 md:grid-cols-3">
        {steps.map(s => (
          <div key={s.n} className="group relative rounded-2xl border bg-card p-6 shadow-soft transition-all hover:-translate-y-1 hover:shadow-elevated">
            <div className="flex items-center justify-between">
              <span className="font-display text-5xl font-bold gradient-text">{s.n}</span>
              <span className="grid h-10 w-10 place-items-center rounded-xl bg-primary/10 text-primary"><s.icon className="h-5 w-5" /></span>
            </div>
            <h3 className="mt-6 text-lg font-semibold">{s.title}</h3>
            <p className="mt-1 text-sm text-muted-foreground">{s.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

function Categories() {
  return (
    <section id="categories" className="bg-muted/30 py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6">
        <div className="flex flex-wrap items-end justify-between gap-4">
          <div>
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">دسته‌بندی‌های پرکاربرد</h2>
            <p className="mt-2 text-muted-foreground">هر کاری که نیاز دارید، یک انجام‌دهنده برایش هست.</p>
          </div>
          <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline">
            مشاهده بازار <ArrowLeft className="h-4 w-4" />
          </Link>
        </div>
        <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {cats.map(c => (
            <div key={c.label} className="group flex items-center gap-4 rounded-2xl border bg-card p-5 shadow-soft transition-all hover:border-primary/40 hover:shadow-elevated">
              <span className="grid h-12 w-12 place-items-center rounded-xl gradient-brand text-white shadow-soft transition-transform group-hover:scale-110">
                <c.icon className="h-5 w-5" />
              </span>
              <div>
                <div className="font-semibold">{c.label}</div>
                <div className="text-sm text-muted-foreground">{c.desc}</div>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function Benefits() {
  const items = [
    { icon: Shield,         title: "پرداخت امن امانی",          desc: "وجه شما تا تایید نهایی در پلتفرم نگهداری می‌شود." },
    { icon: Star,           title: "انجام‌دهنده‌های تایید شده",  desc: "احراز هویت با کارت ملی و امتیازدهی توسط کاربران." },
    { icon: Zap,            title: "تطبیق سریع",                desc: "بیشتر درخواست‌ها در کمتر از ۳۰ دقیقه پاسخ می‌گیرند." },
    { icon: MessagesSquare, title: "گفت‌وگوی درون‌برنامه‌ای",   desc: "هماهنگی، ارسال فایل و پیگیری در لحظه." },
    { icon: Wallet,         title: "پرداخت با تومان",           desc: "شارژ کیف پول و پرداخت با کارت‌های شتاب." },
    { icon: Lock,           title: "حریم خصوصی کامل",            desc: "اطلاعات شما رمزنگاری شده و امن نگه‌داری می‌شود." },
  ];
  return (
    <section id="benefits" className="mx-auto max-w-7xl px-4 py-24 sm:px-6">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">چرا تسک‌بریج؟</h2>
        <p className="mt-3 text-muted-foreground">یک بازار مطمئن که برای زندگی واقعی ساخته شده.</p>
      </div>
      <div className="mt-14 grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {items.map(i => (
          <div key={i.title} className="rounded-2xl border bg-card p-6 shadow-soft">
            <span className="grid h-10 w-10 place-items-center rounded-xl bg-primary/10 text-primary"><i.icon className="h-5 w-5" /></span>
            <h3 className="mt-5 font-semibold">{i.title}</h3>
            <p className="mt-1 text-sm text-muted-foreground">{i.desc}</p>
          </div>
        ))}
      </div>
    </section>
  );
}

function Testimonials() {
  const items = [
    { name: "لیلا هاشمی",  role: "دانشجو در آلمان", quote: "نیاز داشتم ریزنمرات از دانشگاه تهران بگیرم. در ۱۰ دقیقه با یک انجام‌دهنده مطابقت پیدا کردم.", rating: 5 },
    { name: "مارکوس وگل",  role: "مهاجر در امارات", quote: "پرداخت امانی خیال من را راحت کرد. انجام‌دهنده ۳ سند را تایید کرد و وجه فوری آزاد شد.",            rating: 5 },
    { name: "پریا آناند",  role: "بنیان‌گذار",       quote: "ما از تسک‌بریج برای کارهای اداری تیم ریموت استفاده می‌کنیم. گفت‌وگو و رسید فوق‌العاده است.",      rating: 5 },
  ];
  return (
    <section className="bg-muted/30 py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">مورد اعتماد کاربران در سراسر کشور</h2>
        </div>
        <div className="mt-12 grid gap-6 md:grid-cols-3">
          {items.map(t => (
            <figure key={t.name} className="flex flex-col rounded-2xl border bg-card p-6 shadow-soft">
              <div className="flex gap-0.5 text-warning">
                {Array.from({ length: t.rating }).map((_, i) => <Star key={i} className="h-4 w-4 fill-current" />)}
              </div>
              <blockquote className="mt-4 flex-1 text-sm leading-relaxed">«{t.quote}»</blockquote>
              <figcaption className="mt-6 flex items-center gap-3">
                <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">{t.name[0]}</span>
                <div>
                  <div className="text-sm font-semibold">{t.name}</div>
                  <div className="text-xs text-muted-foreground">{t.role}</div>
                </div>
              </figcaption>
            </figure>
          ))}
        </div>
      </div>
    </section>
  );
}

function FAQSection() {
  const [open, setOpen] = useState(0);
  return (
    <section id="faq" className="mx-auto max-w-3xl px-4 py-24 sm:px-6">
      <div className="text-center">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">پرسش‌های متداول</h2>
      </div>
      <div className="mt-10 divide-y rounded-2xl border bg-card shadow-soft">
        {faqs.map((f, i) => (
          <button key={i} onClick={() => setOpen(open === i ? -1 : i)} className="block w-full px-6 py-5 text-right">
            <div className="flex items-center justify-between gap-4">
              <span className="font-semibold">{f.q}</span>
              <ChevronDown className={`h-5 w-5 shrink-0 text-muted-foreground transition-transform ${open === i ? "rotate-180" : ""}`} />
            </div>
            {open === i && <p className="mt-3 text-sm text-muted-foreground">{f.a}</p>}
          </button>
        ))}
      </div>
    </section>
  );
}

function CTA() {
  return (
    <section className="mx-auto max-w-7xl px-4 py-16 sm:px-6">
      <div className="relative overflow-hidden rounded-3xl gradient-brand p-10 text-center text-white shadow-elevated sm:p-16">
        <div className="absolute inset-0 grid-bg opacity-20" />
        <div className="relative">
          <h2 className="font-display text-3xl font-bold sm:text-4xl">آماده انجام کارهایتان هستید؟</h2>
          <p className="mx-auto mt-3 max-w-xl text-white/85">در کمتر از یک دقیقه اولین درخواست خود را ثبت کنید. بدون کارت اعتباری.</p>
          <div className="mt-7 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link to="/register" className="rounded-xl bg-white px-6 py-3 text-sm font-semibold text-primary shadow-soft hover:bg-white/90">
              ثبت درخواست
            </Link>
            <Link to="/register" className="rounded-xl border border-white/30 bg-white/10 px-6 py-3 text-sm font-semibold backdrop-blur hover:bg-white/20">
              انجام‌دهنده شوید
            </Link>
          </div>
        </div>
      </div>
    </section>
  );
}

function Footer() {
  return (
    <footer className="border-t bg-muted/30">
      <div className="mx-auto grid max-w-7xl gap-8 px-4 py-12 sm:grid-cols-2 sm:px-6 lg:grid-cols-4">
        <div>
          <Logo />
          <p className="mt-3 text-sm text-muted-foreground">پلتفرم انجام کارهای اداری و خدماتی، با اعتماد در قلب آن.</p>
        </div>
        <div>
          <h4 className="text-sm font-semibold">محصول</h4>
          <ul className="mt-3 space-y-2 text-sm text-muted-foreground">
            <li><a href="#how" className="hover:text-foreground">نحوه کار</a></li>
            <li><a href="#categories" className="hover:text-foreground">دسته‌بندی‌ها</a></li>
            <li><Link to="/app/tasks" className="hover:text-foreground">بازار درخواست‌ها</Link></li>
          </ul>
        </div>
        <div>
          <h4 className="text-sm font-semibold">شرکت</h4>
          <ul className="mt-3 space-y-2 text-sm text-muted-foreground">
            <li><a href="#" className="hover:text-foreground">درباره ما</a></li>
            <li><a href="#" className="hover:text-foreground">شغل و همکاری</a></li>
            <li><a href="#" className="hover:text-foreground">وبلاگ</a></li>
          </ul>
        </div>
        <div>
          <h4 className="text-sm font-semibold">قانونی</h4>
          <ul className="mt-3 space-y-2 text-sm text-muted-foreground">
            <li><a href="#" className="hover:text-foreground">قوانین و مقررات</a></li>
            <li><a href="#" className="hover:text-foreground">حریم خصوصی</a></li>
            <li><a href="#" className="hover:text-foreground">پشتیبانی</a></li>
          </ul>
        </div>
      </div>
      <div className="border-t py-6 text-center text-xs text-muted-foreground">
        © {toFa(1405)} تسک‌بریج — تمامی حقوق محفوظ است.
      </div>
    </footer>
  );
}
