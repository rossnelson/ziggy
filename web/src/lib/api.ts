import {
  ziggyState,
  feed as mockFeed,
  play as mockPlay,
  pet as mockPet,
  type ZiggyState,
} from './store';
import { get } from 'svelte/store';

export const USE_MOCK = false;
const API_BASE = 'http://localhost:8080';
const POLL_INTERVAL_MS = 2000;

let pollTimer: ReturnType<typeof setInterval> | null = null;

export interface ApiResponse<T> {
  success: boolean;
  data?: T;
  error?: string;
}

interface ApiStateResponse {
  success: boolean;
  data?: ZiggyState;
  error?: string;
}

function syncStateFromApi(response: ApiStateResponse) {
  if (response.success && response.data) {
    ziggyState.set(response.data);
  }
}

async function fetchApi<T>(endpoint: string, options?: RequestInit): Promise<ApiResponse<T>> {
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
    return { success: true, data: data.data ?? data };
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
  const result = await fetchApi<ZiggyState>('/api/state');
  syncStateFromApi(result);
  return result;
}

export async function sendFeed(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockFeed();
    return { success: true, data: get(ziggyState) };
  }
  const result = await fetchApi<ZiggyState>('/api/signal/feed', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendPlay(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockPlay();
    return { success: true, data: get(ziggyState) };
  }
  const result = await fetchApi<ZiggyState>('/api/signal/play', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendPet(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    mockPet();
    return { success: true, data: get(ziggyState) };
  }
  const result = await fetchApi<ZiggyState>('/api/signal/pet', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendWake(): Promise<ApiResponse<ZiggyState>> {
  if (USE_MOCK) {
    const { wake } = await import('./store');
    wake();
    return { success: true, data: get(ziggyState) };
  }
  const result = await fetchApi<ZiggyState>('/api/signal/wake', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function healthCheck(): Promise<boolean> {
  if (USE_MOCK) return true;
  const result = await fetchApi('/api/health');
  return result.success;
}

export function startPolling() {
  if (USE_MOCK || pollTimer) return;
  getState();
  pollTimer = setInterval(() => getState(), POLL_INTERVAL_MS);
}

export function stopPolling() {
  if (pollTimer) {
    clearInterval(pollTimer);
    pollTimer = null;
  }
}
