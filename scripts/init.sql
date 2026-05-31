CREATE TABLE IF NOT EXISTS todo (
    id INTEGER PRIMARY KEY,
    desc TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status in ('todo', 'in-progress', 'done')),
    deadline DATETIME,
    created_at DATETIME NOT NULL,
    updated_at DATETIME
)
