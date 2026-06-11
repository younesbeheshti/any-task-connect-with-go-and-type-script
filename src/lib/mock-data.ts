export type TaskStatus =
  | "posted"               // ثبت شده
  | "awaiting_applicants"  // در انتظار متقاضی
  | "accepted"             // پذیرفته شده
  | "in_progress"          // در حال انجام
  | "completed"            // تکمیل شده
  | "awaiting_verification"// در انتظار تایید
  | "paid"                 // پرداخت شده
  | "cancelled";           // لغو شده

export type Task = {
  id: string;
  title: string;
  description: string;
  category: string;
  city: string;
  budget: number;          // in Toman
  deadline: string;        // Persian display string
  status: TaskStatus;
  postedBy: string;
  postedAgo: string;
  applicants: number;
  attachments?: string[];
};

export const categories = [
  { id: "admin",      label: "اداری",              icon: "FileText" },
  { id: "medical",    label: "پزشکی",              icon: "Stethoscope" },
  { id: "university", label: "دانشگاهی",           icon: "GraduationCap" },
  { id: "legal",      label: "حقوقی",              icon: "Scale" },
  { id: "government", label: "دولتی",              icon: "Landmark" },
  { id: "delivery",   label: "پیک و تحویل",        icon: "Package" },
  { id: "shopping",   label: "خرید",               icon: "ShoppingBag" },
  { id: "other",      label: "سایر",               icon: "Sparkles" },
];

export const cities = [
  "تهران", "مشهد", "اصفهان", "شیراز", "تبریز", "کرج", "اهواز", "قم", "رشت", "کرمان",
];

export const tasks: Task[] = [
  { id: "TB-۱۰۴۲", title: "دریافت ریزنمرات از دانشگاه تهران", description: "نیاز به مراجعه به اداره آموزش، درخواست ریزنمرات مهر شده و ارسال آن به آدرس مشخص.", category: "دانشگاهی", city: "تهران", budget: 850000, deadline: "۴ روز دیگر", status: "awaiting_applicants", postedBy: "سارا محمدی", postedAgo: "۲ ساعت پیش", applicants: 7 },
  { id: "TB-۱۰۴۱", title: "تمدید کارت اقامت در شهرداری", description: "مراجعه به شهرداری منطقه، تحویل مدارک و دریافت رسید رسمی.", category: "دولتی", city: "مشهد", budget: 620000, deadline: "۲ روز دیگر", status: "accepted", postedBy: "امیر احمدی", postedAgo: "۵ ساعت پیش", applicants: 12 },
  { id: "TB-۱۰۳۹", title: "دریافت جواب آزمایش از کلینیک", description: "تحویل جواب آزمایش از کلینیک پارس و بارگذاری اسکن آن.", category: "پزشکی", city: "تهران", budget: 350000, deadline: "فردا", status: "in_progress", postedBy: "لیلا کریمی", postedAgo: "۱ روز پیش", applicants: 4 },
  { id: "TB-۱۰۳۶", title: "تایید و مهر سه سند رسمی در دفترخانه", description: "سه مدرک تحصیلی نیاز به تایید و مهر دارد. تصاویر پیوست شده است.", category: "حقوقی", city: "اصفهان", budget: 2200000, deadline: "۷ روز دیگر", status: "completed", postedBy: "یوسف بهرامی", postedAgo: "۳ روز پیش", applicants: 9 },
  { id: "TB-۱۰۳۱", title: "تحویل بسته در شهر", description: "دریافت بسته از میرداماد و تحویل در پاسداران در همان روز.", category: "پیک و تحویل", city: "تهران", budget: 250000, deadline: "امروز", status: "awaiting_verification", postedBy: "مونا رضایی", postedAgo: "۴ روز پیش", applicants: 3 },
  { id: "TB-۱۰۲۸", title: "خرید دارو از داروخانه", description: "نسخه پیوست شده. هزینه دارو و حق‌الزحمه جداگانه محاسبه می‌شود.", category: "پزشکی", city: "شیراز", budget: 400000, deadline: "۳ روز دیگر", status: "paid", postedBy: "هانی صادقی", postedAgo: "۶ روز پیش", applicants: 6 },
  { id: "TB-۱۰۲۵", title: "ترجمه و ارسال فرم ویزا", description: "ترجمه فرم پیوست‌شده و تحویل آن به دفتر کنسولگری.", category: "اداری", city: "تبریز", budget: 1100000, deadline: "۵ روز دیگر", status: "posted", postedBy: "دیبا پارسا", postedAgo: "۱ روز پیش", applicants: 5 },
  { id: "TB-۱۰۲۱", title: "صدور المثنی کارت ملی", description: "گزارش مفقودی و درخواست صدور المثنی کارت ملی.", category: "دولتی", city: "کرج", budget: 950000, deadline: "۹ روز دیگر", status: "awaiting_applicants", postedBy: "آنا لطیفی", postedAgo: "۸ ساعت پیش", applicants: 2 },
];

export const agents = [
  { id: "ag-01", name: "مهدی توکلی", city: "تهران", rating: 4.9, tasks: 184, badge: "برترین", price: 800000, eta: "امروز", bio: "بیش از ۵ سال تجربه در امور اداری و دانشگاهی." },
  { id: "ag-02", name: "آیلین عثمانی", city: "مشهد", rating: 4.8, tasks: 132, badge: "تایید شده", price: 650000, eta: "فردا", bio: "متخصص امور شهرداری و کارهای دولتی." },
  { id: "ag-03", name: "کریم نظری", city: "اصفهان", rating: 5.0, tasks: 97, badge: "برترین", price: 1100000, eta: "امروز", bio: "تایید رسمی اسناد و امور حقوقی برای دانشجویان." },
];

export const notifications = [
  { id: "n1", type: "task_assigned",  title: "درخواست به مهدی توکلی واگذار شد", desc: "TB-۱۰۴۲ · دریافت ریزنمرات",          time: "۲ دقیقه پیش", unread: true },
  { id: "n2", type: "payment",        title: "پرداخت ۲٬۲۰۰٬۰۰۰ تومان آزاد شد",     desc: "TB-۱۰۳۶ · تایید سه سند",          time: "۱ ساعت پیش", unread: true },
  { id: "n3", type: "application",    title: "درخواست جدید از آیلین عثمانی",        desc: "TB-۱۰۴۱ · تمدید کارت اقامت",      time: "۳ ساعت پیش", unread: true },
  { id: "n4", type: "task_completed", title: "کار تایید شد — پرداخت آزاد گردید",    desc: "TB-۱۰۳۱ · تحویل بسته",            time: "دیروز",      unread: false },
  { id: "n5", type: "message",        title: "پیام جدید از کریم نظری",               desc: "الان در دفترخانه هستم…",           time: "دیروز",      unread: false },
];

export const transactions = [
  { id: "tx-901", date: "۱۴۰۵/۰۳/۱۹", description: "بلوکه در امانت — TB-۱۰۴۲",       amount: -850000,  type: "escrow_in" as const, status: "بلوکه شده" },
  { id: "tx-900", date: "۱۴۰۵/۰۳/۱۸", description: "آزادسازی پرداخت — TB-۱۰۳۶",      amount: -2200000, type: "release" as const,   status: "تکمیل شده" },
  { id: "tx-899", date: "۱۴۰۵/۰۳/۱۷", description: "شارژ کیف پول · کارت بانک ملت",   amount: 5000000,  type: "topup" as const,     status: "تکمیل شده" },
  { id: "tx-898", date: "۱۴۰۵/۰۳/۱۵", description: "بازگشت وجه — TB-۱۰۱۹ لغو شده",    amount: 450000,   type: "refund" as const,    status: "تکمیل شده" },
  { id: "tx-897", date: "۱۴۰۵/۰۳/۱۳", description: "بلوکه در امانت — TB-۱۰۳۱",       amount: -250000,  type: "escrow_in" as const, status: "آزاد شده" },
];

export const messages = [
  { id: "m1", from: "agent", text: "سلام! ریزنمرات رو الان دریافت کردم. دارم به سمت پست می‌رم.", time: "۱۰:۴۲" },
  { id: "m2", from: "me",    text: "عالیه، ممنونم! می‌شه عکس رسید رو بفرستی؟",                  time: "۱۰:۴۳" },
  { id: "m3", from: "agent", text: "حتما — همین الان می‌فرستم.",                                time: "۱۰:۴۳" },
  { id: "m4", from: "agent", text: "📎 رسید_پست_۱۴۰۵.jpg",                                       time: "۱۰:۴۴" },
  { id: "m5", from: "me",    text: "عالی. کد رهگیری رو هم وقتی آماده شد بفرست لطفا.",            time: "۱۰:۴۶" },
];

export const chats = [
  { id: "c1", name: "مهدی توکلی",   last: "حتما — همین الان می‌فرستم.", time: "۱۰:۴۳", unread: 2, online: true,  taskId: "TB-۱۰۴۲" },
  { id: "c2", name: "آیلین عثمانی", last: "ساعت ۹ صبح اونجا هستم.",     time: "دیروز",  unread: 0, online: false, taskId: "TB-۱۰۴۱" },
  { id: "c3", name: "کریم نظری",    last: "همه مدارک تایید شدن.",         time: "دوشنبه", unread: 0, online: true,  taskId: "TB-۱۰۳۶" },
];
