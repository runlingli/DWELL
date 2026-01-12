// src/stores/favouriteStore.ts
import { create } from 'zustand';
import { STORAGE_KEYS } from '@/config';
import * as favoritesApi from '@/api/favorites';

type FavoritesState = {
  favorites: string[];
  isLoading: boolean;
  toggleFavorite: (id: string, userId?: number) => void;
  isFavorite: (id: string) => boolean;
  hydrate: (userId?: number) => Promise<void>;
  syncToBackend: (userId: number) => Promise<void>;
  clearFavorites: () => void;
};

export const useFavoritesStore = create<FavoritesState>((set, get) => ({
  favorites: [],
  isLoading: false,

  toggleFavorite: async (id, userId) => {
    const currentFavorites = get().favorites;
    const isFav = currentFavorites.includes(id);

    // Optimistic update
    const next = isFav
      ? currentFavorites.filter((fid) => fid !== id)
      : [...currentFavorites, id];

    set({ favorites: next });
    localStorage.setItem(STORAGE_KEYS.favorites, JSON.stringify(next));

    // If user is logged in, sync with backend
    if (userId) {
      const postId = parseInt(id, 10);
      if (!isNaN(postId)) {
        try {
          if (isFav) {
            await favoritesApi.removeFavorite(userId, postId);
          } else {
            await favoritesApi.addFavorite(userId, postId);
          }
        } catch (error) {
          // Revert on error
          console.error('Failed to sync favorite:', error);
          set({ favorites: currentFavorites });
          localStorage.setItem(STORAGE_KEYS.favorites, JSON.stringify(currentFavorites));
        }
      }
    }
  },

  isFavorite: (id) => get().favorites.includes(id),

  hydrate: async (userId) => {
    // First, load from localStorage
    const raw = localStorage.getItem(STORAGE_KEYS.favorites);
    const localFavorites: string[] = raw ? JSON.parse(raw) : [];

    if (!userId) {
      // Not logged in, just use localStorage
      set({ favorites: localFavorites });
      return;
    }

    set({ isLoading: true });

    try {
      // Fetch from backend
      const response = await favoritesApi.fetchFavoriteIds(userId);

      if (!response.error && response.data) {
        const backendFavorites = response.data;

        // Merge local and backend favorites (union)
        const merged = [...new Set([...localFavorites, ...backendFavorites])];

        // If there are local favorites not in backend, sync them
        const localOnly = localFavorites.filter(id => !backendFavorites.includes(id));
        if (localOnly.length > 0) {
          const postIds = localOnly.map(id => parseInt(id, 10)).filter(id => !isNaN(id));
          if (postIds.length > 0) {
            await favoritesApi.syncFavorites(userId, postIds);
          }
        }

        set({ favorites: merged });
        localStorage.setItem(STORAGE_KEYS.favorites, JSON.stringify(merged));
      } else {
        // Backend error, use localStorage
        set({ favorites: localFavorites });
      }
    } catch (error) {
      console.error('Failed to hydrate favorites:', error);
      set({ favorites: localFavorites });
    } finally {
      set({ isLoading: false });
    }
  },

  syncToBackend: async (userId) => {
    const localFavorites = get().favorites;
    if (localFavorites.length === 0) return;

    const postIds = localFavorites.map(id => parseInt(id, 10)).filter(id => !isNaN(id));
    if (postIds.length === 0) return;

    try {
      const response = await favoritesApi.syncFavorites(userId, postIds);
      if (!response.error && response.data) {
        set({ favorites: response.data });
        localStorage.setItem(STORAGE_KEYS.favorites, JSON.stringify(response.data));
      }
    } catch (error) {
      console.error('Failed to sync favorites to backend:', error);
    }
  },

  clearFavorites: () => {
    set({ favorites: [] });
    localStorage.removeItem(STORAGE_KEYS.favorites);
  },
}));
