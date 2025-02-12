'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useAuth } from '@/lib/auth/auth-context';
import { ArticleCard } from '@/components/articles/ArticleCard';
import { articlesApi } from '@/lib/api/articles';
import { feedsApi } from '@/lib/api/feeds';
import { Feed } from '@/types';
import Link from 'next/link';

interface Article {
    id: string;
    feedId: string;
    title: string;
    url: string;
    content: string;
    publishedAt: string;
    createdAt: string;
}

export default function FeedArticlesPage() {
    const { feedId } = useParams();
    const router = useRouter();
    const { user, loading } = useAuth();
    const [articles, setArticles] = useState<Article[]>([]);
    const [feed, setFeed] = useState<Feed | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (!loading && !user) {
            router.push('/');
        }
    }, [loading, user, router]);

    useEffect(() => {
        async function fetchArticles() {
            try {
                const data = await articlesApi.listByFeed(feedId as string);
                const feedData = await feedsApi.get(feedId as string);
                setFeed(feedData);
                setArticles(data.articles || []);
            } catch (err: unknown) {
                if (err instanceof Error) {
                    setError(err.message);
                } else {
                    setError("An unknown error occurred.");
                }
            } finally {
                setIsLoading(false);
            }
        }
        if (feedId) {
            fetchArticles();
        }
    }, [feedId]);

    if (loading || isLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-white dark:bg-slate-900">
                <div className="h-12 w-12 animate-spin rounded-full border-4 border-slate-200 border-t-blue-500 dark:border-slate-600"></div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-white dark:bg-slate-900">
                <div className="rounded-lg bg-red-50 p-6 text-center dark:bg-red-900/20">
                    <p className="text-red-600 dark:text-red-400">{error}</p>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-white transition-colors dark:bg-slate-900">
            <div className="mx-auto max-w-6xl p-8">
                <h1 className="mb-6 text-3xl font-bold text-slate-900 dark:text-white">
                    {feed?.title || 'Feed Articles'}
                </h1>
                
                <Link
                    href="/feeds"
                    className="mb-4 inline-block text-blue-500 hover:underline dark:text-blue-400"
                >
                    &larr; Back to Feeds Overview
                </Link>
                {articles.length === 0 ? (
                    <p>No articles found for this feed.</p>
                ) : (
                    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
                        {articles.map((article) => (
                            <ArticleCard key={article.id} article={article} />
                        ))}
                    </div>
                )}
            </div>
        </div>
    );
}