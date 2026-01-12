-- Union query for vehicles, comtrans, and motorcycles
with vs as (        
    select 
        vs.id,
        'car' as type,
        bs.name as brand,
        ms.name as model,
        vs.year,
        vs.price,
        vs.credit,
        vs.status,
        vs.created_at,
        images.images,
        vs.view_count,
        true as my_car,
        vs.crash
    from vehicles vs
    left join brands bs on vs.brand_id = bs.id
    left join models ms on vs.model_id = ms.id
    LEFT JOIN LATERAL (
        SELECT json_agg(img.image) AS images
        FROM (
            SELECT image as image
            FROM images
            WHERE vehicle_id = vs.id
            ORDER BY created_at DESC
        ) img
    ) images ON true
    where vs.user_id = 1 and status = 2
    order by vs.id desc
),
cms as (
    select
        cm.id,
        'comtran' as type,
        cbs.name as brand,
        cms.name as model,
        cm.year,
        cm.price,
        cm.credit,
        cm.status,
        cm.created_at,
        images.images,
        cm.view_count,
        true as my_car,
        cm.crash
    from comtrans cm
    left join com_brands cbs on cbs.id = cm.comtran_brand_id
    left join com_models cms on cms.id = cm.comtran_model_id
    LEFT JOIN LATERAL (
        SELECT json_agg(img.image) AS images
        FROM (
            SELECT image
            FROM comtran_images
            WHERE comtran_id = cm.id
            ORDER BY created_at DESC
        ) img
    ) images ON true
    where cm.user_id = 1 and cm.status = 2
),
mts as (
    select
        mt.id,
        'motorcycle' as type,
        mbs.name as brand,
        mms.name as model,
        mt.year,
        mt.price,
        mt.credit,
        mt.status,
        mt.created_at,
        mt.view_count,
        images.images,
        true as my_car,
        mt.crash
    from motorcycles mt
    left join moto_brands mbs on mbs.id = mt.moto_brand_id
    left join moto_models mms on mms.id = mt.moto_model_id
    LEFT JOIN LATERAL (
        SELECT json_agg(img.image) AS images
        FROM (
            SELECT image
            FROM moto_images
            WHERE moto_id = mt.id
            ORDER BY created_at DESC
        ) img
    ) images ON true
    where mt.user_id = 1 and mt.status = 2
)
-- Union all three CTEs
select 
    id, type, brand, model, 
    year, price, credit, 
    status, created_at, 
    view_count, images, my_car, 
    crash 
from vs
union all
select 
    id, type, brand, model, 
    year, price, credit, 
    status, created_at, 
    view_count, images, my_car, 
    crash 
from cms
union all
select 
    id, type, brand, model, 
    year, price, credit, 
    status, created_at, 
    view_count, images, my_car, 
    crash 
from mts
order by created_at desc;




alter table motorcycles add column "view_count" int not null default 0;
alter table comtrans add column "view_count" int not null default 0;
alter table motorcycles add column "credit" boolean not null default false;
alter table comtrans add column "credit" boolean not null default false;



select 
    c.updated_at,
    u.username,
    p.avatar,
    u.id
from conversations c
join users u on u.id = 
    case 
        when c.user_id_1 = 23 then c.user_id_2 
        else c.user_id_1 
    end
join profiles p on p.user_id = u.id
order by c.updated_at desc;





create table horse_powers (
    "id" serial primary key,
    "name" varchar(255) not null,
    "name_ru" varchar(255) default 'name_ru',
    "name_ae" varchar(255) default 'name_ae',
    "created_at" timestamp default now(),
    unique("name")
);



drop table user_likes;
drop table videos;
drop table images;
drop table vehicles;
drop table generation_modifications ;


create table generation_modifications (
    "id" serial primary key,
    "generation_id" int not null,
    "horse_power_id" int,
    "body_type_id" int not null,
    "engine_id" int not null,
    "fuel_type_id" int not null, 
    "drivetrain_id" int not null,
    "transmission_id" int not null, 
    unique(horse_power_id, generation_id, body_type_id, engine_id, fuel_type_id, drivetrain_id, transmission_id),
    constraint fk_generation_modifications_horse_power_id
        foreign key (horse_power_id)
            references horse_powers(id)
                on delete set null
                on update cascade,
    constraint fk_generation_modifications_generation_id
        foreign key (generation_id)
            references generations(id)
                on delete cascade
                on update cascade,
    constraint fk_generation_modifications_engine_id
        foreign key (engine_id)
            references engines(id)
                on delete cascade
                on update cascade,
    constraint fk_generation_modifications_fuel_type_id
        foreign key (fuel_type_id)
            references fuel_types(id)
                on delete cascade
                on update cascade,
    constraint fk_generation_modifications_drivetrain_id
        foreign key (drivetrain_id)
            references drivetrains(id)
                on delete cascade
                on update cascade,
    constraint fk_generation_modifications_transmission_id
        foreign key (transmission_id)
            references transmissions(id)
                on delete cascade
                on update cascade,
    constraint fk_generation_modifications_body_type_id
        foreign key (body_type_id)
            references body_types(id)
                on delete cascade
                on update cascade
);


create table vehicles (
    "id" serial primary key,
    "user_id" int not null,
    "modification_id" int not null,
    "brand_id" int,
    "region_id" int,
    "city_id" int default 1,
    "model_id" int,
    "ownership_type_id" int not null default 1,
    "owners" int not null default 0,
    "view_count" int not null default 0,
    "year" int not null,
    "popular" int not null default 0,
    "description" text,
    "credit" boolean not null default false,
    "wheel" boolean not null default true, -- true left, false right
    "crash" boolean not null default false,
    "odometer" int not null default 0,
    "vin_code" varchar(255),
    "phone_numbers" varchar(255)[] not null,
    "price" int not null,
    "new" boolean not null default false,
    "color_id" int not null,
    "trade_in" int not null default 1, -- 1. No exchange 2. Equal value 3. More expensive 4. Cheaper 5. Not a car
    "status" int not null default 3, -- 1-pending, 2-not sale (my cars), 3-on sale,
    "updated_at" timestamp default now(),
    "created_at" timestamp default now(),
    constraint fk_vehicles_color_id
        foreign key (color_id)
            references colors(id)
                on delete set null
                on update cascade,
    constraint fk_vehicles_ownership_type_id
        foreign key (ownership_type_id)
            references ownership_types(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_user_id
        foreign key (user_id)
            references users(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_brand_id
        foreign key (brand_id)
            references brands(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_model_id
        foreign key (model_id)
            references models(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_modification_id
        foreign key (modification_id)
            references generation_modifications(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_region_id
        foreign key (region_id)
            references regions(id)
                on delete cascade
                on update cascade,
    constraint fk_vehicles_city_id
        foreign key (city_id)
            references cities(id)
                on delete cascade
                on update cascade
);



CREATE TABLE user_likes (
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    PRIMARY KEY (user_id, vehicle_id),
    constraint fk_user_likes_vehicle_id
        foreign key (vehicle_id)
            references vehicles(id)
                on delete cascade
                on update cascade,
    constraint fk_user_likes_user_id
        foreign key (user_id)
            references users(id)
                on delete cascade
                on update cascade
);




create table images (
    "vehicle_id" int not null,
    "image" varchar(255) not null,
    "created_at" timestamp not null default now(),
    constraint fk_images_vehicle_id
        foreign key (vehicle_id)
            references vehicles(id)
                on delete cascade
                on update cascade
);




create table videos (
    "vehicle_id" int not null,
    "video" varchar(255) not null,
    "created_at" timestamp not null default now(),
    constraint fk_videos_vehicle_id
        foreign key (vehicle_id)
            references vehicles(id)
                on delete cascade
                on update cascade
);




-- 06.12.2025
-- Update all brands' model_count based on the actual number of models per brand
UPDATE brands
SET model_count = COALESCE(sub.model_count, 0)
FROM (
    SELECT brand_id, COUNT(*) AS model_count
    FROM models
    GROUP BY brand_id
) AS sub
WHERE brands.id = sub.brand_id;
-- Optionally set brands with no models to 0
UPDATE brands
SET model_count = 0
WHERE id NOT IN (SELECT DISTINCT brand_id FROM models);


delete from brands;
delete from engines;
delete from body_types;
delete from fuel_types;
delete from drivetrains;
delete from transmissions;
delete from horse_powers;
delete from generations;
delete from generation_modifications;
delete from models;
delete from vehicles;
delete from images;



insert into profiles (user_id, username, registered_by) values (87, 'admin', 'admin');
insert into profiles (user_id, username, registered_by) values (84, 'admin2', 'admin');


UPDATE conversations 
		SET 
			user_1_unread_messages = CASE 
				WHEN user_id_1 = 2 THEN 0 
				ELSE user_1_unread_messages 
			END,
			user_2_unread_messages = CASE 
				WHEN user_id_1 != 2 THEN 0 
				ELSE user_2_unread_messages 
			END,
			updated_at = NOW()
		WHERE id = 1;



select 
			vs.id,
			bs.name,
			rs.name,
			cs.name,
			cls.name,
			ms.name,
			ts.name,
			es.name,
			ds.name,
			bts.name,
			fts.name,
			vs.year,
			vs.price,
			vs.odometer,
			vs.vin_code,
			vs.credit,
			vs.new,
			vs.status,
			vs.created_at,
			vs.trade_in,
			vs.owners,
			vs.crash,
			vs.updated_at,
			images.images,
			videos.videos,
			vs.phone_numbers,
			vs.view_count,
			json_build_object(
				'id', pf.user_id,
				'username', pf.username,
				'avatar', '` + r.config.IMAGE_BASE_URL + `' || pf.avatar,
				'role_id', u.role_id,
				'contacts', pf.contacts
			) as owner,
			vs.description
		from vehicles vs
		left join generation_modifications gms on gms.id = vs.modification_id
		left join colors cls on vs.color_id = cls.id
		left join profiles pf on pf.user_id = vs.user_id
		left join users u on u.id = vs.user_id
		left join brands bs on vs.brand_id = bs.id
		left join regions rs on vs.region_id = rs.id
		left join cities cs on vs.city_id = cs.id
		left join models ms on vs.model_id = ms.id
		left join transmissions ts on gms.transmission_id = ts.id
		left join engines es on gms.engine_id = es.id
		left join drivetrains ds on gms.drivetrain_id = ds.id
		left join body_types bts on gms.body_type_id = bts.id
		left join fuel_types fts on gms.fuel_type_id = fts.id
		LEFT JOIN LATERAL (
			SELECT json_agg(img.image) AS images
			FROM (
				SELECT image as image
				FROM images
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) img
		) images ON true
		LEFT JOIN LATERAL (
			SELECT json_agg(v.video) AS videos
			FROM (
				SELECT video as video
				FROM videos
				WHERE vehicle_id = vs.id
				ORDER BY created_at DESC
			) v
		) videos ON true
		where vs.status = 3
		order by vs.id desc;


 select
    vs.id,
    bs.name as brand,
    rs.name as region,
    cs.name as city,
    cls.name as color,
    ms.name as model,
    ts.name as transmission,
    es.name as engine,
    ds.name as drive,
    bts.name as body_type,
    fts.name as fuel_type,
    vs.year,
    vs.price,
    vs.odometer,
    vs.vin_code,
    vs.credit,
    vs.new,
    vs.status,
    vs.created_at,
    vs.trade_in,
    vs.owners,
    vs.crash,
    vs.updated_at,
    images.images,
    videos.videos,
    vs.phone_numbers,
    vs.view_count,
    CASE
            WHEN vs.user_id = 1 THEN TRUE
            ELSE FALSE
    END AS my_car,
    json_build_object(
            'id', pf.user_id,
            'username', pf.username,
            'avatar', 'https://api.mashynbazar.com/api/v1' || pf.avatar,
            'role_id', u.role_id,
            'contacts', pf.contacts
    ) as owner,
    vs.description,
    CASE
            WHEN ul.vehicle_id IS NOT NULL THEN true
            ELSE false
    END AS liked
from vehicles vs
left join generation_modifications gms on gms.id = vs.modification_id
left join colors cls on vs.color_id = cls.id
left join profiles pf on pf.user_id = vs.user_id
left join users u on u.id = vs.user_id
left join brands bs on vs.brand_id = bs.id
left join regions rs on vs.region_id = rs.id
left join cities cs on vs.city_id = cs.id
left join models ms on vs.model_id = ms.id
left join transmissions ts on gms.transmission_id = ts.id
left join engines es on gms.engine_id = es.id
left join drivetrains ds on gms.drivetrain_id = ds.id
left join body_types bts on gms.body_type_id = bts.id
left join fuel_types fts on gms.fuel_type_id = fts.id
left join user_likes ul on ul.vehicle_id = vs.id AND ul.user_id = 1
LEFT JOIN LATERAL (
    SELECT json_agg(img.image) AS images
    FROM (
            SELECT 'https://api.mashynbazar.com/api/v1' || image as image
            FROM images
            WHERE vehicle_id = vs.id
            ORDER BY created_at DESC
    ) img
) images ON true
LEFT JOIN LATERAL (
    SELECT json_agg(v.video) AS videos
    FROM (
            SELECT 'https://api.mashynbazar.com/api/v1' || video as video
            FROM videos
            WHERE vehicle_id = vs.id
            ORDER BY created_at DESC
    ) v
) videos ON true
where vs.status = 3 and vs.id > 9999

order by vs.id desc
limit 50;

