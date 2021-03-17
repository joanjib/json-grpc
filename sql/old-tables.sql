-- domains		
create domain amount_type_calc as
	numeric(11,2)	check (value >= 0)				;	-- max amount : 999999999.99 and bigger than -1
create domain amount_type as
	numeric(11,2)	not null check (value >= 0)				;	-- max amount : 999999999.99 and bigger than -1
create domain discount_type as
	numeric(5,2)	check (value <=100.00 and value >=0.00);
create domain name_type as
	varchar(128)	not null;
create domain fiscal_identity_type as
	char(9)			not null;
create domain creation_time_type as
	timestamp		not null default now();

-- types
create type sell_order_state		as enum ('ongoing','reversed','locked','committed');
create type invoice_state			as enum ('financing search','rejected','financed');

-- tables
create table issuer (
	id 						int generated always as identity primary key			,   -- identity
	fiscal_identity			fiscal_identity_type									,
	name					name_type												,
	surname					name_type												,	
	balance					amount_type
);

create table investor (
	id 						int generated always as identity primary key			,   -- identity
	fiscal_identity			fiscal_identity_type									,
	name					name_type												,
	surname					name_type												,
	balance					amount_type
);

create table invoice (
	id 						int generated always as identity primary key			,   -- identity
	issuer_id				int references issuer not null							,
	invoice_amount			amount_type												,
	creation_time			creation_time_type										,
	state					invoice_state
);

create table sell_order (
	id 						int generated always as identity primary key			,   -- identity
	invoice_id				int references invoice not null							,
	order_size				amount_type												,
	order_amount			amount_type												,
	-- financed amounts - begin
	finan_size				amount_type_calc	default 0							,
	finan_amount			amount_type_calc	default 0							,
	-- financed amounts - end
	discount				discount_type generated always as (100 - (order_amount/order_size)*100) stored,
	creation_time			creation_time_type										,
	state					sell_order_state not null

);

create table ledger (
	id 						int generated always as identity primary key			,   -- identity
	investor_id				int references investor not null						,
	sell_order_id			int references sell_order not null						,
	creation_time			creation_time_type										,
	investor_size			amount_type												,
	investor_amount			amount_type												,
	investor_balance		amount_type												,
	investor_discount		discount_type 		generated always as (100 - (investor_amount/investor_size)*100) stored,	
	expected_profit			amount_type_calc   	generated always as (investor_size - investor_amount) stored
);
