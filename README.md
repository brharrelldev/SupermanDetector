# Superman Detector
### Overview
Exercise to detect ip addresses based on login info.   This is working with a sample dataset that I built.  The program will check for a SQL database in the beginnning and create one if one doesn't current exists.

This program consist of 2 parts.  Client and Server.
 
The Server is the "webservice" piece, and this does the processing, including calculating the Haversine.   
 
The client is a wrapper to a HTTP Client.   This will bring will simply connect to the server, and specify logins which we want to check

### Server

Server is a subcommand, below is an example of how to use it

`./super-cli --db-file=logins.db --dataset=data1.csv  server --port=3000 --logindb=logins.db   --mmdb=GeoLite2-City.mmdb`

Flags:

````
--loginsdb    Local database that has login info and ip address
--port        Port server will run on
--mmdb        GeoLite2 database with geolocation information about IP address
````

### Client

```
--login       Username located in database, this is the name we want to look up
--addr        Address up super-cli server
```

Dockerfile works and compiles.  You can incorporate this into your workflow.

 
 
 
