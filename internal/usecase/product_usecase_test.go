package usecase_test

import (
	"testing"

	"go-arch/internal/entity"
	"go-arch/internal/repository"
	"go-arch/internal/usecase"
)

type mockProductRepo struct {
	products []*entity.Product
	total    int64
}

func (m *mockProductRepo) FindAll(filter repository.ProductFilter) ([]*entity.Product, int64, error) {
	return m.products, m.total, nil
}

func (m *mockProductRepo) FindBySlug(slug string) (*entity.Product, error) {
	if len(m.products) > 0 {
		return m.products[0], nil
	}
	return nil, nil
}

func (m *mockProductRepo) FindByID(id int64) (*entity.Product, error) {
	if len(m.products) > 0 {
		return m.products[0], nil
	}
	return nil, nil
}

func TestProductUsecase_Pagination(t *testing.T) {
	mockRepo := &mockProductRepo{
		products: []*entity.Product{
			{ID: 1, Name: "Product 1", Slug: "product-1"},
			{ID: 2, Name: "Product 2", Slug: "product-2"},
		},
		total: 34,
	}

	uc := usecase.NewProductUsecase(mockRepo)

	products, pagination, err := uc.GetAll(1, 10, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}

	if pagination.CurrentPage != 1 {
		t.Errorf("expected current_page 1, got %d", pagination.CurrentPage)
	}

	if pagination.PerPage != 10 {
		t.Errorf("expected per_page 10, got %d", pagination.PerPage)
	}

	if pagination.TotalData != 34 {
		t.Errorf("expected total_data 34, got %d", pagination.TotalData)
	}

	if pagination.TotalPages != 4 {
		t.Errorf("expected total_pages 4 (ceil of 34/10), got %d", pagination.TotalPages)
	}
}
