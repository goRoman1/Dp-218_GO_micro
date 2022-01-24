CREATE TABLE IF NOT EXISTS roles
(
    id         smallint PRIMARY KEY,
    name       VARCHAR(50) UNIQUE,
    is_admin    boolean,
    is_user     boolean,
    is_supplier boolean
    );

CREATE TABLE IF NOT EXISTS users
(
    id          serial PRIMARY KEY,
    login_email  VARCHAR(100) UNIQUE NOT NULL,
    is_blocked   boolean,
    user_name    VARCHAR(100),
    user_surname VARCHAR(100),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    role_id      int                 NOT NULL,
    password_hash VARCHAR(512),

    FOREIGN KEY (role_id) REFERENCES roles (id)
    );

CREATE TABLE IF NOT EXISTS login_status
(
    user_id    int PRIMARY KEY,
    logged_in  boolean,
    date_time  TIMESTAMP NOT NULL,
    ip_address VARCHAR(40),

    FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE IF NOT EXISTS contact_types
(
    id   smallserial PRIMARY KEY,
    name VARCHAR(50)
    );

CREATE TABLE IF NOT EXISTS contacts
(
    id          serial PRIMARY KEY,
    type_id      int NOT NULL,
    user_id      int NOT NULL,
    contact_info VARCHAR(200),

    FOREIGN KEY (type_id) REFERENCES contact_types (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE IF NOT EXISTS accounts
(
    id       serial PRIMARY KEY,
    name     VARCHAR(100),
    number   VARCHAR(100) UNIQUE NOT NULL,
    owner_id int                 NOT NULL,

    FOREIGN KEY (owner_id) REFERENCES users (id)
    );

CREATE TABLE IF NOT EXISTS payment_types
(
    id   serial PRIMARY KEY,
    name VARCHAR(100) UNIQUE
    );

CREATE TABLE IF NOT EXISTS supplier_commissions
(
    id                  serial PRIMARY KEY,
    commission_percent  NUMERIC(4, 2),
    user_id             int NOT NULL,

    FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE IF NOT EXISTS supplier_prices
(
    id               serial PRIMARY KEY,
    price            NUMERIC(15, 2),
    payment_type_id  smallint NOT NULL,
    user_id          int      NOT NULL,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE IF NOT EXISTS scooter_models
(
    id               smallserial PRIMARY KEY,
    payment_type_id  smallint     NOT NULL,
    model_name       VARCHAR(100) NOT NULL,
    max_weight       NUMERIC(5, 2),
    speed            smallint     NOT NULL,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id)
    );

CREATE TABLE IF NOT EXISTS scooters
(
    id            serial PRIMARY KEY,
    model_id      smallint            NOT NULL,
    owner_id      int                 NOT NULL,
    serial_number VARCHAR(100) UNIQUE NOT NULL,

    FOREIGN KEY (model_id) REFERENCES scooter_models (id),
    FOREIGN KEY (owner_id) REFERENCES users (id)
    );


CREATE TABLE IF NOT EXISTS scooter_stations
(
    id          serial PRIMARY KEY,
    name        VARCHAR(100),
    is_active   boolean,
    latitude      NUMERIC(16, 14),
    longitude     NUMERIC(16, 14)
    );

CREATE TABLE IF NOT EXISTS locations
(
    id          serial PRIMARY KEY,
    latitude      NUMERIC(16, 14),
    longitude     NUMERIC(16, 14),
    label         VARCHAR(100)
);

CREATE TABLE IF NOT EXISTS scooter_statuses
(
    scooter_id     int PRIMARY KEY,
    battery_remain NUMERIC(5, 2),
    station_id     int,
    latitude      NUMERIC(16, 14),
    longitude     NUMERIC(16, 14),
    can_be_rent   boolean,

    FOREIGN KEY (scooter_id)  REFERENCES scooters (id),
    FOREIGN KEY (station_id)  REFERENCES scooter_stations (id)
    );

CREATE TABLE IF NOT EXISTS problem_types
(
    id   smallint PRIMARY KEY,
    name VARCHAR(150) UNIQUE NOT NULL
    );

CREATE TABLE IF NOT EXISTS problems
(
    id            bigserial PRIMARY KEY,
    user_id       int       NOT NULL,
    type_Id       smallint  NOT NULL,
    scooter_id    int,
    date_reported TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description   text      NOT NULL,
    is_solved     boolean,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (type_id) REFERENCES problem_types (id)
    --FOREIGN KEY (scooter_id) REFERENCES scooters (id)
    );

CREATE TABLE IF NOT EXISTS solutions
(
    problem_id   bigint PRIMARY KEY,
    date_solved  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    description text      NOT NULL,

    FOREIGN KEY (problem_id) REFERENCES problems (id)
    );

CREATE TABLE IF NOT EXISTS scooter_statuses_in_rent
(
    id         bigserial PRIMARY KEY,
    station_id  int,
    date_time   TIMESTAMP NOT NULL,
    latitude      NUMERIC(16, 14),
    longitude     NUMERIC(16, 14),

    FOREIGN KEY (station_id) REFERENCES Scooter_Stations (id)
    );

CREATE TABLE IF NOT EXISTS orders
(
    id             bigserial PRIMARY KEY,
    user_id        int NOT NULL,
    scooter_id     int NOT NULL,
    status_start_id bigint,
    status_end_id  bigint,
    distance       NUMERIC(12, 2),
    amount_cents   bigint,

    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (scooter_id) REFERENCES scooters (id),
    FOREIGN KEY (status_start_id) REFERENCES scooter_statuses_in_rent (id),
    FOREIGN KEY (status_end_id) REFERENCES scooter_statuses_in_rent (id)
    );

CREATE TABLE IF NOT EXISTS account_transactions
(
    id              bigserial PRIMARY KEY,
    date_time       TIMESTAMP NOT NULL,
    payment_type_id smallint  NOT NULL,
    account_from_id int,
    account_to_id   int,
    order_id        bigint,
    amount_cents    bigint,

    FOREIGN KEY (payment_type_id) REFERENCES payment_types (id)
    --    FOREIGN KEY (account_from_id) REFERENCES accounts (id),
--    FOREIGN KEY (account_To_id) REFERENCES accounts (id),
--    FOREIGN KEY (order_id) REFERENCES orders (id)
    );

BEGIN;
/*
INSERT INTO problem_types(name) VALUES('Payment problem');
INSERT INTO problem_types(name) VALUES('Scooter problem');
INSERT INTO problem_types(name) VALUES('Other problem');
 */

INSERT INTO payment_types(name) VALUES('comission');
INSERT INTO payment_types(name) VALUES('simple income');
INSERT INTO payment_types(name) VALUES('simple outcome');
INSERT INTO payment_types(name) VALUES('Xiaomi М365 Mi Scooter');/**/
INSERT INTO payment_types(name) VALUES('Kugoo G2 Pro');/**/

INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(1, 'admin role', true, false, false);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(2, 'user role', false, true, false);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(3, 'supplier role', false, false, true);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(5, 'supplier+user role', false, false, true);
INSERT INTO roles(id, name, is_admin, is_user, is_supplier) VALUES(7, 'super_admin role', true, true, true);

INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('guru_admin@guru.com', false, 'Guru', 'Sadh', 7);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('VikaP@mail.com', false, 'Vika', 'Petrova', 1);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('IraK@mail.com', true, 'Ira', 'Petrova', 1);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('PetrPetroff@mail.com', false, 'Petr', 'Petrov', 3);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('UserChan@mail.com', false, 'Jackie', 'Chan', 2);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('UserB@mail.com', true, 'Beyonce', 'Ivanova', 2);
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id) VALUES('telo@mail.com', false, 'Goga', 'Boba', 2);

INSERT INTO scooter_stations(name, is_active, latitude, longitude ) VALUES ('Pobeda3', true, 48.42367000000000, 35.04436000000000);
INSERT INTO scooter_stations(name, is_active, latitude, longitude ) VALUES ('Dafi Mall', true, 48.42210000000000, 35.01960000000000);
INSERT INTO scooter_stations(name, is_active, latitude, longitude ) VALUES ('Private Sector', false, 48.42543000000000, 35.02183000000000);
INSERT INTO scooter_stations(name, is_active, latitude, longitude ) VALUES ('Getto', true, 48.41943000000000, 35.02293000000000);

INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id, password_hash) VALUES('gtr@gmail.com', false, 'Gregor', 'Tyson', 7, '$2a$10$Le9uo/qFrA.EPFh5d1Z5Wu1EaNCVMkeV1dOT/q86ZZ.obCeSY/472');
INSERT INTO users(login_email, is_blocked, user_name, user_surname, role_id, password_hash) VALUES('roma@gmail.com', false, 'Roman', 'Amelchenko', 3, '$2a$10$Le9uo/qFrA.EPFh5d1Z5Wu1EaNCVMkeV1dOT/q86ZZ.obCeSY/472');

INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed) VALUES(4, 'Xiaomi М365 Mi Scooter', 125,25);/**/
INSERT INTO scooter_models(payment_type_id, model_name, max_weight, speed) VALUES(5, 'Kugoo G2 Pro', 150, 35);/**/

INSERT INTO supplier_prices(price, payment_type_id, user_id) VALUES(50,4,9);/**/
INSERT INTO supplier_prices( price, payment_type_id, user_id) VALUES(60,5,9);/**/

INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(1, 9, '100000');/**/
INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(1, 9, '100001');/**/
INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(1, 9, '100002');/**/
INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(2, 9, '200000');/**/
INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(2, 9, '200001');/**/
INSERT INTO scooters(model_id, owner_id, serial_number) VALUES(2, 9, '200002');/**/

INSERT INTO locations(latitude, longitude, label) VALUES(48.00000000000000, 35.00000000000000, 'Pobeda');/**/



/* If the scooter has a status. it will not be visible on the init page. In the future, you will need to change the Check for checking the can_be_rent field
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(1, 77, 48.41452620789186, 35.01444471956219, true);/**/
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(2, 58, 48.43452620789186, 35.01444471956219, true);/**/
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(3, 100, 48.43452620789186, 35.01444471956219, true);/**/
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(4, 100, 48.43452620789186, 35.01444471956219, true);/**/
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(5, 40, 48.43452620789186, 35.01444471956219, true);/**/
INSERT INTO scooter_statuses(scooter_id, battery_remain, latitude, longitude, can_be_rent) VALUES(6, 100, 48.43452620789186, 35.01444471956219, true);/**/
 */

INSERT INTO accounts(name, number, owner_id) VALUES('Main account', '111222333444', 9);
INSERT INTO accounts(name, number, owner_id) VALUES('One more account', '55555666666', 9);

INSERT INTO account_transactions(date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents) VALUES(current_timestamp, 2, 0, 1, 0, 99999);
INSERT INTO account_transactions(date_time, payment_type_id, account_from_id, account_to_id, order_id, amount_cents) VALUES(current_timestamp, 3, 1, 0, 0, 11111);

INSERT INTO problem_types(id, name) VALUES (1, 'General');
INSERT INTO problem_types(id, name) VALUES (2, 'Payment issues');
INSERT INTO problem_types(id, name) VALUES (3, 'Scooter issues');
INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) VALUES(1, 1, 0, 'Bad service', false);
INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) VALUES(1, 2, 0, 'Wrong sum calculated', false);
INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) VALUES(2, 2, 0, 'Cant pay for service', false);
INSERT INTO problems(user_id, type_Id, scooter_id, description, is_solved) VALUES(3, 3, 1, 'Battery failed and scooter suddenly stopped', false);
COMMIT;