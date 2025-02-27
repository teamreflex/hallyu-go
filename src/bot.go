package bot

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/gorm"
)

func Start() {
    fmt.Println("Starting bot...")

    env, err := GetEnv()
    if err != nil {
        fmt.Printf("Failed to get environment variables: %s\n", err)
        os.Exit(1)
    }

    db := GetDatabase()
    db.AutoMigrate(&Product{})

    ticker := time.NewTicker(time.Second * time.Duration(env.LOOP_INTERVAL))
    defer ticker.Stop()

    client := http.Client{
        Timeout: time.Second * 2,
    }

    go func() {
        for range ticker.C {
             _, err := HandleProducts(db, &client)
             if err != nil {
                 fmt.Println(err.Error())
             }
        }
    }()

    select {}
}

var searchTerms = []string{"artms", "polaroid"}
func HandleProducts(db *gorm.DB, client *http.Client) (bool, error) {
    // fetch products from remote
    products, err := GetProducts(client)
    if err != nil {
        return false, err
    }

    // filter new products by search terms
    filteredProducts := make([]RawProduct, 0)
    for _, product := range products {
        for _, searchTerm := range searchTerms {
            if strings.Contains(strings.ToLower((product.Title)), searchTerm) {
                filteredProducts = append(filteredProducts, product)
                break
            }
        }
    }

    if (len(filteredProducts) == 0) {
        return true, nil
    }
    
    // get ids of remote products
    remoteIds := make([]uint, len(filteredProducts))
    for i, product := range filteredProducts {
        remoteIds[i] = product.ID
    }

    // get products from database
    var databaseProducts []Product
    result := db.Model(&Product{}).Where("remote_id IN ?", remoteIds).Order("created_at DESC").Limit(50).Find(&databaseProducts)
    if result.Error != nil {
        fmt.Printf("Failed to get products from database: %s\n", result.Error)
        databaseProducts = make([]Product, 0)
    }

    // get ids of database products
    databaseIds := make([]uint, len(databaseProducts))
    for i, product := range databaseProducts {
        databaseIds[i] = product.RemoteID
    }

    // get new products that aren't in the database
    newProducts := make([]RawProduct, 0, len(filteredProducts))
    for _, product := range filteredProducts {
        if !slices.Contains(databaseIds, product.ID) {
            newProducts = append(newProducts, product)
        }
    }

    if (len(newProducts) == 0) {
        return true, nil
    }

    // post new products to discord and write to database
    for _, rawProduct := range newProducts {
        fmt.Printf("Posting product: %s\n", rawProduct.Title)
        newProduct := Product{
            RemoteID: rawProduct.ID,
            Title: rawProduct.Title,
            Handle: rawProduct.Handle,
        }

        // post to discord
        _, err := PostToDiscord(rawProduct, client)
        if err != nil {
            fmt.Println(err.Error())
        }
        
        // write to database
        db.Create(&newProduct)
    }

    return true, nil
}