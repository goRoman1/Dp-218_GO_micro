# Scooter Service

The application is similar to uber service. Where you as a:
* supplier  - can get your scooters for rent.
* user - can take a scooter for rent.

# Run

The application runs with docker.
```
docker-compose up --build -d
```

You should configure your run environment in the IDE you use.  

Server runs on http://localhost:8080/  



# The main endpoints:
```/``` - the main page.  
```/login``` - sign-in or sign-up  
```/customer/map``` - here a user can see the nearest stations to his location.  
```/start-trip/1``` - here a user can see all the scooters on the chosen station, choose the destination station and start a trip. Index in a sub-domain depends on the chosen station.  
````/users```` - the list of all users and their statuses.  
```/stations``` - the list of all stations.  
```/scooters``` - the list of scooters in json format.
```/models``` - add scooter models/scooters
```/init``` - place scooters on stations

# How to start the trip

On the page ```http://localhost:8080/customer/map``` you can choose a departure station.  

On click "show station" button you will move to the ```http://localhost:8080/start-trip/{station_id}``` page.
Which shows you all available scooters on the chosen station. Here you also choose the destination station.  

"Start trip" button will start your trip with chosen scooter to the chosen station.

Information about trips will be written to the database table - "Orders".


