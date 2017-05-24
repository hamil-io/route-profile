/*
Interpolate takes a line geometry and a sample resolution and returns a new
geometry with points sampled at the specified resolution.
*/
CREATE OR REPLACE FUNCTION public.interpolate(line geometry, sample numeric)
RETURNS geometry
LANGUAGE plpgsql
AS $function$
DECLARE
    points geometry;
BEGIN
    WITH linemesure AS
       (SELECT ST_AddMeasure(line, 0, ST_Length(line)) as linem,
               generate_series(0, ST_Length(line)::numeric, sample) as i)
       SELECT ST_MakeLine(ST_GeometryN(ST_LocateAlong(linem, i), 1)) INTO points FROM linemesure;
    RETURN points;
END;
$function$;
