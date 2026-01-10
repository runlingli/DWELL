// src/components/layout/Header.tsx
import React from 'react';
import { useAuthStore, useUIStore, type ViewType } from '@/stores';
import { Button } from '@/ui';

export const Header: React.FC = () => {
  const { currentUser, logout } = useAuthStore();
  const { view, navigate, resetToHome, openAuthModal, openCreateModal } = useUIStore();

  const handleNavClick = (targetView: ViewType) => {
    if (targetView === 'profile' && !currentUser) {
      openAuthModal();
      return;
    }
    navigate(targetView);
  };

  const handleLogout = async () => {
    await logout();
    navigate('discover');
  };

  return (
    <header className="fixed top-0 left-0 right-0 bg-[#f3e9d2]/90 backdrop-blur-md z-[600] border-b border-[#4a586e]/10">
      <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
        <button
          onClick={resetToHome}
          className="font-serif text-3xl tracking-tighter text-[#4a586e] transition-opacity hover:opacity-60"
        >
          TinyRent
        </button>

        <nav className="hidden md:flex items-center gap-12 text-[10px] font-bold uppercase tracking-[0.3em]">
          <button
            onClick={() => handleNavClick('discover')}
            className={`${
              view === 'discover' ? 'text-[#4a586e]' : 'text-[#7e918b]'
            } hover:text-[#4a586e] transition-colors underline-offset-[12px] ${
              view === 'discover' ? 'underline decoration-1' : ''
            }`}
          >
            Discover
          </button>
          <button
            onClick={() => handleNavClick('profile')}
            className={`${
              view === 'profile' ? 'text-[#4a586e]' : 'text-[#7e918b]'
            } hover:text-[#4a586e] transition-colors underline-offset-[12px] ${
              view === 'profile' ? 'underline decoration-1' : ''
            }`}
          >
            Profile
          </button>
        </nav>

        <div className="flex items-center gap-4">
          {currentUser ? (
            <div className="flex items-center gap-4">
              <Button variant="outline" className="hidden sm:block !py-2 !px-4" onClick={openCreateModal}>
                Post
              </Button>

              <button
                onClick={handleLogout}
                className="w-10 h-10 border border-[#4a586e] flex items-center justify-center hover:bg-[#4a586e] hover:text-[#f3e9d2] transition-all group text-[#4a586e]"
              >
                <span className="text-[10px] font-bold group-hover:hidden uppercase tracking-tighter">ME</span>
                <svg
                  className="hidden group-hover:block w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                  />
                </svg>
              </button>
            </div>
          ) : (
            <Button onClick={openAuthModal} variant="primary" className="!py-2 !px-4">
              Sign In
            </Button>
          )}
        </div>
      </div>
    </header>
  );
};
