ALTER TABLE general_settings
    DROP COLUMN IF EXISTS company_logo_mobile,
    DROP COLUMN IF EXISTS company_logo_footer,
    DROP COLUMN IF EXISTS company_email,
    DROP COLUMN IF EXISTS company_phone,
    DROP COLUMN IF EXISTS company_address,
    DROP COLUMN IF EXISTS about_us_short,
    DROP COLUMN IF EXISTS about_us_long;
