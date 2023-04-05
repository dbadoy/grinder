# grinder
grinder is a server that collects metadata about contracts deployed on the ethereum network. 
All it does is fetch blocks from the ethereum network in order, check to see if any of the internal transactions are deployment transactions, and if so, extract the metadata from the contract's bytecode and store it. Before we talk about metadata, let's look at the example below.<br>
We have a requirement to collect all ERC-20 contracts that are currently deployed on the ethereum network, so the usual approach is to inspect the block to find the deployment transaction, check if it meets the ERC-20 interface, and store the contract address in the DB.
However, when implementing a server like this, there are limitations when it comes to scaling. For example, if you also need a list of deployed ERC-721 contracts, you have to cycle through the first block to the most recent block once more. A contract that meets interface A, a contract that meets interface B... This process must be done every time a requirement is added, which is a huge burden.
<br>
The key idea is not to check and store which interfaces a contract satisfies during the ingest phase, but to store only the metadata about the contract and determine if the interface is satisfied when a request comes in.
<br>
The goal is to extract a list of methodIDs and eventIDs from the contract's bytecode and store them in a inverted index DB to fulfill the request "Give me a list of contracts that meet the interfaces I gave you". 
<br>
If you want to find items containing specific keywords, a inverted index DB is ideal, but too many of them(specific keyword) can cause performance degradation. We need to benchmark it to make sure it is suitable for production.
