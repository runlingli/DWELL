import { create } from 'zustand';

type FavoritesState = {
  favorites: string[];
  toggleFavorite: (id: string) => void;
  isFavorite: (id: string) => boolean;
  hydrate: () => void;
};

const STORAGE_KEY = 'dwell_favorites';

export const useFavoritesStore = create<FavoritesState>((set, get) => ({
  favorites: [],

  toggleFavorite: (id) =>
    set((state) => {
      const next = state.favorites.includes(id)
        ? state.favorites.filter((fid) => fid !== id)
        : [...state.favorites, id];

      localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
      return { favorites: next };
    }),

  isFavorite: (id) => get().favorites.includes(id),

  hydrate: () => {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (raw) {
      set({ favorites: JSON.parse(raw) });
    }
  },
}));
