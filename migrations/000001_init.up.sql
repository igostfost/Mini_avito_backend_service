
CREATE TABLE tags (
                      tag_id SERIAL PRIMARY KEY
);

INSERT INTO tags (tag_id)
SELECT generate_series(1, 998);

CREATE TABLE banners (
                         banner_id SERIAL PRIMARY KEY,
                         feature_id INTEGER NOT NULL,
                         title VARCHAR(255) NOT NULL,
                         text TEXT NOT NULL,
                         url VARCHAR(255) NOT NULL,
                         is_active BOOLEAN NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);

CREATE TABLE banner_tags (
                             banner_id INTEGER,
                             tag_id INTEGER,
                             PRIMARY KEY (banner_id, tag_id),
                             FOREIGN KEY (banner_id) REFERENCES banners(banner_id) ON DELETE CASCADE,
                             FOREIGN KEY (tag_id) REFERENCES tags(tag_id) ON DELETE CASCADE
);


CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       is_admin BOOLEAN NOT NULL DEFAULT FALSE
);


CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;


CREATE TRIGGER update_banner_updated_at_trigger
    BEFORE UPDATE ON banners
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
