// src/hooks/useListingForm.ts
import { useState, useEffect } from 'react';
import type { Listing } from '@/types';
import { NEIGHBORHOODS } from '@/config';

export interface ListingFormData {
  title: string;
  price: string;
  type: Listing['type'];
  neighborhood: string;
  address: string;
  bedrooms: string;
  bathrooms: string;
  description: string;
  availableFrom: string;
  availableTo: string;
  imageUrl: string;
  lat: number;
  lng: number;
  radius: number;
}

const DEFAULT_FORM_DATA: ListingFormData = {
  title: '',
  price: '',
  type: 'Apartment',
  neighborhood: NEIGHBORHOODS[0],
  address: '',
  bedrooms: '1',
  bathrooms: '1',
  description: '',
  availableFrom: '',
  availableTo: '',
  imageUrl: '',
  lat: 38.5449,
  lng: -121.7405,
  radius: 300,
};

export function useListingForm(initialData: Listing | null | undefined, isOpen: boolean) {
  const [formData, setFormData] = useState<ListingFormData>(DEFAULT_FORM_DATA);

  useEffect(() => {
    if (initialData) {
      setFormData({
        title: initialData.title,
        price: initialData.price.toString(),
        type: initialData.type,
        neighborhood: initialData.neighborhood,
        address: initialData.location || '',
        bedrooms: initialData.bedrooms.toString(),
        bathrooms: initialData.bathrooms.toString(),
        description: initialData.description,
        availableFrom: new Date(initialData.availableFrom).toISOString().split('T')[0],
        availableTo: new Date(initialData.availableTo).toISOString().split('T')[0],
        imageUrl: initialData.imageUrl,
        lat: initialData.coordinates.lat,
        lng: initialData.coordinates.lng,
        radius: initialData.radius || 300,
      });
    } else {
      setFormData(DEFAULT_FORM_DATA);
    }
  }, [initialData, isOpen]);

  const updateField = <K extends keyof ListingFormData>(field: K, value: ListingFormData[K]) => {
    setFormData((prev) => ({ ...prev, [field]: value }));
  };

  const handleFromDateChange = (newFromDate: string) => {
    setFormData((prev) => {
      const updatedToDate = prev.availableTo && newFromDate > prev.availableTo ? '' : prev.availableTo;
      return { ...prev, availableFrom: newFromDate, availableTo: updatedToDate };
    });
  };

  const setCoordinates = (lat: number, lng: number) => {
    setFormData((prev) => ({ ...prev, lat, lng }));
  };

  const buildListingData = (existingData?: Listing | null): Listing => ({
    id: existingData?.id || Math.random().toString(36).substr(2, 9),
    title: formData.title,
    price: Number(formData.price),
    type: formData.type,
    neighborhood: formData.neighborhood,
    location: formData.address,
    coordinates: { lat: formData.lat, lng: formData.lng },
    radius: formData.radius,
    bedrooms: Number(formData.bedrooms),
    bathrooms: Number(formData.bathrooms),
    imageUrl:
      formData.imageUrl ||
      'https://images.unsplash.com/photo-1493809842364-78817add7ffb?q=80&w=2000&auto=format&fit=crop',
    description: formData.description || 'A newly listed minimalist property.',
    availableFrom: new Date(formData.availableFrom).getTime(),
    availableTo: new Date(formData.availableTo).getTime(),
    createdAt: existingData?.createdAt || Date.now(),
    author: existingData?.author || { name: 'Current User' },
  });

  return {
    formData,
    updateField,
    handleFromDateChange,
    setCoordinates,
    buildListingData,
  };
}
