import React from 'react';
import type { Listing } from '../types/types';

interface ListingCardProps {
  listing: Listing;
  isFavorite: boolean;
  onToggleFavorite: (id: string) => void;
  onClick: (listing: Listing) => void;
}

export const ListingCard: React.FC<ListingCardProps> = ({ 
  listing, 
  isFavorite, 
  onToggleFavorite,
  onClick
}) => {
  return (
    <div className="group relative border-b border-[#4a586e]/10 pb-8 mb-8 last:border-0 last:mb-0 lg:border-0 lg:pb-0 lg:mb-0 cursor-pointer">
      <div 
        className="relative overflow-hidden aspect-[4/5] bg-[#7e918b]/10"
        onClick={() => onClick(listing)}
      >
        <img 
          src={listing.imageUrl} 
          alt={listing.title}
          className="w-full h-full object-cover grayscale brightness-90 sepia-[.2] transition-all duration-700 group-hover:grayscale-0 group-hover:brightness-100 group-hover:sepia-0 group-hover:scale-105"
        />
        <div className="absolute inset-0 bg-black/0 group-hover:bg-black/5 transition-colors duration-500" />
        <div className="absolute bottom-4 left-4 z-10 bg-[#f3e9d2] px-2 py-1 text-[8px] font-bold uppercase tracking-widest text-[#4a586e]">
          {listing.neighborhood}
        </div>
      </div>

      <button 
        onClick={(e) => {
          e.stopPropagation();
          onToggleFavorite(listing.id);
        }}
        className="absolute top-4 right-4 p-2 bg-[#f3e9d2]/80 backdrop-blur-sm transition-colors hover:bg-[#4a586e] hover:text-[#f3e9d2] text-[#4a586e] z-10"
      >
        <svg 
          className={`w-5 h-5 ${isFavorite ? 'fill-current' : 'fill-none'}`} 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path 
            strokeLinecap="round" 
            strokeLinejoin="round" 
            strokeWidth={1} 
            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" 
          />
        </svg>
      </button>
      
      <div className="mt-6 flex justify-between items-start" onClick={() => onClick(listing)}>
        <div>
          <h3 className="font-serif text-xl mb-1 text-[#4a586e] group-hover:underline underline-offset-4 decoration-1 tracking-tight">{listing.title}</h3>
          <p className="text-[#7e918b] text-[10px] uppercase tracking-widest font-bold">{listing.location || listing.neighborhood}</p>
        </div>
        <div className="text-right">
          <p className="font-bold text-lg text-[#4a586e]">${listing.price}</p>
          <p className="text-[#7e918b] text-[9px] uppercase tracking-widest font-bold">per month</p>
        </div>
      </div>
      
      <div className="mt-4 flex gap-4 text-[9px] uppercase tracking-widest text-[#4a586e]/60 border-t border-[#4a586e]/10 pt-4" onClick={() => onClick(listing)}>
        <span className="font-bold">{listing.bedrooms} Bedrooms</span>
        <span className="font-bold">{listing.bathrooms} Bathrooms</span>
        <span className="text-[#9bb794] font-black">{listing.type}</span>
      </div>
    </div>
  );
};