// src/components/listings/detail/ListingDetail.tsx
import React from 'react';
import type { Listing, User } from '@/types';
import { getUserDisplayName } from '@/utils';
import { useDeleteConfirmation } from '@/hooks';
import { DetailTopBar } from './DetailTopBar';
import { DetailContent } from './DetailContent';

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
  onDelete,
}) => {
  const { isConfirming, handleClick, cancel } = useDeleteConfirmation(() => onDelete(listing.id));

  if (!isOpen) return null;

  const isAuthor = currentUser ? getUserDisplayName(currentUser) === listing.author.name : false;

  return (
    <div className="fixed inset-0 z-[1000] flex items-center justify-center p-4 md:p-10 lg:p-16">
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-[#4a586e]/60 backdrop-blur-sm animate-in fade-in duration-500"
        onClick={onClose}
      />

      {/* Detail Container */}
      <div className="relative w-full max-w-6xl h-full max-h-[90vh] bg-[#f3e9d2] border border-[#4a586e] shadow-2xl animate-in fade-in zoom-in-95 duration-500 flex flex-col overflow-hidden">
        <DetailTopBar
          listing={listing}
          isFavorite={isFavorite}
          isAuthor={isAuthor}
          isConfirmingDelete={isConfirming}
          onClose={onClose}
          onEdit={onEdit}
          onDeleteClick={handleClick}
          onCancelDelete={cancel}
          onToggleFavorite={onToggleFavorite}
        />

        <DetailContent listing={listing} />
      </div>
    </div>
  );
};
