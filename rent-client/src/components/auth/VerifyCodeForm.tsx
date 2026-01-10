// src/components/auth/VerifyCodeForm.tsx
import React, { useRef, useEffect, useState } from 'react';
import { Button } from '../UI';

interface VerifyCodeFormProps {
  email: string;
  isLoading: boolean;
  isSignUpFlow: boolean;
  onSubmit: (code: string) => void;
  onResend: () => void;
  onBack: () => void;
}

export const VerifyCodeForm: React.FC<VerifyCodeFormProps> = ({
  email,
  isLoading,
  isSignUpFlow,
  onSubmit,
  onResend,
  onBack,
}) => {
  const [otp, setOtp] = useState(['', '', '', '', '', '']);
  const [resendCooldown, setResendCooldown] = useState(0);
  const otpRefs = useRef<(HTMLInputElement | null)[]>([]);

  // Focus first input on mount
  useEffect(() => {
    const timer = setTimeout(() => {
      otpRefs.current[0]?.focus();
    }, 150);
    return () => clearTimeout(timer);
  }, []);

  // Resend cooldown timer
  useEffect(() => {
    if (resendCooldown > 0) {
      const timer = setTimeout(() => setResendCooldown((c) => c - 1), 1000);
      return () => clearTimeout(timer);
    }
  }, [resendCooldown]);

  const handleOtpChange = (value: string, index: number) => {
    if (!/^\d*$/.test(value)) return;

    const newOtp = [...otp];
    newOtp[index] = value.slice(-1);
    setOtp(newOtp);

    // Auto-advance to next input
    if (value && index < 5) {
      otpRefs.current[index + 1]?.focus();
    }

    // Auto-submit when all digits entered
    if (value && index === 5 && newOtp.every((d) => d)) {
      onSubmit(newOtp.join(''));
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent, index: number) => {
    if (e.key === 'Backspace' && !otp[index] && index > 0) {
      otpRefs.current[index - 1]?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const pastedData = e.clipboardData.getData('text').replace(/\D/g, '').slice(0, 6);
    if (pastedData.length === 6) {
      const newOtp = pastedData.split('');
      setOtp(newOtp);
      otpRefs.current[5]?.focus();
      onSubmit(pastedData);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const code = otp.join('');
    if (code.length === 6) {
      onSubmit(code);
    }
  };

  const handleResend = () => {
    if (resendCooldown > 0) return;
    onResend();
    setResendCooldown(60); // 60 second cooldown
  };

  return (
    <div>
      <form onSubmit={handleSubmit} className="space-y-8 py-2">
        <div className="text-center space-y-2">
          <p className="text-[12px] text-[#7e918b] uppercase tracking-[0.3em] leading-relaxed">
            Verification code sent to
          </p>
          <p className="text-[#4a586e] font-bold text-xs uppercase tracking-tighter">
            {email || 'your email address'}
          </p>
        </div>

        <div className="flex justify-between gap-2 max-w-xs mx-auto" onPaste={handlePaste}>
          {otp.map((digit, idx) => (
            <input
              key={idx}
              ref={(el) => {
                otpRefs.current[idx] = el;
              }}
              type="text"
              inputMode="numeric"
              maxLength={1}
              value={digit}
              onChange={(e) => handleOtpChange(e.target.value, idx)}
              onKeyDown={(e) => handleKeyDown(e, idx)}
              disabled={isLoading}
              className="w-10 h-14 md:w-11 md:h-16 bg-transparent border border-[#4a586e]/20 text-center text-2xl font-bold focus:border-[#4a586e] focus:outline-none transition-all text-[#4a586e] rounded-none disabled:opacity-50"
            />
          ))}
        </div>

        <p className="text-center text-[9px] text-[#7e918b] uppercase tracking-[0.2em]">
          Didn't receive code?{' '}
          <button
            type="button"
            onClick={handleResend}
            disabled={resendCooldown > 0 || isLoading}
            className="text-[#4a586e] font-bold hover:underline disabled:opacity-50 disabled:no-underline"
          >
            {resendCooldown > 0 ? `Resend in ${resendCooldown}s` : 'Resend'}
          </button>
        </p>

        <Button type="submit" className="w-full" disabled={isLoading || otp.some((d) => !d)}>
          {isLoading
            ? 'VERIFYING...'
            : isSignUpFlow
            ? 'FINISH REGISTRATION'
            : 'VERIFY & CONTINUE'}
        </Button>
      </form>

      <button
        type="button"
        onClick={onBack}
        className="w-full text-center text-[8px] uppercase tracking-widest text-[#7e918b] hover:text-[#4a586e] font-bold transition-colors mt-4"
      >
        ← Back to Sign In
      </button>
    </div>
  );
};
