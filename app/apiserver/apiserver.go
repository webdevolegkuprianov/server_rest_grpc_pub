package apiserver

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/model"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/store/sqlstore"

	logger "github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/logger"
	pb "github.com/webdevolegkuprianov/server_rest_grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Start(config *model.Service) error {

	dbPostgres, err := newDbPostgres(config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	defer dbPostgres.Close()

	//grpc
	//conn
	newGrpc, err := newConn(config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	stream1, err := newStream1(newGrpc, config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	stream2, err := newStream2(newGrpc, config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	stream3, err := newStream3(newGrpc, config)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	store := sqlstore.New(dbPostgres, stream1, stream2, stream3)

	server := newServer(store, config)

	//setup HTTPS server
	srv := &http.Server{
		Addr:    config.Spec.Ports.RestServerGrpcClient.BindAddr,
		Handler: server.router,
	}

	return srv.ListenAndServe()
}

//start grpc
//new conn
func newConn(config *model.Service) (*grpc.ClientConn, error) {

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	}

	conn, err := grpc.Dial(config.Spec.Ports.GrpcServer.BindAddr, opts...)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return conn, nil

}

//stream1
func newStream1(conn *grpc.ClientConn, config *model.Service) (pb.Stream_Stream1Client, error) {

	// create stream
	client := pb.NewStreamClient(conn)
	stream, err := client.Stream1(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	return stream, nil
}

//stream2
func newStream2(conn *grpc.ClientConn, config *model.Service) (pb.Stream_Stream2Client, error) {

	// create stream
	client := pb.NewStreamClient(conn)
	stream, err := client.Stream2(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	return stream, nil
}

//stream3
func newStream3(conn *grpc.ClientConn, config *model.Service) (pb.Stream_Stream3Client, error) {

	// create stream
	client := pb.NewStreamClient(conn)
	stream, err := client.Stream3(context.Background())
	if err != nil {
		logger.ErrorLogger.Println(err)
	}

	return stream, nil
}

//connect to postgres
func newDbPostgres(conf *model.Service) (*pgxpool.Pool, error) {

	config, _ := pgx.ParseConfig("")
	config.Host = conf.Spec.DB.Host
	config.Port = conf.Spec.DB.Port
	config.User = conf.Spec.DB.User
	config.Password = conf.Spec.DB.Password
	config.Database = conf.Spec.DB.Database
	config.LogLevel = pgx.LogLevelDebug
	config.Logger = logrusadapter.NewLogger(logger.PgLog())
	config.TLSConfig = nil

	poolConfig, _ := pgxpool.ParseConfig("")
	poolConfig.ConnConfig = config
	poolConfig.MaxConnLifetime = time.Duration(conf.Spec.DB.MaxConnLifetime) * time.Minute
	poolConfig.MaxConnIdleTime = time.Duration(conf.Spec.DB.MaxConnIdletime) * time.Minute
	poolConfig.MaxConns = conf.Spec.DB.MaxConns
	poolConfig.MinConns = conf.Spec.DB.MinConns
	poolConfig.HealthCheckPeriod = time.Duration(conf.Spec.DB.HealthCheckPeriod) * time.Minute

	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		logger.ErrorLogger.Println(err)
		return nil, err
	}

	return conn, nil
}
