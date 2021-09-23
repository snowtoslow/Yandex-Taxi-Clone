A. Authentication service( role based auth + jwt token) + Handle of user service
    
    1. Sign in
    2. Sign up
    3. Authorization (to gives permission to access a resource for a user based on provided role)
    4. CRUD for users;
    5. SetUserStatus => will change statuses of a user;
    **User which is driver can have a status like: working, paused, not_working

B. Order service:

    1. CRUD operations for an order;
    2. List Orders by UserID;
    3. Every order should have a status like: waiting, in progress, finished


C. Car service:

    1. CRUD operations for a car model;
    2. Possibility to book a car if it is free and is in nearest area with a user which performs an order
    3.Every car should have a status like: busy, pause, available;
    4. SetCarStatus;
D. Notification service:
    
    1. CRUD operations for a notifications;
    2. Notications should have types: create_order, order_is_taken

Based on role of a user, it can perform different types of operations which is going to be described:
   ROLES:

A. Customer User:

      1. Where he is located
      2. Where he want to arrive
      3. Choose type of car which he wants to order(Comfort, Economy, something in the middle) => *here can be specified a special classification on *RELATED TO CAR SRV
      4. Should have a history of made orders -> *RELATED TO ORDER SRV
      5. Can have a rating as a customer

B. Driver User
   
    1.A user which has driver role has the possibility to view available orders;
    2. Can accept or decline an order which is made by a simple user;
    3. Can access more specific information related to a car;
    4. Should have  completed orders history -> *RELATED TO ORDER SRV
    5. Can have rating as a driver( worker for company)

C. Car owner:
    
    1. User which can create new cars in our application -> *RELATED TO CARS SRV
    2. Can set cars for a DRIVER-USER;




Workflow:

      After a user login, or enter our application, now it can create an order, 
      where he should specify FROM(where he is located at the moment)(*order)
      and TO(where he want to arrive)(*order), 
      (4)after this a query is performed to list all available cars now(*cars) 
      with users  with “working status”(in the best scenario which is 
      located somewhere near the FROM location, 
      but not mandatory for now), and find the user which is the driver 
      for this car. A driver(1) should be notified(2), when a notification 
      status is changed, Car status should be changed to “busy”, 
      order status should be changed to “in progress”. When a driver finish the 
      order, it should send a notification (*notification) with finished text, 
      after this car status(*cars) should be changed to “available” and order 
      status(*order) to finished.


Workflow refactored:

      After created a new account or logs in, he is able to create 
      orders(*order service), for creating an order a user should first of all 
      create a notification(*notification service), notification should be created with a 
      status: “create_order_notification”, after a notification gets created, 
      a call is performed to car service to find all cars which have status “available” 
      based on an algorithm(from discussion) , also we should be careful that the user
      which is the driver of a car should have a status “working”. The found users should 
      receive the notification, and if one of them accepts the notification, an order 
      with status “in progress” is created(*order service), car status 
      is updated to “busy”(*car service). When the road is finished, 
      the user driver should send a request, 
      which firstly is going to update the order to status order status to 
      “finished”(*order service), the car status should be changed 
      to “available”(*car service) and notification should be 
      deleted(*notification service)

*also we should be careful, if a user has is in pause or not_working the car should have also a special status for this;


Technologies to be used:

      1. Golang;
      2. Rust;
      3. Redis(for cache) => gateway;
      4. PostgresSQL;
      5. RavenDB or Couchbase( just to try);
      6. RethinkDB(real time database) for notification service if it will be separated;
      7. gRPC;
      8. HTTP to gateway;
      …. so on :)
