# Quote Driver

The quote driver implements quote storage and retrieval.

## Adding quotes
Quotes can be added in three ways (must address the bot by nick):
- `botNick: qadd <quote>`
- `botNick: quote add <quote>`
- `botNick: add quote <quote>`

## Removing quotes
Quotes are deleted by quote ID (must address the bot by nick):
- `botNick: qdel <qid>`
- `botNick: quote del <qid>`
- `botNick: del quote <qid>`

## Retrieving quotes by ID
A specific quote may be looked up by its ID:
- `botNick: quote #<qid>`

## Retrieving quotes
Quotes may also be retrieved at random, with an optional case-insensitive regular expression:
- `botNick: quote <regex>`
- `botNick: quote`

## Plugin Syntax
The quote driver provides plugin functionality for the factoid driver, allowing you to embed quotes in factoid results. Add a factoid containing a `<plugin=quote>` directive:
- `botNick: quote fact := <plugin=quote>`
- `botNick: quote fact := <plugin=quote #quote ID>`
- `botNick: quote fact := <plugin=quote regex>`
