// src/components/auth/ForgotPasswordForm.tsx
import React from 'react';
import { Input, Button } from '../../UI';

interface ForgotPasswordFormProps {
  email: string;
  isLoading: boolean;
  onEmailChange: (email: string) => void;
  onSubmit: () => void;
  onBack: () => void;
}

export const ForgotPasswordForm: React.FC<ForgotPasswordFormProps> = ({
  email,
  isLoading,
  onEmailChange,
  onSubmit,
  onBack,
}) => {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  return (
    <div>
      <p className="text-[11px] text-[#7e918b] uppercase tracking-[0.2em] mb-8 leading-relaxed">
        Enter your email address...
      </p>

      <form onSubmit={handleSubmit} className="space-y-6">
        <Input
          placeholder="EMAIL ADDRESS"
          type="email"
          value={email}
          onChange={(e) => onEmailChange(e.target.value)}
          required
          disabled={isLoading}
        />

        <Button type="submit" className="w-full mt-4" disabled={isLoading}>
          {isLoading ? 'SENDING...' : 'SEND CODE'}
        </Button>
      </form>

      <button
        type="button"
        onClick={onBack}
        className="w-full text-center text-[8px] uppercase tracking-widest text-[#7e918b] hover:text-[#4a586e] font-bold transition-colors mt-6"
      >
        ← Back to Sign In
      </button>
    </div>
  );
};
