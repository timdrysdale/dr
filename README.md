# dr
digital resource with category, id, description and expiry 

## Background
Digital resources such as descriptions of experiments, user interfaces, and hardware often require to be dynamically updated, but some are multi-use and some are single-use, and a finite lifetime (maybe even a few minutes) is assumed, so a typical storage methods are not sufficiently expressive. Some come close, such as redis and ledis which include TTL. However, it's not wise to tie business rules to implementation specifics, and in any case, we need an additional feature which is a description that can be returned when it is desired to preview a single-use resource (e.g. for selecting between resources according to some criteria (beyond simple category) that are expressed in the description, such as location, ping time, reliability, cost etc.

## Business rules
Resources can be single-use or multi-use. Single-use resources are just that - the resource (a string) can only be revealed once. 
