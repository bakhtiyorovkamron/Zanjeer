	create or  replace procedure set_location(
		imei varchar,
		longitude varchar[],
		latitude varchar[]
	)
	language plpgsql
	as $$
	begin

	delete from devices_location where devices_location.imei=imei;

	insert into devices_location(imei,longitude,latitude) values(imei,longitude,latitude);

	commit;
	end;
	$$;