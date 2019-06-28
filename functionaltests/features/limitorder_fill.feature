
Feature: Limit Order Fill

  Scenario: Buy and sell limit orders and expect fills
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder     | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel |             |
      | NewOrder     | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel |             |
    Then await 4 executions
    And executions should be:
      | Event        | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  |
      | NewOrderAck  | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New       | New       |        | Not Null | Not Null |
      | NewOrderAck  | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 0       |           | 0      | New       | New       |        | Not Null | Not Null |
      | Filled       | John_01   | ABV        | Buy  | Limit   | John_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null |
      | Filled       | Bill_01   | ABV        | Sell | Limit   | Bill_01o | 1.03  | 1000 |          | GoodTillCancel | 1000    | 1.03      | 1000   | Filled    | Trade     |        | Not Null | Not Null |
