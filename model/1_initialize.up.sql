create table grabber_threads (
  url text primary key,
  created timestamp not null default now()
);

create table grabber_images (
  md5 text primary key,
  filepath text not null default ''::text,
  thumbpath text not null default ''::text,
  postid integer not null default 0,
  tags text not null default ''::text,
  rating char not null default 'q',
  parent_md5 text not null default ''::text
);

CREATE SEQUENCE grabber_threads_images_porder_seq START 1;
CREATE SEQUENCE grabber_threads_images_norder_seq START 1;

create table grabber_threads_images (
  thread_url text,
  image_md5 text,
  porder integer not null default nextval('grabber_threads_images_porder_seq'),
  norder integer not null default nextval('grabber_threads_images_norder_seq'),
  pgroup integer not null default 0,
  primary key(thread_url, image_md5)
);