// src/pages/discover/DiscoverPage.tsx
import React, { useMemo, useEffect } from 'react';
import { useUIStore, useListingsStore, useFavoritesStore, type SortOption } from '@/stores';
import { ListingCard, MapView } from '@/components';
import { Select } from '@/ui';

const SORT_OPTIONS: { value: SortOption; label: string }[] = [
  { value: 'newest', label: 'Newest' },
  { value: 'price-low', label: 'Price: Low-High' },
  { value: 'price-high', label: 'Price: High-Low' },
];

export const DiscoverPage: React.FC = () => {
  const { listings, isLoading, fetchListings } = useListingsStore();
  const { favorites, toggleFavorite } = useFavoritesStore();
  const {
    sortBy,
    setSortBy,
    filterStartDate,
    filterEndDate,
    setFilterStartDate,
    setFilterEndDate,
    clearFilters,
    selectListing,
    mapCenter,
  } = useUIStore();

  useEffect(() => {
    fetchListings();
  }, [fetchListings]);

  const filteredAndSortedListings = useMemo(() => {
    let result = [...listings];

    if (filterStartDate && filterEndDate) {
      const startTs = new Date(filterStartDate).getTime();
      const endTs = new Date(filterEndDate).getTime();
      result = result.filter((l) => l.availableFrom <= startTs && l.availableTo >= endTs);
    }

    result.sort((a, b) => {
      if (sortBy === 'price-low') return a.price - b.price;
      if (sortBy === 'price-high') return b.price - a.price;
      return b.createdAt - a.createdAt;
    });

    return result;
  }, [listings, sortBy, filterStartDate, filterEndDate]);

  return (
    <main className="flex-grow pt-32 pb-20 px-6 max-w-7xl mx-auto w-full">
      <div className="mb-16">
        <div className="mb-5">
          <h1 className="font-serif text-6xl md:text-8xl mb-6 tracking-tighter text-[#4a586e] leading-none">
            Short Period Residences
          </h1>
        </div>

        {/* Split View Map & Filters */}
        <div className="grid grid-cols-1 lg:grid-cols-12 lg:max-h-[400px] gap-0 border border-[#4a586e]/10 bg-white/10 backdrop-blur-sm overflow-hidden mb-16">
          <div className="lg:col-span-8 border-b lg:border-b-0 lg:border-r border-[#4a586e]/10">
            <MapView
              listings={filteredAndSortedListings}
              onMarkerClick={selectListing}
              center={mapCenter}
            />
          </div>

          <div className="lg:col-span-4 p-8 md:p-10 flex flex-col justify-between">
            <div>
              <div className="mb-10">
                <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block mb-6">
                  Curation Controls
                </span>
                <Select
                  label="Sequence"
                  options={SORT_OPTIONS}
                  value={sortBy}
                  onChange={(val) => setSortBy(val as SortOption)}
                  className="w-full"
                />
              </div>

              <div className="space-y-10">
                <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block">
                  Timeline Filtering
                </span>

                <div className="space-y-8">
                  <div className="relative">
                    <span className="absolute -top-5 left-0 text-[8px] uppercase tracking-widest text-[#7e918b] font-bold">
                      Inhabitation Start
                    </span>
                    <input
                      type="date"
                      value={filterStartDate}
                      onChange={(e) => setFilterStartDate(e.target.value)}
                      className="w-full bg-transparent border-b border-[#4a586e]/20 py-3 outline-none text-[11px] uppercase text-[#4a586e] font-bold tracking-widest focus:border-[#4a586e] transition-colors"
                    />
                  </div>

                  <div className="relative">
                    <span className="absolute -top-5 left-0 text-[8px] uppercase tracking-widest text-[#7e918b] font-bold">
                      Conclusion Date
                    </span>
                    <input
                      type="date"
                      value={filterEndDate}
                      onChange={(e) => setFilterEndDate(e.target.value)}
                      min={filterStartDate}
                      className="w-full bg-transparent border-b border-[#4a586e]/20 py-3 outline-none text-[11px] uppercase text-[#4a586e] font-bold tracking-widest focus:border-[#4a586e] transition-colors"
                    />
                  </div>
                </div>
              </div>
            </div>

            <div className="mt-12 pt-8 border-t border-[#4a586e]/10 flex items-center justify-between">
              <span className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/40">
                {isLoading ? 'Loading...' : `${filteredAndSortedListings.length} Curated results`}
              </span>
              {(filterStartDate || filterEndDate) && (
                <button
                  onClick={clearFilters}
                  className="text-[9px] uppercase font-bold text-[#4a586e] hover:opacity-50 underline decoration-1 underline-offset-4"
                >
                  Reset All
                </button>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Listing Grid */}
      {isLoading ? (
        <div className="h-64 border border-dashed border-[#4a586e]/20 flex flex-col items-center justify-center gap-4 bg-white/10">
          <div className="animate-pulse">
            <p className="text-[#7e918b] uppercase tracking-widest text-[10px] font-bold">
              Loading residences...
            </p>
          </div>
        </div>
      ) : filteredAndSortedListings.length > 0 ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-12">
          {filteredAndSortedListings.map((listing) => (
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
            No residences match your search.
          </p>
        </div>
      )}
    </main>
  );
};
