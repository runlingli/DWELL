// stores/useUserStore.ts
import { create } from 'zustand';

export interface User {
  first_name: string;
  last_name: string;
  email: string;
}

interface UserState {
  user: User | null;

  // 设置完整 user
  setUser: (user: User) => void;

  // 更新部分字段
  updateUser: (partialUser: Partial<User>) => void;

  // 清空 user
  clearUser: () => void;
}

export const useUserStore = create<UserState>((set) => ({
  user: null,

  setUser: (user) =>
    set(() => ({
      user,
    })),

  updateUser: (partialUser) =>
    set((state) => ({
      user: state.user
        ? { ...state.user, ...partialUser }
        : null,
    })),

  clearUser: () =>
    set(() => ({
      user: null,
    })),
}));
