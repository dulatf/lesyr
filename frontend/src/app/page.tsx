'use client';

import { GithubLoginButton } from '@/components/auth/GithubLoginButton';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth/auth-context';

export default function Home() {
  const router = useRouter();
  const { user, loading } = useAuth();

  useEffect(() => {
    if (!loading) {
      if (user) {
        router.push('/feeds');
      }
    }
  }, [loading, user, router]);

  if (loading || user) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-white dark:bg-slate-900">
        <div className="h-12 w-12 animate-spin rounded-full border-4 border-slate-200 border-t-blue-500 dark:border-slate-600"></div>
      </div>
    );
  }

  return (
    <main className="flex min-h-screen flex-col items-center justify-center bg-white dark:bg-slate-900 p-8">
      <div className="flex flex-col items-center gap-8 text-center">
        <h1 className="text-5xl font-bold tracking-tight text-zinc900">
          Welcome to Lesyr
        </h1>
        <p className="max-w-md text-xl text-zinc-400">
          Your personal RSS feed reader. Stay updated with your favorite content in one place.
        </p>
        <GithubLoginButton className="mt-4" />
      </div>
    </main>
  );
}