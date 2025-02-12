import { Feed } from '@/types'
import Link from 'next/link'

interface FeedCardProps {
  feed: Feed
}

export function FeedCard({ feed }: FeedCardProps) {
  return (
    <Link
      href={`/feeds/${feed.id}`}
      className="group rounded-lg bg-slate-100 p-6 transition-all hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700"
    >
      <h3 className="mb-2 text-lg font-medium text-slate-900 group-hover:text-blue-600 dark:text-white dark:group-hover:text-blue-400">
        {feed.title}
      </h3>
      <p className="text-sm text-slate-400 dark:text-slate-600">
        {feed.url}
      </p>
      {feed.description && (
        <p className="line-clamp-2 text-sm text-slate-600 dark:text-slate-400">
          {feed.description}
        </p>
      )}
    </Link>
  )
}