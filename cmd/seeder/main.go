package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"go-arch/internal/infrastructure/pgsql"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
)

func findEnvFile() string {
	possiblePaths := []string{
		".env",
		"../.env",
		"../../.env",
	}
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	return ".env"
}

type CategorySeed struct {
	Name        string
	Slug        string
	Description string
}

type ProductSeed struct {
	Name            string
	Slug            string
	CategoryName    string
	Description     string
	Image           string
	Price           *float64
	Dimension       *string
	Weight          *float64
	Stock           int
	IsActive        bool
	IsFeatured      bool
	Galleries       []string
	MetaTitle       string
	MetaDescription string
	MetaKeywords    string
	CanonicalURL    string
	OGTitle         string
	OGDescription   string
	OGImage         string
}

type GallerySeed struct {
	Src       string
	Title     string
	Category  string
	AltText   string
	SortOrder int
}

func ptrFloat(f float64) *float64 {
	return &f
}

func ptrString(s string) *string {
	return &s
}

func getBaseURL() string {
	appURL := os.Getenv("APP_URL")
	if appURL != "" {
		return strings.TrimRight(appURL, "/")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return fmt.Sprintf("http://localhost:%s", port)
}

func toUploadURL(baseURL, p string) string {
	if strings.HasPrefix(p, "http://") || strings.HasPrefix(p, "https://") {
		return p
	}
	if strings.HasPrefix(p, "/images/products/") {
		p = strings.Replace(p, "/images/products/", "/uploads/products/", 1)
	}
	if !strings.HasPrefix(p, "/") {
		p = "/" + p
	}
	return baseURL + p
}

func main() {
	envFile := findEnvFile()
	if err := godotenv.Load(envFile); err != nil {
		log.Printf("⚠️  Warning: Error loading .env file from %s: %v", envFile, err)
	} else {
		log.Printf("✅ Env loaded successfully from %s", envFile)
	}

	baseURL := getBaseURL()
	log.Printf("🌐 Using Base URL: %s", baseURL)

	db, err := pgsql.Init()
	if err != nil {
		log.Fatalf("❌ Failed to connect PostgreSQL: %v", err)
	}
	defer db.Close()

	log.Println("🌱 Starting Database Seeding for Categories, Products, Galleries & SEO...")

	tx, err := db.Beginx()
	if err != nil {
		log.Fatalf("❌ Failed to start transaction: %v", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Fatalf("❌ Panic occurred during seeding: %v", r)
		}
	}()

	// 1. Seed Categories
	categoryMap, err := seedCategories(tx)
	if err != nil {
		tx.Rollback()
		log.Fatalf("❌ Failed seeding categories: %v", err)
	}

	// 2. Seed Products and Product Galleries
	if err := seedProducts(tx, categoryMap, baseURL); err != nil {
		tx.Rollback()
		log.Fatalf("❌ Failed seeding products: %v", err)
	}

	// 3. Seed General Galleries
	if err := seedGalleries(tx, baseURL); err != nil {
		tx.Rollback()
		log.Fatalf("❌ Failed seeding general galleries: %v", err)
	}

	// 4. Seed Dynamic Site Configurations
	if err := seedConfigs(tx); err != nil {
		tx.Rollback()
		log.Fatalf("❌ Failed seeding configs: %v", err)
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("❌ Failed to commit transaction: %v", err)
	}

	log.Println("🎉 Database Seeding completed successfully!")
}

func seedCategories(tx *sqlx.Tx) (map[string]int64, error) {
	log.Println("📁 Seeding Categories...")

	categories := []CategorySeed{
		{Name: "Minimalis", Slug: "minimalis", Description: "Roster beton dengan pola minimalis, bersih, dan modern. Sangat cocok untuk arsitektur kontemporer."},
		{Name: "Klasik", Slug: "klasik", Description: "Roster beton berdesain klasik yang timeless, memberikan kesan mewah dan kokoh pada fasad bangunan."},
		{Name: "Motif Klasik", Slug: "motif-klasik", Description: "Roster beton dengan motif kotak dan aksen klasik tradisional khas bangunan elegan."},
		{Name: "Geometri", Slug: "geometri", Description: "Roster beton bermotif geometris presisi (segitiga, lingkaran, kubus) yang artistik dan dinamis."},
		{Name: "Premium", Slug: "premium", Description: "Roster beton kualitas terbaik dengan finishing halus super presisi dan daya tahan maksimal."},
		{Name: "Tradisional", Slug: "tradisional", Description: "Roster beton bernuansa etnik nusantara dan ukiran tradisional yang kaya nilai seni budaya."},
		{Name: "Motif Floral", Slug: "motif-floral", Description: "Roster beton berornamen bunga dan tanaman yang anggun, memberikan nuansa natural dan asri."},
		{Name: "Motif Alam", Slug: "motif-alam", Description: "Roster beton dengan motif dedaunan dan elemen alam untuk sirkulasi udara optimal."},
		{Name: "Ukuran Besar", Slug: "ukuran-besar", Description: "Roster beton ukuran besar khusus untuk pagar perimeter, pabrik, gedung, dan area luas."},
		{Name: "Ukuran Kecil", Slug: "ukuran-kecil", Description: "Roster beton ukuran kompak untuk ventilasi kamar mandi, jendela, dan ornamen detail."},
		{Name: "Ukuran Medium", Slug: "ukuran-medium", Description: "Roster beton ukuran standar serbaguna untuk partisi ruangan dan dinding semi-terbuka."},
		{Name: "Elegan", Slug: "elegan", Description: "Roster beton dengan sentuhan estetika mewah untuk interior dan eksterior hunian idaman."},
		{Name: "Motif Unik", Slug: "motif-unik", Description: "Koleksi roster beton dengan desain custom dan pola unik yang eksklusif."},
	}

	categoryMap := make(map[string]int64)

	for _, cat := range categories {
		var id int64
		query := `
			INSERT INTO categories (name, slug, description, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (name) DO UPDATE 
			SET slug = EXCLUDED.slug, description = EXCLUDED.description, updated_at = NOW()
			RETURNING id
		`
		err := tx.QueryRow(query, cat.Name, cat.Slug, cat.Description).Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("error inserting category %s: %w", cat.Name, err)
		}
		categoryMap[cat.Name] = id
	}

	log.Printf("✅ %d Categories seeded successfully", len(categories))
	return categoryMap, nil
}

func seedProducts(tx *sqlx.Tx, categoryMap map[string]int64, baseURL string) error {
	log.Println("🧱 Seeding Products, Product Galleries & Polymorphic SEO...")

	// Clear existing product galleries, products, and product seos for fresh seed
	if _, err := tx.Exec("DELETE FROM seos WHERE reference_type = 'product'"); err != nil {
		return fmt.Errorf("failed to clear seos for products: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM product_galleries"); err != nil {
		return fmt.Errorf("failed to clear product_galleries: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM products"); err != nil {
		return fmt.Errorf("failed to clear products: %w", err)
	}

	rawProducts := []struct {
		ID          int64
		Name        string
		Slug        string
		Image       string
		Category    string
		Description string
		Gallery     []string
	}{
		{
			ID:          1,
			Name:        "Roster Motif Kotak",
			Slug:        "roster-motif-kotak",
			Image:       "/images/products/galprod-17.jpeg",
			Category:    "Motif Klasik",
			Description: "Roster beton dengan motif kotak klasik yang cocok untuk berbagai jenis bangunan. Dibuat dengan material berkualitas tinggi dan proses produksi yang modern, roster ini menawarkan daya tahan yang luar biasa dan estetika yang menarik.",
			Gallery:     []string{"/images/products/galprod-17.jpeg", "/images/products/galprod-18.jpeg", "/images/products/galprod-19.jpeg", "/images/products/galprod-20.jpeg"},
		},
		{
			ID:          2,
			Name:        "Roster Motif Bunga",
			Slug:        "roster-motif-bunga",
			Image:       "/images/products/galprod-18.jpeg",
			Category:    "Motif Floral",
			Description: "Roster beton dengan motif bunga yang indah dan elegan. Memberikan sentuhan keindahan alam pada bangunan Anda. Perfect untuk pagar, ventilasi, dan dekorasi rumah.",
			Gallery:     []string{"/images/products/galprod-18.jpeg", "/images/products/galprod-22.jpeg", "/images/products/galprod-23.jpeg", "/images/products/galprod-24.jpeg"},
		},
		{
			ID:          3,
			Name:        "Roster Minimalis Modern",
			Slug:        "roster-minimalis-modern",
			Image:       "/images/products/galprod-19.jpeg",
			Category:    "Minimalis",
			Description: "Desain roster minimalis yang modern dan clean, ideal untuk rumah dengan konsep modern atau skandinavia. Memberikan kesan luas dan nyaman.",
			Gallery:     []string{"/images/products/galprod-19.jpeg", "/images/products/galprod-25.jpeg", "/images/products/galprod-26.jpeg", "/images/products/galprod-27.jpeg"},
		},
		{
			ID:          4,
			Name:        "Roster Geometri",
			Slug:        "roster-geometri",
			Image:       "/images/products/galprod-20.jpeg",
			Category:    "Geometri",
			Description: "Motif geometri yang unik dan menarik. Cocok untuk bangunan dengan desain arsitektur modern dan artistik.",
			Gallery:     []string{"/images/products/galprod-20.jpeg", "/images/products/galprod-32.jpeg", "/images/products/galprod-33.jpeg", "/images/products/galprod-34.jpeg"},
		},
		{
			ID:          5,
			Name:        "Roster Klasik Elegan",
			Slug:        "roster-klasik-elegan",
			Image:       "/images/products/galprod-22.jpeg",
			Category:    "Klasik",
			Description: "Roster beton dengan desain klasik yang elegan dan timeless. Memberikan nuansa mewah pada rumah Anda.",
			Gallery:     []string{"/images/products/galprod-22.jpeg", "/images/products/galprod-35.jpeg", "/images/products/galprod-36.jpeg", "/images/products/galprod-37.jpeg"},
		},
		{
			ID:          6,
			Name:        "Roster Premium",
			Slug:        "roster-premium",
			Image:       "/images/products/galprod-23.jpeg",
			Category:    "Premium",
			Description: "Kualitas terbaik dengan finishing premium, cocok untuk proyek-proyek besar dan bangunan mewah.",
			Gallery:     []string{"/images/products/galprod-23.jpeg", "/images/products/galprod-39.jpeg", "/images/products/galprod-40.jpeg", "/images/products/galprod-44.jpeg"},
		},
		{
			ID:          7,
			Name:        "Roster Tradisional",
			Slug:        "roster-tradisional",
			Image:       "/images/products/galprod-24.jpeg",
			Category:    "Tradisional",
			Description: "Motif tradisional yang kaya akan budaya dan seni lokal. Perfect untuk rumah dengan konsep adat Jawa.",
			Gallery:     []string{"/images/products/galprod-24.jpeg", "/images/products/galprod-47.jpeg", "/images/products/galprod-48.jpeg", "/images/products/galprod-50.jpeg"},
		},
		{
			ID:          8,
			Name:        "Roster Kotak Besar",
			Slug:        "roster-kotak-besar",
			Image:       "/images/products/galprod-25.jpeg",
			Category:    "Ukuran Besar",
			Description: "Roster beton dengan ukuran besar, cocok untuk area yang luas seperti pagar gedung atau pabrik.",
			Gallery:     []string{"/images/products/galprod-25.jpeg", "/images/products/galprod-51.jpeg", "/images/products/galprod-53.jpeg", "/images/products/galprod-54.jpeg"},
		},
		{
			ID:          9,
			Name:        "Roster Motif Daun",
			Slug:        "roster-motif-daun",
			Image:       "/images/products/galprod-26.jpeg",
			Category:    "Motif Alam",
			Description: "Inspirasi dari alam dengan motif daun yang natural dan menawan.",
			Gallery:     []string{"/images/products/galprod-26.jpeg", "/images/products/galprod-55.jpeg", "/images/products/galprod-56.jpeg", "/images/products/galprod-57.jpeg"},
		},
		{
			ID:          10,
			Name:        "Roster Modern Minimalis",
			Slug:        "roster-modern-minimalis",
			Image:       "/images/products/galprod-27.jpeg",
			Category:    "Minimalis",
			Description: "Desain modern minimalis dengan garis-garis yang bersih dan rapi.",
			Gallery:     []string{"/images/products/galprod-27.jpeg", "/images/products/galprod-58.jpeg", "/images/products/galprod-59.jpeg", "/images/products/galprod-17.jpeg"},
		},
		{
			ID:          11,
			Name:        "Roster Kotak Kecil",
			Slug:        "roster-kotak-kecil",
			Image:       "/images/products/galprod-28.jpeg",
			Category:    "Ukuran Kecil",
			Description: "Ukuran kecil yang cocok untuk area terbatas seperti jendela atau ventilasi kamar tidur.",
			Gallery:     []string{"/images/products/galprod-28.jpeg", "/images/products/galprod-29.jpeg", "/images/products/galprod-31.jpeg", "/images/products/galprod-32.jpeg"},
		},
		{
			ID:          12,
			Name:        "Roster Motif Bintang",
			Slug:        "roster-motif-bintang",
			Image:       "/images/products/galprod-29.jpeg",
			Category:    "Motif Unik",
			Description: "Motif bintang yang unik dan menarik, cocok untuk dekorasi anak-anak atau tempat bermain.",
			Gallery:     []string{"/images/products/galprod-29.jpeg", "/images/products/galprod-33.jpeg", "/images/products/galprod-34.jpeg", "/images/products/galprod-35.jpeg"},
		},
		{
			ID:          13,
			Name:        "Roster Premium Klasik",
			Slug:        "roster-premium-klasik",
			Image:       "/images/products/galprod-31.jpeg",
			Category:    "Premium",
			Description: "Kualitas premium dengan desain klasik yang elegan.",
			Gallery:     []string{"/images/products/galprod-31.jpeg", "/images/products/galprod-36.jpeg", "/images/products/galprod-37.jpeg", "/images/products/galprod-38.jpeg"},
		},
		{
			ID:          14,
			Name:        "Roster Geometri Modern",
			Slug:        "roster-geometri-modern",
			Image:       "/images/products/galprod-32.jpeg",
			Category:    "Geometri",
			Description: "Motif geometri dengan sentuhan modern.",
			Gallery:     []string{"/images/products/galprod-32.jpeg", "/images/products/galprod-39.jpeg", "/images/products/galprod-40.jpeg", "/images/products/galprod-44.jpeg"},
		},
		{
			ID:          15,
			Name:        "Roster Kotak Medium",
			Slug:        "roster-kotak-medium",
			Image:       "/images/products/galprod-33.jpeg",
			Category:    "Ukuran Medium",
			Description: "Ukuran medium yang fleksibel untuk berbagai kebutuhan.",
			Gallery:     []string{"/images/products/galprod-33.jpeg", "/images/products/galprod-47.jpeg", "/images/products/galprod-48.jpeg", "/images/products/galprod-50.jpeg"},
		},
		{
			ID:          16,
			Name:        "Roster Motif Segitiga",
			Slug:        "roster-motif-segitiga",
			Image:       "/images/products/galprod-34.jpeg",
			Category:    "Geometri",
			Description: "Motif segitiga yang dinamis dan energik.",
			Gallery:     []string{"/images/products/galprod-34.jpeg", "/images/products/galprod-51.jpeg", "/images/products/galprod-53.jpeg", "/images/products/galprod-54.jpeg"},
		},
		{
			ID:          17,
			Name:        "Roster Elegan",
			Slug:        "roster-elegan",
			Image:       "/images/products/galprod-35.jpeg",
			Category:    "Elegan",
			Description: "Desain elegan yang menambah keindahan bangunan Anda.",
			Gallery:     []string{"/images/products/galprod-35.jpeg", "/images/products/galprod-55.jpeg", "/images/products/galprod-56.jpeg", "/images/products/galprod-57.jpeg"},
		},
		{
			ID:          18,
			Name:        "Roster Minimalis Putih",
			Slug:        "roster-minimalis-putih",
			Image:       "/images/products/galprod-36.jpeg",
			Category:    "Minimalis",
			Description: "Warna putih yang bersih dengan desain minimalis.",
			Gallery:     []string{"/images/products/galprod-36.jpeg", "/images/products/galprod-58.jpeg", "/images/products/galprod-59.jpeg", "/images/products/galprod-17.jpeg"},
		},
		{
			ID:          19,
			Name:        "Roster Klasik Tradisional",
			Slug:        "roster-klasik-tradisional",
			Image:       "/images/products/galprod-37.jpeg",
			Category:    "Tradisional",
			Description: "Kombinasi desain klasik dengan sentuhan tradisional.",
			Gallery:     []string{"/images/products/galprod-37.jpeg", "/images/products/galprod-18.jpeg", "/images/products/galprod-19.jpeg", "/images/products/galprod-20.jpeg"},
		},
		{
			ID:          20,
			Name:        "Roster Motif Lingkaran",
			Slug:        "roster-motif-lingkaran",
			Image:       "/images/products/galprod-38.jpeg",
			Category:    "Geometri",
			Description: "Motif lingkaran yang lembut dan elegan.",
			Gallery:     []string{"/images/products/galprod-38.jpeg", "/images/products/galprod-22.jpeg", "/images/products/galprod-23.jpeg", "/images/products/galprod-24.jpeg"},
		},
		{
			ID:          21,
			Name:        "Roster Premium Modern",
			Slug:        "roster-premium-modern",
			Image:       "/images/products/galprod-39.jpeg",
			Category:    "Premium",
			Description: "Kualitas premium dengan desain modern terbaru.",
			Gallery:     []string{"/images/products/galprod-39.jpeg", "/images/products/galprod-25.jpeg", "/images/products/galprod-26.jpeg", "/images/products/galprod-27.jpeg"},
		},
		{
			ID:          22,
			Name:        "Roster Kotak Elegan",
			Slug:        "roster-kotak-elegan",
			Image:       "/images/products/galprod-40.jpeg",
			Category:    "Elegan",
			Description: "Motif kotak dengan sentuhan elegan.",
			Gallery:     []string{"/images/products/galprod-40.jpeg", "/images/products/galprod-32.jpeg", "/images/products/galprod-33.jpeg", "/images/products/galprod-34.jpeg"},
		},
		{
			ID:          23,
			Name:        "Roster Motif Alam",
			Slug:        "roster-motif-alam",
			Image:       "/images/products/galprod-44.jpeg",
			Category:    "Motif Alam",
			Description: "Inspirasi dari alam dengan motif yang natural.",
			Gallery:     []string{"/images/products/galprod-44.jpeg", "/images/products/galprod-35.jpeg", "/images/products/galprod-36.jpeg", "/images/products/galprod-37.jpeg"},
		},
		{
			ID:          24,
			Name:        "Roster Tradisional Klasik",
			Slug:        "roster-tradisional-klasik",
			Image:       "/images/products/galprod-47.jpeg",
			Category:    "Tradisional",
			Description: "Desain tradisional klasik yang kaya akan budaya.",
			Gallery:     []string{"/images/products/galprod-47.jpeg", "/images/products/galprod-39.jpeg", "/images/products/galprod-40.jpeg", "/images/products/galprod-44.jpeg"},
		},
		{
			ID:          25,
			Name:        "Roster Modern Geometri",
			Slug:        "roster-modern-geometri",
			Image:       "/images/products/galprod-48.jpeg",
			Category:    "Geometri",
			Description: "Motif geometri dengan nuansa modern.",
			Gallery:     []string{"/images/products/galprod-48.jpeg", "/images/products/galprod-47.jpeg", "/images/products/galprod-48.jpeg", "/images/products/galprod-50.jpeg"},
		},
		{
			ID:          26,
			Name:        "Roster Minimalis Elegan",
			Slug:        "roster-minimalis-elegan",
			Image:       "/images/products/galprod-50.jpeg",
			Category:    "Minimalis",
			Description: "Kombinasi minimalis dan elegan.",
			Gallery:     []string{"/images/products/galprod-50.jpeg", "/images/products/galprod-51.jpeg", "/images/products/galprod-53.jpeg", "/images/products/galprod-54.jpeg"},
		},
		{
			ID:          27,
			Name:        "Roster Premium Klasik V2",
			Slug:        "roster-premium-klasik-2",
			Image:       "/images/products/galprod-51.jpeg",
			Category:    "Premium",
			Description: "Kualitas premium dengan desain klasik yang timeless.",
			Gallery:     []string{"/images/products/galprod-51.jpeg", "/images/products/galprod-55.jpeg", "/images/products/galprod-56.jpeg", "/images/products/galprod-57.jpeg"},
		},
		{
			ID:          28,
			Name:        "Roster Motif Bunga Modern",
			Slug:        "roster-motif-bunga-modern",
			Image:       "/images/products/galprod-53.jpeg",
			Category:    "Motif Floral",
			Description: "Motif bunga dengan sentuhan modern.",
			Gallery:     []string{"/images/products/galprod-53.jpeg", "/images/products/galprod-58.jpeg", "/images/products/galprod-59.jpeg", "/images/products/galprod-17.jpeg"},
		},
		{
			ID:          29,
			Name:        "Roster Kotak Premium",
			Slug:        "roster-kotak-premium",
			Image:       "/images/products/galprod-54.jpeg",
			Category:    "Premium",
			Description: "Motif kotak dengan kualitas premium.",
			Gallery:     []string{"/images/products/galprod-54.jpeg", "/images/products/galprod-18.jpeg", "/images/products/galprod-19.jpeg", "/images/products/galprod-20.jpeg"},
		},
		{
			ID:          30,
			Name:        "Roster Geometri Unik",
			Slug:        "roster-geometri-unik",
			Image:       "/images/products/galprod-55.jpeg",
			Category:    "Geometri",
			Description: "Motif geometri yang unik dan menarik.",
			Gallery:     []string{"/images/products/galprod-55.jpeg", "/images/products/galprod-22.jpeg", "/images/products/galprod-23.jpeg", "/images/products/galprod-24.jpeg"},
		},
		{
			ID:          31,
			Name:        "Roster Minimalis Putih V2",
			Slug:        "roster-minimalis-putih-2",
			Image:       "/images/products/galprod-56.jpeg",
			Category:    "Minimalis",
			Description: "Warna putih bersih dengan desain minimalis.",
			Gallery:     []string{"/images/products/galprod-56.jpeg", "/images/products/galprod-25.jpeg", "/images/products/galprod-26.jpeg", "/images/products/galprod-27.jpeg"},
		},
		{
			ID:          32,
			Name:        "Roster Tradisional Elegan",
			Slug:        "roster-tradisional-elegan",
			Image:       "/images/products/galprod-57.jpeg",
			Category:    "Tradisional",
			Description: "Desain tradisional dengan sentuhan elegan.",
			Gallery:     []string{"/images/products/galprod-57.jpeg", "/images/products/galprod-32.jpeg", "/images/products/galprod-33.jpeg", "/images/products/galprod-34.jpeg"},
		},
		{
			ID:          33,
			Name:        "Roster Motif Kotak Modern",
			Slug:        "roster-motif-kotak-modern",
			Image:       "/images/products/galprod-58.jpeg",
			Category:    "Motif Klasik",
			Description: "Motif kotak dengan desain modern.",
			Gallery:     []string{"/images/products/galprod-58.jpeg", "/images/products/galprod-35.jpeg", "/images/products/galprod-36.jpeg", "/images/products/galprod-37.jpeg"},
		},
		{
			ID:          34,
			Name:        "Roster Premium Terbaru",
			Slug:        "roster-premium-terbaru",
			Image:       "/images/products/galprod-59.jpeg",
			Category:    "Premium",
			Description: "Kualitas premium dengan desain terbaru.",
			Gallery:     []string{"/images/products/galprod-59.jpeg", "/images/products/galprod-39.jpeg", "/images/products/galprod-40.jpeg", "/images/products/galprod-44.jpeg"},
		},
	}

	for _, item := range rawProducts {
		catID, exists := categoryMap[item.Category]
		var catIDPtr *int64
		if exists {
			catIDPtr = &catID
		}

		uploadImg := toUploadURL(baseURL, item.Image)
		price := 15000.0 + float64(item.ID%5)*5000.0
		dim := "20x20x10 cm"
		weight := 3.5
		isFeatured := item.ID <= 6

		metaTitle := fmt.Sprintf("%s | Jual Roster Beton Berkualitas Purwakarta", item.Name)
		metaDesc := fmt.Sprintf("Jual %s berkualitas unggulan dari CV Roster Purwakarta. %s Cocok untuk ventilasi rumah, pagar, dan fasad bangunan modern.", item.Name, item.Description)
		metaKeywords := fmt.Sprintf("roster beton, %s, roster %s, roster purwakarta, roster beton plered, pagar roster, ventilasi beton minimalis", strings.ToLower(item.Name), strings.ToLower(item.Category))
		canonicalURL := fmt.Sprintf("https://rosterbetonpurwakarta.com/products/%s", item.Slug)
		ogTitle := fmt.Sprintf("%s - CV Roster Purwakarta", item.Name)
		ogDesc := item.Description
		ogImage := uploadImg

		var productID int64
		insertProductQuery := `
			INSERT INTO products (
				id, name, slug, category_id, description, image,
				price, dimension, weight, stock, is_active, is_featured,
				created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5, $6,
				$7, $8, $9, $10, $11, $12,
				NOW(), NOW()
			)
			RETURNING id
		`
		err := tx.QueryRow(
			insertProductQuery,
			item.ID, item.Name, item.Slug, catIDPtr, item.Description, uploadImg,
			price, dim, weight, 100, true, isFeatured,
		).Scan(&productID)
		if err != nil {
			return fmt.Errorf("error inserting product %s: %w", item.Name, err)
		}

		// Insert product galleries
		for idx, gImg := range item.Gallery {
			uploadGImg := toUploadURL(baseURL, gImg)
			altText := fmt.Sprintf("%s - Galeri Foto %d Roster Beton Purwakarta", item.Name, idx+1)
			insertGalleryQuery := `
				INSERT INTO product_galleries (product_id, image_url, alt_text, sort_order, created_at)
				VALUES ($1, $2, $3, $4, NOW())
			`
			if _, err := tx.Exec(insertGalleryQuery, productID, uploadGImg, altText, idx+1); err != nil {
				return fmt.Errorf("error inserting product gallery for product %d: %w", productID, err)
			}
		}

		// Insert Polymorphic SEO record for product
		insertSEOQuery := `
			INSERT INTO seos (
				reference_type, reference_id, meta_title, meta_description, meta_keywords,
				canonical_url, og_title, og_description, og_image, og_type,
				created_at, updated_at
			) VALUES (
				'product', $1, $2, $3, $4,
				$5, $6, $7, $8, 'website',
				NOW(), NOW()
			)
		`
		if _, err := tx.Exec(insertSEOQuery, productID, metaTitle, metaDesc, metaKeywords, canonicalURL, ogTitle, ogDesc, ogImage); err != nil {
			return fmt.Errorf("error inserting seo for product %d: %w", productID, err)
		}
	}

	// Reset sequence
	if _, err := tx.Exec("SELECT setval('products_id_seq', (SELECT MAX(id) FROM products));"); err != nil {
		log.Printf("⚠️  Warning setting sequence products_id_seq: %v", err)
	}
	if _, err := tx.Exec("SELECT setval('product_galleries_id_seq', (SELECT COALESCE(MAX(id), 1) FROM product_galleries));"); err != nil {
		log.Printf("⚠️  Warning setting sequence product_galleries_id_seq: %v", err)
	}

	log.Printf("✅ %d Products and Galleries seeded successfully with full URLs", len(rawProducts))
	return nil
}

func seedGalleries(tx *sqlx.Tx, baseURL string) error {
	log.Println("🖼️  Seeding General Galleries (Projects, Factory, Products)...")

	if _, err := tx.Exec("DELETE FROM galleries"); err != nil {
		return fmt.Errorf("failed to clear galleries: %w", err)
	}

	rawGalleries := []GallerySeed{
		{Src: "/images/products/galprod-17.jpeg", Title: "Roster Beton Motif", Category: "Produk", SortOrder: 1},
		{Src: "/images/products/galprod-18.jpeg", Title: "Produk Unggulan", Category: "Produk", SortOrder: 2},
		{Src: "/images/products/galprod-19.jpeg", Title: "Proyek Pemasangan", Category: "Proyek", SortOrder: 3},
		{Src: "/images/products/galprod-20.jpeg", Title: "Motif Kotak", Category: "Produk", SortOrder: 4},
		{Src: "/images/products/galprod-22.jpeg", Title: "Motif Bunga", Category: "Produk", SortOrder: 5},
		{Src: "/images/products/galprod-23.jpeg", Title: "Roster Minimalis", Category: "Produk", SortOrder: 6},
		{Src: "/images/products/galprod-24.jpeg", Title: "Roster Modern", Category: "Produk", SortOrder: 7},
		{Src: "/images/products/galprod-25.jpeg", Title: "Roster Klasik", Category: "Produk", SortOrder: 8},
		{Src: "/images/products/galprod-26.jpeg", Title: "Roster Premium", Category: "Produk", SortOrder: 9},
		{Src: "/images/products/galprod-27.jpeg", Title: "Detail Produk", Category: "Produk", SortOrder: 10},
		{Src: "/images/products/galprod-28.jpeg", Title: "Proyek Rumah", Category: "Proyek", SortOrder: 11},
		{Src: "/images/products/galprod-29.jpeg", Title: "Pemasangan Roster", Category: "Proyek", SortOrder: 12},
		{Src: "/images/products/galprod-31.jpeg", Title: "Koleksi Motif", Category: "Produk", SortOrder: 13},
		{Src: "/images/products/galprod-32.jpeg", Title: "Produksi Roster", Category: "Produksi", SortOrder: 14},
		{Src: "/images/products/galprod-33.jpeg", Title: "Quality Control", Category: "Produksi", SortOrder: 15},
		{Src: "/images/products/galprod-34.jpeg", Title: "Pengiriman", Category: "Layanan", SortOrder: 16},
		{Src: "/images/products/galprod-35.jpeg", Title: "Proyek Komersial", Category: "Proyek", SortOrder: 17},
		{Src: "/images/products/galprod-36.jpeg", Title: "Detail Pemasangan", Category: "Proyek", SortOrder: 18},
		{Src: "/images/products/galprod-37.jpeg", Title: "Motif Custom", Category: "Produk", SortOrder: 19},
		{Src: "/images/products/galprod-38.jpeg", Title: "Proyek Skala Besar", Category: "Proyek", SortOrder: 20},
		{Src: "/images/products/galprod-39.jpeg", Title: "Hasil Akhir", Category: "Proyek", SortOrder: 21},
		{Src: "/images/products/galprod-40.jpeg", Title: "Konsultasi", Category: "Layanan", SortOrder: 22},
		{Src: "/images/products/galprod-44.jpeg", Title: "Produk Terbaru", Category: "Produk", SortOrder: 23},
		{Src: "/images/products/galprod-47.jpeg", Title: "Motif Geometri", Category: "Produk", SortOrder: 24},
		{Src: "/images/products/galprod-48.jpeg", Title: "Proyek Perumahan", Category: "Proyek", SortOrder: 25},
		{Src: "/images/products/galprod-50.jpeg", Title: "Motif Unik", Category: "Produk", SortOrder: 26},
		{Src: "/images/products/galprod-51.jpeg", Title: "Produk Finishing", Category: "Produksi", SortOrder: 27},
		{Src: "/images/products/galprod-53.jpeg", Title: "Pilihan Warna", Category: "Produk", SortOrder: 28},
		{Src: "/images/products/galprod-54.jpeg", Title: "Proyek Villa", Category: "Proyek", SortOrder: 29},
		{Src: "/images/products/galprod-55.jpeg", Title: "Tim Kami", Category: "Perusahaan", SortOrder: 30},
		{Src: "/images/products/galprod-56.jpeg", Title: "Lokasi Pabrik", Category: "Perusahaan", SortOrder: 31},
		{Src: "/images/products/galprod-57.jpeg", Title: "Proses Produksi", Category: "Produksi", SortOrder: 32},
		{Src: "/images/products/galprod-58.jpeg", Title: "Produk Siap Kirim", Category: "Layanan", SortOrder: 33},
		{Src: "/images/products/galprod-59.jpeg", Title: "Testimoni Pelanggan", Category: "Perusahaan", SortOrder: 34},
	}

	for _, g := range rawGalleries {
		uploadSrc := toUploadURL(baseURL, g.Src)
		altText := fmt.Sprintf("%s - %s CV Roster Purwakarta", g.Title, g.Category)
		query := `
			INSERT INTO galleries (id, src, title, category, alt_text, sort_order, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		`
		if _, err := tx.Exec(query, g.SortOrder, uploadSrc, g.Title, g.Category, altText, g.SortOrder); err != nil {
			return fmt.Errorf("error inserting gallery %s: %w", g.Title, err)
		}
	}

	if _, err := tx.Exec("SELECT setval('galleries_id_seq', (SELECT MAX(id) FROM galleries));"); err != nil {
		log.Printf("⚠️  Warning setting sequence galleries_id_seq: %v", err)
	}

	log.Printf("✅ %d General Galleries seeded successfully with full URLs", len(rawGalleries))
	return nil
}

func seedConfigs(tx *sqlx.Tx) error {
	log.Println("⚙️  Creating 'configs' table & seeding configuration records...")

	createTableQuery := `
		CREATE TABLE IF NOT EXISTS configs (
			id BIGSERIAL PRIMARY KEY,
			key VARCHAR(100) UNIQUE NOT NULL,
			value TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
		CREATE INDEX IF NOT EXISTS idx_configs_key ON configs(key);
	`
	if _, err := tx.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to ensure configs table exists: %w", err)
	}

	configs := []struct {
		Key   string
		Value string
	}{
		{
			Key:   "company_name",
			Value: "CV Roster Purwakarta",
		},
		{
			Key:   "company_tagline",
			Value: "Penyedia roster beton berkualitas tinggi di Plered, Purwakarta, Jawa Barat.",
		},
		{
			Key:   "email",
			Value: "info@rosterbetonpurwakarta.com",
		},
		{
			Key:   "phone",
			Value: "0812-3456-7890",
		},
		{
			Key:   "whatsapp",
			Value: "081234567890",
		},
		{
			Key:   "address",
			Value: "Jl. Raya Plered No. 45, Kecamatan Plered, Kabupaten Purwakarta, Jawa Barat 41161",
		},
		{
			Key:   "operational_weekday",
			Value: "Senin – Jumat: 08.00 – 17.00 WIB",
		},
		{
			Key:   "operational_saturday",
			Value: "Sabtu: 08.00 – 15.00 WIB",
		},
		{
			Key:   "operational_sunday",
			Value: "Minggu & Hari Libur: WhatsApp Tetap Aktif",
		},
		{
			Key:   "google_maps_url",
			Value: "https://www.google.com/maps/embed?pb=!1m18!1m12!1m3!1d3963.278534430604!2d107.4491153749726!3d-6.511323763642216!2m3!1f0!2f0!3f0!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x2e69123456789abcd!2sPlered%2C%20Purwakarta%20Regency%2C%20West%20Java!5e0!3m2!1sen!2sid!4v1234567890123!5m2!1sen!2sid",
		},
		{
			Key:   "google_maps_direction_url",
			Value: "https://maps.google.com/?q=Plered+Purwakarta",
		},
		{
			Key:   "meta_description",
			Value: "Pabrik dan supplier roster beton minimalis berkualitas di Plered Purwakarta. Beragam motif elegan untuk rumah & gedung.",
		},
	}

	for _, cfg := range configs {
		query := `
			INSERT INTO configs (key, value, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
			ON CONFLICT (key) DO UPDATE 
			SET value = EXCLUDED.value, updated_at = NOW()
		`
		if _, err := tx.Exec(query, cfg.Key, cfg.Value); err != nil {
			return fmt.Errorf("error inserting config %s: %w", cfg.Key, err)
		}
	}

	log.Printf("✅ %d Configurations seeded successfully", len(configs))
	return nil
}

// Ensure database/sql package is utilized
var _ = sql.ErrNoRows
