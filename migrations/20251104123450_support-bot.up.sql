create table user_info_posts (
    user_id bigint primary key,
    forward_post_id bigint not null
);

comment on column user_info_posts.user_id IS 'ID пользователя';
comment on column user_info_posts.forward_post_id IS 'ID перепоста в группе информации о пользователе';
