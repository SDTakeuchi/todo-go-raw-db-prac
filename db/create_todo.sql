CREATE TABLE IF NOT EXISTS todos (
    id SERIAL NOT NULL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    title varchar(64) NOT NULL,
    memo varchar(1024) NOT NULL,
    is_done boolean NOT NULL,
    due_date TIMESTAMPTZ NOT NULL
);

insert into todos (title, memo, is_done, due_date) values('eat lunch', 'at 8 am', false, '2004-10-19 10:23:54');

insert into todos (title, memo, is_done, due_date) values('do homework', 'at 10 am', false, '2014-10-19 10:23:54');
