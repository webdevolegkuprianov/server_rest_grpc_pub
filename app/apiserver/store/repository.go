package store

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/model"
)

//user repository
type UserRepository interface {
	//auth methods
	FindUser(string, string) (*model.User1, error)
	FindUserid(uint64) error
	//jwt methods
	CreateToken(uint64, *model.Service) (string, time.Time, error)
	ExtractTokenMetadata(*http.Request, *model.Service) (*model.AccessDetails, error)
	VerifyToken(*http.Request, *model.Service) (*jwt.Token, error)
	ExtractToken(*http.Request) string
}

//data repository
type DataRepository interface {
	SendStreamMessages1([]byte) error
	SendStreamMessages2([]byte) error
	SendStreamMessages3([]byte) error
}
