// Mirrors backend TaskResponse (internal/task/handler/dto.go): flat strings for
// category/city, API-status values, and a denormalized applicantsCount.
export type ApiTask = {
  id: string;
  title: string;
  description: string;
  category: string;
  city: string;
  budget: number;
  currency?: string;
  status: string; // API status: posted | awaiting_applicants | accepted | in_progress | completed | awaiting_verification | paid | cancelled
  deadline: string | null;
  requesterId: string;
  assignedAgentId: string | null;
  applicantsCount: number;
  attachmentUrls?: string[];
  createdAt: string;
  updatedAt?: string;
};

export type ApiCategory = { id: string; title: string };
export type ApiCity = { id: string; title: string };
