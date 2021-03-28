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

func (s *server) RemoveClient(ctx context.Context, in *pb.Id) (*pb.Empty, error) {
	id := uint(in.GetId())
	res:= s.db.WithContext(ctx).Delete(&models.Client{},id)
	return  &pb.Empty{},res.Error
}

func (s *server) ListClients(in *pb.IsInvestor,stream pb.ArexServices_ListClientsServer) error {
	var toRet []models.Client
	res:= s.db.Where("is_investor = ?", in.GetIsInvestor()).Find(&toRet)
	if res.Error != nil {
		return res.Error
	}
	for _,e:=range toRet {
		stream.Send(e.CastgRPC())
	}
	return nil
}

func (s *server) StartInvoiceFinancing(ctx context.Context, in *pb.InvoiceFinancing) (*pb.Id, error) {
	invoice		:= models.CastInvoice	(in.GetInvoice()	)
	sellOrder	:= models.CastSellOrder	(in.GetSellOrder()	)

	e := s.db.Transaction(func(tx *gorm.DB) error {
		var r *gorm.DB

		ErrorTrack :		// handling of errors.
		if r != nil { return r.Error }	// rolling back the transaction

		r = tx.WithContext(ctx).Create(invoice	)
		sellOrder.InvoiceID = invoice.ID
		r = tx.WithContext(ctx).Create(sellOrder)

		return nil
	})
	return  &pb.Id{Id:uint64(sellOrder.ID)},e
}
// macro expansion : the function is exapanded as another function in the .go files
//<<func (s *server) ListBids		(in *pb.Empty,stream pb.ArexServices_ListBidsServer ) error {>>
//<<func (s *server) ListInvoices	(in *pb.Empty,stream pb.ArexServices_ListInvoicesServer ) error {>>
func (s *server) ListSellOrders	(in *pb.Empty,stream pb.ArexServices_ListSellOrdersServer	) error {
//<<var toRet []models.Ledger>>
//<<var toRet []models.Invoice>>
	var toRet []models.SellOrder
	res:= s.db.Find(&toRet)
	if res.Error != nil {
		return res.Error
	}
	for _,e:=range toRet {
		stream.Send(e.CastgRPC())
	}
	return nil
}
//<<end>>

type ID struct {
	ID		uint
}
func (s *server) AddBid(ctx context.Context, in *pb.Ledger) (*pb.Id, error) {
	var r ID
	l := models.CastLedger(in)
	res := s.db.Raw("select add_bid(?,?,?,?)",l.InvestorID,l.SellOrderID,l.Size,l.Amount).Scan(&r)

	if res.Error != nil {
		return nil,res.Error
	}

	return  &pb.Id{Id:uint64(r.ID)},nil
}


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
