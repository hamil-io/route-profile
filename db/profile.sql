/*
Profile will return the Z components of the points along a drape line.
Profile should normally take in the output of the drape function. For the
default SRTM data the output will be in meters.
*/
CREATE OR REPLACE FUNCTION public.profile(line geometry)
RETURNS TABLE(
   z double precision
)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
        WITH points3d AS (SELECT (ST_DumpPoints(line)).geom AS geom)
        SELECT ST_Z(geom) FROM points3d;
END;
$function$;

CREATE OR REPLACE FUNCTION public.points(line geometry)
RETURNS TABLE(
   p text
)
LANGUAGE plpgsql
AS $function$
BEGIN
    RETURN QUERY
        WITH points3d AS (SELECT (ST_DumpPoints(line)).geom AS geom)
        SELECT ST_AsText(geom) FROM points3d;
END;
$function$;
