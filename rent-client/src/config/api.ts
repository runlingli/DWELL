// src/config/api.ts
// API configuration and endpoints

export const API_CONFIG = {
  baseUrl: import.meta.env.VITE_API_URL || 'http://localhost:8080',
  timeout: 10000,
  retryAttempts: 3,
} as const;

export const API_ENDPOINTS = {
  // Auth
  auth: {
    login: '/auth/login',
    logout: '/auth/logout',
    register: '/auth/register',
    profile: '/auth/profile',
    verifyEmail: '/auth/verify-email',
    forgotPassword: '/auth/forgot-password',
    resetPassword: '/auth/reset-password',
    resendCode: '/auth/resend-code',
  },
  // Posts/Listings
  posts: {
    base: '/posts',
    byId: (id: string) => `/posts/${id}`,
    byAuthor: (authorId: number) => `/posts/author/${authorId}`,
  },
} as const;

// Token storage keys
export const STORAGE_KEYS = {
  accessToken: 'access_token',
  refreshToken: 'refresh_token',
  user: 'user',
  favorites: 'favorites',
} as const;
