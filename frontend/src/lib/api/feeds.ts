import api from './axios';
import type { Feed } from '@/types';

export const feedsApi = {
  list: () => api.get<{ feeds: Feed[] }>('/feeds').then((res) => res.data),
  create: (url: string) => api.post<Feed>('/feeds', { url }).then((res) => res.data),
  get: (id: string) => api.get<Feed>(`/feeds/${id}`).then((res) => res.data),
  delete: (id: string) => api.delete(`/feeds/${id}`),
};