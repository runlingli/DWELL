// src/api/broker.ts
// Generic broker client for non-auth API calls
import axios from 'axios';

const BROKER_URL = import.meta.env.VITE_BROKER_URL;

const brokerClient = axios.create({
  baseURL: BROKER_URL,
  withCredentials: true,
  headers: {
    'Content-Type': 'application/json',
  },
});

/**
 * Generic POST to broker /handle endpoint
 * For auth-specific calls, prefer using the functions in ./auth.ts
 */
export async function postToBroker(payload: unknown) {
  const res = await brokerClient.post('/handle', payload);
  return res.data;
}
