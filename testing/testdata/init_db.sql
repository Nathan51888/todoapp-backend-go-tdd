-- public.todos definition

-- Drop table

-- DROP TABLE public.todos;

CREATE TABLE public.todos (
    id serial4 NOT NULL,
    title varchar NULL,
    completed varchar NULL,
    CONSTRAINT todos_pk PRIMARY KEY (id)
);

-- INSERT INTO public.todos
-- (title, completed, id)
-- VALUES ('Test_todo', 'false', nextval('todos_id_seq'::regclass));
