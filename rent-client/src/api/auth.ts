// src/api/auth.ts
import axios, { AxiosError } from 'axios';
import type { User } from '../types/types';

const BROKER_URL = import.meta.env.VITE_BROKER_URL;

const authClient = axios.create({
  baseURL: BROKER_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response types
export interface AuthResponse {
  error: boolean;
  message: string;
  data?: User;
}

export interface ApiError {
  message: string;
  status?: number;
}

// Helper to extract error message
const extractError = (err: unknown): ApiError => {
  if (axios.isAxiosError(err)) {
    const axiosErr = err as AxiosError<{ message?: string }>;
    return {
      message: axiosErr.response?.data?.message || axiosErr.message,
      status: axiosErr.response?.status,
    };
  }
  if (err instanceof Error) {
    return { message: err.message };
  }
  return { message: String(err) };
};

/**
 * Login with email and password
 */
export async function login(email: string, password: string): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'auth',
      auth: { email, password },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Send verification code to email (for signup)
 */
export async function sendVerificationCode(email: string): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'verify',
      verify: { email },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Register a new user with verification code
 */
export async function register(
  email: string,
  password: string,
  firstName: string,
  lastName: string,
  verificationCode: string
): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'register',
      register: {
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        verification_code: verificationCode,
      },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Fetch current user profile (uses cookies for auth)
 */
export async function fetchProfile(): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'resource',
      resource: 'profile',
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Send password reset code to email
 * Note: Backend may not have this implemented yet
 */
export async function sendPasswordResetCode(email: string): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'forgot-password',
      forgot_password: { email },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Reset password with verification code
 * Note: Backend may not have this implemented yet
 */
export async function resetPassword(
  email: string,
  verificationCode: string,
  newPassword: string
): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/handle', {
      action: 'reset-password',
      reset_password: {
        email,
        verification_code: verificationCode,
        new_password: newPassword,
      },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Resend verification code
 */
export async function resendVerificationCode(email: string): Promise<AuthResponse> {
  // Uses same endpoint as sendVerificationCode
  return sendVerificationCode(email);
}

/**
 * Get Google OAuth login URL
 */
export function getGoogleLoginUrl(): string {
  return `${BROKER_URL}/oauth/google/login`;
}

/**
 * Check if user has a refresh token (basic login check)
 */
export function hasRefreshToken(): boolean {
  return document.cookie.includes('refresh_token');
}
