// src/stores/UIStore.ts
import { create } from 'zustand';
import type { Listing } from '@/types';

export type ViewType = 'discover' | 'profile';
export type ProfileTab = 'favorites' | 'posts';
export type SortOption = 'newest' | 'price-low' | 'price-high';

export interface UIState {
  // State
  view: ViewType;
  profileTab: ProfileTab;
  sortBy: SortOption;
  selectedListing: Listing | null;
  showAuthModal: boolean;
  showCreateModal: boolean;
  listingToEdit: Listing | null;
  mapCenter: { lat: number; lng: number } | undefined;
  filterStartDate: string;
  filterEndDate: string;

  // Navigation
  navigate: (view: ViewType) => void;
  setProfileTab: (tab: ProfileTab) => void;
  setSortBy: (sort: SortOption) => void;
  resetToHome: () => void;

  // Auth modal
  openAuthModal: () => void;
  closeAuthModal: () => void;

  // Create/edit modal
  openCreateModal: () => void;
  openEditModal: (listing: Listing) => void;
  closeCreateModal: () => void;

  // Listing selection
  selectListing: (listing: Listing) => void;
  clearSelectedListing: () => void;

  // Filters
  setFilterStartDate: (date: string) => void;
  setFilterEndDate: (date: string) => void;
  clearFilters: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  view: 'discover',
  profileTab: 'favorites',
  sortBy: 'newest',
  selectedListing: null,
  showAuthModal: false,
  showCreateModal: false,
  listingToEdit: null,
  mapCenter: undefined,
  filterStartDate: '',
  filterEndDate: '',

  // Navigation
  navigate: (view: ViewType) =>
    set({
      view,
      selectedListing: null,
    }),

  setProfileTab: (tab: ProfileTab) => set({ profileTab: tab }),

  setSortBy: (sort: SortOption) => set({ sortBy: sort }),

  resetToHome: () =>
    set({
      view: 'discover',
      sortBy: 'newest',
      selectedListing: null,
      filterStartDate: '',
      filterEndDate: '',
    }),

  // Auth modal
  openAuthModal: () => set({ showAuthModal: true }),
  closeAuthModal: () => set({ showAuthModal: false }),

  // Create/edit modal
  openCreateModal: () =>
    set({
      showCreateModal: true,
      listingToEdit: null,
    }),

  openEditModal: (listing: Listing) =>
    set({
      showCreateModal: true,
      listingToEdit: listing,
    }),

  closeCreateModal: () =>
    set({
      showCreateModal: false,
      listingToEdit: null,
    }),

  // Listing selection
  selectListing: (listing: Listing) =>
    set({
      selectedListing: listing,
      mapCenter: listing.coordinates,
    }),

  clearSelectedListing: () =>
    set({
      selectedListing: null,
    }),

  // Filters
  setFilterStartDate: (date: string) =>
    set((state) => ({
      filterStartDate: date,
      filterEndDate: state.filterEndDate && date > state.filterEndDate ? '' : state.filterEndDate,
    })),

  setFilterEndDate: (date: string) => set({ filterEndDate: date }),

  clearFilters: () =>
    set({
      filterStartDate: '',
      filterEndDate: '',
    }),
}));
