package main
// issuer remote tests file
import (
	"io"
	"testing"
	"github.com/stretchr/testify/assert"

	pb "num/arexservices"
    "num/utils"
)

// macro expansions:
//<<func generateListSellOrders	(t *testing.T,stream pb.ArexServices_ListSellOrdersClient	) []pb.SellOrder{>>
//<<func generateListInvoices 	(t *testing.T,stream pb.ArexServices_ListInvoicesClient		) []pb.Invoice	{>>
func	 generateListClients	(t *testing.T,stream pb.ArexServices_ListClientsClient		) []pb.Client	{

//<<var list []pb.SellOrder>>
//<<var list []pb.Invoice>>
	var list []pb.Client

    for {
        e, err := stream.Recv()
        if err == io.EOF {
            break
        }
        assert.Nil(t,err,"Error receiving  issuer")
        list = append(list,*e)
    }
	return list
}
//<<end>>

func TestCRUDClients (t *testing.T) {
	conn,c,ctx,cancel := utils.InitClient()
    defer conn.Close()
    defer cancel()
	// client
	idToRm, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "1",Name: "J1",Surname:"I1",Balance:"1000",IsInvestor:false})
	assert.Nil(t,err,"Error in addition of issuer")

    //id := r.GetId()     // saving the id of the new issuer    

	// client
	id, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "2",Name: "J2",Surname:"I2",Balance:"2000",IsInvestor:false})
	assert.Nil(t,err,"Error in addition of issuer")

	// investor
	_, err = c.AddClient(ctx, &pb.Client{FiscalIdentity: "2",Name: "J2",Surname:"I2",Balance:"2000",IsInvestor:true})
	assert.Nil(t,err,"Error in addition of issuer")

    //id = r.GetId()     // saving the id of the new issuer    

	stream, err := c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the two issuers inserted")

	issuersList := generateListClients(t,stream)

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

	investorsList := generateListClients(t,stream)

	assert.Equal(t,"2"			,investorsList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"  )
    assert.Equal(t,"J2"			,investorsList[0].GetName(),              "Name not equal to J1"		    )
    assert.Equal(t,"I2"			,investorsList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,investorsList[0].GetBalance(),			"Balance numeral is not 2000"		)
    assert.Equal(t,true		    ,investorsList[0].GetIsInvestor(),		"Is investor is not true "			)

	// removing the first client
	_, err = c.RemoveClient(ctx, idToRm)
	assert.Nil(t,err,"Error in addition of issuer")

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the issuer")

	issuersList = generateListClients(t,stream )

    assert.Equal(t,"2"			,issuersList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[0].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,issuersList[0].GetBalance(),			"Balance numeral is not 1000"		)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)

	//Invoice financing process testing:

	// invoice to be financed of 250 €.
	invoice				:= &pb.Invoice{ClientId:id.Id,Amount:"250"}
	sellOrder			:= &pb.SellOrder{Size:"250",Amount:"200"}
	invoiceFinancing	:= &pb.InvoiceFinancing{SellOrder:sellOrder,Invoice:invoice}

	_,err = c.StartInvoiceFinancing(ctx,invoiceFinancing)
    assert.Nil(t,err,"Error starting the financing process")

	streamInv,err := c.ListInvoices(ctx,&pb.Empty{})
    assert.Nil(t,err)
	invoicesList := generateListInvoices(t,streamInv )

	assert.Equal(t,1					,len(invoicesList)				)
	assert.Equal(t,id.GetId()			,invoicesList[0].GetClientId()	)
    assert.Equal(t,"250.00"				,invoicesList[0].GetAmount()	)
    assert.Equal(t,"financing search"	,invoicesList[0].GetState()		)

	streamSO,err := c.ListSellOrders(ctx,&pb.Empty{})
    assert.Nil(t,err)
	soList := generateListSellOrders(t,streamSO )

	assert.Equal(t,1						,len(soList)				)
    assert.Equal(t,invoicesList[0].GetId()	,soList[0].GetInvoiceId()	)
    assert.Equal(t,"250.00"					,soList[0].GetSize()		)
    assert.Equal(t,"200.00"					,soList[0].GetAmount()		)
    assert.Equal(t,"ongoing"				,soList[0].GetState()		)
}


