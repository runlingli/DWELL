import React, { useState, useMemo, useEffect } from 'react';
import { INITIAL_LISTINGS } from './constants';
import { Listing, User } from './types';
import { Button, Select } from './components/UI';
import { ListingCard } from './components/ListingCard';
import { AuthModal } from './components/AuthModal';
import { CreateListingModal } from './components/CreateListingModal';
import { ListingDetail } from './components/ListingDetail';
import { MapView } from './components/MapView';

type SortOption = 'newest' | 'price-low' | 'price-high';
type ViewType = 'all' | 'profile';
type ProfileTab = 'favorites' | 'posts';

const SORT_OPTIONS = [
  { value: 'newest', label: 'Newest' },
  { value: 'price-low', label: 'Price: Low-High' },
  { value: 'price-high', label: 'Price: High-Low' },
];

const App: React.FC = () => {
  const [listings, setListings] = useState<Listing[]>(INITIAL_LISTINGS);
  const [favorites, setFavorites] = useState<string[]>([]);
  const [currentUser, setCurrentUser] = useState<User | null>(null);
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [listingToEdit, setListingToEdit] = useState<Listing | null>(null);
  const [view, setView] = useState<ViewType>('all');
  const [profileTab, setProfileTab] = useState<ProfileTab>('favorites');
  const [sortBy, setSortBy] = useState<SortOption>('newest');
  const [selectedListing, setSelectedListing] = useState<Listing | null>(null);
  const [mapCenter, setMapCenter] = useState<{ lat: number; lng: number } | undefined>();
  
  // Filter States
  const [filterStartDate, setFilterStartDate] = useState<string>('');
  const [filterEndDate, setFilterEndDate] = useState<string>('');

  useEffect(() => {
    const savedFavs = localStorage.getItem('dwell_favorites');
    if (savedFavs) setFavorites(JSON.parse(savedFavs));
    
    const savedUser = localStorage.getItem('dwell_user');
    if (savedUser) setCurrentUser(JSON.parse(savedUser));
  }, []);

  useEffect(() => {
    localStorage.setItem('dwell_favorites', JSON.stringify(favorites));
  }, [favorites]);

  useEffect(() => {
    if (currentUser) {
      localStorage.setItem('dwell_user', JSON.stringify(currentUser));
    } else {
      localStorage.removeItem('dwell_user');
    }
  }, [currentUser]);

  const toggleFavorite = (id: string) => {
    setFavorites(prev => 
      prev.includes(id) ? prev.filter(fid => fid !== id) : [...prev, id]
    );
  };

  const filteredAndSortedListings = useMemo(() => {
    let result = [...listings];

    if (view === 'profile') {
      if (profileTab === 'favorites') {
        result = result.filter(l => favorites.includes(l.id));
      } else {
        result = result.filter(l => l.author.name === currentUser?.name);
      }
    } else {
      // Apply date filtering
      if (filterStartDate && filterEndDate) {
        const startTs = new Date(filterStartDate).getTime();
        const endTs = new Date(filterEndDate).getTime();
        result = result.filter(l => l.availableFrom <= startTs && l.availableTo >= endTs);
      }
    }

    result.sort((a, b) => {
      if (sortBy === 'price-low') return a.price - b.price;
      if (sortBy === 'price-high') return b.price - a.price;
      return b.createdAt - a.createdAt;
    });

    return result;
  }, [listings, favorites, view, profileTab, sortBy, currentUser, filterStartDate, filterEndDate]);

  const handleSaveListing = (listingData: Listing) => {
    if (listingToEdit) {
      setListings(prev => prev.map(l => l.id === listingData.id ? listingData : l));
      if (selectedListing?.id === listingData.id) {
        setSelectedListing(listingData);
      }
    } else {
      const listingWithAuthor = {
        ...listingData,
        author: {
          name: currentUser?.name || 'Anonymous',
          avatar: currentUser?.avatar
        }
      };
      setListings(prev => [listingWithAuthor, ...prev]);
    }
    setListingToEdit(null);
    setShowCreateModal(false);
  };

  const handleDeleteListing = (id: string) => {
    setListings(prev => prev.filter(l => l.id !== id));
    setSelectedListing(null);
  };

  const handleListingClick = (listing: Listing) => {
    setSelectedListing(listing);
    setMapCenter(listing.coordinates);
  };

  const handleNavClick = (targetView: ViewType) => {
    if (targetView === 'profile' && !currentUser) {
      setShowAuthModal(true);
      return;
    }
    setView(targetView);
    setSelectedListing(null);
  };

  const startCreate = () => {
    setListingToEdit(null);
    setShowCreateModal(true);
  };

  const startEdit = (listing: Listing) => {
    setListingToEdit(listing);
    setShowCreateModal(true);
  };

  const handleFilterStartDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const val = e.target.value;
    setFilterStartDate(val);
    if (filterEndDate && val > filterEndDate) {
      setFilterEndDate('');
    }
  };

  return (
    <div className={`min-h-screen flex flex-col bg-[#f3e9d2] text-[#4a586e] ${selectedListing ? 'overflow-hidden' : ''}`}>
      {/* Header */}
      <header className="fixed top-0 left-0 right-0 bg-[#f3e9d2]/90 backdrop-blur-md z-[600] border-b border-[#4a586e]/10">
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <button 
            onClick={() => { setView('all'); setSortBy('newest'); setSelectedListing(null); }}
            className="font-serif text-3xl tracking-tighter text-[#4a586e] transition-opacity hover:opacity-60"
          >
            DWELL.
          </button>
          
          <nav className="hidden md:flex items-center gap-12 text-[10px] font-bold uppercase tracking-[0.3em]">
            <button 
              onClick={() => handleNavClick('all')}
              className={`${view === 'all' ? 'text-[#4a586e]' : 'text-[#7e918b]'} hover:text-[#4a586e] transition-colors underline-offset-[12px] ${view === 'all' ? 'underline decoration-1' : ''}`}
            >
              Discover
            </button>
            <button 
              onClick={() => handleNavClick('profile')}
              className={`${view === 'profile' ? 'text-[#4a586e]' : 'text-[#7e918b]'} hover:text-[#4a586e] transition-colors underline-offset-[12px] ${view === 'profile' ? 'underline decoration-1' : ''}`}
            >
              Profile
            </button>
          </nav>

          <div className="flex items-center gap-4">
            {currentUser ? (
              <div className="flex items-center gap-4">
                <Button variant="outline" className="hidden sm:block !py-2 !px-4" onClick={startCreate}>Post</Button>
                <button 
                  onClick={() => { setCurrentUser(null); setView('all'); }}
                  className="w-10 h-10 border border-[#4a586e] flex items-center justify-center hover:bg-[#4a586e] hover:text-[#f3e9d2] transition-all group text-[#4a586e]"
                >
                  <span className="text-[10px] font-bold group-hover:hidden uppercase tracking-tighter">ME</span>
                  <svg className="hidden group-hover:block w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                  </svg>
                </button>
              </div>
            ) : (
              <Button onClick={() => setShowAuthModal(true)} variant="primary" className="!py-2 !px-4">Sign In</Button>
            )}
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-grow pt-32 pb-20 px-6 max-w-7xl mx-auto w-full">
        {view === 'all' ? (
          <div className="mb-16">
            <div className="mb-5">
              <h1 className="font-serif text-6xl md:text-8xl mb-6 tracking-tighter text-[#4a586e] leading-none">Davis Residences.</h1>
            </div>

            {/* Split View Map & Filters */}
            <div className="grid grid-cols-1 lg:grid-cols-12 lg:max-h-[400px] gap-0 border border-[#4a586e]/10 bg-white/10 backdrop-blur-sm overflow-hidden mb-16">
              {/* Left Column: Map */}
              <div className="lg:col-span-8 border-b lg:border-b-0 lg:border-r border-[#4a586e]/10">
                <MapView 
                  listings={filteredAndSortedListings} 
                  onMarkerClick={handleListingClick} 
                  center={mapCenter}
                />
              </div>

              {/* Right Column: Filters Sidebar */}
              <div className="lg:col-span-4 p-8 md:p-10 flex flex-col justify-between">
                <div>
                  <div className="mb-10">
                    <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block mb-6">Curation Controls</span>
                    <Select 
                      label="Sequence"
                      options={SORT_OPTIONS}
                      value={sortBy}
                      onChange={(val) => setSortBy(val as SortOption)}
                      className="w-full"
                    />
                  </div>

                  <div className="space-y-10">
                    <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block">Timeline Filtering</span>
                    
                    <div className="space-y-8">
                      <div className="relative">
                        <span className="absolute -top-5 left-0 text-[8px] uppercase tracking-widest text-[#7e918b] font-bold">Inhabitation Start</span>
                        <input 
                          type="date" 
                          value={filterStartDate} 
                          onChange={handleFilterStartDateChange} 
                          className="w-full bg-transparent border-b border-[#4a586e]/20 py-3 outline-none text-[11px] uppercase text-[#4a586e] font-bold tracking-widest focus:border-[#4a586e] transition-colors" 
                        />
                      </div>
                      
                      <div className="relative">
                        <span className="absolute -top-5 left-0 text-[8px] uppercase tracking-widest text-[#7e918b] font-bold">Conclusion Date</span>
                        <input 
                          type="date" 
                          value={filterEndDate} 
                          onChange={e => setFilterEndDate(e.target.value)} 
                          min={filterStartDate} 
                          className="w-full bg-transparent border-b border-[#4a586e]/20 py-3 outline-none text-[11px] uppercase text-[#4a586e] font-bold tracking-widest focus:border-[#4a586e] transition-colors" 
                        />
                      </div>
                    </div>
                  </div>
                </div>

                <div className="mt-12 pt-8 border-t border-[#4a586e]/10 flex items-center justify-between">
                  <span className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/40">
                    {filteredAndSortedListings.length} Curated results
                  </span>
                  {(filterStartDate || filterEndDate) && (
                    <button 
                      onClick={() => { setFilterStartDate(''); setFilterEndDate(''); }} 
                      className="text-[9px] uppercase font-bold text-[#4a586e] hover:opacity-50 underline decoration-1 underline-offset-4"
                    >
                      Reset All
                    </button>
                  )}
                </div>
              </div>
            </div>
          </div>
        ) : (
          <div className="mb-16">
            <div className="flex items-center gap-8 mb-12">
               <div className="w-24 h-24 bg-[#4a586e] flex items-center justify-center text-[#f3e9d2] text-4xl font-serif">
                 {currentUser?.name?.charAt(0) || 'U'}
               </div>
               <div>
                 <h1 className="font-serif text-5xl text-[#4a586e] tracking-tighter mb-1">{currentUser?.name}</h1>
                 <p className="text-[#7e918b] text-xs tracking-[0.3em] font-bold">{currentUser?.email}</p>
               </div>
            </div>
            <div className="flex gap-12 border-b border-[#4a586e]/10 mb-12">
               <button onClick={() => setProfileTab('favorites')} className={`pb-4 text-[10px] font-bold uppercase tracking-widest transition-colors ${profileTab === 'favorites' ? 'text-[#4a586e] border-b border-[#4a586e]' : 'text-[#7e918b] hover:text-[#4a586e]'}`}>Favorites ({favorites.length})</button>
               <button onClick={() => setProfileTab('posts')} className={`pb-4 text-[10px] font-bold uppercase tracking-widest transition-colors ${profileTab === 'posts' ? 'text-[#4a586e] border-b border-[#4a586e]' : 'text-[#7e918b] hover:text-[#4a586e]'}`}>My Posts ({listings.filter(l => l.author.name === currentUser?.name).length})</button>
            </div>
          </div>
        )}

        {filteredAndSortedListings.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-12">
            {filteredAndSortedListings.map(listing => (
              <ListingCard 
                key={listing.id} 
                listing={listing} 
                isFavorite={favorites.includes(listing.id)}
                onToggleFavorite={toggleFavorite}
                onClick={handleListingClick}
              />
            ))}
          </div>
        ) : (
          <div className="h-64 border border-dashed border-[#4a586e]/20 flex flex-col items-center justify-center gap-4 bg-white/10">
            <p className="text-[#7e918b] uppercase tracking-widest text-[10px] font-bold">No residences match your search.</p>
          </div>
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-[#4a586e]/10 px-6 py-16 max-w-7xl mx-auto w-full flex flex-col md:flex-row justify-between items-center gap-8">
        <div className="font-serif text-3xl tracking-tighter text-[#4a586e]">DWELL</div>
      </footer>

      {/* Overlays */}
      {selectedListing && (
        <ListingDetail 
          listing={selectedListing} isOpen={!!selectedListing} onClose={() => setSelectedListing(null)} 
          isFavorite={favorites.includes(selectedListing.id)} onToggleFavorite={toggleFavorite} 
          currentUser={currentUser} onEdit={startEdit} onDelete={handleDeleteListing}
        />
      )}
      <AuthModal isOpen={showAuthModal} onClose={() => setShowAuthModal(false)} onLogin={(user) => setCurrentUser(user)} />
      <CreateListingModal isOpen={showCreateModal} onClose={() => { setShowCreateModal(false); setListingToEdit(null); }} onSave={handleSaveListing} initialData={listingToEdit} />
    </div>
  );
};

export default App;