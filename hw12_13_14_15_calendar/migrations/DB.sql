create table events (
    ID text PRIMARY KEY,
    Title text,
    StartDate date,
    Details text,
    UserID bigint
);
create index ind1 on events (StartDate);

create table shed_send_id (
    event_id text
);

create table send_events_stat (
    event_ID text,
    send_date TIMESTAMP default CURRENT_TIMESTAMP
);
