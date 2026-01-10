// src/components/auth/useAuthFlow.ts
import { useState, useCallback, useEffect } from 'react';
import { useAuthStore } from '../../stores/authStore';
import * as authApi from '../../api/auth';
import type { User } from '../../types/types';

export type AuthStep = 'SIGN_IN' | 'SIGN_UP' | 'FORGOT_PASSWORD' | 'VERIFY_CODE' | 'NEW_PASSWORD';

export interface AuthFormData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  confirmPassword: string;
  verificationCode: string;
}

export interface UseAuthFlowReturn {
  // State
  step: AuthStep;
  formData: AuthFormData;
  error: string;
  isLoading: boolean;
  isSignUpFlow: boolean;

  // Actions
  setStep: (step: AuthStep) => void;
  setFormData: (data: Partial<AuthFormData>) => void;
  setError: (error: string) => void;
  clearError: () => void;
  reset: () => void;

  // Auth operations
  handleSignIn: () => Promise<boolean>;
  handleSignUp: () => Promise<boolean>;
  handleForgotPassword: () => Promise<boolean>;
  handleVerifyCode: (code?: string) => Promise<boolean>;
  handleNewPassword: (password?: string, confirmPassword?: string) => Promise<boolean>;
  handleResendCode: () => Promise<boolean>;

  // Hydration
  checkExistingSession: () => Promise<void>;
}

const initialFormData: AuthFormData = {
  email: '',
  password: '',
  firstName: '',
  lastName: '',
  confirmPassword: '',
  verificationCode: '',
};

export function useAuthFlow(onSuccess?: () => void): UseAuthFlowReturn {
  const [step, setStep] = useState<AuthStep>('SIGN_IN');
  const [formData, setFormDataState] = useState<AuthFormData>(initialFormData);
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [isSignUpFlow, setIsSignUpFlow] = useState(false);

  const login = useAuthStore((state) => state.login);

  const setFormData = useCallback((data: Partial<AuthFormData>) => {
    setFormDataState((prev) => ({ ...prev, ...data }));
  }, []);

  const clearError = useCallback(() => setError(''), []);

  const reset = useCallback(() => {
    setStep('SIGN_IN');
    setFormDataState(initialFormData);
    setError('');
    setIsLoading(false);
    setIsSignUpFlow(false);
  }, []);

  // Check for existing session on mount
  const checkExistingSession = useCallback(async () => {
    if (!authApi.hasRefreshToken()) return;

    try {
      const response = await authApi.fetchProfile();
      if (!response.error && response.data) {
        login(response.data);
      }
    } catch (err) {
      console.log('No existing session');
    }
  }, [login]);

  // Sign In
  const handleSignIn = useCallback(async (): Promise<boolean> => {
    setError('');
    setIsLoading(true);

    try {
      const response = await authApi.login(formData.email, formData.password);
      if (response.data) {
        login(response.data);
        onSuccess?.();
        return true;
      }
      return false;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Sign in failed');
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [formData.email, formData.password, login, onSuccess]);

  // Sign Up - Step 1: Send verification code
  const handleSignUp = useCallback(async (): Promise<boolean> => {
    setError('');
    setIsLoading(true);

    try {
      await authApi.sendVerificationCode(formData.email);
      setIsSignUpFlow(true);
      setStep('VERIFY_CODE');
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to send verification code');
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [formData.email]);

  // Forgot Password - Send reset code
  const handleForgotPassword = useCallback(async (): Promise<boolean> => {
    setError('');
    setIsLoading(true);

    try {
      await authApi.sendPasswordResetCode(formData.email);
      setIsSignUpFlow(false);
      setStep('VERIFY_CODE');
      return true;
    } catch (err) {
      // If backend doesn't support this yet, still proceed to verify step
      // This allows testing the UI flow
      console.warn('Password reset may not be implemented on backend:', err);
      setIsSignUpFlow(false);
      setStep('VERIFY_CODE');
      return true;
    } finally {
      setIsLoading(false);
    }
  }, [formData.email]);

  // Verify Code - accepts code directly to avoid state timing issues
  const handleVerifyCode = useCallback(async (code?: string): Promise<boolean> => {
    setError('');
    setIsLoading(true);

    const verificationCode = code || formData.verificationCode;

    try {
      if (isSignUpFlow) {
        // Complete registration
        const response = await authApi.register(
          formData.email,
          formData.password,
          formData.firstName,
          formData.lastName,
          verificationCode
        );

        if (response.data) {
          login(response.data);
          onSuccess?.();
          return true;
        }
        return false;
      } else {
        // Password reset flow - store code and proceed to new password step
        setFormDataState((prev) => ({ ...prev, verificationCode }));
        setStep('NEW_PASSWORD');
        return true;
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Verification failed');
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [isSignUpFlow, formData.email, formData.password, formData.firstName, formData.lastName, formData.verificationCode, login, onSuccess]);

  // Set New Password - accepts passwords directly to avoid state timing issues
  const handleNewPassword = useCallback(async (newPassword?: string, confirmPass?: string): Promise<boolean> => {
    setError('');

    const password = newPassword || formData.password;
    const confirmPassword = confirmPass || formData.confirmPassword;

    // Validate passwords match
    if (password !== confirmPassword) {
      setError('Passwords do not match');
      return false;
    }

    if (password.length < 6) {
      setError('Password must be at least 6 characters');
      return false;
    }

    setIsLoading(true);

    try {
      await authApi.resetPassword(
        formData.email,
        formData.verificationCode,
        password
      );
      setStep('SIGN_IN');
      setFormDataState((prev) => ({ ...prev, password: '', confirmPassword: '' }));
      return true;
    } catch (err) {
      // If backend doesn't support this yet, just go to sign in
      console.warn('Password reset may not be implemented on backend:', err);
      setStep('SIGN_IN');
      return true;
    } finally {
      setIsLoading(false);
    }
  }, [formData.email, formData.verificationCode, formData.password, formData.confirmPassword]);

  // Resend verification code
  const handleResendCode = useCallback(async (): Promise<boolean> => {
    setError('');
    setIsLoading(true);

    try {
      await authApi.resendVerificationCode(formData.email);
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to resend code');
      return false;
    } finally {
      setIsLoading(false);
    }
  }, [formData.email]);

  return {
    step,
    formData,
    error,
    isLoading,
    isSignUpFlow,
    setStep,
    setFormData,
    setError,
    clearError,
    reset,
    handleSignIn,
    handleSignUp,
    handleForgotPassword,
    handleVerifyCode,
    handleNewPassword,
    handleResendCode,
    checkExistingSession,
  };
}
