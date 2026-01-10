// src/stores/authStore.ts
import { create } from 'zustand';
import type { User } from '@/types';
import { fetchProfile, logout as logoutApi } from '@/api/auth';
import { STORAGE_KEYS } from '@/config';

type AuthState = {
  currentUser: User | null;
  isLoading: boolean;
  login: (user: User) => void;
  logout: () => Promise<void>;
  hydrate: () => Promise<void>;
};

// Check if auth cookies exist
const hasAuthCookies = (): boolean => {
  return document.cookie.includes('access_token') || document.cookie.includes('refresh_token');
};

export const useAuthStore = create<AuthState>((set) => ({
  currentUser: null,
  isLoading: false,

  login: (user) => {
    localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(user));
    set({ currentUser: user });
  },

  logout: async () => {
    try {
      await logoutApi();
    } catch (err) {
      console.error('Logout API error:', err);
    }
    localStorage.removeItem(STORAGE_KEYS.user);
    set({ currentUser: null });
  },

  hydrate: async () => {
    // First, try to load from localStorage (fast)
    const raw = localStorage.getItem(STORAGE_KEYS.user);
    if (raw) {
      try {
        const user = JSON.parse(raw);
        set({ currentUser: user });
        return;
      } catch {
        localStorage.removeItem(STORAGE_KEYS.user);
      }
    }

    // If no localStorage but has cookies (e.g., after OAuth redirect), fetch profile
    if (hasAuthCookies()) {
      set({ isLoading: true });
      try {
        const response = await fetchProfile();
        if (!response.error && response.data) {
          const user = response.data;
          localStorage.setItem(STORAGE_KEYS.user, JSON.stringify(user));
          set({ currentUser: user, isLoading: false });
        } else {
          set({ isLoading: false });
        }
      } catch (err) {
        console.error('Failed to fetch profile during hydration:', err);
        set({ isLoading: false });
      }
    }
  },
}));
