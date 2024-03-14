	create or  replace procedure set_location(
		device_imei varchar,
		longitude varchar[],
		latitude varchar[]
	)
	language plpgsql
	as $$
	begin

	delete from devices_location where devices_location.imei=device_imei;

	insert into devices_location(id,imei,longitude,latitiude) values(gen_random_uuid(),device_imei,longitude,latitude);

	commit;
	end;
	$$;