/*
Segments take a geometry and converts each pair of points along the line into
segments. Segments have a start and end point, a length, and a theta component.
Theta is the angle from the equator for the segment. Segments will return one
less segment than the number of points in the line. This function is typically
used with interpolate to return a list of uniform length segments.
*/
CREATE OR REPLACE FUNCTION public.segments(line geometry)
RETURNS TABLE(
    n integer,
    theta double precision,
    length double precision,
    x double precision,
    y double precision,
    start geometry,
    stop geometry
)
LANGUAGE plpgsql
AS $function$
BEGIN
    FOR i in 1..ST_NPoints(line) LOOP
        n := i;
        IF n < ST_NPoints(line) THEN
            start := ST_PointN(line, n);
            stop := ST_PointN(line, n+1);
            length := ST_DistanceSphere(start, stop);
            theta := ST_Azimuth(start, stop);
            x := cos(theta);
            y := sin(theta);
            RETURN NEXT;
        END IF;
    END LOOP;
    RETURN;
END;
$function$


