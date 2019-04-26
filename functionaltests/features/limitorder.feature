
Feature: Limit Orders

  Scenario: Buy and sell limit orders
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID | Price | Qty  | ExpireOn | TimeinForce |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01 | 1.03  | 1000 |          |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01 | 1.05  | 2000 |          |             |
    Then await 2 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID | Price | Qty  | ExpireOn | TimeinForce | Filled | Status | Reason | ExecID   | OrderID  |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01 | 1.03  | 1000 |          |             | 0      | New    |        | Not Null | Not Null |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01 | 1.05  | 2000 |          |             | 0      | New    |        | Not Null | Not Null |

