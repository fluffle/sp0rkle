# sp0rkle User Guide: How to use the Bot

This guide provides an overview of the commands and features available to users of the sp0rkle IRC bot.

## 1. Factoids

Factoids allow the bot to remember and recall information.

### Adding Factoids
There are two ways to teach the bot something:
- `botNick: <key> := <value>` - Stores `<value>` for `<key>`.
- `botNick: <key> :is <value>` - Stores `<key> is <value>` for `<key>`.

### Retrieving Factoids
- Simply type the `<key>` in the channel, and the bot will respond with a random value associated with that key.
- `botNick: literal <key>` - Prints all known values for a key. If there are many, the bot will ask you to do this in a private message.

### Managing Factoids
These commands operate on the **last factoid triggered** in the channel:
- `botNick: forget that` or `delete that` - Deletes the last triggered value.
- `botNick: replace that with <new value>` - Replaces the last value.
- `botNick: chance of that is <N>%` - Sets the probability (0-100%) that the bot will respond to this key.

## 2. Karma

Karma tracks the "score" of things or people.

- `<thing>++` - Increases the karma of `<thing>`.
- `<thing>--` - Decreases the karma of `<thing>`.
- `(a thing with spaces)++` - Use parentheses for things with spaces.
- `!karma <thing>` - Checks the current karma score of `<thing>`.

## 3. Quotes

The quote driver stores and retrieves funny or memorable lines.

- `!qadd <quote>` - Adds a new quote.
- `!quote <regex>` - Finds quotes matching the regular expression.
- `!quote #<ID>` - Retrieves a specific quote by ID.
- `!qdel #<ID>` - Deletes a quote (requires appropriate permissions).

## 4. Reminders and Tells

sp0rkle can remind you of things later or leave messages for absent users.

### Reminders
- `!remind <nick> <message> in <duration>` (e.g., `!remind me coffee in 10m`).
- `!remind <nick> <message> at <time>` (e.g., `!remind me meeting at 2pm`).
- `!snooze [duration]` - Delays a recently triggered reminder.

### Tells
- `!tell <nick> <message>` - Stores a message that the bot will deliver when `<nick>` next speaks or joins the channel.

### Timezones
- `!my timezone is <zone>` - Sets your preferred timezone (e.g., `Europe/London`) for reminders.
- `!forget my timezone` - Clears your timezone setting.

## 5. Decisions

Need help making a choice?
- `!choose <option1> or <option2> or ...` - The bot will pick one for you.
- `!shuffle <option1>, <option2>, ...` - The bot will randomize the order.

## 6. Miscellaneous

- `!calc <expression>` - Basic calculator.
- `!urlfind <regex>` - Search for URLs previously mentioned in the channel.
- `!ignore <nick>` / `!unignore <nick>` - Make the bot ignore or stop ignoring a specific user.
