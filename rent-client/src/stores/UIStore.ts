// src/store/uiStore.ts
import { create } from 'zustand';
import type { Listing } from '../types/types';

export type ViewType = 'discover' | 'profile';
export type ProfileTab = 'favorites' | 'posts';
export type SortOption = 'newest' | 'price-low' | 'price-high';

export interface UIState {
  // ===== state =====
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

  // ===== navigation =====
  navigate: (view: ViewType) => void;
  setProfileTab: (tab: ProfileTab) => void;
  setSortBy: (sort: SortOption) => void;
  resetToHome: () => void;

  // ===== auth modal =====
  openAuthModal: () => void;
  closeAuthModal: () => void;

  // ===== create / edit modal =====
  openCreateModal: () => void;
  openEditModal: (listing: Listing) => void;
  closeCreateModal: () => void;

  // ===== listing selection =====
  selectListing: (listing: Listing) => void;
  clearSelectedListing: () => void;

  // ===== filters =====
  setFilterStartDate: (date: string) => void;
  setFilterEndDate: (date: string) => void;
  clearFilters: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  // ===== state =====
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

  // ===== navigation =====
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

  // ===== auth modal =====
  openAuthModal: () => set({ showAuthModal: true }),
  closeAuthModal: () => set({ showAuthModal: false }),

  // ===== create / edit modal =====
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

  // ===== listing selection =====
  selectListing: (listing: Listing) =>
    set({
      selectedListing: listing,
      mapCenter: listing.coordinates,
    }),

  clearSelectedListing: () =>
    set({
      selectedListing: null,
    }),

  // ===== filters =====
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
