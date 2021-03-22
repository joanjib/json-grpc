package main

import (
	"os"
    "context"
    "log"
    "net"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
    pb "num/arexservices"

    "num/db"
    "num/models"
    "gorm.io/gorm"


)

var initTypes bool

const (
    port = ":50051"
)
type server struct {
    pb.UnimplementedArexServicesServer
	db		*gorm.DB
}

func (s *server) AddClient(ctx context.Context, in *pb.Client) (*pb.Id, error) {
	c := models.CastClient(in)
	res:= s.db.WithContext(ctx).Create(c)
	return  &pb.Id{Id:uint64(c.ID)},res.Error
}

//func (s *server) ListIssuers(in *pb.Empty,stream pb.ArexServices_ListIssuersServer) error {
//	return issuer.ListIssuers(stream,s.db)
//}
func newServer() *server {
	// BEGIN -  database instanciation
    if db.Db == nil  {
        panic("not initialitzed Db object")
    }
	db := db.Db
	// END   -  database instanciation

	// BEGIN - database objects initialization
    var err error
	var typesDomains []byte
	if initTypes {
		typesDomains,err = ioutil.ReadFile("../sql/types-domains.sql")
		if err != nil {
			panic(err)
		}
	}

    var r *gorm.DB

    ErrorTrack :
    if r != nil {
		panic (r)
    }
    ErrorTrack_err :
    if err != nil {
		panic (err)
    }

	if initTypes {
		r = db.Exec  (string(typesDomains) )
	}
    // migrations:
    err = db.AutoMigrate(&models.Client{}   )
    err = db.AutoMigrate(&models.Invoice{}  )
    err = db.AutoMigrate(&models.SellOrder{})
    err = db.AutoMigrate(&models.Ledger{}   )

	// END   - database objects initialization

	return &server{db : db}
}



func main () {
	argsWithoutProg := os.Args[1:]


	switch len(argsWithoutProg) {
	case 1	: initTypes = true
	case 0  : initTypes = false
	default : panic("Too much arguments when init the server")
	}

    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterArexServicesServer(s, newServer())
	reflection.Register(s)
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
