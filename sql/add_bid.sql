create or replace function add_bid(
	investor_id			int				,
	sell_order_id		int				,
    investor_size       amount_type  	,  
    investor_amount     amount_type  	
)
returns int
language plpgsql
as $$
declare
	temp_balance 	amount_type 	= 0;
	-- id of the ledger to return
	ret_id 			int;
	so_id			int;
	so_state		sell_order_state;
	so_size			amount_type 	= 0;
	so_amount		amount_type 	= 0;
	so_fin_size		amount_type 	= 0;
	so_fin_amount	amount_type 	= 0;
	so_discount 	discount_type 	= 0;
	bid_discount	discount_type	= (100 - (investor_amount/investor_size)*100);
begin
	-- sell order basic information retrive
	select 	sell_order_id	,state		,order_size 	,order_amount	,discount		,finan_size		,finan_amount
	into	so_id			,so_state	,so_size		,so_amount		,so_discount	,so_fin_size	,so_fin_amount
	from 	sell_order 
	where 	sell_order.id  = sell_order_id;

	-- discount control
	if 	(bid_discount < so_discount						)		-- not an acceptable discount
	then
		raise exception 'discount of the bid not enought %', bid_discount;
	end if;
	
	-- sell order financed, but need to be recalculated the bid. Due to time limitations it's done here the adjustment.
	if  ((so_fin_size  + investor_size) > so_size	)   
	then
		investor_size	= so_size - so_fin_size;
		investor_amount	= investor_size - (investor_size * (bid_discount/100));
	end if;

	-- updating the investor balance
	update 	investor
	set 	balance = balance - investor_amount
	where 	investor.id = investor_id
	returning balance into temp_balance;

	-- updating the sell order financing fields
	update 	sell_order
	set 	finan_size 	= finan_size 	+ investor_size,
			finan_amount= finan_amount	+ investor_amount
	where 	sell_order.id = sell_order_id;

	-- actual ledger insertion.
	insert into ledger (investor_id,sell_order_id,investor_size,investor_amount,investor_balance)
				values (investor_id,sell_order_id,investor_size,investor_amount,temp_balance    )
	returning id into ret_id;
	
	--  the sell order is finally financed
	if  ((so_fin_size  + investor_size) = so_size	)   
	then
		-- updating the balance of the issuer
		update issuer
		set balance = (select sum(ledger.investor_amount) from ledger where ledger.sell_order_id = so_id)
		where issuer.id = (	select id 
							from issuer 
							where issuer.id = (	select issuer_id 
												from invoice 
												where invoice.id = (select invoice_id from sell_order where sell_order.id = so_id)
											  )
						  );
		-- updating the state of the sell order
		update sell_order
		set state = 'committed'
		where sell_order.id= sell_order_id;
		
		-- updating the state of the invoice.
		update invoice
		set state = 'financed'
		where invoice.id = (select invoice_id from sell_order where sell_order.id = so_id);

	end if;

	return ret_id;

end; 
$$
