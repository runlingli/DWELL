// src/api/posts.ts
import axios, { AxiosError } from 'axios';
import type { Listing } from '../types/types';

const BROKER_URL = import.meta.env.VITE_BROKER_URL;

const postsClient = axios.create({
  baseURL: BROKER_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Response types
export interface PostsResponse {
  error: boolean;
  message?: string;
  data?: Listing[];
}

export interface PostResponse {
  error: boolean;
  message?: string;
  data?: Listing;
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
 * Fetch all posts/listings (RESTful GET /posts)
 */
export async function fetchPosts(): Promise<PostsResponse> {
  try {
    const res = await postsClient.get<PostsResponse>('/posts');
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Create a new post/listing (RESTful POST /posts)
 */
export async function createPost(listing: Omit<Listing, 'id' | 'createdAt' | 'author'>, authorId: number): Promise<PostResponse> {
  try {
    const res = await postsClient.post<PostResponse>('/posts', {
      title: listing.title,
      price: listing.price,
      location: listing.location,
      neighborhood: listing.neighborhood,
      lat: listing.coordinates.lat,
      lng: listing.coordinates.lng,
      radius: listing.radius,
      type: listing.type,
      imageUrl: listing.imageUrl,
      additionalImages: listing.additionalImages || [],
      description: listing.description,
      bedrooms: listing.bedrooms,
      bathrooms: listing.bathrooms,
      availableFrom: listing.availableFrom,
      availableTo: listing.availableTo,
      authorId: authorId,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Update an existing post/listing (RESTful PUT /posts/{id})
 */
export async function updatePost(listing: Listing, authorId: number): Promise<PostResponse> {
  try {
    const postId = typeof listing.id === 'number' ? listing.id : parseInt(listing.id);
    const res = await postsClient.put<PostResponse>(`/posts/${postId}`, {
      title: listing.title,
      price: listing.price,
      location: listing.location,
      neighborhood: listing.neighborhood,
      lat: listing.coordinates.lat,
      lng: listing.coordinates.lng,
      radius: listing.radius,
      type: listing.type,
      imageUrl: listing.imageUrl,
      additionalImages: listing.additionalImages || [],
      description: listing.description,
      bedrooms: listing.bedrooms,
      bathrooms: listing.bathrooms,
      availableFrom: listing.availableFrom,
      availableTo: listing.availableTo,
      authorId: authorId,
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}

/**
 * Delete a post/listing (RESTful DELETE /posts/{id})
 */
export async function deletePost(id: string, authorId: number): Promise<PostResponse> {
  try {
    const postId = typeof id === 'number' ? id : parseInt(id);
    const res = await postsClient.delete<PostResponse>(`/posts/${postId}`, {
      data: { authorId: authorId },
    });
    return res.data;
  } catch (err) {
    const error = extractError(err);
    return { error: true, message: error.message };
  }
}
