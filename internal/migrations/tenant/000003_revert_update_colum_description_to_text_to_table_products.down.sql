UPDATE products 
SET description = LEFT(description, 200) 
WHERE LENGTH(description) > 200;

-- Una vez limpios, el ALTER TABLE funcionará perfectamente:
ALTER TABLE products MODIFY COLUMN description VARCHAR(200);