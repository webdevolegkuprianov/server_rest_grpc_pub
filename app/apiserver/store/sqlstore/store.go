package sqlstore

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/store"
	pb "github.com/webdevolegkuprianov/server_rest_grpc/proto"
)

//Store
type Store struct {
	dbPostgres        *pgxpool.Pool
	clientStreamGrpc1 pb.Stream_Stream1Client
	clientStreamGrpc2 pb.Stream_Stream2Client
	clientStreamGrpc3 pb.Stream_Stream3Client
	userRepository    *UserRepository
	dataRepository    *DataRepository
}

//New_db
func New(db *pgxpool.Pool,
	stream1 pb.Stream_Stream1Client,
	stream2 pb.Stream_Stream2Client,
	stream3 pb.Stream_Stream3Client,
) *Store {
	return &Store{
		dbPostgres:        db,
		clientStreamGrpc1: stream1,
		clientStreamGrpc2: stream2,
		clientStreamGrpc3: stream3,
	}
}

//User
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}

//Data
func (s *Store) Data() store.DataRepository {
	if s.dataRepository != nil {
		return s.dataRepository
	}

	s.dataRepository = &DataRepository{
		store: s,
	}

	return s.dataRepository
}
