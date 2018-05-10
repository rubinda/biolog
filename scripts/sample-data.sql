-- Some example data for testing in the database. IDs always start with 1.

-- [conservation_status]
INSERT INTO conservation_status VALUES (DEFAULT, 'EX', 'Extinct', 'Izumrle');
INSERT INTO conservation_status VALUES (DEFAULT, 'EW', 'Extinct in the Wild', 'V naravi izumrle');
INSERT INTO conservation_status VALUES (DEFAULT, 'CR', 'Critically Endangered', 'Skrajno ogrožene');
INSERT INTO conservation_status VALUES (DEFAULT, 'EN', 'Endangered', 'Ogrožene');
INSERT INTO conservation_status VALUES (DEFAULT, 'VU', 'Vulnerable', 'Ranljive');
INSERT INTO conservation_status VALUES (DEFAULT, 'NT', 'Near Threatened', 'Potencialno ogrožene');
INSERT INTO conservation_status VALUES (DEFAULT, 'CD', 'Conservation Dependent', 'Varstveno odvisne');
INSERT INTO conservation_status VALUES (DEFAULT, 'LC', 'Least Concern', 'Najmanj ogrožene');
INSERT INTO conservation_status VALUES (DEFAULT, 'DD', 'Data deficient', 'Premalo podatkov');
INSERT INTO conservation_status VALUES (DEFAULT, 'NE', 'Not evaluated', 'Neopredeljene');

-- [external_auth_provider]
INSERT INTO external_auth_provider VALUES(DEFAULT, 'Google');

-- [species]
INSERT INTO species VALUES (DEFAULT, 'Passer domesticus', 'Animalia', 'Passeridae', 'Aves',
    'Chordata', 'Passeriformes', 'Passer', 'Passer domesticus (Linnaeus, 1758)', 'Passer domesticus', 8, 5231190);

-- [biolog_user] 
INSERT INTO biolog_user (id, external_id, display_name, given_name, family_name, email, public_observations, picture, external_auth_provider)
    VALUES(DEFAULT, '7464723854823589876345', 'Marjetka Kostanjsek', 'Marjetka', 'Kostanjesek', 'marjetka@fakemail.com', TRUE, 'https://doesnt.exist.com/path/to/picture.png', 1);
INSERT INTO biolog_user (id, external_id, display_name, given_name, family_name, email, public_observations, picture, external_auth_provider)
    VALUES(DEFAULT, '4896207315489620731522', 'Silvika Brezovec', 'Silvika', 'Brezovec', 'silvika@fakemail.com', TRUE, 'https://doesnt.exist.com/path/to/picture.png', 1);
INSERT INTO biolog_user (id, external_id, display_name, given_name, family_name, email, public_observations, picture, external_auth_provider)
    VALUES(DEFAULT, '5916087423591608742356', 'Isabela Trtovnik', 'Isabela', 'Trtovnik', 'isabela@fakemail.com', TRUE, 'https://doesnt.exist.com/path/to/picture.png', 1);
INSERT INTO biolog_user (id, external_id, display_name, given_name, family_name, email, public_observations, picture, external_auth_provider)
    VALUES(DEFAULT, '8427613950842761395075', 'Jagoda Mlinaric', 'Jagoda', 'Mlinaric', 'jagoda@fakemail.com', TRUE, 'https://doesnt.exist.com/path/to/picture.png', 1);
INSERT INTO biolog_user (id, external_id, display_name, given_name, family_name, email, public_observations, picture, external_auth_provider)
    VALUES(DEFAULT, '8146079532814607953221', 'Klementina Koblaric', 'Klementina', 'Koblaric', 'klementina@fakemail.com', TRUE, 'https://doesnt.exist.com/path/to/picture.png', 1);

-- [observation]
INSERT INTO observation(id, sighting_time, sighting_location, quantity, public_visibility, biolog_user, species)
    VALUES (DEFAULT, now(), ST_GeomFromText('POINT(-71.060316 48.432044)'), 6, TRUE, 10000003, 1);