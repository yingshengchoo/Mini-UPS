# Differentiation Features
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
