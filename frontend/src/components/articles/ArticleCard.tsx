'use client';

interface Article {
    id: string;
    feedId: string;
    title: string;
    url: string;
    content: string;
    publishedAt: string;
    createdAt: string;
}

interface ArticleCardProps {
    article: Article;
}

export function ArticleCard({ article }: ArticleCardProps) {
    return (
        <div className="group rounded-lg bg-slate-100 p-6 transition-all hover:bg-slate-200 dark:bg-slate-800 dark:hover:bg-slate-700">
            <h2 className="mb-2 text-xl font-semibold text-slate-900 group-hover:text-blue-600 dark:text-white dark:group-hover:text-blue-400">
                <a
                    href={article.url}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="mt-4 inline-block text-blue-500 hover:underline"
                >
                    {article.title}
                </a>
            </h2>
            <p className="text-sm text-slate-600 dark:text-slate-300">
                {article.content.slice(0, 100)}...
            </p>
            <p className="mt-2 text-xs text-slate-500 dark:text-slate-400">
                {new Date(article.publishedAt).toLocaleDateString()}
            </p>
        </div>
    );
}