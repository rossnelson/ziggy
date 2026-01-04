import { writable } from 'svelte/store';
import { ziggyState, syncCooldownTimestamp, type ZiggyState } from './store';

const API_BASE = 'http://localhost:8080';

let eventSource: EventSource | null = null;

// Chat stores for SSE updates
export const chatMessages = writable<ChatMessage[]>([]);
export const chatLoading = writable(false);
export const mysteryStatus = writable<MysteryStatus | null>(null);

// Config store
export const aiEnabled = writable<boolean | null>(null);

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
    syncCooldownTimestamp();
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
  const result = await fetchApi<ZiggyState>('/api/state');
  syncStateFromApi(result);
  return result;
}

export async function sendFeed(): Promise<ApiResponse<ZiggyState>> {
  const result = await fetchApi<ZiggyState>('/api/signal/feed', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendPlay(): Promise<ApiResponse<ZiggyState>> {
  const result = await fetchApi<ZiggyState>('/api/signal/play', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendPet(): Promise<ApiResponse<ZiggyState>> {
  const result = await fetchApi<ZiggyState>('/api/signal/pet', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function sendWake(): Promise<ApiResponse<ZiggyState>> {
  const result = await fetchApi<ZiggyState>('/api/signal/wake', { method: 'POST' });
  syncStateFromApi(result);
  return result;
}

export async function healthCheck(): Promise<boolean> {
  const result = await fetchApi('/api/health');
  return result.success;
}

interface ConfigResponse {
  aiEnabled: boolean;
}

export async function getConfig(): Promise<void> {
  const result = await fetchApi<ConfigResponse>('/api/config');
  if (result.success && result.data) {
    aiEnabled.set(result.data.aiEnabled);
  }
}

interface ChatHistoryResponse {
  messages: ChatMessage[];
  mysteryStatus?: MysteryStatus;
}

interface SSEEvent {
  type: 'state' | 'chat';
  data: ZiggyState | ChatHistoryResponse;
}

export function startSSE() {
  if (eventSource) return;

  eventSource = new EventSource(`${API_BASE}/api/events`);

  eventSource.onmessage = (event) => {
    try {
      const parsed: SSEEvent = JSON.parse(event.data);

      if (parsed.type === 'state') {
        ziggyState.set(parsed.data as ZiggyState);
        syncCooldownTimestamp();
      } else if (parsed.type === 'chat') {
        const chatData = parsed.data as ChatHistoryResponse;
        chatMessages.set(chatData.messages ?? []);
        mysteryStatus.set(chatData.mysteryStatus ?? null);
      }
    } catch (err) {
      console.error('SSE parse error:', err);
    }
  };

  eventSource.onerror = () => {
    console.error('SSE connection error, reconnecting...');
    stopSSE();
    setTimeout(startSSE, 2000);
  };
}

export function stopSSE() {
  if (eventSource) {
    eventSource.close();
    eventSource = null;
  }
}

// Chat types
export interface ChatMessage {
  id: string;
  role: 'user' | 'ziggy';
  content: string;
  timestamp: string;
}

export interface ChatHistory {
  messages: ChatMessage[];
  activeMystery?: Mystery;
  mysteryProgress: number;
}

export interface Mystery {
  id: string;
  title: string;
  description: string;
  track: string;
  concept?: string;
  hints: string[];
  summary?: string;
  docsUrl?: string;
}

export interface MysteryStatus {
  active: boolean;
  mystery?: Mystery;
  progress: number;
  hintsGiven: string[];
  totalHints: number;
}

// Chat API functions
export async function getChatHistory(): Promise<ApiResponse<ChatHistory>> {
  return fetchApi<ChatHistory>('/api/chat/history');
}

export async function sendChatMessage(content: string): Promise<ApiResponse<ChatHistory>> {
  return fetchApi<ChatHistory>('/api/chat/message', {
    method: 'POST',
    body: JSON.stringify({ content }),
  });
}

export async function getMysteryStatus(): Promise<ApiResponse<MysteryStatus>> {
  return fetchApi<MysteryStatus>('/api/chat/mystery');
}

export async function startMystery(mysteryId: string, track: string = 'fun'): Promise<ApiResponse<void>> {
  return fetchApi<void>('/api/chat/mystery/start', {
    method: 'POST',
    body: JSON.stringify({ mysteryId, track }),
  });
}

export async function getAvailableMysteries(track: string = 'fun'): Promise<ApiResponse<Mystery[]>> {
  return fetchApi<Mystery[]>(`/api/chat/mysteries?track=${track}`);
}
