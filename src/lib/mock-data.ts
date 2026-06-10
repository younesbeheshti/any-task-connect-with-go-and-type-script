export type TaskStatus = "open" | "assigned" | "in_progress" | "completed" | "verified" | "paid";

export type Task = {
  id: string;
  title: string;
  description: string;
  category: string;
  city: string;
  budget: number;
  deadline: string;
  status: TaskStatus;
  postedBy: string;
  postedAgo: string;
  applicants: number;
  attachments?: string[];
};

export const categories = [
  { id: "admin", label: "Administrative", icon: "FileText" },
  { id: "medical", label: "Medical", icon: "Stethoscope" },
  { id: "university", label: "University", icon: "GraduationCap" },
  { id: "legal", label: "Legal", icon: "Scale" },
  { id: "government", label: "Government", icon: "Landmark" },
  { id: "delivery", label: "Pickup & Delivery", icon: "Package" },
  { id: "shopping", label: "Shopping", icon: "ShoppingBag" },
  { id: "other", label: "Other", icon: "Sparkles" },
];

export const cities = ["Tehran", "Istanbul", "Dubai", "Cairo", "Madrid", "Berlin", "London", "Paris", "New York"];

export const tasks: Task[] = [
  { id: "TB-1042", title: "Pick up university transcripts from Tehran University", description: "Need someone to visit the registrar office, request official sealed transcripts, and ship them to Berlin via DHL.", category: "University", city: "Tehran", budget: 85, deadline: "in 4 days", status: "open", postedBy: "Sara M.", postedAgo: "2h ago", applicants: 7 },
  { id: "TB-1041", title: "Submit residence renewal at municipality office", description: "Visit the local municipality, hand-deliver documents, and obtain a stamped receipt.", category: "Government", city: "Istanbul", budget: 60, deadline: "in 2 days", status: "assigned", postedBy: "Omar A.", postedAgo: "5h ago", applicants: 12 },
  { id: "TB-1039", title: "Collect medical test results from clinic", description: "Pickup lab results from Pars Clinic and securely upload scanned PDFs.", category: "Medical", city: "Tehran", budget: 35, deadline: "tomorrow", status: "in_progress", postedBy: "Lina K.", postedAgo: "1d ago", applicants: 4 },
  { id: "TB-1036", title: "Notarize and apostille three documents", description: "Three educational documents need notary + apostille. Photos attached.", category: "Legal", city: "Dubai", budget: 220, deadline: "in 7 days", status: "completed", postedBy: "Yusuf B.", postedAgo: "3d ago", applicants: 9 },
  { id: "TB-1031", title: "Deliver gift package across town", description: "Pickup wrapped package from Mirdamad and deliver to Pasdaran. Same-day.", category: "Delivery", city: "Tehran", budget: 25, deadline: "today", status: "verified", postedBy: "Mona R.", postedAgo: "4d ago", applicants: 3 },
  { id: "TB-1028", title: "Buy specific medication from pharmacy", description: "Prescription required. Will be sent securely. Reimbursement + service fee.", category: "Medical", city: "Cairo", budget: 40, deadline: "in 3 days", status: "paid", postedBy: "Hany S.", postedAgo: "6d ago", applicants: 6 },
  { id: "TB-1025", title: "Translate and submit visa form", description: "Translate the attached form and submit at consulate office.", category: "Administrative", city: "Madrid", budget: 110, deadline: "in 5 days", status: "open", postedBy: "Diego P.", postedAgo: "1d ago", applicants: 5 },
  { id: "TB-1021", title: "Replace lost ID card at registry", description: "Submit lost ID report and request a duplicate.", category: "Government", city: "Berlin", budget: 95, deadline: "in 9 days", status: "open", postedBy: "Anna L.", postedAgo: "8h ago", applicants: 2 },
];

export const agents = [
  { id: "ag-01", name: "Mahdi T.", city: "Tehran", rating: 4.9, tasks: 184, badge: "Top Rated", price: 80, eta: "Today", bio: "5+ years helping expats with admin & university tasks." },
  { id: "ag-02", name: "Aylin O.", city: "Istanbul", rating: 4.8, tasks: 132, badge: "Verified", price: 65, eta: "Tomorrow", bio: "Bilingual TR/EN, specializes in municipality paperwork." },
  { id: "ag-03", name: "Karim N.", city: "Dubai", rating: 5.0, tasks: 97, badge: "Top Rated", price: 110, eta: "Today", bio: "Notarization & legal runs for international students." },
];

export const notifications = [
  { id: "n1", type: "task_assigned", title: "Task assigned to Mahdi T.", desc: "TB-1042 · Pick up university transcripts", time: "2m ago", unread: true },
  { id: "n2", type: "payment", title: "Payment of $220 released", desc: "TB-1036 · Notarize and apostille three documents", time: "1h ago", unread: true },
  { id: "n3", type: "application", title: "New application from Aylin O.", desc: "TB-1041 · Submit residence renewal", time: "3h ago", unread: true },
  { id: "n4", type: "task_completed", title: "Task verified — funds released", desc: "TB-1031 · Deliver gift package", time: "yesterday", unread: false },
  { id: "n5", type: "message", title: "New message from Karim N.", desc: "I'm at the notary office now…", time: "yesterday", unread: false },
];

export const transactions = [
  { id: "tx-901", date: "Jun 9, 2026", description: "Escrow funded — TB-1042", amount: -85, type: "escrow_in" as const, status: "Locked" },
  { id: "tx-900", date: "Jun 8, 2026", description: "Payment released — TB-1036", amount: -220, type: "release" as const, status: "Completed" },
  { id: "tx-899", date: "Jun 7, 2026", description: "Wallet top-up · Visa **24", amount: 500, type: "topup" as const, status: "Completed" },
  { id: "tx-898", date: "Jun 5, 2026", description: "Refund — TB-1019 cancelled", amount: 45, type: "refund" as const, status: "Completed" },
  { id: "tx-897", date: "Jun 3, 2026", description: "Escrow funded — TB-1031", amount: -25, type: "escrow_in" as const, status: "Released" },
];

export const messages = [
  { id: "m1", from: "agent", text: "Hi! I just picked up the transcripts. Heading to DHL now.", time: "10:42" },
  { id: "m2", from: "me", text: "Amazing, thank you! Can you share a photo of the receipt?", time: "10:43" },
  { id: "m3", from: "agent", text: "Sure — sending now.", time: "10:43" },
  { id: "m4", from: "agent", text: "📎 receipt_dhl_2026.jpg", time: "10:44" },
  { id: "m5", from: "me", text: "Perfect. Tracking number when ready please.", time: "10:46" },
];

export const chats = [
  { id: "c1", name: "Mahdi T.", last: "Sure — sending now.", time: "10:43", unread: 2, online: true, taskId: "TB-1042" },
  { id: "c2", name: "Aylin O.", last: "I'll be there at 9am sharp.", time: "Yesterday", unread: 0, online: false, taskId: "TB-1041" },
  { id: "c3", name: "Karim N.", last: "All documents notarized.", time: "Mon", unread: 0, online: true, taskId: "TB-1036" },
];
