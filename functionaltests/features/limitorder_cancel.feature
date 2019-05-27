
Feature: Limit Order Cancel

  Scenario: Buy and sell limit orders and cancel
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel |             |
      | Cancel       | John_01   | ABV        | Buy  | Limit   | John_02o | 1.03  | 1000 |          | GoodTillCancel | John_01o    |
      | Cancel       | Bill_01   | ABV        | Sell | Limit   | Bill_02o | 1.05  | 2000 |          | GoodTillCancel | Bill_01o    |
    Then await 4 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status    | Reason | ExecID   | OrderID  |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New       |        | Not Null | Not Null |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.05  | 2000 |          | GoodTillCancel | 0       |           | 0      | New       |        | Not Null | Not Null |
      | Cancelled    | John_01   | ABV        | Buy  | Limit   | John_02o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | Cancelled |        | Not Null | Not Null |
      | Cancelled    | Bill_01   | ABV        | Sell | Limit   | Bill_02o | 1.05  | 2000 |          | GoodTillCancel | 0       |           | 0      | Cancelled |        | Not Null | Not Null |
