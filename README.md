# Mini-UPS

## structure

```bash
mini-ups
├── config
│   ├── config.go   # config reader
│   └── config.yaml # config file
├── controller      # controller(entry)
│   └── user_contoller.go
├── dao             # database access object(operate data here)
│   └── user_dao.go
├── db              # database
│   └── db.go
├── Dockerfile
├── frontend        # frontend html,css,js (divided by module)
│   ├── home
│   │   ├── home.html
│   │   └── home.js
│   ├── login
│   │   ├── login.html
│   │   └── login.js
│   ├── register
│   │   ├── register.html
│   │   └── register.js
│   └── style.css   # shared css
├── go.mod          # go project module config
├── go.sum          
├── main.go         # program entry
├── model           # database model
│   └── user.go
├── router          # router (register url here)
│   └── router.go
├── service         # called by controller to do core logic, can call dao
│   └── user_service.go
└── util            # any util you want
    └── util.go

```

## Quick Start
> before you start, you may change the config.yaml
```yaml
...

amazon:
  host: // your amazon host
  port: // your amazon host

ups:
  host: // your ups host
  port: 8080

...
```
> deploy
```bash
cd docker_deploy

// deploy
sudo docker-compose up --build

// shut down
sudo docker-compose down
```

# Additional Features
On top of basic features of UPS like tracking packages and view your own package info, and changing package destination. We implemented other features that help differentiate our project form our comptetition (other IG grous). These features are all customer-focused, meaning that they aim to create a faster,more flexible, and transparent delivery experience.

### Truck Logistics Algorithm
We used a `FIFO` algorithm to ensure users receive their packages as quickly as possible and in a fair way. 

If no trucks are available at the time of request, the system system automically assigns the package to the next truck that becomes idle. This provides fairness to each user when using our UPS service.

### Package Prioritization
In the real world, users may have important packages that needs to arrive at the destination as soon as possible. That being saiw, wee created
a feature that allows users to press a button to `prioritize` certain pacakgges. This algorithm 

### Share Link
- user can easily copy a `share link` of their package page by clicking the button "share link" at right-up of the package page, and then share to others

### Stateless Sessions
- use `cookies & sessions` to keep user logined

- use `stateless` sessions, which means all info except private key are stored in client side, so it is adaptive to distrubited deployment



# world_simulator_exec
This is the executable file wrapped in a docker for world_simulator, background world for ECE 568 final project.
No source code included.

To run the world simulator without flakiness, type: "sudo docker-compose build" in the same directory with the yml file.
Then type "sudo docker-compose up". 

To run the world simulator with a different flakiness, go to "docker-compose.yml" file, change the command from bash -c "./wait-for-it.sh mydb:5432 --strict -- ./server 12345 23456 0" into bash -c "./wait-for-it.sh mydb:5432 --strict -- ./server 12345 23456 <flakiness_num>". Then do as last paragragh said.

Note flakiness ranges from 0 to 99. When flakiness equels 0, it mean the world will not deliverately drop any request it receives. As flakiness grows, the possibility that the world randomly drops requests will be larger. You can view this behavior as in real life "error in communication". That's also why we are having ack number for each requests. 
