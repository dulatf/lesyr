import api from './axios';
import type { Article } from '@/types';

export const articlesApi = {
  listByFeed: (feedId: string) => 
    api.get<{ articles: Article[] }>(`/feeds/${feedId}/articles`).then((res) => res.data),
  getUnread: () => 
    api.get<{ articles: Article[] }>('/articles/unread').then((res) => res.data),
  markAsRead: (id: string) => 
    api.post(`/articles/${id}/read`),
};