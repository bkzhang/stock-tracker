CREATE TABLE IF NOT EXISTS user (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(20) NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS stock (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  symbol VARCHAR(10) NOT NULL UNIQUE,
  date DATETIME,
  timezone VARCHAR(20),
  high DOUBLE,
  low DOUBLE,
  open DOUBLE,
  close DOUBLE,
  volume INT UNSIGNED,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_stock (
  user_id INT UNSIGNED NOT NULL,
  stock_id INT UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE, 
  CONSTRAINT FOREIGN KEY (stock_id) REFERENCES stock (id) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (user_id, stock_id)
);

create TABLE IF NOT EXISTS ownedstock (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  symbol VARCHAR(10) NOT NULL,
  date DATETIME,
  timezone VARCHAR(20),
  PRICE DOUBLE NOT NULL,
  SHARES INT UNSIGNED NOT NULL,
  PRIMARY KEY (id)
);

create TABLE IF NOT EXISTS user_stock_ownedstock (
  user_id INT UNSIGNED NOT NULL,
  stock_id INT UNSIGNED NOT NULL,
  ownedstock_id BIGINT UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT FOREIGN KEY (stock_id) REFERENCES stock (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT FOREIGN KEY (ownedstock_id) REFERENCES ownedstock (id) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (user_id, stock_id, ownedstock_id)
);

INSERT INTO user (username) VALUES ('John');
SET @john_id = LAST_INSERT_ID();

INSERT INTO user (username) VALUES ('Matt');
SET @matt_id = LAST_INSERT_ID();

INSERT INTO stock (symbol) VALUES ('SHOP');
SET @shop_id = LAST_INSERT_ID();

INSERT INTO stock (symbol) VALUES ('MSFT');
SET @msft_id = LAST_INSERT_ID();

INSERT INTO stock (symbol) VALUES ('FB');
SET @fb_id = LAST_INSERT_ID();

INSERT INTO stock(symbol) VALUES('AAPL');
SET @aapl_id = LAST_INSERT_ID();

INSERT INTO user_stock (user_id, stock_id)  VALUES (@john_id, @shop_id);
INSERT INTO user_stock (user_id, stock_id)  VALUES (@john_id, @msft_id);
INSERT INTO user_stock (user_id, stock_id)  VALUES (@john_id, @aapl_id);
INSERT INTO user_stock (user_id, stock_id)  VALUES (@matt_id, @msft_id);
INSERT INTO user_stock (user_id, stock_id)  VALUES (@matt_id, @fb_id);

INSERT INTO ownedstock (symbol, price, shares) VALUES ('AAPL', 193.29, 20);
set @ownedmsft_id = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, price, shares) VALUES ('SHOP', 134.48, 50);
set @ownedshop_id = LAST_INSERT_ID();

INSERT INTO user_stock_ownedstock (user_id, stock_id, ownedstock_id) VALUES (@john_id, @msft_id, @ownedmsft_id);
INSERT INTO user_stock_ownedstock (user_id, stock_id, ownedstock_id) VALUES (@john_id, @shop_id, @ownedshop_id);
