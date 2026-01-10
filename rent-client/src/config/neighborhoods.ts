// src/config/neighborhoods.ts
// Geographic neighborhood definitions for the application

export const NEIGHBORHOODS = [
  'Downtown Davis',
  'North Davis',
  'South Davis',
  'West Davis',
  'East Davis',
  'Wildhorse',
  'Mace Ranch',
  'Old North Davis',
] as const;

export type Neighborhood = (typeof NEIGHBORHOODS)[number];

// Default map center (Davis, CA)
export const DEFAULT_MAP_CENTER = {
  lat: 38.5449,
  lng: -121.7405,
} as const;

export const DEFAULT_MAP_ZOOM = 13;
