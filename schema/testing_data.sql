-- each test user's plain password is 0000 --
insert into users(id, email, password_digest, name)
values('97327413-6b65-486f-b299-91be0871f898', 'kenny@example.com', '$2a$10$gVtjNk4YL.O4I//ZBtvfN.YEebwR1Ci3.5OBHan4PWFzniSFqpzce', 'kenny');

insert into users(id, email, password_digest, name)
values('eb3c75df-b0df-4e06-a02f-e2ba77eba68a', 'nicole@example.com', '$2a$10$6tsb.2dRzV5gSTEJmtwkgeKpPIMO0VbMv2E6hP9xuAytwFlf0trVm', 'nicole');

insert into users(id, email, password_digest, name)
values('80695811-0bf2-44fd-980d-1635de7734a8', 'jack@example.com', '$2a$10$WkWwIpCbMyB1A2OuMC9LI.4LtQZtxNb1djcYqzeP0IayazJQgVkHG', 'jack');

-- one pair
insert into pairs(user_id_one, user_id_two)
values('97327413-6b65-486f-b299-91be0871f898', 'eb3c75df-b0df-4e06-a02f-e2ba77eba68a')