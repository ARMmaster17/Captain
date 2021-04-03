CREATE TABLE IF NOT EXISTS airspace (
    AirspaceID integer NOT NULL PRIMARY KEY UNIQUE,
    HumanName varchar(255) NOT NULL UNIQUE,
    NetName varchar(255) NOT NULL UNIQUE
);

INSERT INTO airspace (AirspaceID, HumanName, NetName) VALUES (0, 'Default Airspace', 'default');