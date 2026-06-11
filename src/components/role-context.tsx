import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

export type Role = "requester" | "agent" | "admin";

type Ctx = { role: Role; setRole: (r: Role) => void };
const RoleContext = createContext<Ctx>({ role: "requester", setRole: () => {} });

export function RoleProvider({ children }: { children: ReactNode }) {
  const [role, setRoleState] = useState<Role>("requester");
  useEffect(() => {
    const stored = typeof window !== "undefined" ? localStorage.getItem("tb-role") : null;
    if (stored === "requester" || stored === "agent" || stored === "admin") setRoleState(stored);
  }, []);
  const setRole = (r: Role) => {
    setRoleState(r);
    if (typeof window !== "undefined") localStorage.setItem("tb-role", r);
  };
  return <RoleContext.Provider value={{ role, setRole }}>{children}</RoleContext.Provider>;
}

export const useRole = () => useContext(RoleContext);

export const ROLE_LABELS: Record<Role, string> = {
  requester: "درخواست‌دهنده",
  agent: "انجام‌دهنده",
  admin: "مدیر سامانه",
};
