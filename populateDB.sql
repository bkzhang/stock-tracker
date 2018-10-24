CREATE TABLE IF NOT EXISTS user (
  id INT UNSIGNED NOT NULL AUTO_INCREMENT,
  username VARCHAR(20) NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS quote (
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

create TABLE IF NOT EXISTS ownedstock (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  symbol VARCHAR(10) NOT NULL,
  date DATETIME,
  timezone VARCHAR(20),
  price DOUBLE NOT NULL,
  shares INT UNSIGNED NOT NULL,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS user_ownedstock (
  user_id INT UNSIGNED NOT NULL,
  ownedstock_id BIGINT UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (user_id) REFERENCES user (id) ON DELETE CASCADE ON UPDATE CASCADE, 
  CONSTRAINT FOREIGN KEY (ownedstock_id) REFERENCES ownedstock (id) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (user_id, ownedstock_id)
);


create TABLE IF NOT EXISTS ownedstock_quote (
  ownedstock_id BIGINT UNSIGNED NOT NULL,
  quote_id INT UNSIGNED NOT NULL,
  CONSTRAINT FOREIGN KEY (ownedstock_id) REFERENCES ownedstock (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT FOREIGN KEY (quote_id) REFERENCES quote (id) ON DELETE CASCADE ON UPDATE CASCADE,
  PRIMARY KEY (ownedstock_id, quote_id)
);

INSERT INTO user (username) VALUES ('John');
SET @john_id = LAST_INSERT_ID();

INSERT INTO user (username) VALUES ('Matt');
SET @matt_id = LAST_INSERT_ID();

INSERT INTO quote (symbol) VALUES ('SHOP');
SET @shop_id = LAST_INSERT_ID();

INSERT INTO quote (symbol) VALUES ('MSFT');
SET @msft_id = LAST_INSERT_ID();

INSERT INTO quote (symbol) VALUES ('FB');
SET @fb_id = LAST_INSERT_ID();

INSERT INTO quote (symbol) VALUES('AAPL');
SET @aapl_id = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES ('AAPL', '2015-01-02 05:25:00', 'US/EASTERN', 193.29, 20);
set @john_aapl_id1 = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES ('FB', '2017-09-10 10:04:00', 'US/EASTERN', 249.23, 5);
set @john_fb_id1 = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES ('MSFT', '2014-08-02 12:11:00', 'US/EASTERN', 103.50, 10);
set @john_msft_id1 = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES ('MSFT', '2009-03-09 02:01:00', 'US/EASTERN', 49.98, 120);
set @matt_msft_id1 = LAST_INSERT_ID();

INSERT INTO ownedstock (symbol, date, timezone, price, shares) VALUES ('SHOP', '2016-04-09 11:45:00', 'US/EASTERN', 134.48, 50);
set @matt_shop_id1 = LAST_INSERT_ID();

INSERT INTO user_ownedstock (user_id, ownedstock_id)  VALUES (@john_id, @john_aapl_id1);
INSERT INTO user_ownedstock (user_id, ownedstock_id)  VALUES (@john_id, @john_fb_id1);
INSERT INTO user_ownedstock (user_id, ownedstock_id)  VALUES (@john_id, @john_msft_id1);
INSERT INTO user_ownedstock (user_id, ownedstock_id)  VALUES (@matt_id, @matt_msft_id1);
INSERT INTO user_ownedstock (user_id, ownedstock_id)  VALUES (@matt_id, @matt_shop_id1);

INSERT INTO ownedstock_quote (ownedstock_id, quote_id) VALUES (@john_aapl_id1, @aapl_id);
INSERT INTO ownedstock_quote (ownedstock_id, quote_id) VALUES (@john_fb_id1, @fb_id);
INSERT INTO ownedstock_quote (ownedstock_id, quote_id) VALUES (@john_msft_id1, @msft_id);
INSERT INTO ownedstock_quote (ownedstock_id, quote_id) VALUES (@matt_msft_id1, @msft_id);
INSERT INTO ownedstock_quote (ownedstock_id, quote_id) VALUES (@matt_msft_id1, @shop_id);
