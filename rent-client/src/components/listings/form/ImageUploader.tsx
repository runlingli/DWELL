// src/components/listings/form/ImageUploader.tsx
import React, { useRef } from 'react';

interface ImageUploaderProps {
  imageUrl: string;
  onImageChange: (imageUrl: string) => void;
}

export const ImageUploader: React.FC<ImageUploaderProps> = ({ imageUrl, onImageChange }) => {
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onloadend = () => {
        onImageChange(reader.result as string);
      };
      reader.readAsDataURL(file);
    }
  };

  return (
    <div className="space-y-4">
      <label className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/60 block">
        Cover Aesthetic
      </label>
      <div
        onClick={() => fileInputRef.current?.click()}
        className="h-40 md:h-60 w-full border border-[#4a586e]/20 flex flex-col items-center justify-center cursor-pointer hover:bg-neutral-50 transition-colors overflow-hidden group relative"
      >
        {imageUrl ? (
          <>
            <img
              src={imageUrl}
              className="w-full h-full object-cover grayscale brightness-90 group-hover:brightness-100 transition-all"
              alt="Preview"
            />
            <div className="absolute inset-0 bg-black/20 opacity-0 group-hover:opacity-100 transition-opacity flex items-center justify-center">
              <span className="text-[10px] font-bold uppercase tracking-[0.3em] text-white">
                Change Image
              </span>
            </div>
          </>
        ) : (
          <>
            <svg
              className="w-10 h-10 text-[#4a586e]/20 mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            <span className="text-[10px] font-bold uppercase tracking-widest text-[#4a586e]/40">
              Upload Media
            </span>
          </>
        )}
        <input
          type="file"
          ref={fileInputRef}
          onChange={handleFileChange}
          className="hidden"
          accept="image/*"
        />
      </div>
    </div>
  );
};
