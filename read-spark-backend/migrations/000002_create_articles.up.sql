CREATE TABLE articles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(500) NOT NULL,
    summary TEXT,
    content TEXT NOT NULL,
    translation TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('news', 'fiction', 'exam', 'graded')),
    difficulty VARCHAR(10) NOT NULL CHECK (difficulty IN ('A1', 'A2', 'B1', 'B2', 'C1', 'C2')),
    word_count INT DEFAULT 0,
    audio_url VARCHAR(500),
    cover_image VARCHAR(500),
    is_premium BOOLEAN DEFAULT true,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_difficulty ON articles(difficulty);
CREATE INDEX idx_articles_published ON articles(published_at);

ALTER TABLE articles ADD COLUMN search_vector tsvector
    GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(summary, '')), 'B') ||
        setweight(to_tsvector('english', coalesce(content, '')), 'C')
    ) STORED;

CREATE INDEX idx_articles_search ON articles USING GIN(search_vector);
