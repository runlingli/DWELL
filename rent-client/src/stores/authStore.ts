import { create } from 'zustand';
import type { User } from '../types/types';

type AuthState = {
  currentUser: User | null;
  login: (user: User) => void;
  logout: () => void;
  hydrate: () => void;
};

const STORAGE_KEY = 'dwell_user';

export const useAuthStore = create<AuthState>((set) => ({
  currentUser: null,

  login: (user) => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(user));
    set({ currentUser: user });
  },

  logout: () => {
    localStorage.removeItem(STORAGE_KEY);
    set({ currentUser: null });
  },

  hydrate: () => {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      set({ currentUser: JSON.parse(raw) });
    }
  },
}));
