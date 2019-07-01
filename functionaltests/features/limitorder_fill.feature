
Feature: Limit Order Fill

  Scenario: Buy and sell both filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |
      | Filled       | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |
      | Filled       | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |
    And order state should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |


  Scenario: Buy and sell; buy partially filled and sell filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1100      |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1000      |
      | PartiallyFilled | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel | 1000    | 1.03      | 1000   | PartiallyFilled  | Trade     |        | Not Null | Not Null | 100       |
      | Filled          | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled           | Trade     |        | Not Null | Not Null | 0         |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | PartiallyFilled | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel |         |           | 1000   | PartiallyFilled  |           |        | Not Null | Not Null | 100       |             |           |           |           |


  Scenario: Buy and sell; buy filled and sell partially filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1100      |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1000      |
      | Filled          | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled           | Trade     |        | Not Null | Not Null | 0         |
      | PartiallyFilled | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel | 1000    | 1.03      | 1000   | PartiallyFilled  | Trade     |        | Not Null | Not Null | 100       |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | PartiallyFilled | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel |         |           | 1000   | PartiallyFilled  |           |        | Not Null | Not Null | 100       |              |           |           |           |



  Scenario: Sell and buy; sell partially filled and buy filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel |             |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1000      |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1100      |
      | Filled          | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled           | Trade     |        | Not Null | Not Null | 0         |
      | PartiallyFilled | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel | 1000    | 1.03      | 1000   | PartiallyFilled  | Trade     |        | Not Null | Not Null | 100      |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | PartiallyFilled | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1100 |          | GoodTillCancel |         |           | 1000   | PartiallyFilled  |           |        | Not Null | Not Null | 100       |             |           |           |           |


  Scenario: Sell and buy; sell filled and buy partially filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1000      |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel | 0       |           | 0      | New              | New       |        | Not Null | Not Null | 1100      |
      | PartiallyFilled | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel | 1000    | 1.03      | 1000   | PartiallyFilled  | Trade     |        | Not Null | Not Null | 100         |
      | Filled          | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled           | Trade     |        | Not Null | Not Null | 0      |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | PartiallyFilled | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1100 |          | GoodTillCancel |         |           | 1000   | PartiallyFilled  |           |        | Not Null | Not Null | 100       |             |           |           |           |

