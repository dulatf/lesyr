import { useQuery } from '@tanstack/react-query';
import { feedsApi } from '@/lib/api/feeds';
import type { Feed } from '@/types';

export function useFeeds() {
  return useQuery<{ feeds: Feed[] }>({
    queryKey: ['feeds'],
    queryFn: () => feedsApi.list(),
  });
}