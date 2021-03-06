package main
// issuer remote tests file
import (
	"io"
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"

	pb "num/arexservices"
    "num/utils"
)

// macro expansions:
//<<func generateListBids		(t *testing.T,stream pb.ArexServices_ListBidsClient			) []pb.Ledger	{>>
//<<func generateListSellOrders	(t *testing.T,stream pb.ArexServices_ListSellOrdersClient	) []pb.SellOrder{>>
//<<func generateListInvoices 	(t *testing.T,stream pb.ArexServices_ListInvoicesClient		) []pb.Invoice	{>>
func	 generateListClients	(t *testing.T,stream pb.ArexServices_ListClientsClient		) []pb.Client	{

//<<var list []pb.Ledger>>
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
    assert.Equal(t,"1000.00"    ,issuersList[0].GetBalance(),			"Balance is not 1000"				)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)

    assert.Equal(t,"2"			,issuersList[1].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[1].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[1].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,issuersList[1].GetBalance(),			"Balance is not 2000"				)
    assert.Equal(t,false	    ,issuersList[1].GetIsInvestor(),		"Is investor is not false"			)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err,"Error listint the two issuers inserted")

	investorsList := generateListClients(t,stream)

	assert.Equal(t,"2"			,investorsList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"  )
    assert.Equal(t,"J2"			,investorsList[0].GetName(),              "Name not equal to J1"		    )
    assert.Equal(t,"I2"			,investorsList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,investorsList[0].GetBalance(),			  "Balance is not 2000"				)
    assert.Equal(t,true		    ,investorsList[0].GetIsInvestor(),		  "Is investor is not true "		)

	// removing the first client
	_, err = c.RemoveClient(ctx, idToRm)
	assert.Nil(t,err,"Error in addition of issuer")

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the issuer")

	issuersList = generateListClients(t,stream )

    assert.Equal(t,"2"			,issuersList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[0].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2000.00"    ,issuersList[0].GetBalance(),			"Balance is not 2000"				)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)

	//Invoice financing process testing:

	// invoice to be financed of 250 ???.
	invoice				:= &pb.Invoice{ClientId:id.Id,Amount:"250"}
	sellOrder			:= &pb.SellOrder{Size:"250",Amount:"200"}
	invoiceFinancing	:= &pb.InvoiceFinancing{SellOrder:sellOrder,Invoice:invoice}

	soId,err := c.StartInvoiceFinancing(ctx,invoiceFinancing)			// storing the so id for using it at the add bit test
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
    assert.Equal(t,"0.00"					,soList[0].GetFinanSize()	)
    assert.Equal(t,"0.00"					,soList[0].GetFinanAmount()	)
    assert.Equal(t,"ongoing"				,soList[0].GetState()		)

	// adding 3 investors for the add bid tests:
	// investor
	i1, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "111",Name: "J2",Surname:"I2",Balance:"2000",IsInvestor:true})
	assert.Nil(t,err)
	i2, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "222",Name: "J2",Surname:"I2",Balance:"20",IsInvestor:true})
	assert.Nil(t,err)
	i3, err := c.AddClient(ctx, &pb.Client{FiscalIdentity: "333",Name: "J2",Surname:"I2",Balance:"4000",IsInvestor:true})
	assert.Nil(t,err)

	// we has the soId too

	_,err = c.AddBid(ctx,&pb.Ledger{InvestorId:i1.GetId(),SellOrderId:soId.GetId(),Size:"50",Amount:"40"})
	assert.Nil(t,err)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err)

	investorsList = generateListClients(t,stream)
	assert.Equal(t,"111"			,investorsList[1].GetFiscalIdentity()	)
    assert.Equal(t,"1960.00"	    ,investorsList[1].GetBalance()			)
    assert.Equal(t,true			    ,investorsList[1].GetIsInvestor()		)

	streamBids,err := c.ListBids(ctx,&pb.Empty{})
    assert.Nil(t,err)
	bidsList := generateListBids(t,streamBids )

	assert.Equal(t,1						,len(bidsList)					)
    assert.Equal(t,i1.GetId()				,bidsList[0].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[0].GetSellOrderId()	)
    assert.Equal(t,"50.00"					,bidsList[0].GetSize()			)
    assert.Equal(t,"40.00"					,bidsList[0].GetAmount()		)
    assert.Equal(t,"1960.00"				,bidsList[0].GetBalance()		)
    assert.Equal(t,"20.00"					,bidsList[0].GetDiscount()		)
    assert.Equal(t,"10.00"					,bidsList[0].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[0].GetIsAdjusted()	)

	_,err = c.AddBid(ctx,&pb.Ledger{InvestorId:i2.GetId(),SellOrderId:soId.GetId(),Size:"50",Amount:"40"})
    assert.NotNil(t,err)
	//fmt.Println(err)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
	assert.Nil(t,err)

	investorsList = generateListClients(t,stream)
	assert.Equal(t,"222"			,investorsList[2].GetFiscalIdentity()	)
    assert.Equal(t,"20.00"			,investorsList[2].GetBalance()			)
    assert.Equal(t,true			    ,investorsList[2].GetIsInvestor()		)

	streamBids,err = c.ListBids(ctx,&pb.Empty{})
    assert.Nil(t,err)
	bidsList = generateListBids(t,streamBids )
	// not inserted in the ledger
	assert.Equal(t,1						,len(bidsList)					)
    assert.Equal(t,i1.GetId()				,bidsList[0].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[0].GetSellOrderId()	)
    assert.Equal(t,"50.00"					,bidsList[0].GetSize()			)
    assert.Equal(t,"40.00"					,bidsList[0].GetAmount()		)
    assert.Equal(t,"1960.00"				,bidsList[0].GetBalance()		)
    assert.Equal(t,"20.00"					,bidsList[0].GetDiscount()		)
    assert.Equal(t,"10.00"					,bidsList[0].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[0].GetIsAdjusted()	)

	// this bit also has to fail : discount to low
	_,err = c.AddBid(ctx,&pb.Ledger{InvestorId:i3.GetId(),SellOrderId:soId.GetId(),Size:"300",Amount:"290"})
	assert.NotNil(t,err)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err)

	investorsList = generateListClients(t,stream)
	assert.Equal(t,"333"			,investorsList[3].GetFiscalIdentity()	)
    assert.Equal(t,"4000.00"	    ,investorsList[3].GetBalance()			)
    assert.Equal(t,true			    ,investorsList[3].GetIsInvestor()		)

	streamBids,err = c.ListBids(ctx,&pb.Empty{})
    assert.Nil(t,err)
	bidsList = generateListBids(t,streamBids )

	assert.Equal(t,1						,len(bidsList)					)
    assert.Equal(t,i1.GetId()				,bidsList[0].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[0].GetSellOrderId()	)
    assert.Equal(t,"50.00"					,bidsList[0].GetSize()			)
    assert.Equal(t,"40.00"					,bidsList[0].GetAmount()		)
    assert.Equal(t,"1960.00"				,bidsList[0].GetBalance()		)
    assert.Equal(t,"20.00"					,bidsList[0].GetDiscount()		)
    assert.Equal(t,"10.00"					,bidsList[0].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[0].GetIsAdjusted()	)

	// this bit also has not to fail
	_,err = c.AddBid(ctx,&pb.Ledger{InvestorId:i3.GetId(),SellOrderId:soId.GetId(),Size:"100",Amount:"75"})
	assert.Nil(t,err)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err)

	investorsList = generateListClients(t,stream)
	assert.Equal(t,"333"			,investorsList[3].GetFiscalIdentity()	)
    assert.Equal(t,"3925.00"	    ,investorsList[3].GetBalance()			)
    assert.Equal(t,true			    ,investorsList[3].GetIsInvestor()		)

	streamBids,err = c.ListBids(ctx,&pb.Empty{})
    assert.Nil(t,err)
	bidsList = generateListBids(t,streamBids )

	assert.Equal(t,2						,len(bidsList)					)
    assert.Equal(t,i1.GetId()				,bidsList[0].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[0].GetSellOrderId()	)
    assert.Equal(t,"50.00"					,bidsList[0].GetSize()			)
    assert.Equal(t,"40.00"					,bidsList[0].GetAmount()		)
    assert.Equal(t,"1960.00"				,bidsList[0].GetBalance()		)
    assert.Equal(t,"20.00"					,bidsList[0].GetDiscount()		)
    assert.Equal(t,"10.00"					,bidsList[0].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[0].GetIsAdjusted()	)
    assert.Equal(t,i3.GetId()				,bidsList[1].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[1].GetSellOrderId()	)
    assert.Equal(t,"100.00"					,bidsList[1].GetSize()			)
    assert.Equal(t,"75.00"					,bidsList[1].GetAmount()		)
    assert.Equal(t,"3925.00"				,bidsList[1].GetBalance()		)
    assert.Equal(t,"25.00"					,bidsList[1].GetDiscount()		)
    assert.Equal(t,"25.00"					,bidsList[1].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[1].GetIsAdjusted()	)

	// checking the sell order state:

	streamSO,err = c.ListSellOrders(ctx,&pb.Empty{})
    assert.Nil(t,err)
	soList = generateListSellOrders(t,streamSO )

	assert.Equal(t,1						,len(soList)				)
    assert.Equal(t,invoicesList[0].GetId()	,soList[0].GetInvoiceId()	)
    assert.Equal(t,"250.00"					,soList[0].GetSize()		)
    assert.Equal(t,"200.00"					,soList[0].GetAmount()		)
    assert.Equal(t,"150.00"					,soList[0].GetFinanSize()	)
    assert.Equal(t,"115.00"					,soList[0].GetFinanAmount()	)
    assert.Equal(t,"ongoing"				,soList[0].GetState()		)


	// Final bid
	_,err = c.AddBid(ctx,&pb.Ledger{InvestorId:i3.GetId(),SellOrderId:soId.GetId(),Size:"150",Amount:"100"})
	assert.Nil(t,err)

	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:true})
    assert.Nil(t,err)

	investorsList = generateListClients(t,stream)
	assert.Equal(t,"333"			,investorsList[3].GetFiscalIdentity()	)
    assert.Equal(t,"3858.33"	    ,investorsList[3].GetBalance()			)
    assert.Equal(t,true			    ,investorsList[3].GetIsInvestor()		)

	streamBids,err = c.ListBids(ctx,&pb.Empty{})
    assert.Nil(t,err)
	bidsList = generateListBids(t,streamBids )

	assert.Equal(t,4						,len(bidsList)					)
    assert.Equal(t,i1.GetId()				,bidsList[0].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[0].GetSellOrderId()	)
    assert.Equal(t,"50.00"					,bidsList[0].GetSize()			)
    assert.Equal(t,"40.00"					,bidsList[0].GetAmount()		)
    assert.Equal(t,"1960.00"				,bidsList[0].GetBalance()		)
    assert.Equal(t,"20.00"					,bidsList[0].GetDiscount()		)
    assert.Equal(t,"10.00"					,bidsList[0].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[0].GetIsAdjusted()	)
    assert.Equal(t,i3.GetId()				,bidsList[1].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[1].GetSellOrderId()	)
    assert.Equal(t,"100.00"					,bidsList[1].GetSize()			)
    assert.Equal(t,"75.00"					,bidsList[1].GetAmount()		)
    assert.Equal(t,"3925.00"				,bidsList[1].GetBalance()		)
    assert.Equal(t,"25.00"					,bidsList[1].GetDiscount()		)
    assert.Equal(t,"25.00"					,bidsList[1].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[1].GetIsAdjusted()	)
    assert.Equal(t,i3.GetId()				,bidsList[2].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[2].GetSellOrderId()	)
    assert.Equal(t,"150.00"					,bidsList[2].GetSize()			)
    assert.Equal(t,"100.00"					,bidsList[2].GetAmount()		)
    assert.Equal(t,"0.00"					,bidsList[2].GetBalance()		) // not actually recordered because is adjustment.
    assert.Equal(t,"33.33"					,bidsList[2].GetDiscount()		)
    assert.Equal(t,"50.00"					,bidsList[2].GetExpectedProfit())
    assert.Equal(t,true						,bidsList[2].GetIsAdjusted()	)
    assert.Equal(t,i3.GetId()				,bidsList[3].GetInvestorId()	)
    assert.Equal(t,soId.GetId()				,bidsList[3].GetSellOrderId()	)
    assert.Equal(t,"100.00"					,bidsList[3].GetSize()			)
    assert.Equal(t,"66.67"					,bidsList[3].GetAmount()		)
    assert.Equal(t,"3858.33"				,bidsList[3].GetBalance()		)
    assert.Equal(t,"33.33"					,bidsList[3].GetDiscount()		)
    assert.Equal(t,"33.33"					,bidsList[3].GetExpectedProfit())
    assert.Equal(t,false					,bidsList[3].GetIsAdjusted()	)

	// checking the sell order at the end of the financing process:

	streamSO,err = c.ListSellOrders(ctx,&pb.Empty{})
    assert.Nil(t,err)
	soList = generateListSellOrders(t,streamSO )

	assert.Equal(t,1						,len(soList)				)
    assert.Equal(t,invoicesList[0].GetId()	,soList[0].GetInvoiceId()	)
    assert.Equal(t,"250.00"					,soList[0].GetSize()		)
    assert.Equal(t,"200.00"					,soList[0].GetAmount()		)
    assert.Equal(t,"250.00"					,soList[0].GetFinanSize()	)
    assert.Equal(t,"181.67"					,soList[0].GetFinanAmount()	)
    assert.Equal(t,"committed"				,soList[0].GetState()		)

	// invoice final state:

	streamInv,err = c.ListInvoices(ctx,&pb.Empty{})
    assert.Nil(t,err)
	invoicesList = generateListInvoices(t,streamInv )

	assert.Equal(t,1					,len(invoicesList)				)
	assert.Equal(t,id.GetId()			,invoicesList[0].GetClientId()	)
    assert.Equal(t,"250.00"				,invoicesList[0].GetAmount()	)
    assert.Equal(t,"financed"			,invoicesList[0].GetState()		)

	// issuer checking:
	stream, err = c.ListClients(ctx,&pb.IsInvestor{IsInvestor:false})
    assert.Nil(t,err,"Error listint the two issuers inserted")

	issuersList = generateListClients(t,stream)

    assert.Equal(t,"2"			,issuersList[0].GetFiscalIdentity(),    "Fiscal identity not equal to 2"    )
    assert.Equal(t,"J2"			,issuersList[0].GetName(),              "Name not equal to J2"		        )
    assert.Equal(t,"I2"			,issuersList[0].GetSurname(),           "Surname not equal to I2"			)
    assert.Equal(t,"2181.67"    ,issuersList[0].GetBalance(),			"Balance is not 2181.67"			)
    assert.Equal(t,false	    ,issuersList[0].GetIsInvestor(),		"Is investor is not false"			)
	fmt.Println("End tests")
}


