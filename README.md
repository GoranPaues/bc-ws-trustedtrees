# oow-18
Prep repository for OOW-18 Code ONE ws

Current situation: A HLF network of 3 parties:
1. Claudio's Controlling (Founder) - Owners of the Trusted Trees brand - The Trusted Tree branding guarantees the lumber origin as well as the worker terms and it also guarantees that more trees are planted than cut. This brand attracts a lot of investors.
2. Franco's Forestry (Participant) - Plant & grow trees on site in Senegal under Trusted Tree umbrella.
3. Leo's Lumber (Participant) - Cut trees into lumber under the Trusted Trees brand.

Whats going on:
A second tree planting participant, Papé's Plantation, wants to join the network. It is our job to spin up their chaincode and adapt the code to this participant.

Suggested workflow:
- Log in to OABCS using the credentials provided.
- Locate the Papés Plantation instance.
- Create a channel with the Founder, Trusted Trees.
- Download the chaincode from GitHub in go of NodeJS - your choice!
- Upload, install and instantiate the chaincode on your channel.
- Invoke the chaincode using REST (PostMan or cUrl recommended).
- Update the code with a new price per tree for this participant.
- Upgrade the chaincode with the new version.

Challenge:
- Claudio's controlling could really benefit from a "Query all" function. Please help them out!

Optional:
- Send a pull request with your new code on GitHub.
