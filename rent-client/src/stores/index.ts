// src/stores/index.ts
// Zustand stores barrel export

export { useAuthStore } from './authStore';

export { useUIStore } from './UIStore';
export type { ViewType, ProfileTab, SortOption, UIState } from './UIStore';

export { useListingsStore } from './listingStore';

export { useFavoritesStore } from './favouriteStore';
