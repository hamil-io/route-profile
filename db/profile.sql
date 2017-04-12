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
