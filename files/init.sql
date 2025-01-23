-- System Region Table
CREATE TABLE region
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

-- User Profile Tables
CREATE TABLE "user"
(
    id SERIAL PRIMARY KEY
    -- other user details omitted as per diagram
);

CREATE TABLE user_address
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL REFERENCES "user" (id),
    detail     VARCHAR(255) NOT NULL,
    geohash    VARCHAR(12),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Product Category Tables
CREATE TABLE trash_category_group
(
    key CHAR(50) PRIMARY KEY
);

CREATE TABLE trash_category
(
    id                 SERIAL PRIMARY KEY,
    name               VARCHAR(255) NOT NULL,
    parent_category_id INT,
    "group"            CHAR(50) REFERENCES trash_category_group (key),
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by         INT
);

CREATE TABLE trash_category_detail
(
    id           SERIAL PRIMARY KEY,
    category_id  INT NOT NULL REFERENCES trash_category (id),
    description  TEXT,
    image_url_md VARCHAR(255),
    image_url_sm VARCHAR(255),
    image_url_lg VARCHAR(255),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by   INT
);

-- Order Handler Tables
CREATE TABLE partner
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    owner_id   INT,
    status     VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE partner_detail
(
    id           SERIAL PRIMARY KEY,
    partner_id   INT NOT NULL REFERENCES partner (id),
    description  TEXT,
    address      VARCHAR(255),
    image_url_md VARCHAR(255),
    image_url_sm VARCHAR(255),
    image_url_lg VARCHAR(255),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by   INT
);

CREATE TABLE picker
(
    id         SERIAL PRIMARY KEY,
    user_id    INT,
    partner_id INT REFERENCES partner (id),
    name       VARCHAR(255),
    status     VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE partner_point
(
    id         SERIAL PRIMARY KEY,
    partner_id INT NOT NULL REFERENCES partner (id),
    latitude   DECIMAL(10, 8),
    longitude  DECIMAL(11, 8),
    geohash    VARCHAR(12),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE partner_point_geohash
(
    id               SERIAL PRIMARY KEY,
    partner_point_id INT NOT NULL REFERENCES partner_point (id),
    geohash          VARCHAR(12),
    radius           DECIMAL(10, 2),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE picker_partner_point
(
    id               SERIAL PRIMARY KEY,
    picker_id        INT NOT NULL REFERENCES picker (id),
    partner_point_id INT NOT NULL REFERENCES partner_point (id),
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by       INT
);

CREATE TABLE cost_type
(
    key CHAR(50) PRIMARY KEY
);

CREATE TABLE partner_category
(
    id          SERIAL PRIMARY KEY,
    partner_id  INT NOT NULL REFERENCES partner (id),
    category_id INT NOT NULL REFERENCES trash_category (id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  INT
);

CREATE TABLE partner_category_cost
(
    id                  SERIAL PRIMARY KEY,
    partner_category_id INT      NOT NULL REFERENCES partner_category (id),
    region              INT      NOT NULL REFERENCES region (id),
    cost_type           CHAR(50) NOT NULL REFERENCES cost_type (key),
    cost_per_unit       DECIMAL(10, 2),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by          INT
);

-- Order Tables
CREATE TABLE cart_item
(
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES "user" (id),
    category_id INT NOT NULL REFERENCES trash_category (id),
    quantity    INT NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  INT
);

CREATE TABLE "order"
(
    id         SERIAL PRIMARY KEY,
    address_id INT NOT NULL REFERENCES user_address (id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE order_status
(
    id         SERIAL PRIMARY KEY,
    order_id   INT NOT NULL REFERENCES "order" (id),
    status     VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE order_item
(
    id          SERIAL PRIMARY KEY,
    user_id     INT NOT NULL REFERENCES "user" (id),
    category_id INT NOT NULL REFERENCES trash_category (id),
    quantity    INT NOT NULL,
    total_cost  DECIMAL(10, 2),
    order_id    INT NOT NULL REFERENCES "order" (id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by  INT
);

CREATE TABLE item_image
(
    id         SERIAL PRIMARY KEY,
    item_id    INT NOT NULL REFERENCES order_item (id),
    url        VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by INT
);

CREATE TABLE order_item_cost_detail
(
    id                       SERIAL PRIMARY KEY,
    order_item_id            INT NOT NULL REFERENCES order_item (id),
    partner_category_cost_id INT NOT NULL REFERENCES partner_category_cost (id),
    cost                     DECIMAL(10, 2),
    created_at               TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by               INT
);

-- Add indexes separately
CREATE INDEX idx_partner_point_geohash_partner_point_id ON partner_point_geohash(partner_point_id);
CREATE INDEX idx_partner_point_geohash_geohash ON partner_point_geohash(geohash);
CREATE INDEX idx_partner_point_geohash_radius ON partner_point_geohash(radius);

CREATE INDEX idx_picker_partner_point_picker_id ON picker_partner_point(picker_id);
CREATE INDEX idx_picker_partner_point_partner_point_id ON picker_partner_point(partner_point_id);