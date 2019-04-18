# Barbados Orderbook

Order book matching functional code for a single order book

The order book is a an ordered list of orders. The usual sorting order
is price and then time. That is the oldest order at a price level will
be matched first. Buy and sell orders are ordered differently. As
follows:

*  Buy orders are sorted high to low where high is the top of book 
   where top will be matched first
*  Sell orders are sorted low to high where low is the top of book
   where top will be matched first

Price time does not have to be the only sort order as orders could also 
be ordered by price and size. That is the largest order gets precedence. 
This is not usual in a continuous market but may be applicable 
during an auction when orders are submitted and at auction close
orders are sorted to favour the person with the largest order, thereby
encouraging greater liquidity.

Orders in the context of an order book are both commands and state. 
That is orders such as new or cancel (commands) are submitted that 
may result in trade executons that fill or partially fill (responses)
the order (state). 

An order book needs to accommodate all these ideas:
* Commands that try and change state, for example
  * new
  * cancel
* Responses to commands, for example
  * new order ack
  * new order rejected
* Responses as a result of some other action, for example
  * trade executions when two orders are matched
  * order expiry
* State of current orders, for example
  * order that is not matched or partially matched
* State of completed orders, for example
  * cancelled 
  * rejected
  
#### Using time for sorting orders
At a given price level orders need to be sorted by oldest to newest where
the oldest have precedence over the newest in the matching process. 
However, it is possible that two orders could arrive at the order book
with the same time. Fortunately GO has nano-second precision of time.
Unfortunately this may be too course on the system the
program executes on. 

One option is to use a time UUID that has a sequence
number to preserve ordering. However, the sequence is internal to the UUID
generator and doesn't guarantee ordering across machines so cannot be
supplied to the order book from another machine. So the order
needs to be timestamped and correctly sequenced when it enters the order 
book. A UUID is conveniently arranged to allow ordering comparisons by
just comparing bytes so can be used as the comparitor. Also conveniently
its possible to extract a time form the UUID so timestamps don't need to
duplicated.

One thing missing from the go UUID package is the ability to generate a
UUID from a timestamp. While this is not useful in day to day operations
its useful if orders needs to be created for testing or back testing 
purposes. However, these generated UUIDs cannot be mixed with
UUID's created from a wall clock. 

