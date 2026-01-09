interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  children: React.ReactNode;
  title?: string;
  size?: 'md' | 'lg' | 'xl' | '4xl' | '5xl' |'6xl';
}

const Modal: React.FC<ModalProps> = ({ 
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
	'5xl': 'max-w-5xl',
	'6xl': 'max-w-6xl'
  };

  return (
	<div className="fixed inset-0 z-2000 flex items-center justify-center p-4">
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
		<div className="grow overflow-y-auto overflow-x-hidden custom-scrollbar">
		  {children}
		</div>
	  </div>
	</div>
  );
};

export default Modal;