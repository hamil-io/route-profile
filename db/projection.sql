/*
This is the projection used for the NOAA High Resolution Rapid Refresh rasters.
The system was determined by running gdalinfo on the downloaded GRIB files.
*/
INSERT into spatial_ref_sys
    (srid, auth_name, auth_srid, proj4text, srtext)
    values (
        98411,
        'sr-org',
        8411,
        '+proj=lcc +lat_1=38.5 +lat_2=38.5 +lat_0=38.5 +lon_0=-97.5 +x_0=0 +y_0=0 +ellps=sphere +a=6371229 +b=6371229',
        'GEOGCS["Coordinate System imported from GRIB file", DATUM["unknown", SPHEROID["Sphere",6371229,0]], PRIMEM["Greenwich",0], UNIT["degree",0.0174532925199433]], PROJECTION["Lambert_Conformal_Conic_2SP"], PARAMETER["standard_parallel_1",38.5], PARAMETER["standard_parallel_2",38.5], PARAMETER["latitude_of_origin",38.5], PARAMETER["central_meridian",262.5], PARAMETER["false_easting",0], PARAMETER["false_northing",0]');
