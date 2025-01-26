CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE OR REPLACE FUNCTION nanoid(
    size int DEFAULT 16,
    alphabet text DEFAULT '_-0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ'
)
    RETURNS text
    LANGUAGE plpgsql
    volatile
AS
$$
DECLARE
    idBuilder      text := '';
    counter        int  := 0;
    bytes          bytea;
    alphabetIndex  int;
    alphabetArray  text[];
    alphabetLength int;
    mask           int;
    step           int;
BEGIN
    alphabetArray := regexp_split_to_array(alphabet, '');
    alphabetLength := array_length(alphabetArray, 1);
    mask := (2 << cast(floor(log(alphabetLength - 1) / log(2)) as int)) - 1;
    step := cast(ceil(1.6 * mask * size / alphabetLength) AS int);

    while true
        loop
            bytes := gen_random_bytes(step);
            while counter < step
                loop
                    alphabetIndex := (get_byte(bytes, counter) & mask) + 1;
                    if alphabetIndex <= alphabetLength then
                        idBuilder := idBuilder || alphabetArray[alphabetIndex];
                        if length(idBuilder) = size then
                            return idBuilder;
                        end if;
                    end if;
                    counter := counter + 1;
                end loop;

            counter := 0;
        end loop;
END
$$;

-- System Region Table
CREATE TABLE region
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16)  NOT NULL
);

-- User Profile Tables
CREATE TABLE "user"
(
    id           VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    mobile_code  CHAR(5)     NOT NULL,
    phone_number VARCHAR(20) NOT NULL,
    display_name VARCHAR(64) NOT NULL,
    status       VARCHAR(16) DEFAULT 'created',
    created_at   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_profile
(
    id            SERIAL PRIMARY KEY,
    user_id       VARCHAR(16)  NOT NULL REFERENCES "user" (id),
    avatar_url_md VARCHAR(255) NOT NULL,
    avatar_url_sm VARCHAR(255),
    avatar_url_lg VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_address
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    user_id    VARCHAR(16)  NOT NULL REFERENCES "user" (id),
    detail     VARCHAR(255) NOT NULL,
    geohash    VARCHAR(12),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE TABLE user_auth
(
    id         SERIAL PRIMARY KEY,
    user_id    VARCHAR(16) NOT NULL REFERENCES "user" (id),
    auth_type  CHAR(10)    NOT NULL DEFAULT 'juno',
    created_at TIMESTAMP            DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16) NOT NULL
);

-- Product Category Tables
CREATE TABLE trash_category_group
(
    key CHAR(50) PRIMARY KEY
);

CREATE TABLE trash_category
(
    id                 VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    name               VARCHAR(255) NOT NULL,
    parent_category_id VARCHAR(16) REFERENCES "trash_category" (id),
    "group"            CHAR(50) REFERENCES trash_category_group (key),
    status             CHAR(14)    DEFAULT 'created',
    created_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by         VARCHAR(16)  NOT NULL,
    updated_at         TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    updated_by         VARCHAR(16)  NOT NULL
);

CREATE TABLE trash_category_detail
(
    id           VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    category_id  VARCHAR(16) NOT NULL REFERENCES trash_category (id),
    description  TEXT,
    image_url_md VARCHAR(255),
    image_url_sm VARCHAR(255),
    image_url_lg VARCHAR(255),
    created_at   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(16) NOT NULL
);

-- Order Handler Tables
CREATE TABLE partner
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    owner_id   VARCHAR(16) REFERENCES "user" (id),
    status     VARCHAR(15),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16)  NOT NULL
);

CREATE TABLE partner_detail
(
    id           VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    partner_id   VARCHAR(16) NOT NULL REFERENCES partner (id),
    description  TEXT,
    address      VARCHAR(255),
    image_url_md VARCHAR(255),
    image_url_sm VARCHAR(255),
    image_url_lg VARCHAR(255),
    created_at   TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by   VARCHAR(16) NOT NULL
);

CREATE TABLE picker
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    user_id    VARCHAR(16),
    partner_id VARCHAR(16) REFERENCES partner (id),
    name       VARCHAR(255),
    status     VARCHAR(50),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16) NOT NULL
);

CREATE TABLE partner_point
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    partner_id VARCHAR(16) NOT NULL REFERENCES partner (id),
    latitude   DECIMAL(10, 8),
    longitude  DECIMAL(11, 8),
    geohash    VARCHAR(12),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16) NOT NULL
);

CREATE TABLE partner_point_geohash
(
    id               VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    partner_point_id VARCHAR(16) NOT NULL REFERENCES partner_point (id),
    geohash          VARCHAR(12),
    radius           DECIMAL(10, 2),
    created_at       TIMESTAMP   DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE picker_partner_point
(
    id               VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    picker_id        VARCHAR(16) NOT NULL REFERENCES picker (id),
    partner_point_id VARCHAR(16) NOT NULL REFERENCES partner_point (id),
    created_at       TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by       VARCHAR(16) NOT NULL
);

CREATE TABLE cost_type
(
    key CHAR(50) PRIMARY KEY
);

CREATE TABLE partner_category
(
    id          VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    partner_id  VARCHAR(16) NOT NULL REFERENCES partner (id),
    category_id VARCHAR(16) NOT NULL REFERENCES trash_category (id),
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(16) NOT NULL
);

CREATE TABLE partner_category_cost
(
    id                  VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    partner_category_id VARCHAR(16) NOT NULL REFERENCES partner_category (id),
    region              VARCHAR(16) NOT NULL REFERENCES region (id),
    cost_type           CHAR(50)    NOT NULL REFERENCES cost_type (key),
    cost_per_unit       DECIMAL(10, 2),
    created_at          TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by          VARCHAR(16) NOT NULL
);

-- Order Tables
CREATE TABLE cart_item
(
    id          VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    user_id     VARCHAR(16) NOT NULL REFERENCES "user" (id),
    category_id VARCHAR(16) NOT NULL REFERENCES trash_category (id),
    quantity    INT         NOT NULL,
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(16) NOT NULL
);

CREATE TABLE "order"
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    address_id VARCHAR(16) NOT NULL REFERENCES user_address (id),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16) NOT NULL
);

CREATE TABLE order_status
(
    id         VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    order_id   VARCHAR(16) NOT NULL REFERENCES "order" (id),
    status     VARCHAR(24),
    created_at TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(16) NOT NULL
);

CREATE TABLE order_item
(
    id          VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    user_id     VARCHAR(16) NOT NULL REFERENCES "user" (id),
    category_id VARCHAR(16) NOT NULL REFERENCES trash_category (id),
    quantity    INT         NOT NULL,
    total_cost  DECIMAL(10, 2),
    order_id    VARCHAR(16) NOT NULL REFERENCES "order" (id),
    created_at  TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by  VARCHAR(16) NOT NULL
);

CREATE TABLE item_image
(
    id            VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    cart_item_id  VARCHAR(16)  NOT NULL REFERENCES cart_item (id),
    order_item_id VARCHAR(16)  NOT NULL REFERENCES order_item (id),
    url           VARCHAR(255) NOT NULL,
    created_at    TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by    VARCHAR(16)  NOT NULL
);

CREATE TABLE order_item_cost_detail
(
    id                       VARCHAR(16) DEFAULT nanoid() PRIMARY KEY,
    order_item_id            VARCHAR(16) NOT NULL REFERENCES order_item (id),
    partner_category_cost_id VARCHAR(16) NOT NULL REFERENCES partner_category_cost (id),
    cost                     DECIMAL(10, 2),
    created_at               TIMESTAMP   DEFAULT CURRENT_TIMESTAMP,
    created_by               VARCHAR(16) NOT NULL
);

-- Add indexes separately
CREATE INDEX idx_partner_point_geohash_partner_point_id ON partner_point_geohash (partner_point_id);
CREATE INDEX idx_partner_point_geohash_geohash ON partner_point_geohash (geohash);
CREATE INDEX idx_partner_point_geohash_radius ON partner_point_geohash (radius);

CREATE INDEX idx_picker_partner_point_picker_id ON picker_partner_point (picker_id);
CREATE INDEX idx_picker_partner_point_partner_point_id ON picker_partner_point (partner_point_id);