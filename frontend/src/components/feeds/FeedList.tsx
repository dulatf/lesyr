import { Feed } from '@/types';
import { FeedCard } from './FeedCard';

interface FeedListProps {
  feeds: Feed[];
}

export function FeedList({ feeds }: FeedListProps) {
  if (feeds.length === 0) {
    return (
      <div className="rounded-lg bg-slate-100 p-6 dark:bg-slate-800">
        <p className="text-slate-600 dark:text-slate-400">
          No feeds yet. Add your first RSS feed to get started!
        </p>
      </div>
    )
  }

  return (
    <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
      {feeds.map((feed) => (
        <FeedCard key={feed.id} feed={feed} />
      ))}
    </div>
  )
}
