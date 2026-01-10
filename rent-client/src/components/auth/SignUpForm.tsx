// src/components/auth/SignUpForm.tsx
import React from 'react';
import { Input, Button } from '../../UI';
import { getGoogleLoginUrl } from '../../api/auth';
import googleLoginIcon from '../../assets/google_icon.svg';

interface SignUpFormProps {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  isLoading: boolean;
  onFirstNameChange: (value: string) => void;
  onLastNameChange: (value: string) => void;
  onEmailChange: (value: string) => void;
  onPasswordChange: (value: string) => void;
  onSubmit: () => void;
  onSwitchToSignIn: () => void;
}

export const SignUpForm: React.FC<SignUpFormProps> = ({
  firstName,
  lastName,
  email,
  password,
  isLoading,
  onFirstNameChange,
  onLastNameChange,
  onEmailChange,
  onPasswordChange,
  onSubmit,
  onSwitchToSignIn,
}) => {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="space-y-6">
        <Input
          placeholder="FIRST NAME"
          type="text"
          value={firstName}
          onChange={(e) => onFirstNameChange(e.target.value)}
          required
          disabled={isLoading}
        />

        <Input
          placeholder="LAST NAME"
          type="text"
          value={lastName}
          onChange={(e) => onLastNameChange(e.target.value)}
          required
          disabled={isLoading}
        />

        <Input
          placeholder="EMAIL ADDRESS"
          type="email"
          value={email}
          onChange={(e) => onEmailChange(e.target.value)}
          required
          disabled={isLoading}
        />

        <Input
          placeholder="PASSWORD"
          type="password"
          value={password}
          onChange={(e) => onPasswordChange(e.target.value)}
          required
          disabled={isLoading}
          minLength={6}
        />

        <Button type="submit" className="w-full mt-4" disabled={isLoading}>
          {isLoading ? 'SENDING CODE...' : 'VERIFY EMAIL'}
        </Button>
      </form>

      <div className="relative py-4">
        <div className="absolute inset-0 flex items-center">
          <div className="w-full border-t border-[#4a586e]/10"></div>
        </div>
        <div className="relative flex justify-center text-xs uppercase tracking-widest">
          <span className="bg-[#f3e9d2] px-4 text-[#7e918b]">or</span>
        </div>
      </div>

      <button
        type="button"
        onClick={() => (window.location.href = getGoogleLoginUrl())}
        className="w-full flex items-center justify-center gap-4 border border-[#4a586e] p-4 hover:bg-[#4a586e] hover:text-[#f3e9d2] transition-colors text-[#4a586e]"
        disabled={isLoading}
      >
        <img src={googleLoginIcon} alt="Google" />
        <span className="text-[10px] font-bold uppercase tracking-widest">
          Continue with Google
        </span>
      </button>

      <p className="text-center text-[10px] uppercase tracking-widest text-[#7e918b] mt-4">
        Already have an account?
        <button
          type="button"
          onClick={onSwitchToSignIn}
          className="ml-2 text-[#4a586e] hover:underline font-bold"
        >
          Sign in
        </button>
      </p>
    </div>
  );
};
