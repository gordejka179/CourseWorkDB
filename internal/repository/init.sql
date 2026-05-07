DROP TABLE IF EXISTS BBKAlternative CASCADE;
DROP TABLE IF EXISTS BBKMapping CASCADE;
DROP TABLE IF EXISTS Copy CASCADE;
DROP TABLE IF EXISTS BBKRecord CASCADE;
DROP TABLE IF EXISTS BookAuthor CASCADE;
DROP TABLE IF EXISTS OtherIndex CASCADE;
DROP TABLE IF EXISTS ISBNOther CASCADE;
DROP TABLE IF EXISTS ISBN CASCADE;
DROP TABLE IF EXISTS Reader CASCADE;
DROP TABLE IF EXISTS LibraryBuilding CASCADE;
DROP TABLE IF EXISTS BBKDictionary CASCADE;
DROP TABLE IF EXISTS Publication CASCADE;
DROP TABLE IF EXISTS Author CASCADE;
DROP TABLE IF EXISTS Librarian CASCADE;


CREATE TABLE Publication (
    publicationId SERIAL PRIMARY KEY,
    title VARCHAR NOT NULL,
    publicationYear INT
);


CREATE TABLE ISBN (
    ISBN VARCHAR(18) PRIMARY KEY,
    publicationId INT NOT NULL,
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE
);

CREATE TABLE ISBNOther (
    publicationId INT NOT NULL,
    ISBN VARCHAR(18) NOT NULL,
    PRIMARY KEY (publicationId, ISBN),
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE
);

CREATE TABLE BBKDictionary (
    BBK VARCHAR(100) PRIMARY KEY
);


CREATE TABLE BBKAlternative (
    sourceCode VARCHAR(100),
    targetCode VARCHAR(100),
    PRIMARY KEY (sourceCode, targetCode),
    FOREIGN KEY (sourceCode) REFERENCES BBKDictionary(BBK) ON DELETE CASCADE,
    FOREIGN KEY (targetCode) REFERENCES BBKDictionary(BBK) ON DELETE CASCADE
);

CREATE TABLE BBKMapping (
    fullTableCode VARCHAR(100),
    midTableCode VARCHAR(100),
    PRIMARY KEY (fullTableCode, midTableCode),
    FOREIGN KEY (fullTableCode) REFERENCES BBKDictionary(BBK) ON DELETE CASCADE
);


CREATE TABLE BBKRecord (
    publicationId INT NOT NULL,
    BBK VARCHAR(100) NOT NULL,
    PRIMARY KEY (publicationId, BBK),
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE,
    FOREIGN KEY (BBK) REFERENCES BBKDictionary(BBK) ON DELETE CASCADE
);

CREATE TABLE OtherIndex (
    publicationId INT NOT NULL,
    index VARCHAR NOT NULL,
    PRIMARY KEY (publicationId, index),
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE
);




CREATE TABLE Author (
    authorId   SERIAL PRIMARY KEY,
    birthDate  DATE NOT NULL,
    firstName  VARCHAR(50) NOT NULL,
    lastName   VARCHAR(50) NOT NULL,
    patronymic VARCHAR(50),
    UNIQUE (birthDate, firstName, lastName)
);

CREATE TABLE BookAuthor (
    publicationId INT NOT NULL,
    authorId      INT NOT NULL,
    PRIMARY KEY (publicationId, authorId),
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE,
    FOREIGN KEY (authorId) REFERENCES Author(authorId) ON DELETE CASCADE
);

CREATE TABLE LibraryBuilding (
    libraryBuildingId SERIAL PRIMARY KEY,
    address           VARCHAR NOT NULL UNIQUE,
    description       VARCHAR
);

CREATE TABLE Librarian (
    librarianId   SERIAL PRIMARY KEY,
    staffNum      VARCHAR(10) NOT NULL UNIQUE,
    email         VARCHAR(254) NOT NULL UNIQUE,
    passwordHash  VARCHAR(32) NOT NULL,
    firstName     VARCHAR(100) NOT NULL,
    lastName      VARCHAR(100) NOT NULL,
    patronymic    VARCHAR(100)
);

CREATE TABLE Reader (
    readerId     SERIAL PRIMARY KEY,
    email        VARCHAR(254) NOT NULL UNIQUE,
    libraryCard  VARCHAR(12) NOT NULL UNIQUE,
    passportSeries VARCHAR(4) NOT NULL,
    passportNumber VARCHAR(6) NOT NULL,
    firstName    VARCHAR(100) NOT NULL,
    lastName     VARCHAR(100) NOT NULL,
    patronymic   VARCHAR(100),
    passwordHash VARCHAR(128) NOT NULL,
    UNIQUE (passportSeries, passportNumber)
);

CREATE TABLE Copy (
    copyId         SERIAL PRIMARY KEY,
    inventoryNumber VARCHAR(13) NOT NULL UNIQUE,
    publicationId   INT NOT NULL,
    buildingId      INT NOT NULL,
    readerId        INT,
    librarianId     INT,
    startDate       DATE,
    expiryDate      DATE,
    FOREIGN KEY (publicationId) REFERENCES Publication(publicationId) ON DELETE CASCADE,
    FOREIGN KEY (buildingId) REFERENCES LibraryBuilding(libraryBuildingId) ON DELETE CASCADE,
    FOREIGN KEY (readerId) REFERENCES Reader(readerId) ON DELETE SET NULL,
    FOREIGN KEY (librarianId) REFERENCES Librarian(librarianId) ON DELETE SET NULL
);



INSERT INTO Author (birthDate, firstName, lastName, patronymic) VALUES
('1956-09-16', 'Кирилл', 'Еськов', 'Юрьевич'),

--лидары:
('1970-09-16', 'Иван', 'Иванов', 'Иванович'),
('1971-09-16', 'Пётр', 'Петров', NULL);

INSERT INTO LibraryBuilding (address, description) VALUES
('г. Москва, ул. Б. Дмитровка, д. 5/6, стр. 7', 'Библиотека номер 179'),
('г. Москва, М. Знаменский пер., д. 7/10, стр. 5', 'Библиотека номер 57');

INSERT INTO Librarian (staffNum, email, passwordHash, firstName, lastName, patronymic) VALUES
('LIB0000001', '1@179.ru', MD5('1'), 'л', 'д', 'м'),
('LIB0000002', 'kligunov@179.ru', MD5('1'), 'Клигунов', 'Кирилл', 'Дмитриевич');

INSERT INTO Reader (email, libraryCard, passportSeries, passportNumber, firstName, lastName, patronymic, passwordHash) VALUES
('1@1', '000000000001', '1944', '111111', 'И', 'В', 'Г', MD5('1')),
('reader1@mail.ru', '000000000002', '1945', '111112', 'Семён', 'Георгиевич', 'Чайкин', 'hash_reader1');

INSERT INTO Publication (title, publicationYear) VALUES
--Еськов:
('Удивительная палеонтология: история Земли и жизни на ней', 2008),
('Удивительная палеонтология: история Земли и жизни на ней', 2014),

--лидары (выдуманный пример)
('Алгоритмы, применяемые в лидарах', 2020),

('Лидары', 2025);

INSERT INTO ISBN (ISBN, publicationId) VALUES
--Еськов:
('978-5-93196-711-0', 1),
('978-5-91921-129-7', 2),

--лидары (выдуманный пример)
('978-5-99999-999-9', 3),
('978-5-99999-999-0', 3);


INSERT INTO ISBNOther (publicationId, ISBN) VALUES
(2, '978-5-93196-711-0');

INSERT INTO OtherIndex (publicationId, index) VALUES
(2, '56'),
(2, ' ГРНТИ3543');


INSERT INTO BBKDictionary (BBK) VALUES
--книги на комбинаторику:
('В181'),

('В181.1'),

('В181.11'),

('В181.12'),

('В181.13'),

('В181.14'),

('В181.19'),

--про лидары:

('З81'), 

('З859'),

('З956-5'),

--Еськов:
('Е1');



INSERT INTO BBKRecord (publicationId, BBK) VALUES
(1, 'Е1'),
(2, 'Е1'),

--про лидары:
(3, 'З81'),

(4, 'З956-5');

INSERT INTO BookAuthor (publicationId, authorId) VALUES
(1, 1),
(2, 1),
(3, 2),
(3, 3),

(2, 3);

INSERT INTO Copy (inventoryNumber, publicationId, buildingId, readerId, librarianId, startDate, expiryDate) VALUES
('INV0000000001', 1, 1, NULL, NULL, NULL, NULL),
('INV0000000002', 1, 1, NULL, NULL, NULL, NULL),
('INV0000000003', 1, 1, NULL, NULL, NULL, NULL),
('INV0000000004', 1, 2, NULL, NULL, NULL, NULL),
('INV0000000005', 2, 2, NULL, NULL, NULL, NULL),
('INV0000000006', 3, 2, NULL, NULL, NULL, NULL),
('INV0000000007', 4, 2, NULL, NULL, NULL, NULL);


INSERT INTO BBKAlternative (sourceCode, targetCode) VALUES
--пример с лидарами
('З956-5', 'З859'), 

('З956-5', 'З81');



INSERT INTO BBKMapping (fullTableCode, midTableCode) VALUES
--книги на комбинаторику:
('В181', '22.181'),

('В181.1', '22.181.1'),

('В181.11', '22.181.11'),

('В181.12', '22.181.12'),

('В181.13', '22.181.13'),

('В181.14', '22.181.14'),

('В181.19', '22.181.19'),

--Про лидары:

('З81', '32.81'), 

('З859', '32.859'),

('З956-5', '32.956-5'),

--Еськов:
('Е1', '28.1');








--представления



--индексы






--функции, процедуры


CREATE OR REPLACE FUNCTION createReader(
    newEmail VARCHAR(254),
    newPasswordHash VARCHAR(128),
    newFirstName VARCHAR(100),
    newLastName VARCHAR(100),
    newPassportSeries VARCHAR(4),
    newPassportNumber VARCHAR(6),
    newPatronymic VARCHAR(100) DEFAULT NULL
)
RETURNS TABLE(readerId INTEGER, libraryCard VARCHAR(12))
LANGUAGE plpgsql
AS $$
DECLARE
    lastCard VARCHAR(12);
    nextNum INTEGER; -- читательский номер (число) после инкремента
    prefix CONSTANT VARCHAR(3) := 'LIB';
    newCard VARCHAR(12); -- читательский номер, который получит новый читатель
    newId INTEGER;
BEGIN

    -- Проверка на существующий email
    IF EXISTS (SELECT 1 FROM Reader WHERE email = newEmail) THEN
        RAISE EXCEPTION 'Email already exists' USING ERRCODE = 'EM001';
    END IF;

    -- Проверка на существующий паспорт
    IF EXISTS (SELECT 1 FROM Reader 
               WHERE passportSeries = newPassportSeries 
                 AND passportNumber = newPassportNumber) THEN
        RAISE EXCEPTION 'Passport already exists' USING ERRCODE = 'PS001'; 
    END IF;



    -- Берем самый последний выданный читательский билет
    SELECT Reader.libraryCard INTO lastCard
    FROM Reader
    ORDER BY Reader.readerId DESC
    LIMIT 1;

    IF lastCard IS NULL THEN
        -- Если читателей нет в базе
        nextNum := 1;
    ELSE
        nextNum := CAST(substring(lastCard FROM 4) AS INTEGER) + 1;
    END IF;

    newCard := 'LIB' || nextNum::TEXT;
    

    INSERT INTO Reader (email, libraryCard, passwordHash, firstName, lastName, patronymic, passportSeries, passportNumber)
    VALUES (newEmail, newCard, newPasswordHash, newFirstName, newLastName, newPatronymic, newPassportSeries, newPassportNumber)
    RETURNING Reader.readerId INTO newId;

    RETURN QUERY SELECT newId, newCard;
END;
$$;




CREATE OR REPLACE FUNCTION checkReaderCredentials(
    p_email VARCHAR(254),
    p_password VARCHAR(32)
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 
        FROM Reader 
        WHERE email = p_email AND passwordHash = p_password
    );
END;
$$;

CREATE OR REPLACE FUNCTION checkLibrarianCredentials(
    p_email VARCHAR(254),
    p_password VARCHAR(32)
)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    RETURN EXISTS (
        SELECT 1 
        FROM Librarian 
        WHERE email = p_email AND passwordHash = p_password
    );
END;
$$;


--ищем издания по isbn + смотрим на ISBNOther
CREATE OR REPLACE FUNCTION search_publications_by_isbn(p_isbn VARCHAR)
RETURNS TABLE(
    publicationId INT,
    title VARCHAR,
    publicationYear INT,
    isbns TEXT[],
    BBKs TEXT[],
    otherIndexes TEXT[],
    authors TEXT[]
)
LANGUAGE sql
AS $$

    --находим издание с нужным isbn
    WITH main_pub_id AS(
        SELECT publicationId, isbn
        FROM ISBN
        WHERE ISBN = p_isbn
    ),

    main_pub AS (
        SELECT i.publicationId, i.isbn
        FROM ISBN i
        WHERE i.publicationId IN (SELECT publicationId FROM main_pub_id)
    ),

    --для найденного издания ищем все isbn, с которыми это издание связано
    other_isbn_to_publication AS (
        SELECT isbn_other.ISBN
        FROM ISBNOther isbn_other
        JOIN main_pub mp ON mp.publicationId = isbn_other.PublicationId
    ),

    --для номеров isbn из other_isbn_to_publication ищем, какие к ним относятся издания
    other_publications AS (
        SELECT isbn.publicationId, isbn.ISBN
        FROM ISBN isbn
        JOIN other_isbn_to_publication oisbn ON isbn.ISBN = oisbn.ISBN
    ),

    --теперь собираем вместе все id изданий и их isbn
    all_isbns_for_publication AS (
        SELECT publicationId, isbn 
        FROM main_pub
        UNION
        SELECT DISTINCT publicationId, ISBN FROM other_publications 
    )

    --теперь собираем информацию о всех изданиях
    SELECT 
        p.publicationId,
        p.title,
        p.publicationYear,
        (SELECT array_agg(aip.isbn)
            FROM all_isbns_for_publication aip
            WHERE aip.publicationId = p.publicationId) AS isbns,

        (SELECT array_agg(br.BBK)
            FROM BBKRecord br WHERE br.publicationId = p.publicationId) AS BBKs,

        (SELECT array_agg(oi.index)
            FROM OtherIndex oi WHERE oi.publicationId = p.publicationId) AS otherIndexes,

        (SELECT array_agg(a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
            FROM BookAuthor ba 
            JOIN Author a ON ba.authorId = a.authorId
            WHERE ba.publicationId = p.publicationId) AS authors
    FROM Publication p
    JOIN (SELECT DISTINCT publicationId FROM all_isbns_for_publication) AS aip ON p.publicationId = aip.publicationId
$$;


--ищем издания по автору
CREATE OR REPLACE FUNCTION search_publications_by_author(
    p_lastname VARCHAR,
    p_firstname VARCHAR DEFAULT NULL,
    p_patronymic VARCHAR DEFAULT NULL
)
RETURNS TABLE(
    publicationId INT,
    title VARCHAR,
    publicationYear INT,
    isbns TEXT[],
    bbks TEXT[],
    otherIndexes TEXT[],
    authors TEXT[]
)
LANGUAGE sql
AS $$
    SELECT 
        p.publicationId,
        p.title,
        p.publicationYear,
        COALESCE(
            (SELECT array_agg(i.isbn) FROM ISBN i WHERE i.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS isbns,
        COALESCE(
            (SELECT array_agg(br.BBK) FROM BBKRecord br WHERE br.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS bbks,
        COALESCE(
            (SELECT array_agg(oi.index) FROM OtherIndex oi WHERE oi.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS otherIndexes,
        COALESCE(
            (SELECT array_agg(a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
                FROM BookAuthor ba 
                JOIN Author a ON ba.authorId = a.authorId
                WHERE ba.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS authors
    FROM Publication p
    WHERE EXISTS (
        SELECT 1
        FROM BookAuthor ba
        JOIN Author a ON ba.authorId = a.authorId
        WHERE ba.publicationId = p.publicationId
            AND a.lastName ILIKE '%' || p_lastname || '%'
            AND (p_firstname = '' OR a.firstName ILIKE '%' || p_firstname || '%')
            AND (p_patronymic = '' OR a.patronymic ILIKE '%' || p_patronymic || '%')
    );
$$;


--поиск изданий по названию
CREATE OR REPLACE FUNCTION search_publications_by_title(p_title VARCHAR)
RETURNS TABLE(
    publicationId INT,
    title VARCHAR,
    publicationYear INT,
    isbns TEXT[],
    BBKs TEXT[],
    otherIndexes TEXT[],
    authors TEXT[]
)
LANGUAGE sql
AS $$
    SELECT 
        p.publicationId,
        p.title,
        p.publicationYear,
        (SELECT array_agg(isbn.ISBN)
            FROM ISBN isbn 
            WHERE isbn.publicationId = p.publicationId) AS isbns,

        (SELECT array_agg(br.BBK)
            FROM BBKRecord br 
            WHERE br.publicationId = p.publicationId) AS BBKs,

        (SELECT array_agg(oi.index)
            FROM OtherIndex oi 
            WHERE oi.publicationId = p.publicationId) AS otherIndexes,

        (SELECT array_agg(a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
            FROM BookAuthor ba 
            JOIN Author a ON ba.authorId = a.authorId
            WHERE ba.publicationId = p.publicationId) AS authors
    FROM Publication p
    WHERE p.title ILIKE '%' || COALESCE(p_title, '') || '%'  --ILIKE - для совпадений без учета регистра
$$;


--по массиву из id публикаций получаем инфу про экземпляры
CREATE OR REPLACE FUNCTION get_copies_info_by_ids(p_ids INT[])
RETURNS TABLE(
    copyId INT,
    publicationId INT,
    buildingId INT,
    readerId INT,
    librarianId INT,
    address VARCHAR,
    description VARCHAR
) AS $$
BEGIN
    RETURN QUERY
    SELECT 
        c.copyId,
        c.publicationId, 
        c.buildingId,
        c.readerId, 
        c.librarianId,
        lb.address, 
        lb.description
    FROM Copy c
    JOIN LibraryBuilding lb ON c.buildingId = lb.libraryBuildingId
    WHERE c.publicationId = ANY(p_ids); --функция ANY, для проверки, что совпадает хотя бы с одним значением из массива
END;
$$ LANGUAGE plpgsql;


--Получение полных кодов ббк
CREATE OR REPLACE FUNCTION get_full_codes_by_mid(mid_codes TEXT[])
RETURNS TABLE(full_code TEXT)
LANGUAGE sql
AS $$
    SELECT DISTINCT fullTableCode
    FROM BBKMapping
    WHERE midTableCode = ANY(mid_codes);
$$;

--Получение дополнительных кодов ббк
CREATE OR REPLACE FUNCTION get_alternative_codes_by_source(source_codes TEXT[])
RETURNS TABLE(target_code TEXT)
LANGUAGE sql
AS $$
    SELECT DISTINCT targetCode
    FROM BBKAlternative
    WHERE sourceCode = ANY(source_codes);
$$;

--получение изданий по ббк
CREATE OR REPLACE FUNCTION search_publications_by_bbk(prefixes TEXT[])
RETURNS TABLE(
    publicationId INT,
    title VARCHAR,
    publicationYear INT,
    isbns TEXT[],
    bbks TEXT[],
    otherIndexes TEXT[],
    authors TEXT[]
)
LANGUAGE sql
AS $$
    SELECT 
        p.publicationId,
        p.title,
        p.publicationYear,
        COALESCE(
            (SELECT array_agg(i.isbn) FROM ISBN i WHERE i.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS isbns,
        COALESCE(
            (SELECT array_agg(br.BBK) FROM BBKRecord br WHERE br.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS bbks,
        COALESCE(
            (SELECT array_agg(oi.index) FROM OtherIndex oi WHERE oi.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS otherIndexes,
        COALESCE(
            (SELECT array_agg(a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
             FROM BookAuthor ba 
             JOIN Author a ON ba.authorId = a.authorId
             WHERE ba.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS authors
    FROM Publication p
    WHERE 
        COALESCE(array_length(prefixes, 1), 0) > 0 --1 значит, что массив одномерен
        AND EXISTS (
            SELECT 1
            FROM BBKRecord br
            WHERE br.publicationId = p.publicationId
                AND EXISTS ( --для каждого префикса проверяем, что код ббк для издания начинается с этого префикса
                    SELECT 1
                    FROM unnest(prefixes) AS prefix
                    WHERE br.BBK LIKE prefix || '%'
                )
        )
$$;


CREATE OR REPLACE FUNCTION search_publications_by_other_index(p_index TEXT)
RETURNS TABLE(
    publicationId INT,
    title VARCHAR,
    publicationYear INT,
    isbns TEXT[],
    bbks TEXT[],
    otherIndexes TEXT[],
    authors TEXT[]
)
LANGUAGE sql
AS $$
    SELECT 
        p.publicationId,
        p.title,
        p.publicationYear,
        COALESCE(
            (SELECT array_agg(i.isbn) FROM ISBN i WHERE i.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS isbns,
        COALESCE(
            (SELECT array_agg(br.BBK) FROM BBKRecord br WHERE br.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS bbks,
        COALESCE(
            (SELECT array_agg(oi.index) FROM OtherIndex oi WHERE oi.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS otherIndexes,
        COALESCE(
            (SELECT array_agg(a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
             FROM BookAuthor ba 
             JOIN Author a ON ba.authorId = a.authorId
             WHERE ba.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ) AS authors
    FROM Publication p
    WHERE EXISTS (
        SELECT 1 
        FROM OtherIndex oi
        WHERE oi.publicationId = p.publicationId
          AND oi.index ILIKE '%' || COALESCE(p_index, '') || '%'
    );
$$;


--Бронирование экземпляра по id читателя и id экземпляра. Бронируем на 3 дня
CREATE OR REPLACE FUNCTION reserveCopyByEmail(p_readerId INT, p_copyId INT)
RETURNS BOOLEAN AS $$
DECLARE
    v_readerId INT;
    v_librarianId INT;
    v_publicationId INT;
BEGIN
    -- Проверяем текущее состояние экземпляра
    SELECT readerId, librarianId, publicationId 
    INTO v_readerId, v_librarianId, v_publicationId
    FROM Copy
    WHERE copyId = p_copyId;

    IF NOT FOUND THEN
        RAISE EXCEPTION 'Экземпляр с ID % не найден', p_copyId
        USING ERRCODE = 'BK001';
    END IF;

    -- Если хотя бы один из readerId или librarianId не NULL, то экземпляр уже занят
    IF v_readerId IS NOT NULL OR v_librarianId IS NOT NULL THEN
        RAISE EXCEPTION 'Экземпляр уже занят (выдан или забронирован)'
        USING ERRCODE = 'BK002';
    END IF;

    -- Проверка, что на то, имеет ли читатель уже имеет экземпляр этого издания (на руках или забронирован) 
    IF EXISTS (
        SELECT 1
        FROM Copy
        WHERE publicationId = v_publicationId
          AND readerId = p_readerId
    ) THEN
        RAISE EXCEPTION 'Читатель уже имеет экземпляр этого издания (на руках или забронирован)'
        USING ERRCODE = 'BK003';
    END IF;

    UPDATE Copy
    SET readerId   = p_readerId,
        startDate  = CURRENT_DATE,
        expiryDate = CURRENT_DATE + 3
    WHERE copyId = p_copyId;

    RETURN TRUE;
END;
$$ LANGUAGE plpgsql;


--Получаем бронирования читателя по его id
CREATE OR REPLACE FUNCTION get_current_bookings_by_readerId(p_readerId INT)
RETURNS TABLE(
    copyId INT,
    inventoryNumber VARCHAR,
    title VARCHAR,
    publicationYear INT,
    authors TEXT[],
    isbns TEXT[],
    bbks TEXT[],
    otherIndexes TEXT[],
    buildingId INT,
    buildingAddress VARCHAR,
    expiryDate VARCHAR
) LANGUAGE sql 
AS $$
    SELECT 
        c.copyId,
        c.inventoryNumber,
        p.title,
        p.publicationYear,
        COALESCE(
            (SELECT array_agg(DISTINCT a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
             FROM BookAuthor ba 
             JOIN Author a ON ba.authorId = a.authorId 
             WHERE ba.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT i.isbn) FROM ISBN i WHERE i.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT br.BBK) FROM BBKRecord br WHERE br.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT oi.index) FROM OtherIndex oi WHERE oi.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        lb.libraryBuildingId,
        lb.address,
        c.expiryDate
    FROM Copy c
    JOIN Publication p ON c.publicationId = p.publicationId
    JOIN LibraryBuilding lb ON c.buildingId = lb.libraryBuildingId
    WHERE c.readerId = p_readerId
      AND c.librarianId IS NULL
      AND c.expiryDate >= CURRENT_DATE;
$$;


--Выдача экземляра (выдаём на 30 дней)
CREATE OR REPLACE FUNCTION makeLoan(p_readerId INT, p_librarianId INT, p_copyId INT)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
DECLARE
    v_readerId INT;
    v_librarianId INT;
    v_expiryDate DATE;
    v_publicationId INT;
    v_hasOther BOOLEAN;
BEGIN

    SELECT readerId, librarianId, expiryDate, publicationId
    INTO v_readerId, v_librarianId, v_expiryDate, v_publicationId
    FROM Copy
    WHERE copyId = p_copyId;

    --ВАЖНО:
    --по идее, библиотекарь при выдаче будет сканировать какой-нибудь код на экземпляре
    --и это будет означать, что экземляр есть и ещё никому не выдан (максимум забронирован)
    --но всё равно сделаем проверки


    -- Проверка существования экземпляра
    IF NOT FOUND THEN
        RAISE EXCEPTION 'Экземпляр с id % не найден', p_copyId
        USING ERRCODE = 'ML001';
    END IF;

    -- Проверка, что экземпляр не забронирован другим человеком или забронирован другим человеком, но время брони уже вышло
    IF v_readerId IS NOT NULL AND v_readerId != p_readerId AND v_expiryDate < CURRENT_DATE THEN
        RAISE EXCEPTION 'Экземпляр кем-то забронирован или кому-то выдан'
        USING ERRCODE = 'ML002';
    END IF;


     -- Проверка, нет ли у читателя других экземпляров той же публикации на руках
    SELECT EXISTS (
        SELECT 1
        FROM Copy
        WHERE publicationId = v_publicationId
          AND (readerId = p_readerId AND librarianId IS NOT NULL)
    ) INTO v_hasOther;

    IF v_hasOther THEN
        RAISE SQLSTATE 'ML003' USING MESSAGE = 'У читателя на руках уже есть экземпляр этого издания';
    END IF;

    --выдача
    UPDATE Copy
    SET readerId   = p_readerId,
        librarianId = p_librarianId,
        startDate   = CURRENT_DATE,
        expiryDate  = CURRENT_DATE + INTERVAL '30 days'
    WHERE copyId = p_copyId;

    RETURN TRUE;
END;
$$;


--Получение книг, которые есть у читателя на руках
CREATE OR REPLACE FUNCTION get_current_loans_by_readerId(p_readerId INT)
RETURNS TABLE(
    copyId INT,
    expiryDate DATE,
    inventoryNumber VARCHAR,
    title VARCHAR,
    publicationYear INT,
    authors TEXT[],
    isbns TEXT[],
    bbks TEXT[],
    otherIndexes TEXT[],
    buildingId INT,
    buildingAddress VARCHAR
) LANGUAGE sql AS $$
    SELECT 
        c.copyId,
        c.expiryDate,
        c.inventoryNumber,
        p.title,
        p.publicationYear,
        COALESCE(
            (SELECT array_agg(DISTINCT a.lastName || '|' || a.firstName || COALESCE('|' || a.patronymic, '|'))
                FROM BookAuthor ba 
                JOIN Author a ON ba.authorId = a.authorId 
                WHERE ba.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT i.isbn) FROM ISBN i WHERE i.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT br.BBK) FROM BBKRecord br WHERE br.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        COALESCE(
            (SELECT array_agg(DISTINCT oi.index) FROM OtherIndex oi WHERE oi.publicationId = p.publicationId),
            ARRAY[]::TEXT[]
        ),
        lb.libraryBuildingId,
        lb.address
    FROM Copy c
    JOIN Publication p ON c.publicationId = p.publicationId
    JOIN LibraryBuilding lb ON c.buildingId = lb.libraryBuildingId
    WHERE c.readerId = p_readerId
      AND c.librarianId IS NOT NULL;
$$;

--возврат книги
CREATE OR REPLACE FUNCTION return_copy(p_copy_id INT, p_reader_id INT)
RETURNS BOOLEAN
LANGUAGE plpgsql
AS $$
BEGIN
    UPDATE copy
    SET readerId   = NULL,
        librarianId = NULL,
        startDate   = NULL,
        expiryDate  = NULL
    WHERE copyId = p_copy_id AND readerId = p_reader_id;

    RETURN FOUND; --если строка обновилась, то вернём true
END;
$$;