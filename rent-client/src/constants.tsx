import type { Listing } from './types/types';

const now = new Date();
const nextMonth = new Date(now.getFullYear(), now.getMonth() + 1, 1).getTime();
const sixMonthsLater = new Date(now.getFullYear(), now.getMonth() + 6, 1).getTime();

export const NEIGHBORHOODS = [
  'Downtown Davis',
  'North Davis',
  'South Davis',
  'West Davis',
  'East Davis',
  'Wildhorse',
  'Mace Ranch',
  'Old North Davis'
];

export const INITIAL_LISTINGS: Listing[] = [
  {
    id: '1',
    title: 'Axis',
    price: 1200,
    location: 'Near E St & 2nd',
    neighborhood: 'Downtown Davis',
    coordinates: { lat: 38.5419, lng: -121.7405 },
    radius: 300,
    type: 'Apartment',
    imageUrl: 'https://images.unsplash.com/photo-1502672260266-1c1ef2d93688?q=80&w=2000&auto=format&fit=crop',
    additionalImages: [
      'https://images.unsplash.com/photo-1484154218962-a197022b5858?q=80&w=1000'
    ],
    description: 'Steps away from the UC Davis Arboretum. A minimalist haven with natural light and proximity to the vibrant downtown scene.',
    bedrooms: 2,
    bathrooms: 1,
    createdAt: Date.now() - 86400000,
    availableFrom: nextMonth,
    availableTo: sixMonthsLater,
    author: { name: 'Elena Rossi', avatar: 'https://i.pravatar.cc/150?u=elena' }
  },
  {
    id: '2',
    title: 'Greystone',
    price: 1100,
    location: 'Moore Blvd Area',
    neighborhood: 'North Davis',
    coordinates: { lat: 38.5623, lng: -121.7389 },
    radius: 400,
    type: 'Loft',
    imageUrl: 'https://images.unsplash.com/photo-1493663284031-b7e3aefcae8e?q=80&w=2000&auto=format&fit=crop',
    additionalImages: [
      'https://images.unsplash.com/photo-1556912177-c54030639a6d?q=80&w=1000'
    ],
    description: 'Modern lines meet suburban quietude. This loft features double-height ceilings and easy access to the greenbelt.',
    bedrooms: 2,
    bathrooms: 2,
    createdAt: Date.now() - 172800000,
    availableFrom: nextMonth,
    availableTo: sixMonthsLater + 8640000000,
    author: { name: 'Marcus Chen', avatar: 'https://i.pravatar.cc/150?u=marcus' }
  },
  {
    id: '3',
    title: 'The Green',
    price: 1250,
    location: 'Olive Dr & Richards',
    neighborhood: 'South Davis',
    coordinates: { lat: 38.5385, lng: -121.7345 },
    radius: 250,
    type: 'Studio',
    imageUrl: 'https://plus.unsplash.com/premium_photo-1676968002512-3eac82b1d847?q=80&w=687&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D?q=80&w=2000&auto=format&fit=crop',
    description: 'A stark, beautiful studio designed for the dedicated academic. Quiet, refined, and perfectly positioned near campus.',
    bedrooms: 0,
    bathrooms: 1,
    createdAt: Date.now() - 259200000,
    availableFrom: now.getTime(),
    availableTo: sixMonthsLater,
    author: { name: 'Sarah Miller', avatar: 'https://i.pravatar.cc/150?u=sarah' }
  },
  {
    id: '4',
    title: 'Tanglewood',
    price: 1300,
    location: 'Wildhorse Golf Course',
    neighborhood: 'Wildhorse',
    coordinates: { lat: 38.5750, lng: -121.7100 },
    radius: 500,
    type: 'House',
    imageUrl: 'https://images.unsplash.com/photo-1512917774080-9991f1c4c750?q=80&w=2000&auto=format&fit=crop',
    description: 'Expansive living on the edge of the city. Unobstructed views of the surrounding valley and uncompromising modern architecture.',
    bedrooms: 4,
    bathrooms: 3,
    createdAt: Date.now() - 345600000,
    availableFrom: now.getTime(),
    availableTo: sixMonthsLater,
    author: { name: 'Marcus Chen', avatar: 'https://i.pravatar.cc/150?u=marcus' }
  }
];