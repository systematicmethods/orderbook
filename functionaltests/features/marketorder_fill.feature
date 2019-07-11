
Feature: Market Order Fill

  Scenario: Buy market order
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             |
    Then await 3 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | Filled          | John_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             | 1000    | 1.05      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |              |
      | Filled          | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             | 1000    | 1.05      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |              |


  Scenario: Sell market order
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             |
    Then await 3 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | Filled          | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |              |
      | Filled          | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null | 0         |              |


  Scenario: Buy market order partially filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Market  | John_01m |       | 1100 |          |                |             |
    Then await 3 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status          | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Market  | John_01m |       | 1100 |          |                |             | 0       |           | 0      | New             | New       |        | Not Null | Not Null | 1100      |              |
      | PartiallyFilled | John_01   | ABV        | Buy  | Market  | John_01m |       | 1100 |          |                |             | 1000    | 1.05      | 1000   | PartiallyFilled | Trade     |        | Not Null | Not Null | 0         |              |
      | Filled          | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             | 1000    | 1.05      | 1000   | Filled          | Trade     |        | Not Null | Not Null | 0         |              |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | New             | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |         |           | 0      | New              |           |        | Not Null | Not Null | 1000      |              |           |           |           |


  Scenario: Sell market order partially filled
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1100 |          |                |             |
    Then await 3 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status          | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1100 |          |                |             | 0       |           | 0      | New             | New       |        | Not Null | Not Null | 1000      |              |
      | PartiallyFilled | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1100 |          |                |             | 1000    | 1.03      | 1000   | PartiallyFilled | Trade     |        | Not Null | Not Null | 0         |              |
      | Filled          | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 1000    | 1.03      | 1000   | Filled          | Trade     |        | Not Null | Not Null | 0         |              |
    And order state should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | New             | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |         |           | 0      | New              |           |        | Not Null | Not Null | 1000      |              |           |           |           |


  Scenario: Buy market order can't trade with yourself
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | Bill_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | Bill_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | Rejected        | Bill_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             | 0       |           | 0      | Rejected  | Rejected  |        | Not Null | Not Null | 1000      |              |


  Scenario: Sell market order can't trade with yourself
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder        | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 1000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | NewOrderAck     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             |
    Then await 2 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | New             | John_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             | 0       |           | 0      | New       | New       |        | Not Null | Not Null | 1000      |              |
      | Rejected        | John_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             | 0       |           | 0      | Rejected  | Rejected  |        | Not Null | Not Null | 1000      |              |
