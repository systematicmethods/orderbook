
Feature: Limit Orders

  Scenario: Buy and sell limit orders
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel |
    Then await 2 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status | Reason | ExecID   | OrderID  |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New    |        | Not Null | Not Null |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel | 0       |           | 0      | New    |        | Not Null | Not Null |
