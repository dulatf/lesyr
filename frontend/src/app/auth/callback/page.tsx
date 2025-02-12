'use client';

import { useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

function AuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const token = searchParams.get('token');

    if (!token) {
      console.error('No token provided');
      router.push('/');
      return;
    }

    try {
      // Store the JWT token
      localStorage.setItem('token', token);

      // Parse user data if provided
      const userJson = searchParams.get('user');
      if (userJson) {
        const userData = JSON.parse(decodeURIComponent(userJson));
        // Store user data if needed
        localStorage.setItem('user', JSON.stringify(userData));
      }

      // Show success message and redirect
      router.push('/feeds');
    } catch (error) {
      console.error('Error processing authentication:', error);
      router.push('/');
    }
  }, [router, searchParams]);

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-b from-zinc-900 to-zinc-950">
      <div className="rounded-lg bg-zinc-800 p-8 text-center shadow-xl">
        <div className="mb-4 flex justify-center">
          <div className="h-12 w-12 animate-spin rounded-full border-4 border-zinc-600 border-t-blue-500"></div>
        </div>
        <h2 className="mb-2 text-xl font-semibold text-white">Almost there!</h2>
        <p className="text-zinc-400">Setting up your Lesyr account...</p>
      </div>
    </div>
  );
}

export default function AuthCallback() {
  return (
    <Suspense fallback={<div>Loading...</div>}>
      <AuthCallbackContent />
    </Suspense>
  );
}