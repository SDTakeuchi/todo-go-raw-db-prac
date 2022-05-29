CREATE TABLE IF NOT EXISTS todo (
    id int,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    title varchar(64) NOT NULL,
    memo varchar(1024) NOT NULL,
    is_done boolean NOT NULL,
    due_date timestamp with time zone
);

insert into todo(id, created_at, updated_at, title, memo, is_done, due_date) values(1, '2004-10-19 10:23:54', '2004-10-19 10:23:54', 'eat lunch', 'at 8 am', false, '2004-10-19 10:23:54');

insert into todo(id, created_at, updated_at, title, memo, is_done, due_date) values(2, '2014-10-19 10:23:54', '2014-10-19 10:23:54', 'do homework', 'at 10 am', false, '2014-10-19 10:23:54');
