import { createFileRoute, Link } from "@tanstack/react-router";
import { ArrowRight, Shield, Wallet, MessagesSquare, MapPin, Star, FileText, Stethoscope, GraduationCap, Scale, Landmark, Package, CheckCircle2, Zap, Globe2, Lock, Sparkles, ChevronDown } from "lucide-react";
import { Logo } from "@/components/logo";
import { Button } from "@/components/ui/button";
import { Bridge } from "@/components/brand";
import { useTheme } from "@/components/theme-provider";
import { Sun, Moon } from "lucide-react";
import { useState } from "react";

export const Route = createFileRoute("/")({
  head: () => ({
    meta: [
      { title: "TaskBridge — Get real-world tasks done from anywhere" },
      { name: "description", content: "Find trusted local agents to complete real-world tasks in any city — administrative, medical, university, legal and government. Escrow-protected payments." },
      { property: "og:title", content: "TaskBridge — Get real-world tasks done from anywhere" },
      { property: "og:description", content: "A trusted marketplace connecting requesters and local agents for real-world errands. Escrow-protected, rated, and verified." },
    ],
  }),
  component: Landing,
});

const categories = [
  { icon: FileText, label: "Administrative", desc: "Forms, filings, submissions" },
  { icon: Stethoscope, label: "Medical", desc: "Prescriptions, results, pharmacy" },
  { icon: GraduationCap, label: "University", desc: "Transcripts, registrar, exams" },
  { icon: Scale, label: "Legal", desc: "Notary, apostille, signatures" },
  { icon: Landmark, label: "Government", desc: "Permits, IDs, renewals" },
  { icon: Package, label: "Pickup & Delivery", desc: "Same-city handoffs" },
];

const faqs = [
  { q: "How does escrow work?", a: "Funds are held safely the moment you create a task. We only release them to the agent after you verify the task is complete." },
  { q: "How are agents vetted?", a: "Every agent passes ID verification, address verification, and starts with a probation badge. Ratings and reviews build their reputation over time." },
  { q: "What if something goes wrong?", a: "Our 24/7 dispute team can pause escrow, mediate, and refund within 72 hours. Every task is covered up to $500." },
  { q: "Where is TaskBridge available?", a: "We're live in 12 cities across 9 countries and adding more each month. New agents can apply from anywhere." },
];

function Landing() {
  return (
    <div className="min-h-screen bg-background">
      <Nav />
      <Hero />
      <LogosStrip />
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
          <a href="#how" className="hover:text-foreground">How it works</a>
          <a href="#categories" className="hover:text-foreground">Categories</a>
          <a href="#benefits" className="hover:text-foreground">Benefits</a>
          <a href="#faq" className="hover:text-foreground">FAQ</a>
        </nav>
        <div className="flex items-center gap-2">
          <button onClick={toggle} className="grid h-9 w-9 place-items-center rounded-lg border bg-card hover:bg-accent" aria-label="Toggle theme">
            {theme === "dark" ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </button>
          <Link to="/login" className="hidden text-sm font-medium text-muted-foreground hover:text-foreground sm:inline">Sign in</Link>
          <Link to="/register" className="inline-flex items-center gap-1.5 rounded-lg gradient-brand px-4 py-2 text-sm font-semibold text-white shadow-soft hover:opacity-95">
            Get started <ArrowRight className="h-4 w-4" />
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
      <div className="absolute -left-32 top-20 h-72 w-72 rounded-full bg-primary/20 blur-3xl" />
      <div className="absolute -right-20 top-40 h-72 w-72 rounded-full bg-secondary/20 blur-3xl" />

      <div className="relative mx-auto max-w-7xl px-4 pb-24 pt-16 sm:px-6 lg:pt-24">
        <div className="mx-auto max-w-3xl text-center">
          <div className="inline-flex items-center gap-2 rounded-full border bg-card/60 px-3 py-1 text-xs font-medium text-muted-foreground backdrop-blur">
            <span className="inline-block h-1.5 w-1.5 rounded-full bg-success animate-pulse" />
            Now live in 12 cities · 9 countries
          </div>
          <h1 className="mt-6 text-4xl font-bold leading-[1.05] tracking-tight sm:text-6xl lg:text-7xl">
            Get your tasks done <br className="hidden sm:block" />
            <span className="gradient-text">from anywhere.</span>
          </h1>
          <p className="mx-auto mt-6 max-w-2xl text-base text-muted-foreground sm:text-lg">
            TaskBridge connects you with trusted local agents to complete real-world errands —
            administrative, medical, university, legal and government tasks in any city.
          </p>
          <div className="mt-8 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link to="/register" className="inline-flex w-full items-center justify-center gap-2 rounded-xl gradient-brand px-6 py-3 text-sm font-semibold text-white shadow-glow hover:opacity-95 sm:w-auto">
              Create a task <ArrowRight className="h-4 w-4" />
            </Link>
            <Link to="/register" className="inline-flex w-full items-center justify-center gap-2 rounded-xl border bg-card px-6 py-3 text-sm font-semibold hover:bg-accent sm:w-auto">
              Become an agent
            </Link>
          </div>
          <div className="mt-6 flex items-center justify-center gap-6 text-xs text-muted-foreground">
            <span className="inline-flex items-center gap-1.5"><Shield className="h-3.5 w-3.5 text-primary" /> Escrow protected</span>
            <span className="inline-flex items-center gap-1.5"><Star className="h-3.5 w-3.5 text-warning" /> 4.9 avg rating</span>
            <span className="inline-flex items-center gap-1.5"><Lock className="h-3.5 w-3.5 text-success" /> ID verified</span>
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
        <div className="flex items-center gap-2 border-b px-4 py-3">
          <div className="flex gap-1.5">
            <span className="h-2.5 w-2.5 rounded-full bg-destructive/70" />
            <span className="h-2.5 w-2.5 rounded-full bg-warning/70" />
            <span className="h-2.5 w-2.5 rounded-full bg-success/70" />
          </div>
          <div className="mx-auto rounded-md bg-muted px-3 py-1 text-xs text-muted-foreground">app.taskbridge.io/tasks</div>
        </div>
        <div className="grid gap-4 p-4 sm:grid-cols-3">
          {[
            { city: "Tehran", title: "Pickup university transcripts", price: 85, badge: "Open", color: "bg-primary/10 text-primary" },
            { city: "Istanbul", title: "Residence renewal at municipality", price: 60, badge: "Assigned", color: "bg-secondary/15 text-secondary-foreground" },
            { city: "Dubai", title: "Notarize 3 documents", price: 220, badge: "In Progress", color: "bg-warning/15 text-warning-foreground" },
          ].map((t, i) => (
            <div key={i} className="rounded-xl border bg-background p-4 shadow-soft">
              <div className="flex items-center justify-between text-xs">
                <span className="inline-flex items-center gap-1 text-muted-foreground"><MapPin className="h-3 w-3" />{t.city}</span>
                <span className={`rounded-full px-2 py-0.5 font-medium ${t.color}`}>{t.badge}</span>
              </div>
              <div className="mt-2 font-semibold">{t.title}</div>
              <div className="mt-3 flex items-center justify-between">
                <span className="text-xs text-muted-foreground">3 applicants</span>
                <span className="rounded-md bg-primary/10 px-2 py-0.5 text-sm font-bold text-primary">${t.price}</span>
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

function LogosStrip() {
  return (
    <div className="border-y bg-muted/40">
      <div className="mx-auto flex max-w-7xl flex-wrap items-center justify-center gap-x-10 gap-y-3 px-4 py-6 text-sm font-medium text-muted-foreground sm:px-6">
        <span>Trusted by teams & individuals at</span>
        {["Notion", "Linear", "Vercel", "Figma", "Stripe", "Airbnb"].map(n => (
          <span key={n} className="font-display text-lg font-bold opacity-60">{n}</span>
        ))}
      </div>
    </div>
  );
}

function HowItWorks() {
  const steps = [
    { n: "01", icon: FileText, title: "Create a task", desc: "Describe what you need, set a budget, and where it should happen." },
    { n: "02", icon: Sparkles, title: "Pick the perfect agent", desc: "Review proposals, ratings, and prices. Chat before you commit." },
    { n: "03", icon: Shield, title: "Pay only when done", desc: "Funds sit safely in escrow until you verify the task is complete." },
  ];
  return (
    <section id="how" className="mx-auto max-w-7xl px-4 py-24 sm:px-6">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">How it works</h2>
        <p className="mt-3 text-muted-foreground">From request to receipt — in three steps.</p>
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
            <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">Featured categories</h2>
            <p className="mt-2 text-muted-foreground">Whatever you need done, there's an agent for it.</p>
          </div>
          <Link to="/app/tasks" className="inline-flex items-center gap-1 text-sm font-medium text-primary hover:underline">
            Browse marketplace <ArrowRight className="h-4 w-4" />
          </Link>
        </div>
        <div className="mt-10 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {categories.map(c => (
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
    { icon: Shield, title: "Escrow on every task", desc: "Funds are held by TaskBridge until you confirm the job is done right." },
    { icon: Star, title: "Vetted agents, real ratings", desc: "ID verified, background checked, and rated by people like you." },
    { icon: Zap, title: "Built for speed", desc: "Most tasks get matched in under 30 minutes. Same-day completion common." },
    { icon: MessagesSquare, title: "Built-in chat", desc: "Coordinate, share files, and track progress in real-time." },
    { icon: Wallet, title: "Flexible payments", desc: "Top up your wallet, pay-as-you-go, or invoice for teams." },
    { icon: Globe2, title: "Global coverage", desc: "12 cities across 9 countries — and growing every week." },
  ];
  return (
    <section id="benefits" className="mx-auto max-w-7xl px-4 py-24 sm:px-6">
      <div className="mx-auto max-w-2xl text-center">
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">Why TaskBridge</h2>
        <p className="mt-3 text-muted-foreground">A trustworthy marketplace, built for the way real life works.</p>
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
    { name: "Leyla Hashemi", role: "Student in Berlin", quote: "I needed sealed transcripts from Tehran University. Matched with an agent in 10 minutes — had them in Berlin within a week.", rating: 5 },
    { name: "Marcus Vogel", role: "Expat in Dubai", quote: "The escrow gave me peace of mind. Agent notarized 3 documents and the funds released the moment I verified.", rating: 5 },
    { name: "Priya Ananth", role: "Founder", quote: "We use TaskBridge to handle local admin for our remote team. The chat + receipts are the killer features.", rating: 5 },
  ];
  return (
    <section className="bg-muted/30 py-24">
      <div className="mx-auto max-w-7xl px-4 sm:px-6">
        <div className="mx-auto max-w-2xl text-center">
          <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">Loved by requesters everywhere</h2>
        </div>
        <div className="mt-12 grid gap-6 md:grid-cols-3">
          {items.map(t => (
            <figure key={t.name} className="flex flex-col rounded-2xl border bg-card p-6 shadow-soft">
              <div className="flex gap-0.5 text-warning">
                {Array.from({ length: t.rating }).map((_, i) => <Star key={i} className="h-4 w-4 fill-current" />)}
              </div>
              <blockquote className="mt-4 flex-1 text-sm leading-relaxed">"{t.quote}"</blockquote>
              <figcaption className="mt-6 flex items-center gap-3">
                <span className="grid h-10 w-10 place-items-center rounded-full gradient-brand text-sm font-bold text-white">
                  {t.name.split(" ").map(n => n[0]).join("")}
                </span>
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
        <h2 className="text-3xl font-bold tracking-tight sm:text-4xl">Frequently asked questions</h2>
      </div>
      <div className="mt-10 divide-y rounded-2xl border bg-card shadow-soft">
        {faqs.map((f, i) => (
          <button key={i} onClick={() => setOpen(open === i ? -1 : i)} className="block w-full px-6 py-5 text-left">
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
          <h2 className="font-display text-3xl font-bold sm:text-4xl">Ready to get something done?</h2>
          <p className="mx-auto mt-3 max-w-xl text-white/85">Post your first task in under a minute. No card required to browse.</p>
          <div className="mt-7 flex flex-col items-center justify-center gap-3 sm:flex-row">
            <Link to="/register" className="rounded-xl bg-white px-6 py-3 text-sm font-semibold text-primary shadow-soft hover:bg-white/90">
              Create a task
            </Link>
            <Link to="/register" className="rounded-xl border border-white/30 bg-white/10 px-6 py-3 text-sm font-semibold backdrop-blur hover:bg-white/20">
              Become an agent
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
          <p className="mt-3 max-w-xs text-sm text-muted-foreground">A trusted marketplace for real-world tasks, anywhere.</p>
        </div>
        {[
          { title: "Product", links: ["How it works", "Categories", "Pricing", "Marketplace"] },
          { title: "Company", links: ["About", "Careers", "Press", "Contact"] },
          { title: "Legal", links: ["Privacy", "Terms", "Escrow policy", "Disputes"] },
        ].map(c => (
          <div key={c.title}>
            <div className="text-sm font-semibold">{c.title}</div>
            <ul className="mt-3 space-y-2 text-sm text-muted-foreground">
              {c.links.map(l => <li key={l}><a href="#" className="hover:text-foreground">{l}</a></li>)}
            </ul>
          </div>
        ))}
      </div>
      <div className="border-t">
        <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-3 px-4 py-5 text-xs text-muted-foreground sm:flex-row sm:px-6">
          <span>© 2026 TaskBridge, Inc.</span>
          <span className="inline-flex items-center gap-1"><CheckCircle2 className="h-3.5 w-3.5 text-success" /> SOC 2 · GDPR compliant</span>
        </div>
      </div>
    </footer>
  );
}
