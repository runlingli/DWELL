// src/components/map/MapView.tsx
import React, { useEffect, useRef } from 'react';
import type { Listing } from '@/types';

declare const L: any;

interface MapViewProps {
  listings: Listing[];
  onMarkerClick: (listing: Listing) => void;
  center?: { lat: number; lng: number };
}

export const MapView: React.FC<MapViewProps> = ({ listings, onMarkerClick, center }) => {
  const mapRef = useRef<any>(null);
  const containerRef = useRef<HTMLDivElement>(null);
  const layersRef = useRef<any[]>([]);

  useEffect(() => {
    if (!containerRef.current) return;

    if (!mapRef.current) {
      mapRef.current = L.map(containerRef.current, {
        center: [38.5449, -121.7405], // Davis, CA Center
        zoom: 13,
        zoomControl: false,
        attributionControl: false,
        scrollWheelZoom: true,
      });

      L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png', {
        maxZoom: 19,
      }).addTo(mapRef.current);
    }

    const map = mapRef.current;

    // Clear previous layers
    layersRef.current.forEach((layer) => map.removeLayer(layer));
    layersRef.current = [];

    listings.forEach((listing) => {
      // Create a stylized circular area (approximate boundary)
      const circle = L.circle([listing.coordinates.lat, listing.coordinates.lng], {
        radius: listing.radius || 300,
        color: '#4a586e',
        weight: 1,
        opacity: 0.3,
        fillColor: '#4a586e',
        fillOpacity: 0.05,
        interactive: false,
      }).addTo(map);

      // Create custom price tag icon
      const priceIcon = L.divIcon({
        className: 'price-tag-container',
        html: `<div class="price-tag-marker">$${listing.price.toLocaleString()}</div>`,
        iconSize: [60, 24],
        iconAnchor: [30, 12],
      });

      const priceMarker = L.marker([listing.coordinates.lat, listing.coordinates.lng], {
        icon: priceIcon,
      }).addTo(map);

      priceMarker.on('click', () => onMarkerClick(listing));

      layersRef.current.push(circle, priceMarker);
    });

    if (center && map) {
      map.setView([center.lat, center.lng], 15, { animate: true });
    }

    // Ensure the map resizes correctly if the parent container changes
    const resizeObserver = new ResizeObserver(() => {
      map.invalidateSize();
    });
    resizeObserver.observe(containerRef.current);

    return () => {
      resizeObserver.disconnect();
    };
  }, [listings, center, onMarkerClick]);

  return (
    <div className="w-full h-75 md:h-112.5 lg:h-137.5 relative overflow-hidden group">
      <div ref={containerRef} className="w-full h-full" />
      <div className="absolute top-4 left-4 z-500 bg-[#4a586e] text-[#f3e9d2] px-3 py-1 text-[8px] font-bold uppercase tracking-widest pointer-events-none">
        DAVIS
      </div>
      <div className="absolute bottom-4 right-4 z-500 bg-[#f3e9d2]/80 px-2 py-1 text-[8px] font-bold uppercase tracking-widest text-[#4a586e] opacity-0 group-hover:opacity-100 transition-opacity">
        Interactive Discovery
      </div>
    </div>
  );
};
