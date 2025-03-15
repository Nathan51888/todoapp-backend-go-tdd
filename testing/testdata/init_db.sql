-- public.todo definition

-- Drop table

-- DROP TABLE public.todo;

CREATE TABLE public.todo (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title varchar NULL,
    completed bool NULL,
    CONSTRAINT todo_pk PRIMARY KEY (id),
    CONSTRAINT todo_user_fk FOREIGN KEY (id) REFERENCES public."user" (id)
);

-- INSERT INTO public.todo
-- (title, completed)
-- VALUES ('Test_todo', 'false');


-- public."user" definition

-- Drop table

-- DROP TABLE public."user";

CREATE TABLE public."user" (
    id uuid NOT NULL,
    email varchar NOT NULL,
    "password" varchar NOT NULL,
    CONSTRAINT user_pk PRIMARY KEY (id)
);
