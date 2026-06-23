import { createContext, useContext, useState, type ReactNode } from "react";

export type Role = "requester" | "agent" | "admin";

type Ctx = { role: Role; setRole: (r: Role) => void };
const RoleContext = createContext<Ctx>({ role: "requester", setRole: () => {} });

export function RoleProvider({ children }: { children: ReactNode }) {
  const [role, setRoleState] = useState<Role>(() => {
    if (typeof window === "undefined") return "requester";
    const stored = localStorage.getItem("tb-role");
    if (stored === "requester" || stored === "agent" || stored === "admin") return stored as Role;
    return "requester";
  });
  const setRole = (r: Role) => {
    setRoleState(r);
    localStorage.setItem("tb-role", r);
  };
  return <RoleContext.Provider value={{ role, setRole }}>{children}</RoleContext.Provider>;
}

export const useRole = () => useContext(RoleContext);

export const ROLE_LABELS: Record<Role, string> = {
  requester: "درخواست‌دهنده",
  agent: "انجام‌دهنده",
  admin: "مدیر سامانه",
};
