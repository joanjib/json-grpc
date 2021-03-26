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
	id, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "1",Name: "J1",Surname:"I1",Balance:"1000",IsInvestor:false})
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

    assert.Equal(t,"2"			,issuersList[1].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[1].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[1].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,issuersList[1].GetBalance(),			"Balance numeral is not 1000"		)
    assert.Equal(t,false	    ,issuersList[1].GetIsInvestor(),		"Is investor is not false"			)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err,"Error listint the two issuers inserted")

	investorsList := []pb.Client{}

    for {
        investor, err := stream.Recv()
        if err == io.EOF {
            break
        }
        assert.Nil(t,err,"Error receiving  investor")
        investorsList = append(investorsList,*investor)
    }
	assert.Equal(t,"2"			,investorsList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"  )
    assert.Equal(t,"J2"			,investorsList[0].GetName(),              "Name not equal to J1"		    )
    assert.Equal(t,"I2"			,investorsList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,investorsList[0].GetBalance(),			"Balance numeral is not 2000"		)
    assert.Equal(t,true		    ,investorsList[0].GetIsInvestor(),		"Is investor is not true "			)

	// removing the first client
	_, err = c.RemoveClient(ctx, id)
	assert.Nil(t,err,"Error in addition of issuer")

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the issuer")

    issuersList = []pb.Client{}

    for {
        issuer, err := stream.Recv()
        if err == io.EOF {
            break
        }
        assert.Nil(t,err,"Error receiving  issuer")
        issuersList = append(issuersList,*issuer)
    }

    assert.Equal(t,"2"			,issuersList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[0].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,issuersList[0].GetBalance(),			"Balance numeral is not 1000"		)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)


}


