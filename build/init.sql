CREATE TABLE "users" (
    "id" varchar PRIMARY KEY,
    "value" NUMERIC NOT NULL,
    "currency" varchar NOT NULL DEFAULT 'USD',
    "date_of_creation" timestampz NOT NULL DEFAULT (now())
);

CREATE TABLE "deposits" (
    "id" bigserial PRIMARY KEY,
    "value" NUMERIC NOT NULL,
    "date_of_creation" timestampz NOT NULL DEFAULT (now()),
    "from_user_id" varchar DEFAULT 'atm machine',
    "to_user_id" varchar NOT NULL
);

CREATE TABLE "cashouts" {
    "id" bigserial PRIMARY KEY, 
    "value" NUMERIC NOT NULL,
    "date_of_creation" timestampz NOT NULL DEFAULT (now()),
    "from_user_id" varchar NOT NULL,
    "to_user_id" varchar DEFAULT 'atm_machine'
};

CREATE TABLE "transactions" {
    "id" bigserial PRIMARY KEY, 
    "value" NUMERIC NOt NULL,
    "date_of_creation" timestampz NOT NULL DEFAULT (now()),
    "from_user_id" varchar NOT NULL,
    "to_user_id" varchar NOT NULL
}

ALTER TABLE "deposits" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");

ALTER TABLE "withdrawals" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("from_user_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("to_user_id") REFERENCES "users" ("id");