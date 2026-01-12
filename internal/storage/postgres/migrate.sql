-- 0001	-	add status to users table
alter table users add column status int not null default 1;
alter table temp_users add column status int not null default 1;



alter table temp_users add column "company_name" varchar(100);
alter table temp_users add column "company_type_id" int;
alter table temp_users add column "activity_field_id" int;
alter table temp_users add column "vat_number" varchar(100);
alter table temp_users add column "address" varchar(100);
alter table temp_users add column "license_issue_date" timestamp;
alter table temp_users add column "license_expiry_date" timestamp;



alter table temp_users add constraint fk_temp_users_activity_field_id
        foreign key (activity_field_id)
            references activity_fields(id)
                on delete cascade
                on update cascade;

alter table temp_users add column documents_id int;
alter table temp_users add constraint fk_temp_users_documents_id
        foreign key (documents_id)
            references documents(id)
                on delete cascade
                on update cascade;

alter table documents
    alter column licence_issue_date drop not null,
    alter column licence_expiry_date drop not null;


alter table temp_users rename column license_issue_date to licence_issue_date;
alter table temp_users rename column license_expiry_date to licence_expiry_date;



insert into company_types (name) values ('Company');
insert into company_types (name) values ('Individual');
insert into company_types (name) values ('Other');

insert into activity_fields (name) values ('Transport');
insert into activity_fields (name) values ('Trade');
insert into activity_fields (name) values ('Services');
insert into activity_fields (name) values ('Industry');




(1, 310, 3232, 1, 1, 1, 1, 2020, 50000, 'VIN1234567890', false, 1, false, ARRAY['123456789'], 25000, false, 1, 1, 3),
(1, 310, 3232, 1, 1, 1, 1, 2020, 30000, 'VIN1234567891', false, 1, true, ARRAY['123456790'], 28000, false, 2, 1, 3),
(2, 310, 3232, 1, 1, 2, 1, 2019, 70000, 'VIN1234567892', false, 2, false, ARRAY['123456791'], 22000, false, 3, 2, 3),
(1, 310, 3232, 1, 1, 1, 1, 2020, 15000, 'VIN1234567893', false, 1, true, ARRAY['123456792'], 32000, false, 1, 1, 3),
(2, 310, 3232, 1, 1, 3, 1, 2018, 85000, 'VIN1234567894', true, 2, false, ARRAY['123456793'], 18000, false, 2, 3, 3),
(1, 310, 3233, 1, 1, 1, 1, 2020, 25000, 'VIN1234567895', false, 1, true, ARRAY['123456794'], 35000, false, 3, 1, 3),
(2, 310, 3233, 1, 1, 2, 1, 2020, 40000, 'VIN1234567896', false, 1, false, ARRAY['123456795'], 30000, false, 1, 2, 3),
(1, 310, 3233, 1, 1, 1, 1, 2020, 12000, 'VIN1234567897', false, 1, true, ARRAY['123456796'], 38000, false, 2, 1, 3),
(2, 310, 3233, 1, 1, 3, 1, 2019, 60000, 'VIN1234567898', false, 2, false, ARRAY['123456797'], 27000, false, 3, 4, 3),
(1, 310, 3233, 1, 1, 1, 1, 2020, 8000, 'VIN1234567899', false, 1, true, ARRAY['123456798'], 42000, true, 1, 1, 3),
(2, 310, 3234, 1, 1, 2, 1, 2020, 35000, 'VIN1234567900', false, 1, false, ARRAY['123456799'], 28000, false, 2, 2, 3),
(1, 310, 3234, 1, 1, 1, 1, 2020, 22000, 'VIN1234567901', false, 1, true, ARRAY['123456800'], 33000, false, 3, 1, 3),
(2, 310, 3234, 1, 1, 3, 1, 2019, 75000, 'VIN1234567902', true, 3, false, ARRAY['123456801'], 24000, false, 1, 3, 3),
(1, 310, 3234, 1, 1, 1, 1, 2020, 18000, 'VIN1234567903', false, 1, true, ARRAY['123456802'], 36000, false, 2, 1, 3),
(2, 310, 3234, 1, 1, 2, 1, 2023, 5000, 'VIN1234567904', false, 1, false, ARRAY['123456803'], 45000, true, 3, 5, 3),
(1, 310, 3232, 1, 1, 1, 1, 2020, 95000, 'VIN1234567905', false, 3, false, ARRAY['123456804'], 16000, false, 1, 2, 3),
(2, 310, 3232, 1, 1, 3, 1, 2021, 28000, 'VIN1234567906', false, 1, true, ARRAY['123456805'], 31000, false, 2, 1, 3),
(1, 310, 3232, 1, 1, 1, 1, 2020, 45000, 'VIN1234567907', false, 2, false, ARRAY['123456806'], 26000, false, 3, 4, 3),
(2, 310, 3233, 1, 1, 2, 1, 2022, 16000, 'VIN1234567908', false, 1, true, ARRAY['123456807'], 37000, false, 1, 1, 3),
(1, 310, 3233, 1, 1, 1, 1, 2020, 65000, 'VIN1234567909', true, 2, false, ARRAY['123456808'], 25000, false, 2, 3, 3),
(2, 310, 3233, 1, 1, 3, 1, 2023, 3000, 'VIN1234567910', false, 1, true, ARRAY['123456809'], 48000, true, 3, 1, 3),
(1, 310, 3234, 1, 1, 1, 1, 2020, 52000, 'VIN1234567911', false, 2, false, ARRAY['123456810'], 29000, false, 1, 2, 3),
(2, 310, 3234, 1, 1, 2, 1, 2021, 33000, 'VIN1234567912', false, 1, true, ARRAY['123456811'], 34000, false, 2, 1, 3),
(1, 310, 3234, 1, 1, 1, 1, 2020, 88000, 'VIN1234567913', true, 3, false, ARRAY['123456812'], 20000, false, 3, 3, 3),
(2, 310, 3232, 1, 1, 3, 1, 2022, 14000, 'VIN1234567914', false, 1, true, ARRAY['123456813'], 39000, false, 1, 1, 3),
(1, 310, 3232, 1, 1, 1, 1, 2020, 72000, 'VIN1234567915', false, 2, false, ARRAY['123456814'], 23000, false, 2, 4, 3),
(2, 310, 3233, 1, 1, 2, 1, 2023, 6000, 'VIN1234567916', false, 1, true, ARRAY['123456815'], 46000, true, 3, 1, 3),
(1, 310, 3233, 1, 1, 1, 1, 2020, 38000, 'VIN1234567917', false, 1, false, ARRAY['123456816'], 32000, false, 1, 2, 3),
(2, 310, 3234, 1, 1, 3, 1, 2021, 26000, 'VIN1234567918', false, 1, true, ARRAY['123456817'], 35000, false, 2, 1, 3),
(1, 310, 3234, 1, 1, 1, 1, 2020, 19000, 'VIN1234567919', false, 1, false, ARRAY['123456818'], 41000, false, 3, 5, 3);


alter table profiles add column "company_name" varchar(200);

alter table moto_brands alter column image drop not null;



select 
    u.id,
    u.username,
    u.last_active_date,
    p.avatar,
    json_agg(
        json_build_object(
            'id', m.id,
            'message', m.message,
            'type', m.type,
            'created_at', m.created_at
        )
    ) as messages
from messages m
left join users u on m.sender_id = u.id
left join profiles p on u.id = p.user_id
where m.status = 1
group by u.id, p.avatar;



insert into users (email, username, role_id, password, phone) values ('dealer1@example.com', 'dealer1', 2, 'password1', '1234567891');
insert into users (email, username, role_id, password, phone) values ('logist1@example.com', 'logist1', 3, 'password1', '1234567892');
insert into users (email, username, role_id, password, phone) values ('broker1@example.com', 'broker1', 4, 'password1', '1234567893');
insert into users (email, username, role_id, password, phone) values ('service1@example.com', 'service1', 5, 'password1', '1234567894');
insert into users (email, username, role_id, password, phone) values ('dealer2@example.com', 'dealer2', 2, 'password2', '1234567895');
insert into users (email, username, role_id, password, phone) values ('logist2@example.com', 'logist2', 3, 'password2', '1234567896');
insert into users (email, username, role_id, password, phone) values ('broker2@example.com', 'broker2', 4, 'password2', '1234567897');
insert into users (email, username, role_id, password, phone) values ('service2@example.com', 'service2', 5, 'password2', '1234567898');
insert into users (email, username, role_id, password, phone) values ('dealer3@example.com', 'dealer3', 2, 'password3', '1234567899');
insert into users (email, username, role_id, password, phone) values ('logist3@example.com', 'logist3', 3, 'password3', '1234567900');
insert into users (email, username, role_id, password, phone) values ('broker3@example.com', 'broker3', 4, 'password3', '1234567901');
insert into users (email, username, role_id, password, phone) values ('service3@example.com', 'service3', 5, 'password3', '1234567902');


insert into profiles (user_id, company_name, company_type_id, activity_field_id, vat_number, driving_experience, notification, username, registered_by, google, avatar, banner, contacts, address, coordinates, message, birthday, about_me) 
    values 
    (79, 'Company 1', 1, 1, '1234567890', 10, true, 'dealer1', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567890", "telegram": "1234567890"}'::jsonb, '1234567890', '1234567890', '1234567890', '2025-12-10', '1234567890'),
    (80, 'Company 2', 2, 2, '1234567891', 15, true, 'dealer2', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567891", "telegram": "1234567891"}'::jsonb, '1234567891', '1234567891', '1234567891', '2025-12-10', '1234567891'),
    (81, 'Company 3', 3, 3, '1234567892', 20, true, 'dealer3', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567892", "telegram": "1234567892"}'::jsonb, '1234567892', '1234567892', '1234567892', '2025-12-10', '1234567892'),
    (82, 'Company 4', 1, 1, '1234567893', 25, true, 'broker1', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567893", "telegram": "1234567893"}'::jsonb, '1234567893', '1234567893', '1234567893', '2025-12-10', '1234567893'),
    (83, 'Company 5', 2, 2, '1234567894', 30, true, 'service1', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567894", "telegram": "1234567894"}'::jsonb, '1234567894', '1234567894', '1234567894', '2025-12-10', '1234567894'),
    (84, 'Company 6', 3, 3, '1234567895', 35, true, 'broker2', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567895", "telegram": "1234567895"}'::jsonb, '1234567895', '1234567895', '1234567895', '2025-12-10', '1234567895'),
    (85, 'Company 7', 1, 1, '1234567896', 40, true, 'service2', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567896", "telegram": "1234567896"}'::jsonb, '1234567896', '1234567896', '1234567896', '2025-12-10', '1234567896'),
    (86, 'Company 8', 2, 2, '1234567897', 45, true, 'broker3', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567897", "telegram": "1234567897"}'::jsonb, '1234567897', '1234567897', '1234567897', '2025-12-10', '1234567897'),
    (87, 'Company 9', 3, 3, '1234567898', 50, true, 'service3', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567898", "telegram": "1234567898"}'::jsonb, '1234567898', '1234567898', '1234567898', '2025-12-10', '1234567898'),
    (88, 'Company 10', 1, 1, '1234567899', 55, true, 'dealer4', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567899", "telegram": "1234567899"}'::jsonb, '1234567899', '1234567899', '1234567899', '2025-12-10', '1234567899'),
    (89, 'Company 11', 2, 2, '1234567900', 60, true, 'logist4', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567900", "telegram": "1234567900"}'::jsonb, '1234567900', '1234567900', '1234567900', '2025-12-10', '1234567900'),
    (90, 'Company 12', 3, 3, '1234567901', 65, true, 'broker4', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567901", "telegram": "1234567901"}'::jsonb, '1234567901', '1234567901', '1234567901', '2025-12-10', '1234567901');

update profiles set contacts = '{"whatsapp": "1234567890", "telegram": "1234567890"}'::jsonb where user_id = 253;




alter table vehicles alter column owners drop not null;
alter table vehicles alter column owners set default 0;


alter table users add column temp_password varchar(100);

drop table services;
drop table service_types;

30.10.2025
alter table company_types add column "name_ru" varchar(50) default 'name_ru';
alter table activity_fields add column "name_ru" varchar(50) default 'name_ru';
alter table brands add column "name_ru" varchar(50) default 'name_ru';
alter table models add column "name_ru" varchar(50) default 'name_ru';
alter table body_types add column "name_ru" varchar(50) default 'name_ru';
alter table transmissions add column "name_ru" varchar(50) default 'name_ru';
alter table drivetrains add column "name_ru" varchar(50) default 'name_ru';
alter table fuel_types add column "name_ru" varchar(50) default 'name_ru';
alter table generations add column "name_ru" varchar(50) default 'name_ru';
alter table ownership_types add column "name_ru" varchar(50) default 'name_ru';
alter table colors add column "name_ru" varchar(50) default 'name_ru';
alter table moto_categories add column "name_ru" varchar(50) default 'name_ru';
alter table moto_brands add column "name_ru" varchar(50) default 'name_ru';
alter table moto_models add column "name_ru" varchar(50) default 'name_ru';
alter table moto_parameters add column "name_ru" varchar(50) default 'name_ru';
alter table moto_parameter_values add column "name_ru" varchar(50) default 'name_ru';
alter table com_categories add column "name_ru" varchar(50) default 'name_ru';
alter table com_brands add column "name_ru" varchar(50) default 'name_ru';
alter table com_models add column "name_ru" varchar(50) default 'name_ru';
alter table com_parameters add column "name_ru" varchar(50) default 'name_ru';
alter table com_parameter_values add column "name_ru" varchar(50) default 'name_ru';



alter table company_types add column "name_ae" varchar(50) default 'name_ae';
alter table activity_fields add column "name_ae" varchar(50) default 'name_ae';
alter table brands add column "name_ae" varchar(50) default 'name_ae';
alter table models add column "name_ae" varchar(50) default 'name_ae';
alter table body_types add column "name_ae" varchar(50) default 'name_ae';
alter table transmissions add column "name_ae" varchar(50) default 'name_ae';
alter table drivetrains add column "name_ae" varchar(50) default 'name_ae';
alter table fuel_types add column "name_ae" varchar(50) default 'name_ae';
alter table generations add column "name_ae" varchar(50) default 'name_ae';
alter table ownership_types add column "name_ae" varchar(50) default 'name_ae';
alter table colors add column "name_ae" varchar(50) default 'name_ae';
alter table moto_categories add column "name_ae" varchar(50) default 'name_ae';
alter table moto_brands add column "name_ae" varchar(50) default 'name_ae';
alter table moto_models add column "name_ae" varchar(50) default 'name_ae';
alter table moto_parameters add column "name_ae" varchar(50) default 'name_ae';
alter table moto_parameter_values add column "name_ae" varchar(50) default 'name_ae';
alter table com_categories add column "name_ae" varchar(50) default 'name_ae';
alter table com_brands add column "name_ae" varchar(50) default 'name_ae';
alter table com_models add column "name_ae" varchar(50) default 'name_ae';
alter table com_parameters add column "name_ae" varchar(50) default 'name_ae';
alter table com_parameter_values add column "name_ae" varchar(50) default 'name_ae';
alter table countries add column "name_ae" varchar(50) default 'name_ae';
alter table countries add column "name_ru" varchar(50) default 'name_ru';
alter table countries add column "country_code" varchar(50) default 'country_code';


alter table cities add column "name_ae" varchar(50) default 'name_ae';
alter table cities add column "name_ru" varchar(50) default 'name_ru';
alter table regions add column "name_ae" varchar(50) default 'name_ae';
alter table regions add column "name_ru" varchar(50) default 'name_ru';
alter table engines add column "name_ae" varchar(50) default 'name_ae';
alter table engines add column "name_ru" varchar(50) default 'name_ru';


drop table if exists admins;
alter table users add column "permissions" jsonb default '[]';

insert into users (email, username, role_id, password, phone, permissions) 
values 
    ('admin@admin.com', 'admin', 0, '$2a$10$H6OHFABvTjMScHwB6qIvte4teoXtGP1h/ViqTnVHg1R.iw4yy9xTq', '989898989', '["cars", "motorcycles", "comtrans"]');

insert into users (email, username, role_id, password, phone, permissions) 
values 
    ('admin2@admin.com', 'admin2', 0, '$2a$10$H6OHFABvTjMScHwB6qIvte4teoXtGP1h/ViqTnVHg1R.iw4yy9xTq', '9090909090', '["cars", "motorcycles", "comtrans", "chat"]');

update users set permissions = '["cars", "motorcycles", "chat"]' where id = 165;

insert into profiles (user_id, company_name, company_type_id, activity_field_id, vat_number, driving_experience, notification, username, registered_by, google, avatar, banner, contacts, address, coordinates, message, birthday, about_me) 
    values 
    (165, 'Admin', 1, 1, '1234567890', 10, true, 'admin', 'admin', 'https://google.com', 'https://avatar.com', 'https://banner.com', '{"whatsapp": "1234567890", "telegram": "1234567890"}'::jsonb, '1234567890', '1234567890', '1234567890', '2025-12-10', '1234567890');

create table message_files (
    "id" serial primary key,
    "sender_id" int,
    "file_path" varchar(255) not null,
    "created_at" timestamp without time zone not null default now(),
    constraint fk_message_files_sernder_id
        foreign key (sender_id)
            references users(id)
                on delete set null
                on update set null
);



create table conversations (
    "id" serial primary key,
    "user_id_1" int not null,
    "user_id_2" int not null,
    "created_at" timestamp default now(),
    constraint fk_conversations_user_id_1
        foreign key (user_id_1)
            references users(id)
                on delete cascade
                on update cascade,
    constraint fk_conversations_user_id_2
        foreign key (user_id_2)
            references users(id)
                on delete cascade
                on update cascade,
    unique(user_id_1, user_id_2)
);

-- 20.12.2025
alter table conversations add column "new_messages" int not null default 0;
alter table conversations drop column "ney_message";


drop table messages;
drop table conversations;

create table conversations (
    "id" serial primary key,
    "user_id_1" int not null,
    "user_id_2" int not null,
    "user_1_unread_messages" int not null default 0,
    "user_2_unread_messages" int not null default 0,
    "last_message_id" int,
    "last_message" varchar(500),
    "last_message_type" int not null default 1, -- 1-text, 2-item, 3-video, 4-image,
    "updated_at" timestamp not null default now(),
    "created_at" timestamp not null default now(),
    constraint fk_conversations_user_id_1
        foreign key (user_id_1)
            references users(id)
                on delete cascade
                on update cascade,
    constraint fk_conversations_user_id_2
        foreign key (user_id_2)
            references users(id)
                on delete cascade
                on update cascade,
    constraint unique_conversation_pair 
        unique(user_id_1, user_id_2)
);

-- Indexes for efficient querying of conversations by user with ordering
-- These support: WHERE (user_id_1 = $1 OR user_id_2 = $1) ORDER BY updated_at DESC
CREATE INDEX idx_conversations_user1_updated ON conversations(user_id_1, updated_at DESC);
CREATE INDEX idx_conversations_user2_updated ON conversations(user_id_2, updated_at DESC);


create table messages (
    "id" serial primary key,
    "conversation_id" int not null,
    "sender_id" int not null,
    "status" int not null default 1, -- 1-unread, 2-read
    "message" varchar(500) not null, --  it is an id if type "item".
    "type" int not null default 1, -- 1-text, 2-item, 3-video, 4-image,
    "created_at" timestamp not null,
    constraint fk_messages_sender_id
        foreign key (sender_id)
            references users(id)
                on delete cascade
                on update cascade,
    constraint fk_messages_conversation_id
        foreign key (conversation_id)
            references conversations(id)
                on delete cascade
                on update cascade
);


-- 26.12.2025
