// src/components/listings/form/ListingFormFields.tsx
import React from 'react';
import { Input, Select } from '@/ui';
import { NEIGHBORHOODS } from '@/config';
import type { Listing } from '@/types';
import type { ListingFormData } from '@/hooks';

const TYPE_OPTIONS = [
  { value: 'Apartment', label: 'Apartment' },
  { value: 'House', label: 'House' },
  { value: 'Studio', label: 'Studio' },
  { value: 'Loft', label: 'Loft' },
];

const NEIGHBORHOOD_OPTIONS = NEIGHBORHOODS.map((n) => ({ value: n, label: n }));

interface ListingFormFieldsProps {
  formData: ListingFormData;
  onFieldChange: <K extends keyof ListingFormData>(field: K, value: ListingFormData[K]) => void;
  onFromDateChange: (date: string) => void;
}

export const ListingFormFields: React.FC<ListingFormFieldsProps> = ({
  formData,
  onFieldChange,
  onFromDateChange,
}) => {
  return (
    <div className="space-y-7">
      {/* Identity Section */}
      <div>
        <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">
          Identity
        </label>
        <Input
          placeholder="TITLE"
          value={formData.title}
          onChange={(e) => onFieldChange('title', e.target.value)}
          uppercase
          required
        />

        <div className="grid grid-cols-2 gap-4 mt-6">
          <Select
            label="Neighborhood"
            options={NEIGHBORHOOD_OPTIONS}
            value={formData.neighborhood}
            onChange={(val) => onFieldChange('neighborhood', val)}
          />
          <Select
            label="Property Type"
            options={TYPE_OPTIONS}
            value={formData.type}
            onChange={(val) => onFieldChange('type', val as Listing['type'])}
          />
        </div>
      </div>

      {/* Specifications Section */}
      <div className="space-y-2">
        <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block mb-2">
          Specifications
        </label>
        <div className="grid grid-cols-3 gap-8">
          <div className="flex flex-col">
            <label className="text-[10px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">
              Monthly Price ($)
            </label>
            <Input
              placeholder="e.g. 2400"
              type="number"
              value={formData.price}
              onChange={(e) => onFieldChange('price', e.target.value)}
              required
            />
          </div>
          <div className="flex flex-col">
            <label className="text-[10px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">
              Bedrooms
            </label>
            <Input
              placeholder="e.g. 2"
              type="number"
              value={formData.bedrooms}
              onChange={(e) => onFieldChange('bedrooms', e.target.value)}
            />
          </div>
          <div className="flex flex-col">
            <label className="text-[10px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">
              Bathrooms
            </label>
            <Input
              placeholder="e.g. 1"
              type="number"
              value={formData.bathrooms}
              onChange={(e) => onFieldChange('bathrooms', e.target.value)}
            />
          </div>
        </div>
        <div className="flex flex-col pt-2">
          <label className="text-[10px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">
            Location / Cross Streets
          </label>
          <Input
            placeholder="E.G. NEAR E ST & 2ND"
            value={formData.address}
            onChange={(e) => onFieldChange('address', e.target.value)}
            uppercase
            required
          />
        </div>
      </div>

      {/* Availability Section */}
      <div className="space-y-0.5">
        <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">
          Availability Window
        </label>
        <div className="grid grid-cols-2 gap-8">
          <div>
            <p className="text-[10px] uppercase tracking-widest text-[#7e918b] mb-1 font-bold">From</p>
            <Input
              type="date"
              value={formData.availableFrom}
              onChange={(e) => onFromDateChange(e.target.value)}
              required
            />
          </div>
          <div>
            <p className="text-[10px] uppercase tracking-widest text-[#7e918b] mb-1 font-bold">Until</p>
            <Input
              type="date"
              value={formData.availableTo}
              onChange={(e) => onFieldChange('availableTo', e.target.value)}
              min={formData.availableFrom}
              required
            />
          </div>
        </div>
      </div>

      {/* Description Section */}
      <div className="space-y-4">
        <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">
          The Narrative
        </label>
        <textarea
          placeholder="DESCRIBE THE SPACE AND ITS ESSENCE..."
          className="w-full bg-transparent border border-[#4a586e]/20 p-6 focus:border-[#4a586e] outline-none transition-colors placeholder:text-[#4a586e]/30 text-[11px] h-40 resize-none font-bold uppercase tracking-widest"
          value={formData.description}
          onChange={(e) => onFieldChange('description', e.target.value)}
          required
        />
      </div>
    </div>
  );
};
