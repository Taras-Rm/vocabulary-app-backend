
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name text,
    email text,
    password text,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE user_settings(
    id SERIAL PRIMARY KEY,
    user_id int UNIQUE,
    app_language varchar(2) NOT NULL DEFAULT 'en',

    FOREIGN KEY(user_id)
        REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE collections(
    id SERIAL PRIMARY KEY,
    name text,
    owner_id int,
    created_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_owner
        FOREIGN KEY(owner_id)
            REFERENCES users(id) ON DELETE CASCADE
);