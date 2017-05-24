/*
Wind will calculate the headwing along the segments of a line geometry.
Segments are interpolated uniformly using the sample resolution. Headwind
is calculated from the U and V wind components from the loaded wind raster.
Only wind that is normal to the segment angle is considered headwind.
*/
CREATE OR REPLACE FUNCTION public.wind(line geometry, sample numeric)
RETURNS TABLE(
    length double precision,
    theta double precision,
    headwind double precision,
    wind_u double precision,
    wind_v double precision
)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
    SELECT s.length, degrees(s.theta),
        (ST_Value(rast, 1,  ST_Transform(start, 98411)) * x) + (ST_Value(rast, 2,  ST_Transform(start, 98411)) * y) AS magnitude,
        ST_Value(rast, 1,  ST_Transform(start, 98411)) AS wind_u,
        ST_Value(rast, 2,  ST_Transform(start, 98411)) AS wind_v
    FROM segments(interpolate(line, sample)) as s, wind
    WHERE ST_Intersects(rast, ST_Transform(start, 98411));
END;
$function$;
