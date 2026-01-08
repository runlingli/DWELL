import React, { useState, useRef, useEffect } from 'react';
import { Modal, Button, Input, Select } from './UI';
import { Listing } from '../types';
import { NEIGHBORHOODS } from '../constants';

declare const L: any;

interface CreateListingModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSave: (listing: Listing) => void;
  initialData?: Listing | null;
}

const TYPE_OPTIONS = [
  { value: 'Apartment', label: 'Apartment' },
  { value: 'House', label: 'House' },
  { value: 'Studio', label: 'Studio' },
  { value: 'Loft', label: 'Loft' },
];

const NEIGHBORHOOD_OPTIONS = NEIGHBORHOODS.map(n => ({ value: n, label: n }));

export const CreateListingModal: React.FC<CreateListingModalProps> = ({ isOpen, onClose, onSave, initialData }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const mapContainerRef = useRef<HTMLDivElement>(null);
  const pickerMapRef = useRef<any>(null);
  const pickerCircleRef = useRef<any>(null);

  const [formData, setFormData] = useState({
    title: '',
    price: '',
    type: 'Apartment' as Listing['type'],
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
    radius: 300
  });

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
        radius: initialData.radius || 300
      });
    } else {
      setFormData({
        title: '', price: '', type: 'Apartment', neighborhood: NEIGHBORHOODS[0], address: '', bedrooms: '1', bathrooms: '1',
        description: '', availableFrom: '', availableTo: '', imageUrl: '', lat: 38.5449, lng: -121.7405, radius: 300
      });
    }
  }, [initialData, isOpen]);

  // Handle Location Picker Map
  useEffect(() => {
    if (isOpen && mapContainerRef.current) {
      setTimeout(() => {
        if (!pickerMapRef.current) {
          pickerMapRef.current = L.map(mapContainerRef.current, {
            center: [formData.lat, formData.lng],
            zoom: 14,
            zoomControl: false,
            attributionControl: false
          });

          L.tileLayer('https://{s}.basemaps.cartocdn.com/light_all/{z}/{x}/{y}{r}.png').addTo(pickerMapRef.current);

          pickerMapRef.current.on('click', (e: any) => {
            const { lat, lng } = e.latlng;
            setFormData(prev => ({ ...prev, lat, lng }));
          });
        }

        if (pickerCircleRef.current) {
          pickerMapRef.current.removeLayer(pickerCircleRef.current);
        }

        pickerCircleRef.current = L.circle([formData.lat, formData.lng], {
          radius: formData.radius,
          color: '#4a586e',
          weight: 1,
          fillColor: '#4a586e',
          fillOpacity: 0.2
        }).addTo(pickerMapRef.current);

        pickerMapRef.current.setView([formData.lat, formData.lng], pickerMapRef.current.getZoom());
      }, 100);
    }

    return () => {
      if (!isOpen && pickerMapRef.current) {
        pickerMapRef.current.remove();
        pickerMapRef.current = null;
      }
    };
  }, [isOpen, formData.lat, formData.lng, formData.radius]);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        setFormData(prev => ({ ...prev, imageUrl: reader.result as string }));
      };
      reader.readAsDataURL(file);
    }
  };

  const handleFromDateChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newFromDate = e.target.value;
    setFormData(prev => {
      const updatedToDate = (prev.availableTo && newFromDate > prev.availableTo) ? '' : prev.availableTo;
      return { ...prev, availableFrom: newFromDate, availableTo: updatedToDate };
    });
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSave({
      id: initialData?.id || Math.random().toString(36).substr(2, 9),
      title: formData.title,
      price: Number(formData.price),
      type: formData.type,
      neighborhood: formData.neighborhood,
      location: formData.address,
      coordinates: { lat: formData.lat, lng: formData.lng },
      radius: formData.radius,
      bedrooms: Number(formData.bedrooms),
      bathrooms: Number(formData.bathrooms),
      imageUrl: formData.imageUrl || 'https://images.unsplash.com/photo-1493809842364-78817add7ffb?q=80&w=2000&auto=format&fit=crop',
      description: formData.description || 'A newly listed minimalist property.',
      availableFrom: new Date(formData.availableFrom).getTime(),
      availableTo: new Date(formData.availableTo).getTime(),
      createdAt: initialData?.createdAt || Date.now(),
      author: initialData?.author || { name: 'Current User' } 
    });
    onClose();
  };

  return (
    <Modal size="6xl" isOpen={isOpen} onClose={onClose} title={initialData ? "EDIT LISTING" : "NEW LISTING"}>
      <form onSubmit={handleSubmit} className="pb-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 lg:gap-24">
          
          {/* Left Column: Visuals and Spatial */}
          <div className="space-y-10">
            <div className="space-y-4">
               <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">Cover Aesthetic</label>
               <div 
                onClick={() => fileInputRef.current?.click()}
                className="aspect-[4/3] w-full border border-[#4a586e]/20 flex flex-col items-center justify-center cursor-pointer hover:bg-neutral-50 transition-colors overflow-hidden group relative"
              >
                {formData.imageUrl ? (
                  <>
                    <img src={formData.imageUrl} className="w-full h-full object-cover grayscale brightness-90 group-hover:brightness-100 transition-all" alt="Preview" />
                    <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
                       <span className="text-[10px] font-bold uppercase tracking-[0.3em] text-white">Change Image</span>
                    </div>
                  </>
                ) : (
                  <>
                    <svg className="w-10 h-10 text-[#4a586e]/20 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    <span className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/40">Upload Media</span>
                  </>
                )}
                <input type="file" ref={fileInputRef} onChange={handleFileChange} className="hidden" accept="image/*" />
              </div>
            </div>

            <div className="space-y-4">
              <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">Geographic Boundary</label>
              <div ref={mapContainerRef} className="w-full h-64 border border-[#4a586e]/20 grayscale" />
              <div className="space-y-4 pt-2">
                <div className="flex justify-between items-center">
                  <span className="text-[9px] font-bold uppercase tracking-widest text-[#4a586e]">Approx. Radius</span>
                  <span className="text-[11px] font-bold text-[#4a586e]">{formData.radius}m</span>
                </div>
                <input 
                  type="range" min="100" max="1000" step="50" 
                  value={formData.radius} 
                  onChange={e => setFormData({...formData, radius: Number(e.target.value)})}
                  className="w-full h-[1px] bg-[#4a586e]/20 appearance-none cursor-pointer accent-[#4a586e]"
                />
              </div>
            </div>
          </div>

          {/* Right Column: Information and Specs */}
          <div className="space-y-10">
            <div className="space-y-6">
              <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block mb-[-1.5rem]">Identity</label>
              <Input 
                placeholder="TITLE" 
                value={formData.title} 
                onChange={e => setFormData({...formData, title: e.target.value})}
                uppercase
                required 
              />
              
              <div className="grid grid-cols-2 gap-8">
                <Select
                  label="Neighborhood"
                  options={NEIGHBORHOOD_OPTIONS}
                  value={formData.neighborhood}
                  onChange={(val) => setFormData({...formData, neighborhood: val})}
                />
                <Select
                  label="Property Type"
                  options={TYPE_OPTIONS}
                  value={formData.type}
                  onChange={(val) => setFormData({...formData, type: val as Listing['type']})}
                />
              </div>
            </div>

            <div className="space-y-6">
               <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block mb-2">Specifications</label>
               <div className="grid grid-cols-3 gap-8">
                  <div className="flex flex-col">
                    <label className="text-[8px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">Monthly Price ($)</label>
                    <Input 
                      placeholder="e.g. 2400" 
                      type="number"
                      value={formData.price} 
                      onChange={e => setFormData({...formData, price: e.target.value})}
                      required 
                    />
                  </div>
                  <div className="flex flex-col">
                    <label className="text-[8px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">Bedrooms</label>
                    <Input 
                      placeholder="e.g. 2" 
                      type="number" 
                      value={formData.bedrooms} 
                      onChange={e => setFormData({...formData, bedrooms: e.target.value})} 
                    />
                  </div>
                  <div className="flex flex-col">
                    <label className="text-[8px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">Bathrooms</label>
                    <Input 
                      placeholder="e.g. 1" 
                      type="number" 
                      value={formData.bathrooms} 
                      onChange={e => setFormData({...formData, bathrooms: e.target.value})} 
                    />
                  </div>
               </div>
               <div className="flex flex-col pt-2">
                 <label className="text-[8px] font-bold uppercase tracking-widest text-[#7e918b] mb-1">Location / Cross Streets</label>
                 <Input 
                    placeholder="E.G. NEAR E ST & 2ND" 
                    value={formData.address} 
                    onChange={e => setFormData({...formData, address: e.target.value})}
                    uppercase
                    required 
                  />
               </div>
            </div>

            <div className="space-y-6">
              <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block mb-[-1.5rem]">Availability Window</label>
              <div className="grid grid-cols-2 gap-8">
                <div>
                  <p className="text-[8px] uppercase tracking-widest text-[#7e918b] mb-1 font-bold">From</p>
                  <Input type="date" value={formData.availableFrom} onChange={handleFromDateChange} required />
                </div>
                <div>
                  <p className="text-[8px] uppercase tracking-widest text-[#7e918b] mb-1 font-bold">Until</p>
                  <Input type="date" value={formData.availableTo} onChange={e => setFormData({...formData, availableTo: e.target.value})} min={formData.availableFrom} required />
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">The Narrative</label>
              <textarea 
                placeholder="DESCRIBE THE SPACE AND ITS ESSENCE..."
                className="w-full bg-transparent border border-[#4a586e]/20 p-6 focus:border-[#4a586e] outline-none transition-colors placeholder:text-[#4a586e]/30 text-[11px] h-40 resize-none font-bold uppercase tracking-widest"
                value={formData.description}
                onChange={e => setFormData({...formData, description: e.target.value})}
                required
              />
            </div>

            <div className="pt-8">
               <Button type="submit" className="w-full !py-6 text-sm tracking-[0.4em]">{initialData ? "UPDATE POSTING" : "CONFIRM & PUBLISH"}</Button>
            </div>
          </div>
        </div>
      </form>
    </Modal>
  );
};