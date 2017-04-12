BEGIN;
CREATE TABLE "altitude" ("rid" serial PRIMARY KEY,"rast" raster);
CREATE INDEX ON "altitude" USING gist (st_convexhull("rast"));
END;
