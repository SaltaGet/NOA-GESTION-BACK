ALTER TABLE products
DROP COLUMN IF EXISTS primary_image,
DROP COLUMN IF EXISTS secondary_images,
DROP COLUMN IF EXISTS is_visible;

