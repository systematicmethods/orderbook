
Feature: Market Order Fill

  Scenario: Buy market order rejected as order book empty
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | John_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             |
    Then await 1 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | Rejected        | John_01   | ABV        | Buy  | Market  | John_01m |       | 1000 |          |                |             | 0       |           | 0      | Rejected  | Rejected  |        | Not Null | Not Null | 1000      |              |


  Scenario: Sell market orderas order book empty
    Given An order book for instrument "ABV"
    When users send orders with:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID |
      | NewOrder        | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             |
    Then await 1 executions
    And executions should be:
      | Event           | ClientID  | Instrument | Side | OrdType | ClOrdID  | Price | Qty  | ExpireOn | TimeInForce    | OrigClOrdID | LastQty | LastPrice | CumQty | Status    | ExecType  | Reason | ExecID   | OrderID  | LeavesQty | TransactTime |
      | Rejected        | Bill_01   | ABV        | Sell | Market  | Bill_01m |       | 1000 |          |                |             | 0       |           | 0      | Rejected  | Rejected  |        | Not Null | Not Null | 1000      |              |


