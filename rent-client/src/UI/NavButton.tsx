type NavButtonProps = {
  active?: boolean;
  onClick: () => void;
  children: React.ReactNode;
};

const NavButton: React.FC<NavButtonProps> = ({ active = false, onClick, children }) => (
  <button
    onClick={onClick}
    className={`pb-1 text-[10px] font-bold uppercase tracking-widest transition-colors underline-offset-12 ${
      active
        ? 'text-[#4a586e] underline decoration-1'
        : 'text-[#7e918b] hover:text-[#4a586e]'
    }`}
  >
    {children}
  </button>
);

export default NavButton;