CREATE TABLE currency (
    id uuid NOT NULL PRIMARY KEY,
    code text NOT NULL
);

CREATE TABLE product (
	id uuid NOT NULL PRIMARY KEY,
	name text NOT NULL,
    capacity int NOT NULL,

    price int NOT NULL,
    currency_id uuid NOT NULL,

    FOREIGN KEY (currency_id) REFERENCES currency(id)
);

CREATE TYPE availability_status AS ENUM ('AVAILABLE', 'SOLD_OUT');

CREATE TABLE availability (
    id uuid NOT NULL PRIMARY KEY,
    product_id uuid NOT NULL,

    localDate date NOT NULL,
    status availability_status NOT NULL,
    vacancies int NOT NULL,
    available bool NOT NULL,

    FOREIGN KEY (product_id) REFERENCES product(id)
);

CREATE TYPE booking_status AS ENUM ('RESERVED', 'CONFIRMED');

CREATE TABLE booking (
    id uuid NOT NULL PRIMARY KEY,
    product_id uuid NOT NULL,
    availability_id uuid NOT NULL,

    status booking_status NOT NULL,
    price int NOT NULL,
    currency_id uuid NOT NULL,

    FOREIGN KEY (product_id) REFERENCES product(id),
    FOREIGN KEY (availability_id) REFERENCES availability(id),
    FOREIGN KEY (currency_id) REFERENCES currency(id)
);

CREATE TABLE booking_unit (
    id uuid NOT NULL PRIMARY KEY,
    booking_id uuid NOT NULL,
    ticket text,

    FOREIGN KEY (booking_id) REFERENCES booking(id)
);