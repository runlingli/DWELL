// src/components/listings/CreateListingModal.tsx
import React from 'react';
import { Modal, Button } from '@/ui';
import type { Listing } from '@/types';
import { useListingForm } from '@/hooks';
import { ImageUploader } from './form/ImageUploader';
import { MapPicker } from './form/MapPicker';
import { ListingFormFields } from './form/ListingFormFields';

interface CreateListingModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (listing: Listing) => void;
  initialData?: Listing | null;
}

export const CreateListingModal: React.FC<CreateListingModalProps> = ({
  isOpen,
  onClose,
  onSave,
  initialData,
}) => {
  const { formData, updateField, handleFromDateChange, setCoordinates, buildListingData } =
    useListingForm(initialData, isOpen);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const data = buildListingData(initialData);
    onSave(data);
    onClose();
  };

  return (
    <Modal
      size="5xl"
      isOpen={isOpen}
      onClose={onClose}
      title={initialData ? 'EDIT LISTING' : 'NEW LISTING'}
    >
      <form onSubmit={handleSubmit} className="pb-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-5 lg:gap-10">
          <div className="space-y-10">
            <ImageUploader
              imageUrl={formData.imageUrl}
              onImageChange={(url) => updateField('imageUrl', url)}
            />

            <MapPicker
              lat={formData.lat}
              lng={formData.lng}
              radius={formData.radius}
              isOpen={isOpen}
              onLocationChange={setCoordinates}
              onRadiusChange={(r) => updateField('radius', r)}
            />
          </div>

          <div>
            <ListingFormFields
              formData={formData}
              onFieldChange={updateField}
              onFromDateChange={handleFromDateChange}
            />

            <div className="mt-7">
              <Button type="submit" className="w-full !py-4 text-sm tracking-[0.4em]">
                {initialData ? 'UPDATE POSTING' : 'CONFIRM & PUBLISH'}
              </Button>
            </div>
          </div>
        </div>
      </form>
    </Modal>
  );
};
