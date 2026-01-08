import React, { useState, useRef, useEffect } from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline';
}

export const Button: React.FC<ButtonProps> = ({ 
  children, 
  variant = 'primary', 
  className = '', 
  ...props 
}) => {
  const baseStyles = "px-6 py-4 transition-all duration-300 uppercase tracking-[0.3em] text-[10px] font-bold focus:outline-none border border-[#4a586e]";
  const variants = {
    primary: "bg-[#4a586e] text-[#f3e9d2] hover:bg-transparent hover:text-[#4a586e]",
    secondary: "bg-[#f3e9d2] text-[#4a586e] hover:bg-[#4a586e] hover:text-[#f3e9d2]",
    outline: "bg-transparent text-[#4a586e] border-[#4a586e] hover:bg-[#4a586e] hover:text-[#f3e9d2]"
  };

  return (
    <button 
      className={`${baseStyles} ${variants[variant]} ${className}`}
      {...props}
    >
      {children}
    </button>
  );
};


interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  uppercase?: boolean;
}

export const Input: React.FC<InputProps> = ({ uppercase = false, className = '', ...props }) => (
  <input 
    className={`w-full bg-transparent border-b border-[#4a586e]/20 py-4 px-0 focus:border-[#4a586e] outline-none transition-all placeholder:text-[#4a586e]/30 text-[11px] tracking-widest font-bold text-[#4a586e] ${uppercase ? 'uppercase' : ''} ${className}`}
    {...props}
  />
);

interface SelectOption {
  value: string;
  label: string;
}

interface SelectProps {
  options: SelectOption[];
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  label?: string;
  className?: string;
}

export const Select: React.FC<SelectProps> = ({ options, value, onChange, placeholder, label, className = '' }) => {
  const [isOpen, setIsOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (containerRef.current && !containerRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };
    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  const selectedOption = options.find(opt => opt.value === value);

  return (
    <div className={`relative ${className}`} ref={containerRef}>
      {label && <label className="text-[8px] font-bold uppercase tracking-widest text-[#7e918b] mb-1 block">{label}</label>}
      <button
        type="button"
        onClick={() => setIsOpen(!isOpen)}
        className="w-full bg-transparent border-b border-[#4a586e]/20 py-4 px-0 text-left flex justify-between items-center focus:border-[#4a586e] transition-all outline-none"
      >
        <span className="text-[11px] uppercase tracking-widest font-bold text-[#4a586e]">
          {selectedOption ? selectedOption.label : placeholder}
        </span>
        <svg 
          className={`w-3 h-3 text-[#4a586e] transition-transform duration-300 ${isOpen ? 'rotate-180' : ''}`} 
          fill="none" stroke="currentColor" viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {isOpen && (
        <div className="absolute z-[3000] w-full mt-0 bg-[#f3e9d2] border border-[#4a586e] shadow-2xl animate-in fade-in slide-in-from-top-2 duration-200">
          <div className="max-h-60 overflow-y-auto custom-scrollbar">
            {options.map((option) => (
              <button
                key={option.value}
                type="button"
                className={`w-full text-left px-4 py-3 text-[10px] uppercase tracking-widest font-bold transition-colors border-b border-[#4a586e]/5 last:border-0 ${
                  value === option.value 
                    ? 'bg-[#4a586e] text-[#f3e9d2]' 
                    : 'text-[#4a586e] hover:bg-[#4a586e]/5'
                }`}
                onClick={() => {
                  onChange(option.value);
                  setIsOpen(false);
                }}
              >
                {option.label}
              </button>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
  title?: string;
  size?: 'md' | 'lg' | 'xl' | '4xl' | '6xl';
}

export const Modal: React.FC<ModalProps> = ({ 
  isOpen, 
  onClose, 
  children,
  title,
  size = 'md'
}) => {
  if (!isOpen) return null;

  const sizeClasses = {
    'md': 'max-w-md',
    'lg': 'max-w-lg',
    'xl': 'max-w-xl',
    '4xl': 'max-w-4xl',
    '6xl': 'max-w-6xl'
  };

  return (
    <div className="fixed inset-0 z-[2000] flex items-center justify-center p-4">
      <div 
        className="absolute inset-0 bg-[#4a586e]/60 backdrop-blur-none" 
        onClick={onClose} 
      />
      <div className={`relative bg-[#f3e9d2] border border-[#4a586e] w-full ${sizeClasses[size]} p-8 md:p-12 lg:p-16 animate-in fade-in zoom-in duration-300 overflow-hidden flex flex-col max-h-[90vh]`}>
        <button 
          onClick={onClose}
          className="absolute top-8 right-8 text-[#4a586e] hover:opacity-50 transition-opacity z-10"
        >
          <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        {title && (
          <h2 className="font-serif text-3xl md:text-5xl mb-8 md:mb-12 tracking-tighter text-[#4a586e] uppercase leading-none">{title}</h2>
        )}
        <div className="flex-grow overflow-y-auto custom-scrollbar pr-2">
          {children}
        </div>
      </div>
    </div>
  );
};