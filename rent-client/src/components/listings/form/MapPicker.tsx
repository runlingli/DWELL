// src/components/listings/form/MapPicker.tsx
import React, { useRef, useEffect } from 'react';

declare const L: any;

interface MapPickerProps {
  lat: number;
  lng: number;
  radius: number;
  isOpen: boolean;
  onLocationChange: (lat: number, lng: number) => void;
  onRadiusChange: (radius: number) => void;
}

export const MapPicker: React.FC<MapPickerProps> = ({
  lat,
  lng,
  radius,
  isOpen,
  onLocationChange,
  onRadiusChange,
}) => {
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const pickerMapRef = useRef<any>(null);
  const pickerCircleRef = useRef<any>(null);

  useEffect(() => {
    if (isOpen && mapContainerRef.current) {
      const timer = setTimeout(() => {
        if (!pickerMapRef.current) {
          pickerMapRef.current = L.map(mapContainerRef.current, {
            center: [lat, lng],
            zoom: 14,
            zoomControl: false,
            attributionControl: false,
          });

          L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png').addTo(
            pickerMapRef.current
          );

          pickerMapRef.current.on('click', (e: any) => {
            const { lat: newLat, lng: newLng } = e.latlng;
            onLocationChange(newLat, newLng);
          });
        }

        if (pickerCircleRef.current) {
          pickerMapRef.current.removeLayer(pickerCircleRef.current);
        }

        pickerCircleRef.current = L.circle([lat, lng], {
          radius: radius,
          color: '#4a586e',
          weight: 1,
          fillColor: '#4a586e',
          fillOpacity: 0.2,
        }).addTo(pickerMapRef.current);

        pickerMapRef.current.setView([lat, lng], pickerMapRef.current.getZoom());
      }, 100);

      return () => clearTimeout(timer);
    }

    return () => {
      if (!isOpen && pickerMapRef.current) {
        pickerMapRef.current.remove();
        pickerMapRef.current = null;
      }
    };
  }, [isOpen, lat, lng, radius, onLocationChange]);

  return (
    <div className="space-y-4">
      <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">
        Geographic Boundary
      </label>
      <div ref={mapContainerRef} className="w-full h-40 md:h-60 border border-[#4a586e]/20 grayscale" />
      <div className="space-y-4 pt-2">
        <div className="flex justify-between items-center">
          <span className="text-[9px] font-bold uppercase tracking-widest text-[#4a586e]">
            Approx. Radius
          </span>
          <span className="text-[11px] font-bold text-[#4a586e]">{radius}m</span>
        </div>
        <input
          type="range"
          min="100"
          max="1000"
          step="50"
          value={radius}
          onChange={(e) => onRadiusChange(Number(e.target.value))}
          className="w-full h-[1px] bg-[#4a586e]/20 appearance-none cursor-pointer accent-[#4a586e]"
        />
      </div>
    </div>
  );
};
