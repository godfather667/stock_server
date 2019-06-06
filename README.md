**stock_server**

# How to Build

Create the folder **stock_server** in your **golang src Directory**

insert **stock_server.go, LICENSE, .gitignore, and README.md** into that folder

~~~~~~~~~~~~~~~~~~~~~~~~~~~~
go/src
  :: ::
     stock_server
        .gitignore
        LICENSE
        README.md
        stock_server.go
   :: ::
   :: ::
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Execute **go install**


# How to Use
 Package stock_server implements stock price lookup on most exchanges.

 Execute **stock_server** or simply **go run stock_server.go** in source area.

 Server is on localhost:3000

 Responds to endpoint *"/stock/"*  -- Example: **localhost:3000/stock/**

 Specify the stock:  **symbol=stock, stock, ...**

 Exchanges are optional, append them after stocks: **&Stock_exchange=exchange, exchange, ...**

 If an exchange is not specified the default is **AMEX**.

~~~~~~~~~~~~~~~~~~~~~~~~
 Example:   http://localhost:3000/stock/symbol=MSFT,AAPL,FAX&stock_exchange=NASDAQ,AMEX

 Result of example:
	stock: AAPL	price  153.30	exchange  NASDAQ
	stock: FAX	price  3.99	exchange  AMEX
	stock: MSFT	price  105.68	exchange  NASDAQ
~~~~~~~~~~~~~~~~~~~~~~~~
## Errors:  

If a viable stock name is not found:
~~~~~~~~~~~~~~~~~~~~~~~~
     "Error! The requested stock(s) could not be found."
~~~~~~~~~~~~~~~~~~~~~~~~
*Unrecoverable Errors are handled by the "log" package.*