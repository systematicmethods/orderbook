
Feature: Limit Orders

  Scenario: Buy and sell limit orders
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel |             |
    Then await 2 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status | ExecType | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             | 0       |           | 0      | New    | New      |        | Not Null | Not Null |           |              |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel |             | 0       |           | 0      | New    | New      |        | Not Null | Not Null |           |              |
    And order state should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status           | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime | CreatedOn | UpdatedOn | Timestamp |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |         |           | 0      | New              |           |        | Not Null | Not Null | 1000      |              |           |           |           |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel |             |         |           | 0      | New              |           |        | Not Null | Not Null | 2000      |              |           |           |           |

