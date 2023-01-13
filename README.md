# id_generator

### Finch ID
Finch IDs are 64-bit unique identifiers that can be used in distributed computing. 

#### Properties
- Unique
- Monotonically increasing (guarenteed at millisecond sensitivity, i.e ids created in same millisecond may not be monotonically increasing)
- Virtually impossible to predict or scrape ids
- Cloud compatible

#### Design
Each instance of the id_generator server is assigned a unique generator id. This id is not included in the generated Finch IDs, but is used as a input in the generation algorithm to prevent duplication. Furthermore, all servers share a common secret seed.

##### Schema
<pre>
                      42                           12              10
 +-----------------------------------------+----------------+---------------+
 |     milliseconds since start-epoch      |   random part  |  unique part  |
 +-----------------------------------------+----------------+---------------+
</pre>

The number of bits for the random part and unique part are adjustable.

##### Algorithm
To generate a new id:
1. Retrieve the number of milliseconds that have passed since the *start-epoch*.
2. Create a random 12-bit value to be used as the random component. This random component should be generated using a different method than the shared seed, such as a PRNG with safety measures to prevent duplication within the same millisecond.
3. The *prefix* refers to the Finch ID with the timestamp and random component already included, but the unique part set to 0. Verify that this prefix has not been previously used with the current generator Id. If it has been used, generate a new random component or wait for the next millisecond. It is acceptable for multiple instances to generate the same prefix, as the unique part will ensure that the final Finch IDs are distinct.
4. Create the unique component, which should appear random but is guaranteed to be one-of-a-kind for a specific generator id assigned to an individual instance:
    - Initialize a pseudorandom number generator using the shared seed and the prefix as a nonce.
    - Use the aforementioned random number generator to generate a permutation of the set of 2^10 (1024) generator Ids. As the random number generator was initialized with the shared seed and prefix in a deterministic manner, all servers will produce the same permutation for a specific prefix.
    - Use the permutation as a means of determining the unique component by selecting the element associated with the specific generator Id from the permutation.
5. Assemble the Timestamp, Random, and Unique components along with a leading reserved bit to form a Finch ID.

This algorithm relies on the assumption that it is possible to globally enforce that no more than one instance uses the same unique generator id at any given time.

Due to the restriction that no two instances can run concurrently with overlapping generator ids, it is guaranteed that no two instances will select the same element of the permutation, and as a result, no two instances will select the same unique component from the generated permutation for a given prefix. Since the permutation is created using cryptographic methods based on the seed, it should be impossible to correlate Finch IDs generated by the same instance or even within the same cluster.

#### Performance
On a 1.8GHz and 16GB RAM PC, one instance of id_generator can generate 15K IDs per second.  
The maximum possible number of id_generator instances = 2^(number of bits in unique part) = 1024.  
Hence, can generate around 15 million IDs per second.

___

### Usage
[1] update [zookeeper/setup.go](https://github.com/nova-plixx/id_generator/blob/main/zookeeper/setup.go) with desired
  - zookeeper ensemble connection string
  - shared int64 seed
  - start epoch milli
 
it defaults to 
  - zookeeper ensemble in localhost => *127.0.0.1:2181*
  - shared int64 seed => *9223372036854775783*
  - start epoch milli => *1577836800000* => *2020-01-01 00:00:00.000*
 
[2] run the zookeeper setup file => `go run zookeeper/setup.go`

[3] update [server/main.go](https://github.com/nova-plixx/id_generator/blob/main/server/main.go) with desired
  - zookeeper ensemble connection string

[4] run the server file => `go run server/main.go`

[5] sample client code can be found at [example-client/main.go](https://github.com/nova-plixx/id_generator/blob/main/example-client/main.go)
