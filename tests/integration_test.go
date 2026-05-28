package tests

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/builder"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/cache"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/loader"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/repository"
)

func TestIntegration_ProductDocumentBuilder(
	t *testing.T,
) {

	ctx := context.Background()

	databaseURL :=
		"postgres://postgres:postgres@localhost:5432/wyzauto?sslmode=disable"

	db, err := pgxpool.New(
		ctx,
		databaseURL,
	)

	if err != nil {
		t.Fatalf("connect db: %v", err)
	}

	defer db.Close()

	translationCache := cache.NewTranslationCache(
		5 * time.Minute,
	)

	translationLoader := loader.NewPostgresLoader(
		db,
		translationCache,
	)

	productRepository :=
		repository.NewPostgresProductRepository(
			db,
		)

	documentBuilder :=
		builder.NewProductDocumentBuilder(
			productRepository,
			translationLoader,
		)

	document, err := documentBuilder.Build(
		ctx,
		"11111111-1111-1111-1111-111111111111",
		[]string{
			"en",
			"th",
		},
	)

	if err != nil {
		t.Fatalf("build document: %v", err)
	}

	if document.SKU !=
		"BP-OIL-5W30-1L" {

		t.Fatalf(
			"unexpected sku: %s",
			document.SKU,
		)
	}

	if len(document.ProductName) != 2 {
		t.Fatalf(
			"expected 2 translations",
		)
	}
}

