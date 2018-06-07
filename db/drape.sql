/*
Drape takes a line and returns a 3d line with the Z component taken from the
loaded digital elevation model. Drape will interpolate sample points along
the provided geometry at the specified sample resoultion. The resolution should
be provided in spatial reference units. In our case this is degrees.
*/
CREATE OR REPLACE FUNCTION public.drape(line geometry, sample numeric)
RETURNS geometry
LANGUAGE plpgsql
AS $function$
DECLARE
    geom3d geometry;
BEGIN
 WITH linemeasure AS
    (SELECT ST_AddMeasure(line, 0, ST_Length(line)) as linem,
            generate_series(0, ST_Length(line)::numeric, sample) as i),
  points2d AS
    (SELECT ST_GeometryN(ST_LocateAlong(linem, i), 1) AS geom FROM linemeasure),
  cells AS
    (SELECT p.geom AS geom, ST_Value(altitude.rast, 1, p.geom) AS val
     FROM altitude, points2d p
     WHERE ST_Intersects(altitude.rast, p.geom)),
  points3d AS
    (SELECT ST_SetSRID(ST_MakePoint(ST_X(geom), ST_Y(geom), val), 4326) AS geom FROM cells)
    SELECT ST_MakeLine(geom) INTO geom3d FROM points3d;
    RETURN geom3d;
END;
$function$;

CREATE OR REPLACE FUNCTION public.drape(line geometry, __table text, sample numeric)
RETURNS geometry
LANGUAGE plpgsql
AS $function$
DECLARE
    geom3d geometry;
BEGIN
 EXECUTE format('
 WITH linemeasure AS
    -- Add a measure dimension to extract steps
    (SELECT ST_AddMeasure(%L::geometry, 0, ST_Length(%L::geometry)) as linem,
            generate_series(0, ST_Length(%L::geometry)::numeric, %L) as i),
  points2d AS
    (SELECT ST_GeometryN(ST_LocateAlong(linem, i), 1) AS geom FROM linemeasure),
  cells AS
    -- Get DEM elevation for each
    (SELECT p.geom AS geom, ST_Value(%I.rast, 1, p.geom) AS val
     FROM %I, points2d p
     WHERE ST_Intersects(%I.rast, p.geom)),
    -- Instantiate 3D points
  points3d AS
    (SELECT ST_SetSRID(ST_MakePoint(ST_X(geom), ST_Y(geom), val), 4326) AS geom FROM cells)
    SELECT ST_MakeLine(geom) FROM points3d', line, line, line, sample, __table, __table, __table) INTO geom3d;
    RETURN geom3d;
END;
$function$;
