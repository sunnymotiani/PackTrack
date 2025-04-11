-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    event_id UUID REFERENCES events(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    UNIQUE(event_id, name)
);

-- Items in the packing list
CREATE TABLE IF NOT EXISTS items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES categories(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    quantity INTEGER DEFAULT 1,
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    status TEXT CHECK (status IN ('to_pack', 'packed', 'delivered')) DEFAULT 'to_pack',
    notes TEXT,
    created_at TIMESTAMP DEFAULT now()
);
-- Track item status changes (for activity log / real-time dashboard)
CREATE TABLE IF NOT EXISTS item_status_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_id UUID REFERENCES items(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    old_status TEXT,
    new_status TEXT,
    changed_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS item_status_history;
-- +goose StatementEnd
