/* пока в полуручном режиме*/
sudo systemctl start postgresql
sudo -u postgres psql

create database hw12;
create user otus1 with encrypted password '123456';
grant all privileges on database hw12 to otus1;
create table events (
    ID text primary key,
    Title text,
    StartDate date,
    Details text,
    UserID bigint
);
create index ind1 on event (StartDate);
