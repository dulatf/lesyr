'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { useFeeds } from '@/lib/hooks/useFeeds';
import { FeedList } from '@/components/feeds/FeedList';
import { AddFeedButton } from '@/components/feeds/AddFeedButton';
import { UserAvatar } from '@/components/auth/UserAvatar';
import { useAuth } from '@/lib/auth/auth-context';

export default function FeedsPage() {
  const router = useRouter();
  const { user, loading } = useAuth();
  const { data, isLoading, isError } = useFeeds();
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (!loading && !user) {
      // router.push('/');
    }
  }, [loading, user, router]);

  // Don't render anything during SSR
  if (!isMounted) {
    return null;
  }

  if (loading || isLoading) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-white dark:bg-slate-900">
        <div className="h-12 w-12 animate-spin rounded-full border-4 border-slate-200 border-t-blue-500 dark:border-slate-600"></div>
      </div>
    );
  }

  if (!user) {
    console.log("No user")
    return null;
  }

  if (isError) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-white dark:bg-slate-900">
        <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
          <p className="text-red-600 dark:text-red-400">Failed to load feeds. Please try again later.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-white transition-colors dark:bg-slate-900">
      <div className="mx-auto max-w-6xl p-8">
        <div className="mb-8 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-slate-900 dark:text-white">Your Feeds</h1>
          <div className="flex items-center space-x-4">
            <AddFeedButton />
            <UserAvatar />
          </div>
        </div>
        <FeedList feeds={data?.feeds || []} />
      </div>
    </div>
  );
}