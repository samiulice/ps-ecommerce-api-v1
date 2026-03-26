CREATE TABLE IF NOT EXISTS purchase_attachments (
    id          BIGSERIAL PRIMARY KEY,
    purchase_id BIGINT NOT NULL,
    file_url    VARCHAR(500) NOT NULL,
    file_name   VARCHAR(255),
    file_ext    VARCHAR(20),
    mime_type   VARCHAR(120),
    file_size   BIGINT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT purchase_attachments_purchase_id_fkey
        FOREIGN KEY (purchase_id) REFERENCES purchases(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_purchase_attachments_purchase_id
    ON purchase_attachments(purchase_id);
