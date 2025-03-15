-- public.users definition

-- Drop table

-- DROP TABLE public.users;

CREATE TABLE public.users (
    user_id uuid DEFAULT gen_random_uuid() NOT NULL,
    email varchar NOT NULL,
    user_password varchar NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (user_id)
);

-- public.todo definition

-- Drop table

-- DROP TABLE public.todo;

CREATE TABLE public.todos (
    todo_id uuid DEFAULT gen_random_uuid() NOT NULL,
    title varchar NULL,
    completed bool NULL,
    user_id uuid,
    CONSTRAINT todo_pk PRIMARY KEY (todo_id),
    CONSTRAINT todo_user_fk
    FOREIGN KEY (user_id) REFERENCES public.users (user_id)
);

-- INSERT INTO public.todo
-- (title, completed)
-- VALUES ('Test_todo', 'false');
