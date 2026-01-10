// src/components/listings/detail/DetailContent.tsx
import React from 'react';
import type { Listing } from '@/types';
import { Button } from '@/ui';

interface DetailContentProps {
  listing: Listing;
}

const formatDate = (ts: number) =>
  new Date(ts).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  });

export const DetailContent: React.FC<DetailContentProps> = ({ listing }) => {
  return (
    <div className="grow overflow-y-auto custom-scrollbar">
      <div className="grid grid-cols-1 lg:grid-cols-12 min-h-full">
        {/* Left Column: Visuals */}
        <div className="lg:col-span-6 xl:col-span-7 bg-[#4a586e]/5 flex items-center justify-center p-6 md:p-10 xl:p-16 border-b lg:border-b-0 lg:border-r border-[#4a586e]/10">
          <div className="relative w-full aspect-4/5 shadow-2xl overflow-hidden group">
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
              <span className="text-[9px] font-bold uppercase tracking-[0.4em] text-[#7e918b] block mb-2">
                Residence ID: {listing.id.toUpperCase()}
              </span>
              <h1 className="font-serif text-4xl md:text-5xl xl:text-6xl mb-4 tracking-tighter text-[#4a586e] leading-tight uppercase">
                {listing.title}
              </h1>
              <p className="text-[#4a586e] text-[10px] uppercase tracking-[0.3em] font-bold flex items-center gap-2">
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                  />
                </svg>
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
                <p className="text-lg font-serif italic text-[#4a586e]">
                  {listing.bedrooms}BR / {listing.bathrooms}BA
                </p>
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
                {listing.author.avatar ? (
                  <img
                    src={listing.author.avatar}
                    alt={listing.author.name}
                    className="w-full h-full object-cover"
                  />
                ) : (
                  <div className="text-[#f3e9d2] font-serif text-xl">{listing.author.name.charAt(0)}</div>
                )}
              </div>
              <div>
                <p className="text-[8px] text-[#7e918b] uppercase tracking-widest mb-1 font-bold">Curated By</p>
                <p className="text-[10px] font-black text-[#4a586e] uppercase tracking-widest">
                  {listing.author.name}
                </p>
              </div>
            </div>
          </div>

          <div className="mt-auto pt-8">
            <Button variant="primary" className="w-full !py-6 text-[11px] tracking-[0.4em]">
              Inquire Now
            </Button>
            <p className="text-center text-[7px] uppercase tracking-[0.3em] text-[#7e918b] mt-4 font-bold">
              Professional Response Guaranteed
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
