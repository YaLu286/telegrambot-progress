-- DELETE FROM beers;
-- DELETE FROM users;

-- insert into beers values (1, 'Red River', 'Stamm Brewing', 'APA', 'Ароматный светлый эль', 5.5, 3.73, '/images/redriver.jpg', 280);
-- insert into beers values (2, 'Your Young Lordship', 'Stamm Brewing', 'Lager', 'Светлый лагер с лёгкой горечью', 5.6, 3.86, '/images/item_5357.jpg', 280);
-- insert into beers values (3, 'Blowing Up: Pineapple & Mango', 'Stamm Brewing', 'Sour - Fruited', 'Кислый эль c манго и ананасом. Сочный, яркий, шелковистый.', 5, 4.06, '/images/blowingup.png', 430);
-- insert into beers values (4, 'Aftermath Pale Ale', 'Stamm Brewing', 'NE Pale Ale', 'New Zealand Pale Ale, сброженный тиолизированным штаммом дрожжей, и охмеленный новозеландскими хмелями Riwaka и Motueka.', 5.5, 3.94, '/images/aftermath.jpeg', 430);
-- insert into beers values (5, 'Green Steps', 'Stamm Brewing', 'NE Pale Ale', 'Сочный New England Pale Ale, охмелённый Motueka, Citra и Simcoe.', 5.5, 4.02, '/images/greensteps.jpg', 430);
-- insert into beers values (6, 'Shadowtricks', 'Stamm Brewing', 'NE Pale Ale', ' Свежий New England Pale Ale,охмелённый Riwaka & Mosaic.', 5.5, 4.03, '/images/shadowtricks.jpg', 430);
-- insert into beers values (7, 'Kranhaus', 'Stamm Brewing', 'Lager', 'Пиво в стиле кёльш - насыщенный и свежий вкус, с нотками цитрусовых и легкой горчинкой.', 4.7, 3.63, '/images/kranhaus.jpg', 280);
-- insert into beers values (8, 'Silk Supersonic', 'Stamm Brewing', 'NE Pale Ale', 'Шелковистый новоанглийский пейл эль, охмеленный гранулами Talus, Simcoe, Mosaic и Citra.', 5.5, 4.09, '/images/supersonic.jpg', 430);
-- insert into beers values (9, 'Alpaca Juice', 'Stamm Brewing', 'Sour - Fruited', ' Имперский саур эль с добавлением клубники и черники, сурово настоянный на соке граната.', 7, 4.18, '/images/alpaca.jpg', 450);
-- insert into beers values (10, 'Quartz V2.0', 'Zavod', 'Lager', 'Что может быть яснее, свежее, живительнее и понятнее легкого классического лагера? Наш вариант, чистый как кварц, придется по вкусу любому, ведь это не простой лагер, а охмеленный премиальными сортами американских хмелей Citra и Simcoe.', 5, 3.71 , '/images/quartz.png', 450);
-- insert into beers values (11, 'Empress Maria', 'Zavod', 'Gose', 'Имперское томатное гозе. Больше алкоголя, больше пряностей, больше Вустерского соуса, двойная доза остроты, а ещё новые составляющие - соевый соус, чеснок и немного копчености - не оставят вас в числе равнодушныхю', 6.9, 3.88, '/images/maria.png', 430);
-- insert into beers values (12, 'Ticket to Köln', 'Zavod', 'Lager', 'Легкое тело, баланс горечи и солодовой составляющей, тонкий фруктовый и цветочный аромат; наш кёльш - ваш билетик в Кёльн!', 4.3, 3.71, '/images/ticket.png', 430);

-- insert into  users values (1032058526, true, 'sokol');

-- select * from users;


-- delete from locations;

-- insert into locations values('presnya', 'Добро пожаловать в Прогресс на Пресне!', '+7(925)888-13-16', 'progress.presnya@gmail.com', '/images/presnya.jpg');
-- insert into locations values('sokol', 'Добро пожаловать в Прогресс на Соколе!', '+7(925)433-52-94', 'progress.sokol@gmail.com', '/images/sokol.jpg');
-- insert into locations values('rizhskaya', 'Добро пожаловать в Прогресс на Рижской!', '+7(925)635-70-19', 'progress.rizhskaya@gmail.com', '/images/rizhskaya.jpg');
-- insert into locations values('frunza', 'Добро пожаловать в Прогресс на Фрунзенской!', '+7(903)167-22-53', 'progress.frunza@gmail.com', '/images/frunza.jpg');

-- select * from locations;

select * from user_sessions;

-- DROP table user_sessions;
