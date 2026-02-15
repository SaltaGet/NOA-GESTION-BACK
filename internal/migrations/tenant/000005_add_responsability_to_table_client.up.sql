ALTER TABLE clients
ADD COLUMN IF NOT EXISTS responsability_front_iva VARCHAR(255) DEFAULT 'consumidor_final' NULL;