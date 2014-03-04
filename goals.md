==

## goals:
simple really...
* frontend: 
* * main page 
* * * displays a "most popular", i.e. most points on reddit
* * * also displays a scrolling list of recently submitted maps
* * * above map is a link submission form to send a thread to be processed
* * search page:
* * * a search by "food type", location, and stuff
* backend:
* * link processing server
* * * some entity extraction service will pull out the names of places and save them (apialchemy)
* * * some entity will run searches to get address of places
* * * some entity will geocode addresses to lat/lon (mapquest)
* * * some entity will build GeoJSON and/or kml layers for rendering on map (mapquest again?)
* rendering server
* * select all places for a particular thread and send the output to the client

