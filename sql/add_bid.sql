create or replace function add_bid(
	investor_id			bigint			,
	sell_order_id		bigint			,
    investor_size       amount_type  	,  
    investor_amount     amount_type  	
)
returns bigint
language plpgsql
as $$
declare
	temp_balance 	amount_type 	= 0;
	-- id of the ledger to return
	ret_id 			bigint;
	so_id			bigint;
	so_state		sell_order_state;
	so_size			amount_type 	= 0;
	so_amount		amount_type 	= 0;
	so_fin_size		amount_type 	= 0;
	so_fin_amount	amount_type 	= 0;
	so_discount 	discount_type 	= 0;
	bid_discount	discount_type	= (100 - (investor_amount/investor_size)*100);
	is_adjusted		boolean			= false;
begin

	-- sell order basic information retrive
	select 	sell_orders.id	,state		,size 			,amount			,discount		,finan_size		,finan_amount
	into	so_id			,so_state	,so_size		,so_amount		,so_discount	,so_fin_size	,so_fin_amount
	from 	sell_orders 
	where 	sell_orders.id  = sell_order_id;
	
	if (so_state != 'ongoing')
	then
		raise exception 'sell order state is not ongoing, is %',so_state;
	end if;

	-- discount control
	if 	(bid_discount < so_discount						)		-- not an acceptable discount
	then
		raise exception 'discount of the bid not enought %', bid_discount;
	end if;
	
	-- sell order financed, but need to be recalculated the bid. Due to time limitations it's done here the adjustment.
	if  ((so_fin_size  + investor_size) > so_size	)   
	then
		-- if is adjusted we must insert an entry in the ledger too:
		insert into ledgers(investor_id,sell_order_id,size		   ,amount			,balance	,is_adjusted)
					values (investor_id,sell_order_id,investor_size,investor_amount,temp_balance,true    	);

		investor_size	= so_size - so_fin_size;
		investor_amount	= investor_size - (investor_size * (bid_discount/100));
		is_adjusted 	= true;
	end if;

	-- updating the investor balance
	update 	clients
	set 	balance = balance - investor_amount
	where 	clients.id = investor_id
	returning balance into temp_balance;

	-- updating the sell order financing fields
	update 	sell_orders
	set 	finan_size 	= finan_size 	+ investor_size,
			finan_amount= finan_amount	+ investor_amount
	where 	sell_orders.id = sell_order_id;

	-- actual ledger insertion.
	insert into ledgers(investor_id,sell_order_id,size		   ,amount			,balance	,is_adjusted)
				values (investor_id,sell_order_id,investor_size,investor_amount,temp_balance,false    	)
	returning id into ret_id;

	--  the sell order is finally financed
	if  ((so_fin_size  + investor_size) = so_size	)   
	then
		-- updating the balance of the issuer
		update clients
		set balance = (select sum(ledger.investor_amount) from ledgers where ledgers.sell_order_id = so_id and is_adjusted = false)
		where clients.id = (select id 
							from   clients
							where  clients.id = (	select 	client_id 
													from 	invoices 
													where 	invoices.id = (select invoice_id from sell_orders where sell_orders.id = so_id)
											  )
						  );
		-- updating the state of the sell order
		update sell_orders
		set state = 'committed'
		where sell_orders.id= sell_order_id;
		
		-- updating the state of the invoice.
		update invoices
		set state = 'financed'
		where invoices.id = (select invoice_id from sell_orders where sell_orders.id = so_id);

	end if;

	return ret_id;

end; 
$$
