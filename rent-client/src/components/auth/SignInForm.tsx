// src/components/auth/SignInForm.tsx
import React from 'react';
import { Input, Button } from '../../UI';
import { getGoogleLoginUrl } from '../../api/auth';
import googleLoginIcon from '../../assets/google_icon.svg';

interface SignInFormProps {
  email: string;
  password: string;
  isLoading: boolean;
  onEmailChange: (email: string) => void;
  onPasswordChange: (password: string) => void;
  onSubmit: () => void;
  onForgotPassword: () => void;
  onSwitchToSignUp: () => void;
}

export const SignInForm: React.FC<SignInFormProps> = ({
  email,
  password,
  isLoading,
  onEmailChange,
  onPasswordChange,
  onSubmit,
  onForgotPassword,
  onSwitchToSignUp,
}) => {
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit();
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="space-y-6">
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
        />

        <Button type="submit" className="w-full mt-4" disabled={isLoading}>
          {isLoading ? 'SIGNING IN...' : 'CONTINUE'}
        </Button>
      </form>

      <div className="text-right mt-2">
        <button
          type="button"
          onClick={onForgotPassword}
          className="text-[10px] uppercase tracking-widest text-[#7e918b] hover:text-[#4a586e] transition-colors font-bold border-b border-transparent hover:border-[#4a586e]"
        >
          Forgot Password?
        </button>
      </div>

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
        New to Dwell?
        <button
          type="button"
          onClick={onSwitchToSignUp}
          className="ml-2 text-[#4a586e] hover:underline font-bold"
        >
          Sign up
        </button>
      </p>
    </div>
  );
};
