"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  useCallback,
  ReactNode,
} from "react";
import { User, UserStats } from "@/types";
import { login as apiLogin, register as apiRegister, fetchMe } from "@/lib/api";

interface AuthContextType {
  user: User | null;
  stats: UserStats | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  register: (username: string, password: string) => Promise<void>;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [stats, setStats] = useState<UserStats | null>(null);
  const [loading, setLoading] = useState(true);

  const refreshUser = useCallback(async () => {
    try {
      const token = localStorage.getItem("token");
      if (!token) {
        setUser(null);
        setStats(null);
        return;
      }
      const data = await fetchMe();
      setUser(data.user);
      setStats(data.stats);
    } catch {
      localStorage.removeItem("token");
      setUser(null);
      setStats(null);
    }
  }, []);

  useEffect(() => {
    refreshUser().finally(() => setLoading(false));
  }, [refreshUser]);

  const login = async (username: string, password: string) => {
    const data = await apiLogin(username, password);
    localStorage.setItem("token", data.token);
    setUser(data.user);
    await refreshUser();
  };

  const register = async (username: string, password: string) => {
    const data = await apiRegister(username, password);
    localStorage.setItem("token", data.token);
    setUser(data.user);
    await refreshUser();
  };

  const logout = () => {
    localStorage.removeItem("token");
    setUser(null);
    setStats(null);
  };

  return (
    <AuthContext.Provider
      value={{ user, stats, loading, login, register, logout, refreshUser }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
