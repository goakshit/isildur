create table product (
    id uuid not null primary key
    name varchar
    description varchar
    monthly_price numeric
    instructor_name varchar
);
insert into product values('56f79fee-0cb0-4e87-9bca-7b5811cca4ce', 'YOGA L1', 'Basic yoga lessons', 5, 'A. Dhar');
insert into product values('56f79fee-0cb0-4e87-9bca-7b5811cca4cf', 'YOGA L2', 'Intermediate yoga lessons', 7, 'A. Dhar');
create table subscription (
    id uuid not null primary key
    product_id uuid not null
    duration_in_months smallint
    tax numeric
    total_cost numeric
    status varchar
    start_date timestamptz
    end_date timestamptz
);