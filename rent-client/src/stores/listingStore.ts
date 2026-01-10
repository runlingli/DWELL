// src/stores/listingStore.ts
import { create } from 'zustand';
import type { Listing } from '@/types';
import { INITIAL_LISTINGS } from '@/config';
import * as postsApi from '@/api/posts';

type ListingsState = {
  listings: Listing[];
  isLoading: boolean;
  error: string | null;
  hasFetched: boolean;
  fetchListings: (force?: boolean) => Promise<void>;
  addListing: (listing: Listing, authorId?: number) => Promise<void>;
  updateListing: (listing: Listing, authorId?: number) => Promise<void>;
  deleteListing: (id: string, authorId?: number) => Promise<void>;
  setListings: (listings: Listing[]) => void;
};

export const useListingsStore = create<ListingsState>((set, get) => ({
  listings: INITIAL_LISTINGS,
  isLoading: false,
  error: null,
  hasFetched: false,

  fetchListings: async (force = false) => {
    const state = get();
    if (state.isLoading) return;
    if (state.hasFetched && !force) return;

    set({ isLoading: true, error: null });
    try {
      const response = await postsApi.fetchPosts();
      if (!response.error && response.data) {
        set({ listings: response.data, isLoading: false, hasFetched: true });
      } else {
        set({ isLoading: false, hasFetched: true, error: response.message || 'Failed to fetch listings' });
      }
    } catch (err) {
      console.error('Failed to fetch listings:', err);
      set({
        isLoading: false,
        hasFetched: true,
        error: err instanceof Error ? err.message : 'Failed to fetch listings',
      });
    }
  },

  addListing: async (listing, authorId) => {
    console.log('========== addListing START ==========');
    console.log('Listing to add:', listing);
    console.log('AuthorId:', authorId);

    set({ isLoading: true, error: null });
    try {
      if (authorId) {
        console.log('Creating post on backend with authorId:', authorId);
        const response = await postsApi.createPost(listing, authorId);
        console.log('Backend response:', response);
        if (!response.error && response.data) {
          console.log('SUCCESS: Post created on backend:', response.data);
          set((state) => ({
            listings: [response.data!, ...state.listings],
            isLoading: false,
          }));
          return;
        } else {
          console.log('Backend returned error or no data:', response.message);
        }
      } else {
        console.log('WARNING: No authorId provided - post will only be added locally');
      }
      // Fallback: add locally
      console.log('Fallback: Adding locally only');
      set((state) => ({
        listings: [listing, ...state.listings],
        isLoading: false,
      }));
    } catch (err) {
      console.error('Failed to create listing:', err);
      set((state) => ({
        listings: [listing, ...state.listings],
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to create listing',
      }));
    }
  },

  updateListing: async (listing, authorId) => {
    set({ isLoading: true, error: null });
    try {
      if (authorId) {
        const response = await postsApi.updatePost(listing, authorId);
        if (!response.error && response.data) {
          set((state) => ({
            listings: state.listings.map((l) => (l.id === listing.id ? response.data! : l)),
            isLoading: false,
          }));
          return;
        }
      }
      // Fallback: update locally
      set((state) => ({
        listings: state.listings.map((l) => (l.id === listing.id ? listing : l)),
        isLoading: false,
      }));
    } catch (err) {
      console.error('Failed to update listing:', err);
      set((state) => ({
        listings: state.listings.map((l) => (l.id === listing.id ? listing : l)),
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to update listing',
      }));
    }
  },

  deleteListing: async (id, authorId) => {
    set({ isLoading: true, error: null });
    try {
      if (authorId) {
        const response = await postsApi.deletePost(id, authorId);
        if (!response.error) {
          set((state) => ({
            listings: state.listings.filter((l) => l.id !== id),
            isLoading: false,
          }));
          return;
        }
      }
      // Fallback: delete locally
      set((state) => ({
        listings: state.listings.filter((l) => l.id !== id),
        isLoading: false,
      }));
    } catch (err) {
      console.error('Failed to delete listing:', err);
      set((state) => ({
        listings: state.listings.filter((l) => l.id !== id),
        isLoading: false,
        error: err instanceof Error ? err.message : 'Failed to delete listing',
      }));
    }
  },

  setListings: (listings) => set({ listings }),
}));
