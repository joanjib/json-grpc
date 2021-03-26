package main
// issuer remote tests file
import (
	"io"
	"testing"
	"github.com/stretchr/testify/assert"

	pb "num/arexservices"
    "num/utils"
)

func TestCRUDClients (t *testing.T) {
	conn,c,ctx,cancel := utils.InitClient()
    defer conn.Close()
    defer cancel()
	// client
	_, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "1",Name: "J1",Surname:"I1",Balance:"1000",IsInvestor:false})
	assert.Nil(t,err,"Error in addition of issuer")

    //id := r.GetId()     // saving the id of the new issuer    

	// client
	_, err = c.AddClient(ctx, &pb.Client{FiscalIdentity: "2",Name: "J2",Surname:"I2",Balance:"2000",IsInvestor:false})
	assert.Nil(t,err,"Error in addition of issuer")

	// investor
	_, err = c.AddClient(ctx, &pb.Client{FiscalIdentity: "2",Name: "J2",Surname:"I2",Balance:"2000",IsInvestor:true})
	assert.Nil(t,err,"Error in addition of issuer")

    //id = r.GetId()     // saving the id of the new issuer    

	stream, err := c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the two issuers inserted")

    issuersList := []pb.Client{}

    for {
        issuer, err := stream.Recv()
        if err == io.EOF {
            break
        }
        assert.Nil(t,err,"Error receiving  issuer")
        issuersList = append(issuersList,*issuer)
    }


    assert.Equal(t,"1"			,issuersList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 1"    )
    assert.Equal(t,"J1"			,issuersList[0].GetName(),              "Name not equal to Joan1"           )
    assert.Equal(t,"I1"			,issuersList[0].GetSurname(),           "Surname not equal to Iglesias1"    )
    assert.Equal(t,"1000.00"    ,issuersList[0].GetBalance(),			"Balance numeral is not 1000"		)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)

}


