DROP TABLE IF EXISTS album;
CREATE TABLE weather (
  id         INT AUTO_INCREMENT NOT NULL,
  country    VARCHAR(2),
  date       DATE,
  message    VARCHAR(255),
  PRIMARY KEY (`id`)
);