import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

export type AuthUser = {
  id: string;
  fullName: string;
  phone: string;
  email: string | null;
  role: "requester" | "agent" | "admin";
  avatarUrl: string | null;
  rating: number;
  completedCount: number;
};

type RegisterInput = {
  fullName: string;
  phone: string;
  password: string;
  role: "requester" | "agent";
  email?: string;
};

type AuthState = {
  user: AuthUser | null;
  token: string | null;
  loading: boolean;
  isAuthenticated: boolean;
  login: (phone: string, password: string) => Promise<AuthUser>;
  register: (data: RegisterInput) => Promise<AuthUser>;
  logout: () => void;
};

const AuthContext = createContext<AuthState>({
  user: null,
  token: null,
  loading: true,
  isAuthenticated: false,
  login: async () => { throw new Error("not mounted"); },
  register: async () => { throw new Error("not mounted"); },
  logout: () => {},
});

const STORAGE_KEY = "tb-auth";

type ApiUser = {
  id: string;
  fullName: string;
  phone: string;
  email: string | null;
  role: string;
  avatarUrl: string | null;
  rating: number;
  completedCount: number;
};

type AuthResponse = {
  token: string;
  user: ApiUser;
};

function toAuthUser(u: ApiUser): AuthUser {
  const role = u.role === "agent" ? "agent" : u.role === "admin" ? "admin" : "requester";
  return {
    id: u.id,
    fullName: u.fullName,
    phone: u.phone,
    email: u.email,
    role,
    avatarUrl: u.avatarUrl,
    rating: u.rating,
    completedCount: u.completedCount,
  };
}

const API_BASE = import.meta.env.VITE_API_URL ?? "http://localhost:8000";

async function apiPost<T>(path: string, body: unknown): Promise<T> {
  const res = await fetch(`${API_BASE}/v1${path}`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  if (!res.ok) {
    const err = data?.error;
    // If field-level validation errors exist, show the first one
    const fields = err?.fields as Record<string, string> | undefined;
    const firstField = fields ? Object.values(fields)[0] : undefined;
    throw new Error(firstField ?? err?.message ?? "خطای سرور");
  }
  return data as T;
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthUser | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        const parsed = JSON.parse(stored) as { user: AuthUser; token: string };
        setUser(parsed.user);
        setToken(parsed.token);
      }
    } catch {
      // corrupted storage — ignore
    } finally {
      setLoading(false);
    }
  }, []);

  const persist = (u: AuthUser, t: string) => {
    setUser(u);
    setToken(t);
    localStorage.setItem(STORAGE_KEY, JSON.stringify({ user: u, token: t }));
    localStorage.setItem("tb-role", u.role);
  };

  const login = async (identifier: string, password: string): Promise<AuthUser> => {
    const body = identifier.includes("@")
      ? { email: identifier, password }
      : { phone: identifier, password };
    const data = await apiPost<AuthResponse>("/auth/login", body);
    const u = toAuthUser(data.user);
    persist(u, data.token);
    return u;
  };

  const register = async (input: RegisterInput): Promise<AuthUser> => {
    const data = await apiPost<AuthResponse>("/auth/register", input);
    const u = toAuthUser(data.user);
    persist(u, data.token);
    return u;
  };

  const logout = () => {
    setUser(null);
    setToken(null);
    localStorage.removeItem(STORAGE_KEY);
    localStorage.removeItem("tb-role");
  };

  return (
    <AuthContext.Provider value={{ user, token, loading, isAuthenticated: !!token, login, register, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => useContext(AuthContext);
