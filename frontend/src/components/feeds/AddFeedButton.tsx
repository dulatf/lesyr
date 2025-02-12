'use client';

import { useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { feedsApi } from '@/lib/api/feeds';

export function AddFeedButton() {
  const [isOpen, setIsOpen] = useState(false);
  const [url, setUrl] = useState('');
  const queryClient = useQueryClient();

  const addFeedMutation = useMutation({
    mutationFn: (url: string) => feedsApi.create(url),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['feeds'] });
      setUrl('');
      setIsOpen(false);
    },
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (url.trim()) {
      addFeedMutation.mutate(url);
    }
  };

  return (
    <div className="mb-8">
      {!isOpen ? (
        <button
          onClick={() => setIsOpen(true)}
          className="rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700"
        >
          Add New Feed
        </button>
      ) : (
        <form onSubmit={handleSubmit} className="flex max-w-xl gap-2">
          <input
            type="url"
            value={url}
            onChange={(e) => setUrl(e.target.value)}
            placeholder="Enter RSS feed URL"
            className="flex-1 rounded-lg bg-slate-100 px-4 py-2 text-slate-900 placeholder-slate-500 dark:bg-slate-800 dark:text-white dark:placeholder-slate-400"
            required
          />
          <button
            type="submit"
            disabled={addFeedMutation.isPending}
            className="rounded-lg bg-blue-600 px-4 py-2 text-white transition-colors hover:bg-blue-700 disabled:opacity-50"
          >
            {addFeedMutation.isPending ? 'Adding...' : 'Add Feed'}
          </button>
          <button
            type="button"
            onClick={() => setIsOpen(false)}
            className="rounded-lg bg-slate-200 px-4 py-2 text-slate-700 transition-colors hover:bg-slate-300 dark:bg-slate-700 dark:text-slate-200 dark:hover:bg-slate-600"
          >
            Cancel
          </button>
        </form>
      )}
    </div>
  );
}