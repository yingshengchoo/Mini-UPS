# Dangerlog
1. What if there an overflow occurs due to large number of packages, numerous responses between us and world, large number of trucks/warehouses? We would have an error as the primary key constrain would be violated.

2. What if unauthorized entities use our api calls to change other users's package information or do authorized actions such as sending commands to our trucks or access other uses package information? This is a huge security flaw that can easily be abused
- We can add addtional checks to guarentee that each api call are from valid entities.

3. What if the server abruptly shuts down? We did not implement the code to restore our server back to the correct stsate. Additionally, we would lose information about the package order to be delivered. Simularly, we would lose data in the sendWindow which holds information on the commands we have yet recieved acknowledgement from world.
- we can backup these information in a database

4. Currently, our WorldHandler and AmazonHandler are serialized, meaning they process requests one at a time. If a hacker were to compromise Amazon and flood us with malicious or excessive requests, our system could suffer a denial-of-service (DoS) attack, as it would be forced to process each request sequentially without prioritization or throttling.
- We can add goroutines to make each resonse handler. The hard part about this is to ensure that there are no race conditions.

5. Although our World and Amazon handlers run concurrently, the handling of individual responses is still serialized. This design could become a major bottleneck if we need to scale and support multiple Amazon servers communicating with us simultaneously. To achieve true scalability, we would need to make the response handling itself concurrent.
- same as above.

6. What if a hacker intercepts the messages between world or amazon? Currently, are package ids and sequence numbers are sequential, so attackers can easily extract critical information about our system.
- Randomize package id and seqnum number ids to make it harder for hackers to guess. Maybe consider making package id strings to add more complexity, however, this would require some changes to the world simulator as well.

7. If #2 and #6 both happen, hackers can easily guess all existing packages.



