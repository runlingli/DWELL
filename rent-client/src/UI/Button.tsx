
interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'outline';
}


const Button: React.FC<ButtonProps> = ({ 
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

export default Button;