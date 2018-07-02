CREATE TABLE figures (
id INT(6) UNSIGNED AUTO_INCREMENT PRIMARY KEY,
`Number` VARCHAR(10),
`Name` VARCHAR(255),
`Character`     VARCHAR(255),
`Category`      VARCHAR(255),
`Subcategory`   VARCHAR(255),
Sculptor      VARCHAR(255),
OfficialPrice VARCHAR(255),
PreorderDate  VARCHAR(255),
ReleaseDate   VARCHAR(255),
Reedition1    VARCHAR(255),
Reedition2    VARCHAR(255),
Height        VARCHAR(255),
Weight        VARCHAR(255),
BoxSize       VARCHAR(255),
Observations  TEXT
);
