// src/api/favorites.ts
import axios, { AxiosError } from 'axios';

const BROKER_URL = import.meta.env.VITE_BROKER_URL;

const favoritesClient = axios.create({
  baseURL: BROKER_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response types
export interface FavoriteIDsResponse {
  error: boolean;
  message?: string;
  data?: string[];
}

export interface FavoriteResponse {
  error: boolean;
  message?: string;
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
 * Fetch favorite post IDs for a user
 */
export async function fetchFavoriteIds(userId: number): Promise<FavoriteIDsResponse> {
  try {
    const res = await favoritesClient.get<FavoriteIDsResponse>(`/favorites/${userId}/ids`);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Add a favorite
 */
export async function addFavorite(userId: number, postId: number): Promise<FavoriteResponse> {
  try {
    const res = await favoritesClient.post<FavoriteResponse>('/favorites', {
      userId,
      postId,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Remove a favorite
 */
export async function removeFavorite(userId: number, postId: number): Promise<FavoriteResponse> {
  try {
    const res = await favoritesClient.delete<FavoriteResponse>(`/favorites/${userId}/${postId}`);
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Sync localStorage favorites to backend
 */
export async function syncFavorites(userId: number, postIds: number[]): Promise<FavoriteIDsResponse> {
  try {
    const res = await favoritesClient.post<FavoriteIDsResponse>('/favorites/sync', {
      userId,
      postIds,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}
