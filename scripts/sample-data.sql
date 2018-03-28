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

-- [user] 
INSERT INTO biolog_user VALUES(DEFAULT, 'Malcolm Reynolds', TRUE);
INSERT INTO biolog_user VALUES(DEFAULT, 'Hoban Washburne', FALSE);
INSERT INTO biolog_user VALUES(DEFAULT, 'Zoe Washburne', TRUE);
INSERT INTO biolog_user VALUES(DEFAULT, 'Kaylee Frye', TRUE);
INSERT INTO biolog_user VALUES(DEFAULT, 'Jayne Cobb', FALSE);
INSERT INTO biolog_user VALUES(DEFAULT, 'Derrial Book', FALSE);

-- [external_user]
INSERT INTO external_user VALUES(DEFAULT, 4567891920, 'Zoe', 'Washburne', 
    'zoe.washburne@gmail.com', 'https://vignette.wikia.nocookie.net/firefly/images/1/10/Zoe.jpg', 1, 3);
INSERT INTO external_user VALUES(DEFAULT, 9302175329, 'Kaywinnet', 'Lee Frye', 
    'kaylee4@gmail.com', 'https://vignette.wikia.nocookie.net/firefly/images/4/44/Kaylee_closeup.jpg', 1, 4);
INSERT INTO external_user VALUES(DEFAULT, 8573482930, 'Derrial', 'Book', 
    'shepherd.book@gmail.com', 'https://upload.wikimedia.org/wikipedia/commons/thumb/f/fc/Ron_Glass_Serenity_premiere_1.jpg/220px-Ron_Glass_Serenity_premiere_1.jpg',
    1, 6);

-- [observation]
INSERT INTO observation(id, sighting_time, sighting_location, quantity, public_visibility, biolog_user, species)
    VALUES (DEFAULT, now(), ST_GeomFromText('POINT(-71.060316 48.432044)'), 6, TRUE, 3, 1);