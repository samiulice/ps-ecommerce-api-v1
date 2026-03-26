package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/projuktisheba/pse-api-v1/internal/model"
)

type PurchaseRepo struct {
	db *pgxpool.Pool
}

func NewPurchaseRepo(db *pgxpool.Pool) *PurchaseRepo {
	return &PurchaseRepo{db: db}
}

func (r *PurchaseRepo) Create(ctx context.Context, p *model.Purchase) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO purchases (
			purchase_date, prefix_code, count_id, purchase_code, reference_no,
			purchase_order_id, party_id, state_id, carrier_id, note,
			shipping_charge, is_shipping_charge_distributed, round_off, grand_total,
			change_return, paid_amount, currency_id, exchange_rate, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20
		)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, query,
		p.PurchaseDate, p.PrefixCode, p.CountID, p.PurchaseCode, p.ReferenceNo,
		p.PurchaseOrderID, p.PartyID, p.StateID, p.CarrierID, p.Note,
		p.ShippingCharge, p.IsShippingChargeDistributed, p.RoundOff, p.GrandTotal,
		p.ChangeReturn, p.PaidAmount, p.CurrencyID, p.ExchangeRate, p.CreatedBy, p.UpdatedBy,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)

	if mappedErr := mapPurchaseUniqueViolation(err); mappedErr != nil {
		return mappedErr
	}
	if err != nil {
		return err
	}

	if err := r.replaceItems(ctx, tx, p.ID, p.Items); err != nil {
		return err
	}

	if err := r.insertAttachments(ctx, tx, p.ID, p.Attachments); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PurchaseRepo) FindByID(ctx context.Context, id int64) (*model.Purchase, error) {
	query := `
		SELECT
			p.id, p.purchase_date, p.prefix_code, p.count_id, p.purchase_code, p.reference_no,
			p.purchase_order_id, p.party_id, COALESCE(s.name, ''), p.state_id, p.carrier_id,
			p.note, p.shipping_charge, p.is_shipping_charge_distributed, p.round_off, p.grand_total,
			p.change_return, p.paid_amount, p.currency_id, p.exchange_rate, p.created_by, p.updated_by,
			p.created_at, p.updated_at
		FROM purchases p
		LEFT JOIN suppliers s ON s.id = p.party_id
		WHERE p.id = $1`

	p := &model.Purchase{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.PurchaseDate, &p.PrefixCode, &p.CountID, &p.PurchaseCode, &p.ReferenceNo,
		&p.PurchaseOrderID, &p.PartyID, &p.SupplierName, &p.StateID, &p.CarrierID,
		&p.Note, &p.ShippingCharge, &p.IsShippingChargeDistributed, &p.RoundOff, &p.GrandTotal,
		&p.ChangeReturn, &p.PaidAmount, &p.CurrencyID, &p.ExchangeRate, &p.CreatedBy, &p.UpdatedBy,
		&p.CreatedAt, &p.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("purchase with id %d not found", id)
	}
	if err == nil {
		itemsMap, itemsErr := r.fetchItemsByPurchaseIDs(ctx, []int64{id})
		if itemsErr != nil {
			return nil, itemsErr
		}
		p.Items = itemsMap[id]

		attachmentsMap, attachErr := r.fetchAttachmentsByPurchaseIDs(ctx, []int64{id})
		if attachErr != nil {
			return nil, attachErr
		}
		p.Attachments = attachmentsMap[id]
	}
	return p, err
}

func (r *PurchaseRepo) Update(ctx context.Context, p *model.Purchase) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE purchases SET
			purchase_date = $1,
			prefix_code = $2,
			count_id = $3,
			purchase_code = $4,
			reference_no = $5,
			purchase_order_id = $6,
			party_id = $7,
			state_id = $8,
			carrier_id = $9,
			note = $10,
			shipping_charge = $11,
			is_shipping_charge_distributed = $12,
			round_off = $13,
			grand_total = $14,
			change_return = $15,
			paid_amount = $16,
			currency_id = $17,
			exchange_rate = $18,
			created_by = $19,
			updated_by = $20,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $21
		RETURNING updated_at`

	err = tx.QueryRow(ctx, query,
		p.PurchaseDate, p.PrefixCode, p.CountID, p.PurchaseCode, p.ReferenceNo,
		p.PurchaseOrderID, p.PartyID, p.StateID, p.CarrierID, p.Note,
		p.ShippingCharge, p.IsShippingChargeDistributed, p.RoundOff, p.GrandTotal,
		p.ChangeReturn, p.PaidAmount, p.CurrencyID, p.ExchangeRate, p.CreatedBy, p.UpdatedBy,
		p.ID,
	).Scan(&p.UpdatedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("purchase with id %d not found", p.ID)
	}
	if mappedErr := mapPurchaseUniqueViolation(err); mappedErr != nil {
		return mappedErr
	}
	if err != nil {
		return err
	}

	if err := r.replaceItems(ctx, tx, p.ID, p.Items); err != nil {
		return err
	}

	if len(p.RemoveAttachmentIDs) > 0 {
		if err := r.deleteAttachments(ctx, tx, p.ID, p.RemoveAttachmentIDs); err != nil {
			return err
		}
	}

	if err := r.insertAttachments(ctx, tx, p.ID, p.Attachments); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *PurchaseRepo) Delete(ctx context.Context, id int64) error {
	tag, err := r.db.Exec(ctx, "DELETE FROM purchases WHERE id = $1", id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("purchase with id %d not found", id)
	}
	return nil
}

func (r *PurchaseRepo) ExistsByCode(ctx context.Context, purchaseCode string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM purchases WHERE purchase_code = $1)`
	err := r.db.QueryRow(ctx, query, purchaseCode).Scan(&exists)
	return exists, err
}

func (r *PurchaseRepo) List(ctx context.Context, filter model.PurchaseFilter) ([]model.Purchase, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		conditions = append(conditions, fmt.Sprintf(
			`(p.purchase_code ILIKE $%d OR p.reference_no ILIKE $%d OR s.name ILIKE $%d OR p.count_id ILIKE $%d)`,
			argIndex, argIndex, argIndex, argIndex,
		))
		args = append(args, searchPattern)
		argIndex++
	}

	if filter.HasSupplierSet && filter.PartyID > 0 {
		conditions = append(conditions, fmt.Sprintf("p.party_id = $%d", argIndex))
		args = append(args, filter.PartyID)
		argIndex++
	}

	if filter.FromDate != "" {
		conditions = append(conditions, fmt.Sprintf("p.purchase_date >= $%d", argIndex))
		args = append(args, filter.FromDate)
		argIndex++
	}

	if filter.ToDate != "" {
		conditions = append(conditions, fmt.Sprintf("p.purchase_date <= $%d", argIndex))
		args = append(args, filter.ToDate)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM purchases p
		LEFT JOIN suppliers s ON s.id = p.party_id
		%s
	`, whereClause)

	var total int64
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("failed to count purchases: %w", err)
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (filter.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	query := fmt.Sprintf(`
		SELECT
			p.id, p.purchase_date, p.prefix_code, p.count_id, p.purchase_code, p.reference_no,
			p.purchase_order_id, p.party_id, COALESCE(s.name, ''), p.state_id, p.carrier_id,
			p.note, p.shipping_charge, p.is_shipping_charge_distributed, p.round_off, p.grand_total,
			p.change_return, p.paid_amount, p.currency_id, p.exchange_rate, p.created_by, p.updated_by,
			p.created_at, p.updated_at
		FROM purchases p
		LEFT JOIN suppliers s ON s.id = p.party_id
		%s
		ORDER BY p.id DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to query purchases: %w", err)
	}
	defer rows.Close()

	var purchases []model.Purchase
	for rows.Next() {
		var p model.Purchase
		err := rows.Scan(
			&p.ID, &p.PurchaseDate, &p.PrefixCode, &p.CountID, &p.PurchaseCode, &p.ReferenceNo,
			&p.PurchaseOrderID, &p.PartyID, &p.SupplierName, &p.StateID, &p.CarrierID,
			&p.Note, &p.ShippingCharge, &p.IsShippingChargeDistributed, &p.RoundOff, &p.GrandTotal,
			&p.ChangeReturn, &p.PaidAmount, &p.CurrencyID, &p.ExchangeRate, &p.CreatedBy, &p.UpdatedBy,
			&p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan purchase: %w", err)
		}
		purchases = append(purchases, p)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating purchase rows: %w", err)
	}

	if len(purchases) > 0 {
		ids := make([]int64, 0, len(purchases))
		for i := range purchases {
			ids = append(ids, purchases[i].ID)
		}

		attachmentsMap, err := r.fetchAttachmentsByPurchaseIDs(ctx, ids)
		if err != nil {
			return nil, 0, err
		}

		itemsMap, err := r.fetchItemsByPurchaseIDs(ctx, ids)
		if err != nil {
			return nil, 0, err
		}

		for i := range purchases {
			purchases[i].Items = itemsMap[purchases[i].ID]
			purchases[i].Attachments = attachmentsMap[purchases[i].ID]
		}
	}

	return purchases, total, nil
}

func (r *PurchaseRepo) replaceItems(ctx context.Context, tx pgx.Tx, purchaseID int64, items []model.PurchaseItem) error {
	if _, err := tx.Exec(ctx, `DELETE FROM purchase_items WHERE purchase_id = $1`, purchaseID); err != nil {
		return fmt.Errorf("failed to clear purchase items: %w", err)
	}

	if len(items) == 0 {
		return nil
	}

	query := `
		INSERT INTO purchase_items (
			purchase_id, item_type, product_id, item_name, quantity, unit_price, total_price, note
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at
	`

	for i := range items {
		items[i].PurchaseID = purchaseID
		if err := tx.QueryRow(ctx, query,
			purchaseID,
			items[i].ItemType,
			items[i].ProductID,
			items[i].ItemName,
			items[i].Quantity,
			items[i].UnitPrice,
			items[i].TotalPrice,
			items[i].Note,
		).Scan(&items[i].ID, &items[i].CreatedAt, &items[i].UpdatedAt); err != nil {
			return fmt.Errorf("failed to insert purchase item: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepo) fetchItemsByPurchaseIDs(ctx context.Context, purchaseIDs []int64) (map[int64][]model.PurchaseItem, error) {
	result := make(map[int64][]model.PurchaseItem)
	if len(purchaseIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT id, purchase_id, item_type, product_id, item_name, quantity, unit_price, total_price, note, created_at, updated_at
		FROM purchase_items
		WHERE purchase_id = ANY($1)
		ORDER BY id ASC
	`

	rows, err := r.db.Query(ctx, query, purchaseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch purchase items: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.PurchaseItem
		if err := rows.Scan(
			&item.ID,
			&item.PurchaseID,
			&item.ItemType,
			&item.ProductID,
			&item.ItemName,
			&item.Quantity,
			&item.UnitPrice,
			&item.TotalPrice,
			&item.Note,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan purchase item: %w", err)
		}

		result[item.PurchaseID] = append(result[item.PurchaseID], item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchase items: %w", err)
	}

	return result, nil
}

func (r *PurchaseRepo) insertAttachments(ctx context.Context, tx pgx.Tx, purchaseID int64, attachments []model.PurchaseAttachment) error {
	if len(attachments) == 0 {
		return nil
	}

	query := `
		INSERT INTO purchase_attachments (
			purchase_id, file_url, file_name, file_ext, mime_type, file_size
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	for _, a := range attachments {
		if strings.TrimSpace(a.FileURL) == "" {
			continue
		}

		if _, err := tx.Exec(ctx, query,
			purchaseID,
			a.FileURL,
			a.FileName,
			a.FileExt,
			a.MimeType,
			a.FileSize,
		); err != nil {
			return fmt.Errorf("failed to insert purchase attachment: %w", err)
		}
	}

	return nil
}

func (r *PurchaseRepo) deleteAttachments(ctx context.Context, tx pgx.Tx, purchaseID int64, attachmentIDs []int64) error {
	if len(attachmentIDs) == 0 {
		return nil
	}

	ids := make([]int64, 0, len(attachmentIDs))
	for _, id := range attachmentIDs {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })

	query := `DELETE FROM purchase_attachments WHERE purchase_id = $1 AND id = ANY($2)`
	if _, err := tx.Exec(ctx, query, purchaseID, ids); err != nil {
		return fmt.Errorf("failed to delete purchase attachments: %w", err)
	}

	return nil
}

func (r *PurchaseRepo) fetchAttachmentsByPurchaseIDs(ctx context.Context, purchaseIDs []int64) (map[int64][]model.PurchaseAttachment, error) {
	result := make(map[int64][]model.PurchaseAttachment)
	if len(purchaseIDs) == 0 {
		return result, nil
	}

	query := `
		SELECT id, purchase_id, file_url, file_name, file_ext, mime_type, file_size, created_at
		FROM purchase_attachments
		WHERE purchase_id = ANY($1)
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, query, purchaseIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch purchase attachments: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var a model.PurchaseAttachment
		if err := rows.Scan(
			&a.ID,
			&a.PurchaseID,
			&a.FileURL,
			&a.FileName,
			&a.FileExt,
			&a.MimeType,
			&a.FileSize,
			&a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan purchase attachment: %w", err)
		}

		result[a.PurchaseID] = append(result[a.PurchaseID], a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating purchase attachment rows: %w", err)
	}

	return result, nil
}

func mapPurchaseUniqueViolation(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) || pgErr.Code != "23505" {
		return nil
	}

	switch pgErr.ConstraintName {
	case "purchases_purchase_code_unique":
		return fmt.Errorf("purchase_code already exists")
	default:
		return fmt.Errorf("duplicate purchase record")
	}
}
