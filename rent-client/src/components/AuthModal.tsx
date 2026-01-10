// src/components/AuthModal.tsx
import React, { useEffect } from 'react';
import { Modal } from './UI';
import {
  SignInForm,
  SignUpForm,
  ForgotPasswordForm,
  VerifyCodeForm,
  NewPasswordForm,
  useAuthFlow,
} from './auth';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const TITLES = {
  SIGN_IN: 'SIGN IN',
  SIGN_UP: 'JOIN DWELL',
  FORGOT_PASSWORD: 'RESET ACCESS',
  VERIFY_CODE: 'VERIFY IDENTITY',
  NEW_PASSWORD: 'NEW PASSWORD',
} as const;

export const AuthModal: React.FC<AuthModalProps> = ({ isOpen, onClose }) => {
  const {
    step,
    formData,
    error,
    isLoading,
    isSignUpFlow,
    setStep,
    setFormData,
    clearError,
    reset,
    handleSignIn,
    handleSignUp,
    handleForgotPassword,
    handleVerifyCode,
    handleNewPassword,
    handleResendCode,
    checkExistingSession,
  } = useAuthFlow(onClose);

  // Check for existing session on mount
  useEffect(() => {
    checkExistingSession();
  }, [checkExistingSession]);

  // Clear error when switching steps
  useEffect(() => {
    clearError();
  }, [step, clearError]);

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleVerifySubmit = (code: string) => {
    handleVerifyCode(code);
  };

  const handleNewPasswordSubmit = (password: string, confirmPassword: string) => {
    handleNewPassword(password, confirmPassword);
  };

  return (
    <Modal isOpen={isOpen} onClose={handleClose} title={TITLES[step]}>
      <div>
        {/* Error Display */}
        {error && (
          <p className="text-[#f47979] text-[10px] uppercase tracking-widest mb-4 p-3 bg-[#f47979]/10 border border-[#f47979]/20">
            {error}
          </p>
        )}

        {step === 'SIGN_IN' && (
          <SignInForm
            email={formData.email}
            password={formData.password}
            isLoading={isLoading}
            onEmailChange={(email) => setFormData({ email })}
            onPasswordChange={(password) => setFormData({ password })}
            onSubmit={handleSignIn}
            onForgotPassword={() => setStep('FORGOT_PASSWORD')}
            onSwitchToSignUp={() => setStep('SIGN_UP')}
          />
        )}

        {step === 'SIGN_UP' && (
          <SignUpForm
            firstName={formData.firstName}
            lastName={formData.lastName}
            email={formData.email}
            password={formData.password}
            isLoading={isLoading}
            onFirstNameChange={(firstName) => setFormData({ firstName })}
            onLastNameChange={(lastName) => setFormData({ lastName })}
            onEmailChange={(email) => setFormData({ email })}
            onPasswordChange={(password) => setFormData({ password })}
            onSubmit={handleSignUp}
            onSwitchToSignIn={() => setStep('SIGN_IN')}
          />
        )}

        {step === 'FORGOT_PASSWORD' && (
          <ForgotPasswordForm
            email={formData.email}
            isLoading={isLoading}
            onEmailChange={(email) => setFormData({ email })}
            onSubmit={handleForgotPassword}
            onBack={() => setStep('SIGN_IN')}
          />
        )}

        {step === 'VERIFY_CODE' && (
          <VerifyCodeForm
            email={formData.email}
            isLoading={isLoading}
            isSignUpFlow={isSignUpFlow}
            onSubmit={handleVerifySubmit}
            onResend={handleResendCode}
            onBack={() => setStep('SIGN_IN')}
          />
        )}

        {step === 'NEW_PASSWORD' && (
          <NewPasswordForm
            isLoading={isLoading}
            onSubmit={handleNewPasswordSubmit}
            onBack={() => setStep('SIGN_IN')}
          />
        )}
      </div>
    </Modal>
  );
};
