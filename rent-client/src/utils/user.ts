// src/utils/user.ts
// User-related utility functions
import type { User } from '@/types';

/**
 * Get user ID as number
 */
export const getUserId = (user: User | null): number | undefined => {
  if (!user?.id) return undefined;
  return typeof user.id === 'number' ? user.id : parseInt(user.id, 10);
};

/**
 * Get display name from user
 */
export const getUserDisplayName = (user: User): string => {
  return `${user.first_name} ${user.last_name}`.trim();
};

/**
 * Get user initials for avatar placeholder
 */
export const getUserInitials = (user: User): string => {
  const first = user.first_name?.[0] || '';
  const last = user.last_name?.[0] || '';
  return `${first}${last}`.toUpperCase();
};
