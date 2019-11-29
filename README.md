# dr

![alt text][logo]

![alt text][status]

manage digital resources with expiry, reuse, category, id, and description

## Background
Digital resources such as descriptions of experiments, user interfaces, and hardware often require to be dynamically updated, but some are multi-use and some are single-use, and a finite lifetime (maybe even a few minutes) is assumed, so a typical TTL-aware key-value storage method such as ```ledis``` is pretty close to what we need, but not quite there. I'd like to avoid tying business rules to a specific database implementation. And I want to  an additional feature which is a description that can be returned when it is desired to preview a single-use resource (e.g. for selecting between resources according to some criteria (beyond simple category) that are expressed in the description, such as location, ping time, reliability, cost etc

## Business rules
Resources can be single-use or multi-use. Single-use resources are just that - the resource (a string) can only be revealed once, but in order to wisely select which resource to get, a description of the resource's qualities should be given (e.g. what type of experiment it is, how much it costs, where it is based, etc). 


## APIs

There is no need for long lived connections, so a REST-like API is sufficient ... for now. 

### restapi

A basic REST-like API is intended to offer additional convenience by offering operations with greater power, e.g. submitting multiple resources at once. Security considerations may drive some modifications to this API during the next phase of development.


#### Security
The REST(-ish) API combines user, and admin features. For example, submitting, updating and deleting tokens are admin roles, 
There's a very minor security issue around DELETE - it is a slightly more weaponised command than GET in the wrong hands, yet is usable at the same endpoint as the user-facing commands, and although I semi want to deprecate DELETE asap, there is likely always a human involved in building the UI for experiments that are going to be loaded into this system, and we need a way to be friendly and support an UNDO-like operation. So it probably stays...As for fine-grained security, the api may have to have some understanding of roles, if it proves incompatible with the seceurity system to filter on method as well as endpoint path ... to be continued.  

#### the sneaky Delete()
I swithered over Delete() - I didn't initially include it because it violates the policy of a system that does not rely on two-part state transitions for ordinary running. If an unreliable actor is present, then let each atomic action be sufficient in its own right for running of the system, and let any negative effects of not being around to handle a future part of the interaction fall on the faulty party, as it were. So this implies that you don't submit tokens that are wrong, because it opens the door to the supplier changing their mind, and weakens the trust you might place in a token even if it has a generous TTL. Of course, if the kit behind a token has gone offline then an Update() is advisable so that the change can be inferred from comparing the old and new token of the same ID. HOWEVER, RESTful interfaces have a DELETE and in copying in some mux.Router setup code from ```github.com/timdrysdale/vw``` I realised that there may come a time when a human involved in setting up tokens (e.g. for webpages) might make a mistake and need to. 


## Architectural thoughts

The interface and resource struct  ```dr.go``` are intended to be the abstract, stable(-ish!) representation, at the core of the onion. The common test expresses the expected behaviour, depending on only on those definitions. The implementations of the storage fall in the next layer, and these currently must implement the business rules that are expressed in the tests. It seemed unnecessarily verbose and unperformant to abstract the business rules themselves from the particular storage method, because storage methods implementing TTL or reusable features would be best left to handle that in their optimised and tested way. In any case, the first implementation of the business rules would require testing, and that testing would require a store, so I compromised by writing ```./ram``` which combines the business rules with using native golang structs stored in memory as a reference design with which the generic tests ```./test/test.go``` can be validated.

![alt text][arch]

### Storage for restart
This feature is omitted on the grounds that future usage and immediate testing needs are not predicated upon expectation of having a valuable, long lasting dataset that is difficult to load. Quite the opposite. Anything that is hard to set up, is not going to fit the bill for wider use anyway. Plus, previous experience of server failure mitigation suggests that failure blast radius and system recovery time are both proportional to the mean lifetime of the most-used data in the system. So, you can do a lot worse than design systems with short lifetimes in them, and avoid altogether the issue of trying to failover with already fatally-corrupted data set (not a good day out). Start clean and reconstruct what you need from a trusted corruption source. Short lifetime expectations also make systems more amenable to deployment on spot-priced servers. Bonus 90% compute saving. No one moan about premature optimisation please. 

### Deployment
Given the small size of the initial amount of experiments to be served over the following months, it is a debatable YAGNI point whether the various implementations of the layers of the onion need to be split into their own separate repositories, and whether the storage and the api need to separated so that new apis can be added without restarting the store - which of course only applies to ```./ram``` or some other in-memory embedded database (e.g. ```github.com/boltdb/bolt```) which implies it is only a problem for small scale operation where reloading the existing shortlived data should not be onerous (and provide a sense of how it is to operate with this approach). And in any case, there is nothing stopping said interface from being developed separately and connecting to the existing ```restapi``` - afterall, some sort of store-facing API is needed if the user-facing API is to be put in a separate package.

As background info, in any case I am already considering consolidating ```agg```,```hub```,```rwc```,```rcws``` into ```vw``` to simplify troubleshooting conversations with new users adopting ```vw```, although the barrier to that is re-use in ```crossbar``` and ```hbar``` - but it is not really a genuine reuse, just a convenience for as long as these three codes develop with the same goals in mind (which is temporary situation that is true for now). In any case, even if sharing modules, a smarter solution would be to engage with the dependency management in more recent versions of go (at the cost of then having to diagnose and fix at a distance any errors relating to dependency issues caused by updates that will inevitably come.





[arch]: ./img/arch.png "Onion diagram of the arch centre to outside is dr, test, ram, testapi, main" 
[logo]: ./img/logo.png "logo for dr, clock and files"
[status]: https://img.shields.io/badge/alpha-do%20not%20use-orange "Alpha status, do not use" 