
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name text,
    email text,
    password text,
    created_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE collections(
    id SERIAL PRIMARY KEY,
    name text,
    owner_id int,
    created_at TIMESTAMP WITH TIME ZONE,
    pdf_file_url text,

    CONSTRAINT fk_owner
        FOREIGN KEY(owner_id)
            REFERENCES users(id) ON DELETE CASCADE
);