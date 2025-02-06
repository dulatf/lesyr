-- Insert test user
INSERT INTO users (
    id,
    email,
    provider,
    provider_user_id
) VALUES (
    '11111111-1111-1111-1111-111111111111',
    'test@example.com',
    'github',
    'test123'
);

-- Insert test feeds
INSERT INTO feeds (
    id,
    user_id,
    url,
    title,
    description,
    last_fetched_at
) VALUES
    (
        '22222222-2222-2222-2222-222222222222',
        '11111111-1111-1111-1111-111111111111',
        'https://news.ycombinator.com/rss',
        'Hacker News',
        'Hacker News RSS feed',
        CURRENT_TIMESTAMP
    ),
    (
        '33333333-3333-3333-3333-333333333333',
        '11111111-1111-1111-1111-111111111111',
        'https://lobste.rs/rss',
        'Lobsters',
        'Lobsters RSS feed',
        CURRENT_TIMESTAMP
    );