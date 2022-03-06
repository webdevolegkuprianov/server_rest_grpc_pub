package sqlstore

import (
	logger "github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/logger"
	pb "github.com/webdevolegkuprianov/server_rest_grpc/proto"
)

//Data repository
type DataRepository struct {
	store *Store
}

//send message in stream1
func (r *DataRepository) SendStreamMessages1(m []byte) error {

	bodySend := pb.Request1{
		Message: m,
	}

	if err := r.store.clientStreamGrpc1.Send(&bodySend); err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//send message in stream2
func (r *DataRepository) SendStreamMessages2(m []byte) error {

	bodySend := pb.Request2{
		Message: m,
	}

	if err := r.store.clientStreamGrpc2.Send(&bodySend); err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}

//send message in stream2
func (r *DataRepository) SendStreamMessages3(m []byte) error {

	bodySend := pb.Request3{
		Message: m,
	}

	if err := r.store.clientStreamGrpc3.Send(&bodySend); err != nil {
		logger.ErrorLogger.Println(err)
		return err
	}

	return nil

}
