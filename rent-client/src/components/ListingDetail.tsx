import React, { useState } from 'react';
import { Listing, User } from '../types';
import { Button } from './UI';

interface ListingDetailProps {
  listing: Listing;
  isOpen: boolean;
  onClose: () => void;
  isFavorite: boolean;
  onToggleFavorite: (id: string) => void;
  currentUser: User | null;
  onEdit: (listing: Listing) => void;
  onDelete: (id: string) => void;
}

export const ListingDetail: React.FC<ListingDetailProps> = ({ 
  listing, 
  isOpen, 
  onClose,
  isFavorite,
  onToggleFavorite,
  currentUser,
  onEdit,
  onDelete
}) => {
  const [isConfirmingDelete, setIsConfirmingDelete] = useState(false);

  if (!isOpen) return null;

  const formatDate = (ts: number) => new Date(ts).toLocaleDateString('en-US', { 
    month: 'short', day: 'numeric', year: 'numeric' 
  });

  const isAuthor = currentUser?.name === listing.author.name;

  const handleDeleteClick = () => {
    if (isConfirmingDelete) {
      onDelete(listing.id);
    } else {
      setIsConfirmingDelete(true);
      // Auto-reset after 3 seconds if not clicked again
      setTimeout(() => setIsConfirmingDelete(false), 3000);
    }
  };

  return (
    <div className="fixed inset-0 z-[1000] flex items-center justify-center p-4 md:p-10 lg:p-16">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-[#4a586e]/60 backdrop-blur-sm animate-in fade-in duration-500" 
        onClick={onClose} 
      />
      
      {/* Detail Container */}
      <div className="relative w-full max-w-6xl h-full max-h-[90vh] bg-[#f3e9d2] border border-[#4a586e] shadow-2xl animate-in fade-in zoom-in-95 duration-500 flex flex-col overflow-hidden">
        
        {/* Top Bar */}
        <div className="flex-shrink-0 flex justify-between items-center px-6 md:px-10 py-5 bg-[#f3e9d2] border-b border-[#4a586e]/10 z-10">
          <button 
            onClick={onClose}
            className="text-[#4a586e] hover:opacity-50 transition-opacity flex items-center gap-2 text-[10px] font-bold uppercase tracking-[0.2em]"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
            </svg>
            Back
          </button>
          <div className="flex items-center gap-6 md:gap-12">
            {isAuthor && (
              <div className="flex items-center gap-4 md:gap-8">
                {!isConfirmingDelete && (
                  <button 
                    onClick={() => onEdit(listing)}
                    className="text-[#4a586e] hover:opacity-50 transition-opacity flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest border-b border-[#4a586e]"
                  >
                    Edit
                  </button>
                )}
                <button 
                  onClick={handleDeleteClick}
                  className={`${isConfirmingDelete ? 'bg-red-800 text-white px-3 py-1 animate-pulse' : 'text-red-800 border-b border-red-800/20'} hover:opacity-80 transition-all flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest`}
                >
                  {isConfirmingDelete ? 'Confirm Deletion?' : 'Delete'}
                </button>
                {isConfirmingDelete && (
                  <button 
                    onClick={() => setIsConfirmingDelete(false)}
                    className="text-[#7e918b] text-[10px] font-bold uppercase tracking-widest hover:text-[#4a586e]"
                  >
                    Cancel
                  </button>
                )}
              </div>
            )}
            <button 
              onClick={() => onToggleFavorite(listing.id)}
              className="text-[#4a586e] hover:opacity-50 transition-opacity flex items-center gap-2"
            >
              <span className="text-[10px] font-bold uppercase tracking-widest hidden sm:inline">{isFavorite ? 'Saved' : 'Save'}</span>
              <svg 
                className={`w-5 h-5 ${isFavorite ? 'fill-current' : 'fill-none'}`} 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
            </button>
          </div>
        </div>

        {/* Content Area */}
        <div className="flex-grow overflow-y-auto custom-scrollbar">
          <div className="grid grid-cols-1 lg:grid-cols-12 min-h-full">
            
            {/* Left Column: Visuals */}
            <div className="lg:col-span-6 xl:col-span-7 bg-[#4a586e]/5 flex items-center justify-center p-6 md:p-10 xl:p-16 border-b lg:border-b-0 lg:border-r border-[#4a586e]/10">
              <div className="relative w-full aspect-[4/5] shadow-2xl overflow-hidden group">
                <img 
                  src={listing.imageUrl} 
                  alt={listing.title} 
                  className="w-full h-full object-cover grayscale brightness-90 sepia-[.1] group-hover:grayscale-0 group-hover:sepia-0 transition-all duration-1000" 
                />
                <div className="absolute top-6 left-6 bg-[#f3e9d2] px-3 py-1 text-[8px] font-bold uppercase tracking-widest text-[#4a586e]">
                  {listing.neighborhood}
                </div>
              </div>
            </div>

            {/* Right Column: Information */}
            <div className="lg:col-span-6 xl:col-span-5 bg-[#f3e9d2] p-8 md:p-12 xl:p-16 flex flex-col">
              <div className="flex-grow">
                <div className="mb-12">
                  <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block mb-2">Residence ID: {listing.id.toUpperCase()}</span>
                  <h1 className="font-serif text-4xl md:text-5xl xl:text-6xl mb-4 tracking-tighter text-[#4a586e] leading-tight uppercase">{listing.title}</h1>
                  <p className="text-[#4a586e] text-[10px] uppercase tracking-[0.3em] font-bold flex items-center gap-2">
                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/></svg>
                    {listing.location || listing.neighborhood}
                  </p>
                </div>

                <div className="grid grid-cols-2 gap-y-8 mb-12 border-t border-[#4a586e]/10 pt-8">
                  <div>
                    <p className="text-[9px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Monthly</p>
                    <p className="text-3xl font-serif text-[#4a586e] tracking-tighter">${listing.price}</p>
                  </div>
                  <div>
                    <p className="text-[9px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Category</p>
                    <p className="text-lg font-serif italic text-[#4a586e]">{listing.type}</p>
                  </div>
                  <div>
                    <p className="text-[9px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Config</p>
                    <p className="text-lg font-serif italic text-[#4a586e]">{listing.bedrooms}BR / {listing.bathrooms}BA</p>
                  </div>
                  <div>
                    <p className="text-[9px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Timeline</p>
                    <p className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]">
                      {formatDate(listing.availableFrom)}
                    </p>
                  </div>
                </div>

                <div className="mb-12">
                  <p className="text-[9px] text-[#7e918b] uppercase tracking-widest mb-4 font-bold">Narrative</p>
                  <p className="text-[#4a586e] leading-relaxed text-xl font-light italic font-serif">
                    "{listing.description}"
                  </p>
                </div>

                <div className="flex items-center gap-4 border-t border-[#4a586e]/10 pt-8 mb-8">
                  <div className="w-12 h-12 bg-[#4a586e] flex items-center justify-center grayscale">
                     {listing.author.avatar ? <img src={listing.author.avatar} alt={listing.author.name} className="w-full h-full object-cover" /> : <div className="text-[#f3e9d2] font-serif text-xl">{listing.author.name.charAt(0)}</div>}
                  </div>
                  <div>
                    <p className="text-[8px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Curated By</p>
                    <p className="text-[10px] font-black text-[#4a586e] uppercase tracking-widest">{listing.author.name}</p>
                  </div>
                </div>
              </div>

              <div className="mt-auto pt-8">
                <Button variant="primary" className="w-full !py-6 text-[11px] tracking-[0.4em]">Inquire Now</Button>
                <p className="text-center text-[7px] uppercase tracking-[0.3em] text-[#7e918b] mt-4 font-bold">Professional Response Guaranteed</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};