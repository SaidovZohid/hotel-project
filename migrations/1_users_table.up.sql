create table "users" (
    "id" serial primary key,
    "first_name" varchar(50) not null,
    "last_name" varchar(50) not null,
    "email" varchar(100) not null unique,
    "password" varchar not null,
    "phone_number" varchar(30) unique,
    "type" varchar(20) check("type" in('user', 'superuser', 'manager')) not null,
    "created_at" timestamp default current_timestamp,
    "updated_at" timestamp,
    "deleted_at" timestamp
);

create table "hotels" (
    "id" serial primary key,
    "manager_id" int not null references users(id) unique,
    "hotel_name" varchar(100) not null,
    "description" text not null,
    "address" varchar(100) not null,
    "image_url" varchar not null,
    "num_of_rooms" int not null
);

create table "hotel_images" (
    "id" serial primary key,
    "hotel_id" int not null references hotels(id),
    "image_url" varchar not null,
    "sequence_number" int not null
);

create table "rooms" (
    "id" serial primary key,
    "room_number" int not null,
    "hotel_id" int not null references hotels(id),
    "type" varchar(20) check ("type" in ('single', 'double', 'family')),
    "description" text not null,
    "price_per_night" numeric(18, 2) not null,
    "status" boolean not null default true
);

create table "bookings" (
    "id" serial primary key,
    "check_in" date not null,
    "check_out" date not null,
    "hotel_id" int not null references hotels(id),
    "room_id" int not null references rooms(id),
    "user_id" int not null references users(id),
    "booked_at" timestamp default current_timestamp
);