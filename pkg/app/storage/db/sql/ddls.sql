create table if not exists "user"
(
    id                    char(26)                               not null
        primary key,
    kind                  text                                   not null,
    name                  text                                   not null,
    email                 text
        constraint user_email_unique
            unique,
    phone                 text,
    github_handle         text,
    x_handle              text,
    created_at            timestamp with time zone default now() not null,
    updated_at            timestamp with time zone,
    deleted_at            timestamp with time zone,
    github_remote_id      text
        constraint user_github_remote_id_unique
            unique,
    x_remote_id           text,
    individual_profile_id char(26)
);

create table if not exists profile
(
    id                  char(26)                               not null
        primary key,
    kind                text                                   not null,
    slug                text                                   not null
        constraint profile_slug_unique
            unique,
    profile_picture_uri text,
    title               text                                   not null,
    description         text                                   not null,
    show_stories        boolean                  default false not null,
    show_projects       boolean                  default false not null,
    created_at          timestamp with time zone default now() not null,
    updated_at          timestamp with time zone,
    deleted_at          timestamp with time zone
);

create table if not exists profile_membership
(
    id         char(26)                               not null
        primary key,
    kind       text                                   not null,
    profile_id char(26)                               not null,
    user_id    char(26)                               not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    constraint profile_membership_profile_id_user_id_unique
        unique (profile_id, user_id)
);

create table if not exists session
(
    id                          char(26)                               not null
        primary key,
    status                      text                                   not null,
    oauth_request_state         text                                   not null,
    oauth_request_code_verifier text                                   not null,
    oauth_redirect_uri          text,
    logged_in_user_id           char(26),
    logged_in_at                timestamp with time zone,
    expires_at                  timestamp with time zone,
    created_at                  timestamp with time zone default now() not null,
    updated_at                  timestamp with time zone
);

create index if not exists session_logged_in_user_id_index
    on session (logged_in_user_id);

create table if not exists question
(
    id             char(26)                               not null
        primary key,
    user_id        char(26)                               not null,
    content        text                                   not null,
    is_hidden      boolean                  default false not null,
    created_at     timestamp with time zone default now() not null,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    answered_at    timestamp with time zone,
    answer_uri     text,
    is_anonymous   boolean                  default false not null,
    answer_kind    text,
    answer_content text
);

create table if not exists question_vote
(
    id          char(26)                               not null
        primary key,
    question_id char(26)                               not null,
    user_id     char(26)                               not null,
    score       integer                                not null,
    created_at  timestamp with time zone default now() not null,
    constraint question_vote_question_id_user_id_unique
        unique (question_id, user_id)
);

create table if not exists event
(
    id                char(26)                                       not null
        primary key,
    kind              text                                           not null,
    slug              text                                           not null
        constraint event_slug_unique
            unique,
    event_picture_uri text,
    title             text                                           not null,
    description       text                                           not null,
    time_start        timestamp with time zone                       not null,
    time_end          timestamp with time zone                       not null,
    created_at        timestamp with time zone default now()         not null,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    series_id         char(26),
    status            text                     default 'draft'::text not null,
    attendance_uri    text,
    published_at      timestamp with time zone
);

create table if not exists event_attendance
(
    id         char(26)                               not null
        primary key,
    kind       text                                   not null,
    event_id   char(26)                               not null,
    profile_id char(26)                               not null,
    created_at timestamp with time zone default now() not null,
    updated_at timestamp with time zone,
    deleted_at timestamp with time zone,
    constraint event_attendance_event_id_profile_id_unique
        unique (event_id, profile_id)
);

create table if not exists event_series
(
    id                char(26)                               not null
        primary key,
    slug              text                                   not null
        constraint event_series_slug_unique
            unique,
    event_picture_uri text,
    title             text                                   not null,
    description       text                                   not null,
    created_at        timestamp with time zone default now() not null,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
);

create table if not exists story
(
    id                char(26)                               not null
        primary key,
    kind              text                                   not null,
    status            text                                   not null,
    is_featured       boolean                  default false,
    slug              text                                   not null
        constraint story_slug_unique
            unique,
    story_picture_uri text,
    title             text                                   not null,
    description       text                                   not null,
    author_profile_id char(26),
    content           text,
    published_at      timestamp with time zone,
    created_at        timestamp with time zone default now() not null,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
);
