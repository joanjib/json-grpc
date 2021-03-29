-- domains		
create domain amount_type_calc as
	numeric(11,2)	check (value >= 0)				;	-- max amount : 999999999.99 and bigger than -1
create domain amount_type as
	amount_type_calc not null						;	-- max amount : 999999999.99 and bigger than -1
create domain discount_type as
	numeric(5,2)	check (value <=100.00 and value >=0.00);

-- types
create type sell_order_state		as enum ('ongoing','reversed','locked','committed');
create type invoice_state			as enum ('financing search','rejected','financed');
