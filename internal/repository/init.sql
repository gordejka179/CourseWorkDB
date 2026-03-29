DROP TABLE IF EXISTS Loan CASCADE;
DROP TABLE IF EXISTS Reservation CASCADE;
DROP TABLE IF EXISTS Copy CASCADE;
DROP TABLE IF EXISTS BookAuthor CASCADE;
DROP TABLE IF EXISTS Book CASCADE;
DROP TABLE IF EXISTS Author CASCADE;
DROP TABLE IF EXISTS Publisher CASCADE;
DROP TABLE IF EXISTS LibraryBuilding CASCADE;
DROP TABLE IF EXISTS Reader CASCADE;
DROP TABLE IF EXISTS Librarian CASCADE;

CREATE TABLE Librarian (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    Patronymic VARCHAR(100),
    Is_admin BOOLEAN DEFAULT FALSE
);

CREATE TABLE Author (
    id SERIAL PRIMARY KEY,
    birth_date DATE,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    Patronymic VARCHAR(100),
    description TEXT
);

CREATE TABLE Publisher (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    address VARCHAR(255),
);

CREATE TABLE Book (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    bbk VARCHAR(50),
    udk VARCHAR(50),
    isbn VARCHAR(20) UNIQUE,
    publication_year INT,
    publisher_id INT REFERENCES Publisher(id) ON DELETE SET NULL
);

CREATE TABLE BookAuthor (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES Book(id) ON DELETE CASCADE,
    author_id INT NOT NULL REFERENCES Author(id) ON DELETE CASCADE,
    UNIQUE(book_id, author_id)
);

CREATE TABLE LibraryBuilding (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    address VARCHAR(255) NOT NULL,
    description TEXT
);

CREATE TABLE Copy (
    id SERIAL PRIMARY KEY,
    inventory_number VARCHAR(50) UNIQUE NOT NULL,
    acquisition_date DATE,
    book_id INT NOT NULL REFERENCES Book(id) ON DELETE CASCADE,
    building_id INT NOT NULL REFERENCES LibraryBuilding(id) ON DELETE CASCADE
);

CREATE TABLE Reader (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    library_card_number VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    FirstName VARCHAR(100) NOT NULL,
    LastName VARCHAR(100) NOT NULL,
    Patronymic VARCHAR(100),
    Address VARCHAR(255)
);

CREATE TABLE Reservation (
    id SERIAL PRIMARY KEY,
    reservation_date TIMESTAMP NOT NULL DEFAULT NOW(),
    expiry_date TIMESTAMP NOT NULL,
    reader_id INT NOT NULL REFERENCES Reader(id) ON DELETE CASCADE,
    copy_id INT NOT NULL REFERENCES Copy(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed'))
);

CREATE TABLE Loan (
    id SERIAL PRIMARY KEY,
    loan_date DATE NOT NULL,
    due_date DATE NOT NULL,
    return_date DATE,
    reader_id INT NOT NULL REFERENCES Reader(id) ON DELETE CASCADE,
    copy_id INT NOT NULL REFERENCES Copy(id) ON DELETE CASCADE,
    issued_by_id INT REFERENCES Librarian(id) ON DELETE SET NULL,
    returned_by_id INT REFERENCES Librarian(id) ON DELETE SET NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'returned'))
);

CREATE INDEX idx_book_title ON Book(title);
CREATE INDEX idx_book_bbk ON Book(bbk);
CREATE INDEX idx_book_udk ON Book(udk);
CREATE INDEX idx_book_isbn ON Book(isbn);
CREATE INDEX idx_book_year ON Book(publication_year);
CREATE INDEX idx_author_name ON Author(LastName, FirstName);
CREATE INDEX idx_publisher_name ON Publisher(name);
CREATE INDEX idx_copy_book ON Copy(book_id);
CREATE INDEX idx_copy_building ON Copy(building_id);
CREATE INDEX idx_reservation_copy_status ON Reservation(copy_id, status);
CREATE INDEX idx_loan_copy_status ON Loan(copy_id, status);

INSERT INTO Librarian (email, password_hash, FirstName, LastName, Patronymic, Is_admin) VALUES
('admin@library.ru', 'hashed_password_1', 'Иван', 'Иванов', 'Иванович', true),
('librarian1@library.ru', 'hashed_password_2', 'Мария', 'Петрова', 'Сергеевна', false),
('librarian2@library.ru', 'hashed_password_3', 'Алексей', 'Смирнов', 'Викторович', false);

INSERT INTO LibraryBuilding (name, Address, description) VALUES
('Главное здание', 'ул. Университетская, 1, г. Москва', 'Основной корпус, отдел художественной литературы'),
('Научная библиотека', 'пр. Ленина, 22, г. Москва', 'Научный фонд, редкие книги'),
('Детский отдел', 'ул. Пушкина, 10, г. Москва', 'Литература для детей и подростков');

INSERT INTO Publisher (name, Address) VALUES
('Эксмо', 'г. Москва, ул. Клары Цеткин, 33',),
('АСТ', 'г. Москва, ул. Правды, 15'),
('Наука', 'г. Москва, ул. Профсоюзная, 90'),
('Oxford University Press', 'Oxford, UK');

INSERT INTO Author (birth_date, FirstName, LastName, Patronymic, description) VALUES
('1828-09-09', 'Лев', 'Толстой', 'Николаевич', 'Великий русский писатель'),
('1821-11-11', 'Фёдор', 'Достоевский', 'Михайлович', 'Русский писатель, мыслитель'),
('1903-06-25', 'Джордж', 'Оруэлл', '', 'Английский писатель и публицист'),
('1891-05-15', 'Михаил', 'Булгаков', 'Афанасьевич', 'Русский писатель, драматург'),
('1775-12-16', 'Джейн', 'Остин', '', 'Английская писательница');

INSERT INTO Book (title, bbk, udc, isbn, publication_year, publisher_id) VALUES
('Война и мир', '84(2Рос)1', '821.161.1', '978-5-699-12345-6', 2008, (SELECT id FROM Publisher WHERE name='Эксмо')),
('Преступление и наказание', '84(2Рос)1', '821.161.1', '978-5-17-12345-7', 2015, (SELECT id FROM Publisher WHERE name='АСТ')),
('1984', '84(4Вел)', '821.111', '978-5-699-98765-4', 2010, (SELECT id FROM Publisher WHERE name='Эксмо')),
('Мастер и Маргарита', '84(2Рос)1', '821.161.1', '978-5-17-98765-4', 2012, (SELECT id FROM Publisher WHERE name='АСТ')),
('Гордость и предубеждение', '84(4Вел)', '821.111', '978-0-19-953556-9', 2006, (SELECT id FROM Publisher WHERE name='Oxford University Press'));

INSERT INTO BookAuthor (book_id, author_id) VALUES
((SELECT id FROM Book WHERE isbn='978-5-699-12345-6'), (SELECT id FROM Author WHERE LastName='Толстой')),
((SELECT id FROM Book WHERE isbn='978-5-17-12345-7'), (SELECT id FROM Author WHERE LastName='Достоевский')),
((SELECT id FROM Book WHERE isbn='978-5-699-98765-4'), (SELECT id FROM Author WHERE LastName='Оруэлл')),
((SELECT id FROM Book WHERE isbn='978-5-17-98765-4'), (SELECT id FROM Author WHERE LastName='Булгаков')),
((SELECT id FROM Book WHERE isbn='978-0-19-953556-9'), (SELECT id FROM Author WHERE LastName='Остин'));

INSERT INTO Copy (inventory_number, acquisition_date, book_id, building_id) VALUES
('INV-001', '2020-01-15', (SELECT id FROM Book WHERE isbn='978-5-699-12345-6'), (SELECT id FROM LibraryBuilding WHERE name='Главное здание')),
('INV-002', '2020-02-10', (SELECT id FROM Book WHERE isbn='978-5-699-12345-6'), (SELECT id FROM LibraryBuilding WHERE name='Научная библиотека')),
('INV-003', '2021-03-20', (SELECT id FROM Book WHERE isbn='978-5-17-12345-7'), (SELECT id FROM LibraryBuilding WHERE name='Главное здание')),
('INV-004', '2021-04-05', (SELECT id FROM Book WHERE isbn='978-5-17-12345-7'), (SELECT id FROM LibraryBuilding WHERE name='Детский отдел')),
('INV-005', '2022-05-12', (SELECT id FROM Book WHERE isbn='978-5-699-98765-4'), (SELECT id FROM LibraryBuilding WHERE name='Научная библиотека')),
('INV-006', '2022-06-18', (SELECT id FROM Book WHERE isbn='978-5-699-98765-4'), (SELECT id FROM LibraryBuilding WHERE name='Главное здание')),
('INV-007', '2019-07-22', (SELECT id FROM Book WHERE isbn='978-5-17-98765-4'), (SELECT id FROM LibraryBuilding WHERE name='Главное здание')),
('INV-008', '2019-08-30', (SELECT id FROM Book WHERE isbn='978-5-17-98765-4'), (SELECT id FROM LibraryBuilding WHERE name='Детский отдел')),
('INV-009', '2020-09-14', (SELECT id FROM Book WHERE isbn='978-0-19-953556-9'), (SELECT id FROM LibraryBuilding WHERE name='Научная библиотека')),
('INV-010', '2020-10-25', (SELECT id FROM Book WHERE isbn='978-0-19-953556-9'), (SELECT id FROM LibraryBuilding WHERE name='Главное здание'));

-- Читатели
INSERT INTO Reader (email, library_card_number, password_hash, FirstName, LastName, Patronymic, Address) VALUES
('reader1@mail.ru', 'LIB-001', 'hashed_reader1', 'Ольга', 'Кузнецова', 'Алексеевна', 'г. Москва, ул. Ленина, 10, кв. 5'),
('reader2@mail.ru', 'LIB-002', 'hashed_reader2', 'Дмитрий', 'Соколов', 'Игоревич', 'г. Москва, ул. Мира, 15, кв. 78'),
('reader3@mail.ru', 'LIB-003', 'hashed_reader3', 'Елена', 'Волкова', 'Петровна', 'г. Москва, ул. Горького, 5, кв. 12');
