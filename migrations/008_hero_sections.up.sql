DROP TABLE IF EXISTS hero_sections;

CREATE TABLE hero_sections (
    id SERIAL PRIMARY KEY,
    
    -- Main Banner (Left)
    main_banner VARCHAR(255) DEFAULT '',
    main_title VARCHAR(100) DEFAULT '',
    main_subtitle VARCHAR(150) DEFAULT '',
    
    -- Side Top Banner (Right Top)
    side_top_banner VARCHAR(255) DEFAULT '',
    side_top_title VARCHAR(100) DEFAULT '',
    side_top_tag VARCHAR(50) DEFAULT '',
    
    -- Mini Banner 1 (Right Bottom Left)
    mini_banner_1 VARCHAR(255) DEFAULT '',
    mini_banner_1_title VARCHAR(100) DEFAULT '',
    
    -- Mini Banner 2 (Right Bottom Right)
    mini_banner_2 VARCHAR(255)  DEFAULT '',
    mini_banner_2_title VARCHAR(100) DEFAULT '',

    -- Metadata
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ensure Singleton Row
CREATE UNIQUE INDEX only_one_row ON hero_sections ((1));

-- Seed Initial Row
INSERT INTO hero_sections (id, main_title) VALUES (1, 'Online Super Market');