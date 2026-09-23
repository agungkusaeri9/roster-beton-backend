package pgsql

import (
	"database/sql"
	"fmt"
	"strings"

	"go-arch/internal/entity"
	"go-arch/internal/repository"

	"github.com/jmoiron/sqlx"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepoPg(db *sqlx.DB) repository.ProductRepository {
	return &productRepo{
		db: db,
	}
}

type productRow struct {
	ID          int64           `db:"id"`
	Name        string          `db:"name"`
	Slug        string          `db:"slug"`
	CategoryID  *int64          `db:"category_id"`
	Description string          `db:"description"`
	Image       string          `db:"image"`
	Price       *float64        `db:"price"`
	Dimension   *string         `db:"dimension"`
	Weight      *float64        `db:"weight"`
	Stock       int             `db:"stock"`
	IsActive    bool            `db:"is_active"`
	IsFeatured  bool            `db:"is_featured"`
	CreatedAt   sql.NullTime    `db:"created_at"`
	UpdatedAt   sql.NullTime    `db:"updated_at"`
	// Joined category fields
	CatID       sql.NullInt64   `db:"cat_id"`
	CatName     sql.NullString  `db:"cat_name"`
	CatSlug     sql.NullString  `db:"cat_slug"`
	CatDesc     sql.NullString  `db:"cat_description"`
}

func (r *productRepo) FindAll(filter repository.ProductFilter) ([]*entity.Product, int64, error) {
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit < 1 {
		filter.Limit = 10
	}
	offset := (filter.Page - 1) * filter.Limit

	var conditions []string
	var args []interface{}
	argIdx := 1

	// Filter active products
	conditions = append(conditions, "p.is_active = true")

	// Search filter
	if strings.TrimSpace(filter.Search) != "" {
		conditions = append(conditions, fmt.Sprintf("(p.name ILIKE $%d OR p.description ILIKE $%d)", argIdx, argIdx))
		args = append(args, "%"+strings.TrimSpace(filter.Search)+"%")
		argIdx++
	}

	// Category filter (slug or name)
	if strings.TrimSpace(filter.Category) != "" && strings.ToLower(filter.Category) != "semua" {
		conditions = append(conditions, fmt.Sprintf("(LOWER(c.slug) = LOWER($%d) OR LOWER(c.name) = LOWER($%d))", argIdx, argIdx))
		args = append(args, strings.TrimSpace(filter.Category))
		argIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 1. Count Total
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT p.id)
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
	`, whereClause)

	var total int64
	err := r.db.Get(&total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed counting products: %w", err)
	}

	if total == 0 {
		return []*entity.Product{}, 0, nil
	}

	// 2. Query Products with Category JOIN
	args = append(args, filter.Limit, offset)
	dataQuery := fmt.Sprintf(`
		SELECT 
			p.id, p.name, p.slug, p.category_id, p.description, p.image,
			p.price, p.dimension, p.weight, p.stock, p.is_active, p.is_featured,
			p.created_at, p.updated_at,
			c.id AS cat_id, c.name AS cat_name, c.slug AS cat_slug, c.description AS cat_description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		%s
		ORDER BY p.id ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	var rows []productRow
	err = r.db.Select(&rows, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed fetching products: %w", err)
	}

	if len(rows) == 0 {
		return []*entity.Product{}, total, nil
	}

	// Collect Product IDs for batch loading relations
	productIDs := make([]int64, len(rows))
	products := make([]*entity.Product, len(rows))
	productMap := make(map[int64]*entity.Product)

	for i, row := range rows {
		p := &entity.Product{
			ID:          row.ID,
			Name:        row.Name,
			Slug:        row.Slug,
			CategoryID:  row.CategoryID,
			Description: row.Description,
			Image:       row.Image,
			Price:       row.Price,
			Dimension:   row.Dimension,
			Weight:      row.Weight,
			Stock:       row.Stock,
			IsActive:    row.IsActive,
			IsFeatured:  row.IsFeatured,
			Galleries:   []entity.ProductGallery{},
			Gallery:     []string{},
			CreatedAt:   row.CreatedAt.Time,
			UpdatedAt:   row.UpdatedAt.Time,
		}

		if row.CatID.Valid {
			p.Category = &entity.Category{
				ID:          row.CatID.Int64,
				Name:        row.CatName.String,
				Slug:        row.CatSlug.String,
				Description: row.CatDesc.String,
			}
		}

		products[i] = p
		productIDs[i] = row.ID
		productMap[row.ID] = p
	}

	// 3. Batch load Product Galleries
	if err := r.loadGalleries(productIDs, productMap); err != nil {
		return nil, 0, err
	}

	// 4. Batch load Polymorphic SEO
	if err := r.loadSEO("product", productIDs, productMap); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepo) FindBySlug(slug string) (*entity.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.slug, p.category_id, p.description, p.image,
			p.price, p.dimension, p.weight, p.stock, p.is_active, p.is_featured,
			p.created_at, p.updated_at,
			c.id AS cat_id, c.name AS cat_name, c.slug AS cat_slug, c.description AS cat_description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.slug = $1 AND p.is_active = true
		LIMIT 1
	`
	var row productRow
	err := r.db.Get(&row, query, slug)
	if err != nil {
		return nil, err
	}

	p := &entity.Product{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		CategoryID:  row.CategoryID,
		Description: row.Description,
		Image:       row.Image,
		Price:       row.Price,
		Dimension:   row.Dimension,
		Weight:      row.Weight,
		Stock:       row.Stock,
		IsActive:    row.IsActive,
		IsFeatured:  row.IsFeatured,
		Galleries:   []entity.ProductGallery{},
		Gallery:     []string{},
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	if row.CatID.Valid {
		p.Category = &entity.Category{
			ID:          row.CatID.Int64,
			Name:        row.CatName.String,
			Slug:        row.CatSlug.String,
			Description: row.CatDesc.String,
		}
	}

	productMap := map[int64]*entity.Product{p.ID: p}
	_ = r.loadGalleries([]int64{p.ID}, productMap)
	_ = r.loadSEO("product", []int64{p.ID}, productMap)

	return p, nil
}

func (r *productRepo) FindByID(id int64) (*entity.Product, error) {
	query := `
		SELECT 
			p.id, p.name, p.slug, p.category_id, p.description, p.image,
			p.price, p.dimension, p.weight, p.stock, p.is_active, p.is_featured,
			p.created_at, p.updated_at,
			c.id AS cat_id, c.name AS cat_name, c.slug AS cat_slug, c.description AS cat_description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1 AND p.is_active = true
		LIMIT 1
	`
	var row productRow
	err := r.db.Get(&row, query, id)
	if err != nil {
		return nil, err
	}

	p := &entity.Product{
		ID:          row.ID,
		Name:        row.Name,
		Slug:        row.Slug,
		CategoryID:  row.CategoryID,
		Description: row.Description,
		Image:       row.Image,
		Price:       row.Price,
		Dimension:   row.Dimension,
		Weight:      row.Weight,
		Stock:       row.Stock,
		IsActive:    row.IsActive,
		IsFeatured:  row.IsFeatured,
		Galleries:   []entity.ProductGallery{},
		Gallery:     []string{},
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}

	if row.CatID.Valid {
		p.Category = &entity.Category{
			ID:          row.CatID.Int64,
			Name:        row.CatName.String,
			Slug:        row.CatSlug.String,
			Description: row.CatDesc.String,
		}
	}

	productMap := map[int64]*entity.Product{p.ID: p}
	_ = r.loadGalleries([]int64{p.ID}, productMap)
	_ = r.loadSEO("product", []int64{p.ID}, productMap)

	return p, nil
}

func (r *productRepo) loadGalleries(productIDs []int64, productMap map[int64]*entity.Product) error {
	if len(productIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		SELECT id, product_id, image_url, alt_text, sort_order, created_at
		FROM product_galleries
		WHERE product_id IN (?)
		ORDER BY sort_order ASC
	`, productIDs)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	var galleries []entity.ProductGallery
	err = r.db.Select(&galleries, query, args...)
	if err != nil {
		return fmt.Errorf("failed to fetch product galleries: %w", err)
	}

	for _, g := range galleries {
		if p, ok := productMap[g.ProductID]; ok {
			p.Galleries = append(p.Galleries, g)
			p.Gallery = append(p.Gallery, g.ImageURL)
		}
	}

	return nil
}

func (r *productRepo) loadSEO(refType string, productIDs []int64, productMap map[int64]*entity.Product) error {
	if len(productIDs) == 0 {
		return nil
	}

	query, args, err := sqlx.In(`
		SELECT id, reference_type, reference_id, meta_title, meta_description, meta_keywords,
		       canonical_url, og_title, og_description, og_image, og_type, structured_data,
		       created_at, updated_at
		FROM seos
		WHERE reference_type = ? AND reference_id IN (?)
	`, refType, productIDs)
	if err != nil {
		return err
	}

	query = r.db.Rebind(query)
	var seos []entity.SEO
	err = r.db.Select(&seos, query, args...)
	if err != nil {
		return fmt.Errorf("failed to fetch product seos: %w", err)
	}

	for _, s := range seos {
		if p, ok := productMap[s.ReferenceID]; ok {
			seoCopy := s
			p.SEO = &seoCopy
		}
	}

	return nil
}
