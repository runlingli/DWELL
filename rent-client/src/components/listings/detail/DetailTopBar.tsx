// src/components/listings/detail/DetailTopBar.tsx
import React from 'react';
import type { Listing } from '@/types';

interface DetailTopBarProps {
  listing: Listing;
  isFavorite: boolean;
  isAuthor: boolean;
  isConfirmingDelete: boolean;
  onClose: () => void;
  onEdit: (listing: Listing) => void;
  onDeleteClick: () => void;
  onCancelDelete: () => void;
  onToggleFavorite: (id: string) => void;
}

export const DetailTopBar: React.FC<DetailTopBarProps> = ({
  listing,
  isFavorite,
  isAuthor,
  isConfirmingDelete,
  onClose,
  onEdit,
  onDeleteClick,
  onCancelDelete,
  onToggleFavorite,
}) => {
  return (
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
              onClick={onDeleteClick}
              className={`${
                isConfirmingDelete
                  ? 'bg-red-800 text-white px-3 py-1 animate-pulse'
                  : 'text-red-800 border-b border-red-800/20'
              } hover:opacity-80 transition-all flex items-center gap-2 text-[10px] font-bold uppercase tracking-widest`}
            >
              {isConfirmingDelete ? 'Confirm Deletion?' : 'Delete'}
            </button>
            {isConfirmingDelete && (
              <button
                onClick={onCancelDelete}
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
          <span className="text-[10px] font-bold uppercase tracking-widest hidden sm:inline">
            {isFavorite ? 'Saved' : 'Save'}
          </span>
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
      </div>
    </div>
  );
};
