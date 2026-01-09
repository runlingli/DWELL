import { create } from 'zustand';
import { Listing } from '../types/types';
import { INITIAL_LISTINGS } from '../constants';

type ListingsState = {
  listings: Listing[];
  addListing: (listing: Listing) => void;
  updateListing: (listing: Listing) => void;
  deleteListing: (id: string) => void;
};

export const useListingsStore = create<ListingsState>((set) => ({
  listings: INITIAL_LISTINGS,

  addListing: (listing) =>
    set((state) => ({
      listings: [listing, ...state.listings],
    })),

  updateListing: (listing) =>
    set((state) => ({
      listings: state.listings.map((l) =>
        l.id === listing.id ? listing : l
      ),
    })),

  deleteListing: (id) =>
    set((state) => ({
      listings: state.listings.filter((l) => l.id !== id),
    })),
}));
