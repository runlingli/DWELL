// src/stores/favouriteStore.ts
import { create } from 'zustand';
import { STORAGE_KEYS } from '@/config';

type FavoritesState = {
  favorites: string[];
  toggleFavorite: (id: string) => void;
  isFavorite: (id: string) => boolean;
  hydrate: () => void;
};

export const useFavoritesStore = create<FavoritesState>((set, get) => ({
  favorites: [],

  toggleFavorite: (id) =>
    set((state) => {
      const next = state.favorites.includes(id)
        ? state.favorites.filter((fid) => fid !== id)
        : [...state.favorites, id];

      localStorage.setItem(STORAGE_KEYS.favorites, JSON.stringify(next));
      return { favorites: next };
    }),

  isFavorite: (id) => get().favorites.includes(id),

  hydrate: () => {
    const raw = localStorage.getItem(STORAGE_KEYS.favorites);
    if (raw) {
      set({ favorites: JSON.parse(raw) });
    }
  },
}));
