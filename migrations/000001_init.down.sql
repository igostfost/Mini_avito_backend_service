DROP TRIGGER IF EXISTS update_banner_updated_at_trigger ON banners;


ALTER TABLE banner_tags DROP CONSTRAINT IF EXISTS banner_tags_banner_id_fkey;
ALTER TABLE banner_tags DROP CONSTRAINT IF EXISTS banner_tags_tag_id_fkey;


DROP TABLE IF EXISTS banner_tags;
DROP TABLE IF EXISTS banners;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tags;


DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;