// src/pages/profile/ProfilePage.tsx
import React from 'react';
import { useAuthStore, useUIStore, useFavoritesStore, useListingsStore } from '@/stores';
import { ListingCard } from '@/components';
import { getUserDisplayName } from '@/utils';

export const ProfilePage: React.FC = () => {
  const currentUser = useAuthStore((state) => state.currentUser);
  const listings = useListingsStore((state) => state.listings);
  const { favorites, toggleFavorite } = useFavoritesStore();
  const { profileTab, setProfileTab, openCreateModal, selectListing } = useUIStore();

  if (!currentUser) {
    return (
      <div className="py-32 text-center text-[#4a586e]">
        <p className="text-lg">You must be signed in to view your profile.</p>
      </div>
    );
  }

  const displayName = getUserDisplayName(currentUser);
  const myListings = listings.filter((l) => l.author.name === displayName);
  const favoriteListings = listings.filter((l) => favorites.includes(l.id));
  const listingsToShow = profileTab === 'favorites' ? favoriteListings : myListings;

  return (
    <main className="grow pt-32 pb-20 px-6 max-w-7xl mx-auto w-full">
      {/* User Info */}
      <div className="mb-16">
        <div className="flex items-center gap-8 mb-12">
          <div className="w-24 h-24 bg-[#4a586e] flex items-center justify-center text-[#f3e9d2] text-4xl font-serif">
            {currentUser.first_name.charAt(0).toUpperCase()}
          </div>
          <div>
            <h1 className="font-serif text-5xl text-[#4a586e] tracking-tighter mb-1">{displayName}</h1>
            <p className="text-[#7e918b] text-xs tracking-[0.3em] font-bold">{currentUser.email}</p>
          </div>

          <button
            onClick={openCreateModal}
            className="ml-auto border border-[#4a586e] px-4 py-2 text-[#4a586e] font-bold hover:bg-[#4a586e] hover:text-[#f3e9d2] transition-colors text-[10px] uppercase tracking-widest"
          >
            New Post
          </button>
        </div>

        <div className="flex gap-12 border-b border-[#4a586e]/10 mb-12">
          <button
            onClick={() => setProfileTab('favorites')}
            className={`pb-4 text-[10px] font-bold uppercase tracking-widest transition-colors ${
              profileTab === 'favorites'
                ? 'text-[#4a586e] border-b border-[#4a586e]'
                : 'text-[#7e918b] hover:text-[#4a586e]'
            }`}
          >
            Favorites ({favoriteListings.length})
          </button>
          <button
            onClick={() => setProfileTab('posts')}
            className={`pb-4 text-[10px] font-bold uppercase tracking-widest transition-colors ${
              profileTab === 'posts'
                ? 'text-[#4a586e] border-b border-[#4a586e]'
                : 'text-[#7e918b] hover:text-[#4a586e]'
            }`}
          >
            My Posts ({myListings.length})
          </button>
        </div>
      </div>

      {/* Listing Grid */}
      {listingsToShow.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-12">
          {listingsToShow.map((listing) => (
            <ListingCard
              key={listing.id}
              listing={listing}
              isFavorite={favorites.includes(listing.id)}
              onToggleFavorite={toggleFavorite}
              onClick={selectListing}
            />
          ))}
        </div>
      ) : (
        <div className="h-64 border border-dashed border-[#4a586e]/20 flex flex-col items-center justify-center gap-4 bg-white/10">
          <p className="text-[#7e918b] uppercase tracking-widest text-[10px] font-bold">
            {profileTab === 'favorites' ? 'No saved residences yet.' : 'No posts yet.'}
          </p>
        </div>
      )}
    </main>
  );
};
