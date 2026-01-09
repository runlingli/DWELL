
interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  uppercase?: boolean;
}

export const Input: React.FC<InputProps> = ({ uppercase = false, className = '', ...props }) => (
  <input 
	className={`w-full bg-transparent border-b border-[#4a586e]/20 py-2 px-0 focus:border-[#4a586e] outline-none transition-all placeholder:text-[#4a586e]/30 text-[11px] tracking-widest font-bold text-[#4a586e] ${uppercase ? 'uppercase' : ''} ${className}`}
	{...props}
  />
);