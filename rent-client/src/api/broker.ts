import axios from 'axios';

const BROKER_URL = import.meta.env.VITE_BROKER_URL;


const brokerClient = axios.create({
  baseURL: BROKER_URL,
  withCredentials: true, // 等价于 fetch 的 credentials: "include"
  headers: {
    'Content-Type': 'application/json',
  },
});

export async function postToBroker(payload: unknown) {
  const res = await brokerClient.post('/handle', payload);
  return res.data;
}

// src/api/broker.ts
export async function fetchProfile() {
  const res = await fetch(`${BROKER_URL}/resource/profile`, {
    method: "GET",
    credentials: "include", // 必须带 cookie
  });

  return res.json(); // 返回 { user: { email, firstName, lastName } }
}
