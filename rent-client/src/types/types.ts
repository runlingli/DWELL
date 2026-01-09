export interface Listing {
  id: string;
  title: string;
  price: number;
  location?: string; // Street name or cross streets for approximate display
  neighborhood: string;
  coordinates: {
    lat: number;
    lng: number;
  };
  radius: number; // radius in meters for the approximate area
  type: 'Apartment' | 'House' | 'Studio' | 'Loft';
  imageUrl: string;
  additionalImages?: string[];
  description: string;
  bedrooms: number;
  bathrooms: number;
  createdAt: number;
  availableFrom: number; // timestamp
  availableTo: number;   // timestamp
  author: {
    name: string;
    avatar?: string;
  };
}

export interface User {
  id?: string;
  first_name: string;
  last_name: string;
  email: string;
  avatar?: string;
}

// Helper to get display name from user
export const getUserDisplayName = (user: User): string => {
  return `${user.first_name} ${user.last_name}`.trim();
}

export interface AppState {
  listings: Listing[];
  favorites: string[];
  currentUser: User | null;
}