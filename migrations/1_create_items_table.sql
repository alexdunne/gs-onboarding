CREATE TABLE IF NOT EXISTS items (
    id INT PRIMARY KEY,
    type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL,
    url TEXT NOT NULL,
    score INT DEFAULT 0,
    title TEXT NOT NULL,
    created_by VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL
);
