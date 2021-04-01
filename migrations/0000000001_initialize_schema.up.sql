CREATE TABLE IF NOT EXISTS airspace (
    AirspaceID int,
    HumanName varchar(255),
    NetName varchar(255)
);

INSERT INTO airspace (AirspaceID, HumanName, NetName) VALUES (0, 'Default Airspace', 'default');