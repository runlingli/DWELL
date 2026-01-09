import { useEffect, useRef, useState } from "react";

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
	  {label && <label className="text-[10px] font-bold uppercase tracking-widest text-[#7e918b] mb-1 block">{label}</label>}
	  <button
		type="button"
		onClick={() => setIsOpen(!isOpen)}
		className="w-full bg-transparent border-b border-[#4a586e]/20 py-2 px-0 text-left flex justify-between items-center focus:border-[#4a586e] transition-all outline-none"
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
		<div className="absolute z-3000 w-full mt-0 bg-[#f3e9d2] border border-[#4a586e] shadow-2xl animate-in fade-in slide-in-from-top-2 duration-200">
		  <div className="max-h-60 overflow-y-auto custom-scrollbar">
			{options.map((option) => (
			  <button
				key={option.value}
				type="button"
				className={`w-full text-left px-4 py-2 text-[10px] uppercase tracking-widest font-bold transition-colors border-b border-[#4a586e]/5 last:border-0 ${
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
