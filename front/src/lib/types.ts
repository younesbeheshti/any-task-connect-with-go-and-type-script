export type ApiTask = {
  id: string;
  title: string;
  description: string;
  status: string;
  budget: number;
  deadline: string | null;
  applicantCount: number;
  createdAt: string;
  category: { id: string; title: string } | null;
  city: { id: string; title: string } | null;
  requester?: { id: string; fullName: string } | null;
  assignedAgent?: { id: string; fullName: string } | null;
};

export type ApiCategory = { id: string; title: string };
export type ApiCity = { id: string; title: string };
