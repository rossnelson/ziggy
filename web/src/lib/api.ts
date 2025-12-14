import {
  ziggyState,
  feed as mockFeed,
  play as mockPlay,
  pet as mockPet,
  type ZiggyState,
} from './store';
import { get } from 'svelte/store';

export const USE_MOCK = true;
const API_BASE = 'http://localhost:8080';

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  try {
    const response = await fetch(`${API_BASE}${endpoint}`, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
    });

    if (!response.ok) {
      return { success: false, error: `HTTP ${response.status}` };
    }

    const data = await response.json();
    return { success: true, data };
  } catch (err) {
    return {
      success: false,
      error: err instanceof Error ? err.message : 'Unknown error',
    };
  }
}

export async function getState(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    return { success: true, data: get(ziggyState) };
  }
  return fetchApi<ZiggyState>('/api/state');
}

export async function sendFeed(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockFeed();
    return { success: true, data: get(ziggyState) };
  }
  return fetchApi<ZiggyState>('/api/signal/feed', { method: 'POST' });
}

export async function sendPlay(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockPlay();
    return { success: true, data: get(ziggyState) };
  }
  return fetchApi<ZiggyState>('/api/signal/play', { method: 'POST' });
}

export async function sendPet(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockPet();
    return { success: true, data: get(ziggyState) };
  }
  return fetchApi<ZiggyState>('/api/signal/pet', { method: 'POST' });
}

export async function healthCheck(): Promise<boolean> {
  if (USE_MOCK) return true;
  const result = await fetchApi('/api/health');
  return result.success;
}
