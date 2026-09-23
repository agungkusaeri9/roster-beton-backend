-- 1. Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(120) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

-- 2. Create products table (Clean, tanpa kolom SEO)
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(250) NOT NULL UNIQUE,
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    description TEXT,
    image VARCHAR(500) NOT NULL,
    price NUMERIC(12, 2) DEFAULT 0,
    dimension VARCHAR(100),
    weight NUMERIC(8, 2) DEFAULT 0,
    stock INT DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE,
    is_featured BOOLEAN DEFAULT FALSE,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);

-- 3. Create product galleries table
CREATE TABLE IF NOT EXISTS product_galleries (
    id BIGSERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_product_galleries_product_id ON product_galleries(product_id);

-- 4. Create polymorphic SEO table (Bisa untuk product, article/blog, category, page)
CREATE TABLE IF NOT EXISTS seos (
    id BIGSERIAL PRIMARY KEY,
    reference_type VARCHAR(50) NOT NULL, -- 'product', 'article', 'category', 'page'
    reference_id BIGINT NOT NULL,        -- ID dari entitas terkait
    meta_title VARCHAR(255),
    meta_description TEXT,
    meta_keywords TEXT,
    canonical_url VARCHAR(500),
    og_title VARCHAR(255),
    og_description TEXT,
    og_image VARCHAR(500),
    og_type VARCHAR(50) DEFAULT 'website',
    structured_data TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT uq_seo_reference UNIQUE (reference_type, reference_id)
);

CREATE INDEX IF NOT EXISTS idx_seos_reference ON seos(reference_type, reference_id);

-- 5. Create general galleries table
CREATE TABLE IF NOT EXISTS galleries (
    id BIGSERIAL PRIMARY KEY,
    src VARCHAR(500) NOT NULL,
    title VARCHAR(200) NOT NULL,
    category VARCHAR(100) NOT NULL,
    alt_text VARCHAR(255),
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_galleries_category ON galleries(category);
