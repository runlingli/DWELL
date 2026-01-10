// src/components/auth/NewPasswordForm.tsx
import React, { useState } from 'react';
import { Input, Button } from '../UI';

interface NewPasswordFormProps {
  isLoading: boolean;
  onSubmit: (password: string, confirmPassword: string) => void;
  onBack: () => void;
}

export const NewPasswordForm: React.FC<NewPasswordFormProps> = ({
  isLoading,
  onSubmit,
  onBack,
}) => {
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [validationError, setValidationError] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setValidationError('');

    if (password.length < 6) {
      setValidationError('Password must be at least 6 characters');
      return;
    }

    if (password !== confirmPassword) {
      setValidationError('Passwords do not match');
      return;
    }

    onSubmit(password, confirmPassword);
  };

  return (
    <div>
      <p className="text-[11px] text-[#7e918b] uppercase tracking-[0.2em] mb-8 leading-relaxed">
        Create a new password for your account. Make sure it's at least 6 characters long.
      </p>

      <form onSubmit={handleSubmit} className="space-y-6">
        <Input
          placeholder="NEW PASSWORD"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          disabled={isLoading}
          minLength={6}
        />

        <Input
          placeholder="CONFIRM NEW PASSWORD"
          type="password"
          value={confirmPassword}
          onChange={(e) => setConfirmPassword(e.target.value)}
          required
          disabled={isLoading}
          minLength={6}
        />

        {validationError && (
          <p className="text-[#f47979] text-[10px] uppercase tracking-widest">
            {validationError}
          </p>
        )}

        <Button type="submit" className="w-full mt-4" disabled={isLoading}>
          {isLoading ? 'UPDATING...' : 'UPDATE PASSWORD'}
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
