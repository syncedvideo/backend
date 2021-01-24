CREATE TABLE IF NOT EXISTS sv_user (
    id UUID PRIMARY KEY,
    name VARCHAR,
    color VARCHAR,
    is_admin BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS sv_room (
    id UUID PRIMARY KEY,
    owner_user_id UUID REFERENCES sv_user(id),
    name VARCHAR
);

CREATE TABLE IF NOT EXISTS sv_room_user (
    room_id UUID NOT NULL REFERENCES sv_room(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES sv_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sv_queue_item (
    id UUID PRIMARY KEY,
    room_id UUID NOT NULL REFERENCES sv_room(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES sv_user(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sv_queue_item_vote (
    id UUID PRIMARY KEY,
    queue_item_id UUID NOT NULL REFERENCES sv_queue_item(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES sv_user(id) ON DELETE CASCADE
);
