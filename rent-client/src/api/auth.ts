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
 * Login with email and password (RESTful: POST /auth/login)
 */
export async function login(email: string, password: string): Promise<AuthResponse> {
  try {
    console.log('Attempting login for:', email);
    const res = await authClient.post<AuthResponse>('/auth/login', {
      email,
      password,
    });
    console.log('Login response:', res.data);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    console.error('Login error:', error);
    throw new Error(error.message);
  }
}

/**
 * Send verification code to email (RESTful: POST /auth/verify-email)
 */
export async function sendVerificationCode(email: string): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/auth/verify-email', {
      email,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Register a new user with verification code (RESTful: POST /auth/register)
 */
export async function register(
  email: string,
  password: string,
  firstName: string,
  lastName: string,
  verificationCode: string
): Promise<AuthResponse> {
  try {
    console.log('Attempting registration for:', email);
    const res = await authClient.post<AuthResponse>('/auth/register', {
      email,
      password,
      first_name: firstName,
      last_name: lastName,
      verification_code: verificationCode,
    });
    console.log('Registration response:', res.data);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    console.error('Registration error:', error);
    throw new Error(error.message);
  }
}

/**
 * Fetch current user profile (RESTful: GET /auth/profile)
 */
export async function fetchProfile(): Promise<AuthResponse> {
  try {
    console.log('Fetching user profile...');
    const res = await authClient.get<AuthResponse>('/auth/profile');
    console.log('Profile response:', res.data);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    console.error('Profile fetch error:', error);
    throw new Error(error.message);
  }
}

/**
 * Send password reset code to email (RESTful: POST /auth/forgot-password)
 */
export async function sendPasswordResetCode(email: string): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/auth/forgot-password', {
      email,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    throw new Error(error.message);
  }
}

/**
 * Reset password with verification code (RESTful: POST /auth/reset-password)
 */
export async function resetPassword(
  email: string,
  verificationCode: string,
  newPassword: string
): Promise<AuthResponse> {
  try {
    const res = await authClient.post<AuthResponse>('/auth/reset-password', {
      email,
      verification_code: verificationCode,
      new_password: newPassword,
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
 * Logout - calls backend to clear cookies (RESTful: POST /auth/logout)
 */
export async function logout(): Promise<AuthResponse> {
  try {
    console.log('Logging out...');
    const res = await authClient.post<AuthResponse>('/auth/logout');
    console.log('Logout response:', res.data);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    console.error('Logout error:', error);
    throw new Error(error.message);
  }
}

/**
 * Check if user has a refresh token (basic login check)
 */
export function hasRefreshToken(): boolean {
  return document.cookie.includes('refresh_token');
}
