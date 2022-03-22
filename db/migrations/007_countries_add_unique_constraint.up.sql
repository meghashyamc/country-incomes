ALTER TABLE countries
ADD CONSTRAINT iso_gdp_per_capita_year UNIQUE (iso,gdp_percapita_ppp,year);