package auction

import (
	"context"
	"os"
	"testing"
	"time"

	"fullcycle-auction_go/configuration/database/mongodb"
	"fullcycle-auction_go/internal/entity/auction_entity"
)

func TestAutoCloseAuction(t *testing.T) {
	ctx := context.Background()

	db, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		t.Fatalf("could not connect to mongodb: %v", err)
	}
	if err := db.Collection("auctions").Drop(ctx); err != nil {
		t.Fatalf("could not drop auctions collection: %v", err)
	}
	os.Setenv("AUCTION_INTERVAL", "1s")

	repo := NewAuctionRepository(db)
	a, internalErr := auction_entity.CreateAuction("TestProduct", "TestCat", "Desc", auction_entity.New)
	if internalErr != nil {
		t.Fatalf("could not create auction entity: %v", internalErr)
	}
	if err := repo.CreateAuction(ctx, a); err != nil {
		t.Fatalf("CreateAuction failed: %v", err)
	}

	got, internalErr := repo.FindAuctionById(ctx, a.Id)
	if internalErr != nil {
		t.Fatalf("FindAuctionById failed: %v", internalErr)
	}
	if got.Status != auction_entity.Active {
		t.Errorf("expected Active, got %v", got.Status)
	}
	time.Sleep(2 * time.Second)

	got2, internalErr := repo.FindAuctionById(ctx, a.Id)
	if internalErr != nil {
		t.Fatalf("FindAuctionById failed: %v", internalErr)
	}
	if got2.Status != auction_entity.Completed {
		t.Errorf("expected Completed after auto-close, got %v", got2.Status)
	}
}
