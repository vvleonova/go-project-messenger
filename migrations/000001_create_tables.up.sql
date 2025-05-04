-- user
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    login VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(11) UNIQUE NOT NULL,
    birth_date TIMESTAMP NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_on TIMESTAMP NOT NULL
);
CREATE INDEX ON users using hash (login);
CREATE INDEX ON users using hash (phone);

-- chat
CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    created_on TIMESTAMP NOT NULL
);

-- chat + user
CREATE TABLE IF NOT EXISTS chat_users (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (chat_id, user_id)
);

-- message
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE NOT NULL,
	sender_id UUID NOT NULL,
	receiver_id UUID NOT NULL,
    text TEXT NOT NULL,
    send_at TIMESTAMP NOT NULL,
	read BOOL NOT NULL,
	is_attachment BOOL NOT NULL,
	is_pinned BOOL NOT NULL
);
-- CREATE INDEX ON users using hash (chat_id);
-- CREATE INDEX ON users using hash (sender_id);
-- CREATE INDEX ON users using hash (receiver_id);

-- refresh tokens
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY,
    token TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_on TIMESTAMP NOT NULL,
    revoked BOOL NOT NULL
);
 