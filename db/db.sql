-- DROP TABLE rust_test1.course;
CREATE TABLE rust_test1.course (
	id int4 GENERATED ALWAYS AS IDENTITY NOT NULL,
	teacher_id int4 DEFAULT 0 NULL,
	"name" varchar DEFAULT ''::character varying NULL,
	"time" timestamp DEFAULT now() NULL,
	CONSTRAINT course_pk PRIMARY KEY (id)
);