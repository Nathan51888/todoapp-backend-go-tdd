-- public.todo definition

-- Drop table

-- DROP TABLE public.todo;

CREATE TABLE public.todo (
    id serial4 NOT NULL,
    title varchar NULL,
    completed boolean,
    CONSTRAINT todo_pk PRIMARY KEY (id)
);

-- INSERT INTO public.todo
-- (title, completed, id)
-- VALUES ('Test_todo', 'false', nextval('todo_id_seq'::regclass));
