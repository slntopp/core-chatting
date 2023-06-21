package graph

import (
	"context"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"go.uber.org/zap"
)

func ConnectDb(logger *zap.Logger, host, cred string) driver.Database {
	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://" + cred + "@" + host},
	})
	if err != nil {
		logger.Fatal("Error creating connection to DB", zap.Error(err))
	}

	logger.Debug("Instantiated DB connection", zap.Any("conn", conn))

	logger.Info("Setting up DB client")
	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		logger.Fatal("Error creating driver instance for DB", zap.Error(err))
	}
	logger.Debug("Instantiated DB client", zap.Any("client", c))

	db_connect_attempts := 0
db_connect:
	logger.Info("Trying to connect to DB")
	db, err := c.Database(context.TODO(), DB_NAME)
	if e, isArangoError := driver.AsArangoError(err); isArangoError && e.ErrorMessage == "database not found" {
		logger.Info("DB not found, creating it")
		db, err = c.CreateDatabase(context.TODO(), DB_NAME, nil)
		if err != nil {
			logger.Fatal("Error creating DB", zap.Error(err))
		}
		logger.Info("DB created")
		goto db_connect
	} else if err != nil {
		db_connect_attempts++
		logger.Error("Failed to connect DB", zap.Error(err), zap.Int("attempts", db_connect_attempts), zap.Int("next_attempt", db_connect_attempts*5))
		time.Sleep(time.Duration(db_connect_attempts*5) * time.Second)
		goto db_connect
	}

	return db
}

func GraphGetEnsure(logger *zap.Logger, ctx context.Context, db driver.Database, name string) driver.Graph {
	exists, err := db.GraphExists(ctx, name)
	if err != nil {
		logger.Fatal("Error checking graph existence", zap.Error(err))
	}
	if !exists {
		logger.Info("Creating permissions graph")
		graph, err := db.CreateGraphV2(ctx, name, nil)
		if err != nil {
			logger.Fatal("Error creating permissions graph", zap.Error(err))
		}
		return graph
	}

	graph, err := db.Graph(ctx, name)
	if err != nil {
		logger.Fatal("Error getting graph", zap.Error(err))
	}
	return graph
}

func GraphGetVertexEnsure(logger *zap.Logger, ctx context.Context, db driver.Database, graph driver.Graph, name string) (col driver.Collection) {
	exists, err := db.CollectionExists(ctx, name)
	if err != nil {
		logger.Fatal("Error checking if collection exists", zap.Error(err))
	}
	if !exists {
		col, err = db.CreateCollection(ctx, name, &driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{AllowUserKeys: true, Type: "traditional"},
		})
		if err != nil {
			logger.Fatal("Error creating collection", zap.Error(err))
		}
		return col
	}
	col, err = db.Collection(ctx, name)
	if err != nil {
		logger.Fatal("Error getting collection", zap.Error(err))
	}
	return col
}

func GraphGetEdgeEnsure(logger *zap.Logger, ctx context.Context, graph driver.Graph, name, from, to string) driver.Collection {
	exists, err := graph.EdgeCollectionExists(ctx, name)
	if err != nil {
		logger.Fatal("Error checking if edge collection exists", zap.Error(err))
	}
	if !exists {
		col, err := graph.CreateEdgeCollection(ctx, name, driver.VertexConstraints{
			From: []string{from}, To: []string{to},
		})
		if err != nil {
			logger.Fatal("Error creating collection", zap.Error(err))
		}
		return col
	}

	err = graph.SetVertexConstraints(ctx, name, driver.VertexConstraints{
		From: []string{from}, To: []string{to},
	})
	if err != nil {
		logger.Fatal("Error setting vertex constraints for edge collection",
			zap.String("col", name), zap.String("from", from),
			zap.String("to", to), zap.Error(err),
		)
	}

	col, _, err := graph.EdgeCollection(ctx, name)
	if err != nil {
		logger.Fatal("Error getting collection", zap.Error(err))
	}
	return col
}

func GetEnsureCollection(log *zap.Logger, ctx context.Context, db driver.Database, name string) driver.Collection {
	exists, err := db.CollectionExists(ctx, name)
	if err != nil {
		log.Fatal("Error checking if collection exists", zap.Error(err))
	}
	if !exists {
		col, err := db.CreateCollection(ctx, name, &driver.CreateCollectionOptions{
			KeyOptions: &driver.CollectionKeyOptions{AllowUserKeys: true, Type: "uuid"},
		})
		if err != nil {
			log.Fatal("Error creating collection", zap.Error(err))
		}
		return col
	}

	col, err := db.Collection(ctx, name)
	if err != nil {
		log.Fatal("Error getting collection", zap.Error(err))
	}
	return col
}