CREATE TABLE IF NOT EXISTS social_links (
	id             BIGSERIAL PRIMARY KEY,
	name           VARCHAR(120) NOT NULL,
	icon           VARCHAR(160) NOT NULL,
	link           VARCHAR(500) NOT NULL,
	alt_text       VARCHAR(255) NOT NULL DEFAULT '',
	is_active      BOOLEAN NOT NULL DEFAULT TRUE,
	show_in_topbar BOOLEAN NOT NULL DEFAULT FALSE,
	display_order  INTEGER NOT NULL DEFAULT 0,
	created_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at     TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_social_links_active_topbar
	ON social_links (is_active, show_in_topbar, display_order, id);

CREATE INDEX IF NOT EXISTS idx_social_links_display_order
	ON social_links (display_order, id);