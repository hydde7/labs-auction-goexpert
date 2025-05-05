package auction

import (
	"context"
	"os"
	"time"

	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{
		Collection: database.Collection("auctions"),
	}
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction,
) *internal_error.InternalError {
	intervalStr := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(intervalStr)
	if err != nil || duration <= 0 {
		logger.Error("Invalid AUCTION_INTERVAL, defaulting to 1m", err)
		duration = time.Minute
	}

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}
	_, err = ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	go func(id string) {
		timer := time.NewTimer(duration)
		<-timer.C

		filter := bson.M{"_id": id, "status": auction_entity.Active}
		update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
		if _, err := ar.Collection.UpdateOne(context.Background(), filter, update); err != nil {
			logger.Error("Error trying to close auction", err)
		}
	}(auctionEntity.Id)

	return nil
}
